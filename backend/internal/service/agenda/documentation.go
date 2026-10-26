package agenda

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	agendaDto "github.com/ahmadammarm/amami/backend/internal/dto/agenda"
	"github.com/ahmadammarm/amami/backend/pkg/apperror"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	maxFilesPerUpload = 10
	sniffLength       = 512
)

// allowedFileTypes maps a detected MIME type to the extension it is stored with.
// The type is sniffed from the file content, never taken from the client.
var allowedFileTypes = map[string]string{
	"image/jpeg":      ".jpg",
	"image/png":       ".png",
	"image/webp":      ".webp",
	"application/pdf": ".pdf",
}

// UploadFile is one file from a multipart upload.
type UploadFile struct {
	Name string
	Size int64
	Open func() (io.ReadCloser, error)
}

func (s *agendaService) GetDocumentation(agendaID uuid.UUID) ([]agendaDto.DocumentationResponse, error) {
	if _, err := s.findAgenda(agendaID); err != nil {
		return nil, err
	}

	docs, err := s.repo.GetDocumentation(agendaID)
	if err != nil {
		return nil, apperror.Internal("failed to fetch documentation", err)
	}
	return toDocumentationResponses(docs), nil
}

// UploadDocumentation validates every file before storing any of them, then
// writes the files and their rows. If the database write fails, the stored
// files are removed so no orphans are left on disk.
func (s *agendaService) UploadDocumentation(agendaID uuid.UUID, files []UploadFile, userID uuid.UUID) ([]agendaDto.DocumentationResponse, error) {
	if _, err := s.findAgenda(agendaID); err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, apperror.Validation("at least one file is required")
	}
	if len(files) > maxFilesPerUpload {
		return nil, apperror.Validation(fmt.Sprintf("at most %d files can be uploaded at once", maxFilesPerUpload))
	}
	for _, f := range files {
		if f.Size > s.maxUploadBytes {
			return nil, apperror.Validation(fmt.Sprintf("%s exceeds the %d MB limit", f.Name, s.maxUploadBytes>>20))
		}
	}

	ctx := context.Background()
	docs := make([]domain.AgendaDocumentation, 0, len(files))
	storedKeys := make([]string, 0, len(files))

	for _, f := range files {
		doc, key, err := s.storeFile(ctx, agendaID, f, userID)
		if err != nil {
			s.removeStored(ctx, storedKeys)
			return nil, err
		}
		storedKeys = append(storedKeys, key)
		docs = append(docs, *doc)
	}

	if err := s.repo.CreateDocumentation(docs); err != nil {
		s.removeStored(ctx, storedKeys)
		return nil, apperror.Internal("failed to save documentation", err)
	}
	return toDocumentationResponses(docs), nil
}

func (s *agendaService) storeFile(ctx context.Context, agendaID uuid.UUID, f UploadFile, userID uuid.UUID) (*domain.AgendaDocumentation, string, error) {
	src, err := f.Open()
	if err != nil {
		return nil, "", apperror.Internal("failed to read upload", err)
	}
	defer src.Close()

	buffered := bufio.NewReaderSize(src, sniffLength)
	head, err := buffered.Peek(sniffLength)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, "", apperror.Internal("failed to read upload", err)
	}

	contentType := http.DetectContentType(head)
	ext, ok := allowedFileTypes[contentType]
	if !ok {
		return nil, "", apperror.Validation(fmt.Sprintf("%s is not an allowed file type (jpeg, png, webp, pdf)", f.Name))
	}

	docID := uuid.New()
	key := fmt.Sprintf("agenda/%s/%s%s", agendaID, docID, ext)
	url, err := s.storage.Save(ctx, key, io.LimitReader(buffered, s.maxUploadBytes))
	if err != nil {
		return nil, "", apperror.Internal("failed to store file", err)
	}

	return &domain.AgendaDocumentation{
		ID:         docID,
		AgendaID:   agendaID,
		StorageKey: key,
		FileURL:    url,
		FileName:   f.Name,
		FileType:   contentType,
		FileSize:   f.Size,
		UploadedBy: &userID,
		UploadedAt: time.Now(),
	}, key, nil
}

func (s *agendaService) DeleteDocumentation(agendaID, docID uuid.UUID) error {
	doc, err := s.repo.GetDocumentationByID(agendaID, docID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.NotFound("documentation not found")
		}
		return apperror.Internal("failed to fetch documentation", err)
	}

	if err := s.repo.DeleteDocumentation(doc.ID); err != nil {
		return apperror.Internal("failed to delete documentation", err)
	}
	// The row is gone, so a failed file removal only leaves an orphan on disk.
	if err := s.storage.Delete(context.Background(), doc.StorageKey); err != nil {
		zap.L().Warn("failed to remove documentation file", zap.String("key", doc.StorageKey), zap.Error(err))
	}
	return nil
}

func (s *agendaService) removeStored(ctx context.Context, keys []string) {
	for _, key := range keys {
		if err := s.storage.Delete(ctx, key); err != nil {
			zap.L().Warn("failed to clean up uploaded file", zap.String("key", key), zap.Error(err))
		}
	}
}

func toDocumentationResponses(docs []domain.AgendaDocumentation) []agendaDto.DocumentationResponse {
	res := make([]agendaDto.DocumentationResponse, 0, len(docs))
	for i := range docs {
		res = append(res, agendaDto.ToDocumentationResponse(&docs[i]))
	}
	return res
}
