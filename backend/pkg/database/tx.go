package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// LockForUpdate loads the row of type T with the given primary key and holds
// an exclusive row lock (SELECT ... FOR UPDATE) until the surrounding
// transaction ends. Concurrent writers to the same row wait, so a rule checked
// against the returned row still holds when the transaction writes.
//
// It must be called on a transaction (inside db.Transaction); on a plain
// connection the lock is released as soon as the query finishes.
// A missing or soft-deleted row returns gorm.ErrRecordNotFound.
func LockForUpdate[T any](tx *gorm.DB, id any) (*T, error) {
	return lockRow[T](tx, id, "UPDATE")
}

// LockForShare is LockForUpdate with a shared lock (SELECT ... FOR SHARE):
// other transactions may read and share-lock the row, but cannot change it
// until this transaction ends. Use it to keep a row stable while acting on it
// without modifying it.
func LockForShare[T any](tx *gorm.DB, id any) (*T, error) {
	return lockRow[T](tx, id, "SHARE")
}

func lockRow[T any](tx *gorm.DB, id any, strength string) (*T, error) {
	var row T
	if err := tx.Clauses(clause.Locking{Strength: strength}).First(&row, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ExpectRows reports a write that matched no row (for example because the
// row was deleted meanwhile) as gorm.ErrRecordNotFound instead of success.
func ExpectRows(res *gorm.DB) error {
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
