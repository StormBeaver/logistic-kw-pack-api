package repo

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// lock for tx
func AcquireLockTx(ctx context.Context, tx *sqlx.Tx) (sql.Result, error) {
	res, err := tx.ExecContext(ctx, "SELECT pg_advisory_xact_lock(hashtext('txLock'))")
	return res, err
}

// lock for regular sql requests
func AcquireLock(ctx context.Context, db *sqlx.DB) (sql.Result, error) {
	res, err := db.ExecContext(ctx, "SELECT pg_advisory_lock(hashtext('dbLock'))")
	return res, err
}

// unlock for regular sql requests
func Unlock(ctx context.Context, db *sqlx.DB) (sql.Result, error) {
	res, err := db.ExecContext(ctx, "SELECT pg_advisory_Unlock(hashtext('dbLock'))")
	return res, err
}
