// Package pgtest runs tests against a real PostgreSQL server.
//
// Set AMAMI_TEST_DATABASE_URL to a database the tests may use, for example
// postgres://postgres:postgres@localhost:5432/amami_db?sslmode=disable.
// Each call to New creates a fresh schema, migrates it and drops it when the
// test ends, so existing data is never read or changed. Without the variable,
// tests that call New are skipped.
package pgtest

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/ahmadammarm/amami/backend/internal/domain"
	"github.com/ahmadammarm/amami/backend/pkg/database"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// EnvURL names the variable holding the test database connection string.
const EnvURL = "AMAMI_TEST_DATABASE_URL"

// New returns a connection whose search_path points at a new, migrated schema.
func New(t *testing.T) *gorm.DB {
	t.Helper()

	dsn := os.Getenv(EnvURL)
	if dsn == "" {
		t.Skipf("%s is not set; skipping PostgreSQL test", EnvURL)
	}

	admin, err := open(dsn)
	if err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	schema := "amami_test_" + randomHex(t)
	if err := admin.Exec(fmt.Sprintf("CREATE SCHEMA %q", schema)).Error; err != nil {
		t.Fatalf("create schema: %v", err)
	}

	// public stays on the path so shared extensions such as pg_trgm resolve.
	db, err := open(withSearchPath(dsn, schema+",public"))
	if err != nil {
		t.Fatalf("connect to test schema: %v", err)
	}
	t.Cleanup(func() {
		closeDB(db)
		if err := admin.Exec(fmt.Sprintf("DROP SCHEMA %q CASCADE", schema)).Error; err != nil {
			t.Errorf("drop schema %s: %v", schema, err)
		}
		closeDB(admin)
	})

	if err := database.Migrate(db); err != nil {
		t.Fatalf("migrate test schema: %v", err)
	}
	return db
}

// RunConcurrently calls fn n times on separate goroutines, releasing them all
// at once so they contend for the same rows, and returns each call's error.
func RunConcurrently(n int, fn func(i int) error) []error {
	errs := make([]error, n)
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := range n {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			errs[i] = fn(i)
		}()
	}
	close(start)
	wg.Wait()
	return errs
}

// CountNil reports how many of errs are nil.
func CountNil(errs []error) int {
	n := 0
	for _, err := range errs {
		if err == nil {
			n++
		}
	}
	return n
}

// CreateUser inserts a user (and the role it needs) for rows that reference a creator.
func CreateUser(t *testing.T, db *gorm.DB) *domain.User {
	t.Helper()
	role := &domain.Role{Name: "TEST_" + randomHex(t)}
	mustCreate(t, db, role)
	suffix := randomHex(t)
	user := &domain.User{
		Username:     "user_" + suffix,
		Email:        suffix + "@example.org",
		PasswordHash: "not-a-real-hash",
		RoleID:       role.ID,
	}
	mustCreate(t, db, user)
	return user
}

// CreateJamaah inserts a jamaah, optionally registered as mustahik.
func CreateJamaah(t *testing.T, db *gorm.DB, mustahik bool) *domain.Jamaah {
	t.Helper()
	j := &domain.Jamaah{FullName: "Jamaah " + randomHex(t), IsMustahik: mustahik}
	mustCreate(t, db, j)
	return j
}

func mustCreate(t *testing.T, db *gorm.DB, row any) {
	t.Helper()
	if err := db.Create(row).Error; err != nil {
		t.Fatalf("create %T: %v", row, err)
	}
}

func open(dsn string) (*gorm.DB, error) {
	db, err := database.Open(dsn)
	if err != nil {
		return nil, err
	}
	db.Logger = logger.Default.LogMode(logger.Silent)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	// Enough connections for every goroutine in a concurrency test to hold one.
	sqlDB.SetMaxOpenConns(50)
	return db, nil
}

func closeDB(db *gorm.DB) {
	if sqlDB, err := db.DB(); err == nil {
		sqlDB.Close()
	}
}

// withSearchPath adds search_path to a URL or keyword/value connection string.
func withSearchPath(dsn, path string) string {
	if u, err := url.Parse(dsn); err == nil && strings.Contains(dsn, "://") {
		q := u.Query()
		q.Set("search_path", path)
		u.RawQuery = q.Encode()
		return u.String()
	}
	return dsn + " search_path=" + path
}

func randomHex(t *testing.T) string {
	t.Helper()
	b := make([]byte, 6)
	if _, err := rand.Read(b); err != nil {
		t.Fatalf("random suffix: %v", err)
	}
	return hex.EncodeToString(b)
}
