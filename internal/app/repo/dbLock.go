package repo

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

func AcquireTryLockTx(ctx context.Context, tx *sqlx.Tx, key string) (bool, error) {
	var isAcquired bool
	err := tx.GetContext(ctx, &isAcquired, fmt.Sprintf("SELECT pg_try_advisory_xact_lock(hashtext('%s'))", key))
	return isAcquired, err
}

func AcquireTryLock(ctx context.Context, db *sqlx.DB, key string) (bool, error) {
	var isAcquired bool
	err := db.GetContext(ctx, &isAcquired, fmt.Sprintf("SELECT pg_try_advisory_lock(hashtext('%s'))", key))
	return isAcquired, err
}
