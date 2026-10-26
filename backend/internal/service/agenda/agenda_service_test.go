package agenda

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	agendaDto "github.com/ahmadammarm/amami/backend/internal/dto/agenda"
	agendaRepo "github.com/ahmadammarm/amami/backend/internal/repository/agenda"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fakeRepo struct {
	agendaRepo.AgendaRepository
	agendas map[uuid.UUID]*domain.Agenda
	created []domain.AgendaDocumentation
}

func (r *fakeRepo) GetByID(id uuid.UUID) (*domain.Agenda, error) {
	a, ok := r.agendas[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return a, nil
}

func (r *fakeRepo) DeleteAgenda(id uuid.UUID) error {
	if _, ok := r.agendas[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(r.agendas, id)
	return nil
}

func (r *fakeRepo) CreateDocumentation(docs []domain.AgendaDocumentation) error {
	r.created = append(r.created, docs...)
	return nil
}

type memStorage struct {
	files map[string][]byte
}

func (m *memStorage) Save(_ context.Context, key string, r io.Reader) (string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}
	m.files[key] = data
	return "/uploads/" + key, nil
}

func (m *memStorage) Delete(_ context.Context, key string) error {
	delete(m.files, key)
	return nil
}

func newTestService(agendas ...*domain.Agenda) (*fakeRepo, *memStorage, AgendaService) {
	repo := &fakeRepo{agendas: map[uuid.UUID]*domain.Agenda{}}
	for _, a := range agendas {
		repo.agendas[a.ID] = a
	}
	store := &memStorage{files: map[string][]byte{}}
	return repo, store, NewAgendaService(repo, store, 1<<20)
}

func uploadOf(name string, content []byte) UploadFile {
	return UploadFile{
		Name: name,
		Size: int64(len(content)),
		Open: func() (io.ReadCloser, error) { return io.NopCloser(bytes.NewReader(content)), nil },
	}
}

func TestMissingAgendaIsNotFound(t *testing.T) {
	_, _, svc := newTestService()
	missing := uuid.New()

	_, err := svc.UpdateAgenda(missing, agendaDto.AgendaRequest{Title: "x", Status: domain.AgendaScheduled})
	if !apperror.IsKind(err, apperror.KindNotFound) {
		t.Fatalf("update: got %v, want not found", err)
	}
	if err := svc.DeleteAgenda(missing); !apperror.IsKind(err, apperror.KindNotFound) {
		t.Fatalf("delete: got %v, want not found", err)
	}
	if _, err := svc.UploadDocumentation(missing, []UploadFile{uploadOf("a.pdf", []byte("%PDF-1.4"))}, uuid.New()); !apperror.IsKind(err, apperror.KindNotFound) {
		t.Fatalf("upload: got %v, want not found", err)
	}
}

func TestUploadDocumentation(t *testing.T) {
	agenda := &domain.Agenda{ID: uuid.New(), Title: "Kajian"}
	pdf := []byte("%PDF-1.4\n%fake pdf body")
	png := []byte("\x89PNG\r\n\x1a\n\x00\x00\x00\rIHDR")

	t.Run("stores allowed files under the agenda", func(t *testing.T) {
		repo, store, svc := newTestService(agenda)

		res, err := svc.UploadDocumentation(agenda.ID, []UploadFile{uploadOf("lpj.pdf", pdf), uploadOf("foto.png", png)}, uuid.New())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(res) != 2 || len(repo.created) != 2 || len(store.files) != 2 {
			t.Fatalf("expected 2 stored docs, got res=%d rows=%d files=%d", len(res), len(repo.created), len(store.files))
		}
		if res[0].FileType != "application/pdf" || res[1].FileType != "image/png" {
			t.Fatalf("unexpected detected types: %s, %s", res[0].FileType, res[1].FileType)
		}
		if !bytes.Equal(store.files[repo.created[0].StorageKey], pdf) {
			t.Fatalf("stored content differs from upload")
		}
	})

	t.Run("rejects disguised executables and cleans up earlier files", func(t *testing.T) {
		repo, store, svc := newTestService(agenda)
		exe := []byte("MZ\x90\x00\x03\x00\x00\x00")

		_, err := svc.UploadDocumentation(agenda.ID, []UploadFile{uploadOf("ok.pdf", pdf), uploadOf("photo.jpg", exe)}, uuid.New())
		if !apperror.IsKind(err, apperror.KindValidation) {
			t.Fatalf("got %v, want validation error", err)
		}
		if len(store.files) != 0 || len(repo.created) != 0 {
			t.Fatalf("rejected upload left files=%d rows=%d", len(store.files), len(repo.created))
		}
	})

	t.Run("rejects oversized files before storing", func(t *testing.T) {
		_, store, svc := newTestService(agenda)
		big := UploadFile{Name: "big.pdf", Size: 2 << 20, Open: uploadOf("big.pdf", pdf).Open}

		_, err := svc.UploadDocumentation(agenda.ID, []UploadFile{big}, uuid.New())
		if !apperror.IsKind(err, apperror.KindValidation) {
			t.Fatalf("got %v, want validation error", err)
		}
		if len(store.files) != 0 {
			t.Fatalf("oversized file was stored")
		}
	})

	t.Run("requires at least one file", func(t *testing.T) {
		_, _, svc := newTestService(agenda)
		if _, err := svc.UploadDocumentation(agenda.ID, nil, uuid.New()); !apperror.IsKind(err, apperror.KindValidation) {
			t.Fatalf("got %v, want validation error", err)
		}
	})
}
