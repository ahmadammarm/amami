// Package storage abstracts where uploaded files live, so the local-disk
// implementation can later be swapped for object storage.
package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type Storage interface {
	// Save writes r under key and returns the public URL of the stored file.
	Save(ctx context.Context, key string, r io.Reader) (string, error)
	// Delete removes the file stored under key. Missing files are not an error.
	Delete(ctx context.Context, key string) error
}

// LocalStorage stores files under Root and serves them from URLPrefix
// (main.go mounts Root as a static route at that prefix).
type LocalStorage struct {
	Root      string
	URLPrefix string
}

func NewLocalStorage(root, urlPrefix string) *LocalStorage {
	return &LocalStorage{Root: root, URLPrefix: strings.TrimSuffix(urlPrefix, "/")}
}

func (s *LocalStorage) Save(ctx context.Context, key string, r io.Reader) (string, error) {
	fullPath, err := s.resolve(key)
	if err != nil {
		return "", err
	}
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", fmt.Errorf("create upload directory: %w", err)
	}

	f, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return "", fmt.Errorf("create file: %w", err)
	}
	if _, err := io.Copy(f, readerWithContext(ctx, r)); err != nil {
		f.Close()
		os.Remove(fullPath)
		return "", fmt.Errorf("write file: %w", err)
	}
	if err := f.Close(); err != nil {
		os.Remove(fullPath)
		return "", fmt.Errorf("close file: %w", err)
	}

	return s.URLPrefix + "/" + path.Clean(key), nil
}

func (s *LocalStorage) Delete(_ context.Context, key string) error {
	fullPath, err := s.resolve(key)
	if err != nil {
		return err
	}
	if err := os.Remove(fullPath); err != nil && !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("delete file: %w", err)
	}
	return nil
}

// resolve maps a key to a path inside Root and rejects keys that escape it.
func (s *LocalStorage) resolve(key string) (string, error) {
	cleaned := path.Clean("/" + key)
	if cleaned == "/" {
		return "", errors.New("empty storage key")
	}
	root, err := filepath.Abs(s.Root)
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(root, filepath.FromSlash(cleaned))
	if !strings.HasPrefix(fullPath, root+string(filepath.Separator)) {
		return "", errors.New("storage key escapes the upload directory")
	}
	return fullPath, nil
}

type ctxReader struct {
	ctx context.Context
	r   io.Reader
}

func (c ctxReader) Read(p []byte) (int, error) {
	if err := c.ctx.Err(); err != nil {
		return 0, err
	}
	return c.r.Read(p)
}

func readerWithContext(ctx context.Context, r io.Reader) io.Reader {
	return ctxReader{ctx: ctx, r: r}
}
