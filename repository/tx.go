package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type txRepository struct {
	db *sqlx.DB
}

func NewTxRepository(db *sqlx.DB) *txRepository {
	return &txRepository{db: db}
}

func (repo *txRepository) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	tx, err := repo.db.BeginTxx(ctx, nil)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (repo *txRepository) Commit(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Commit()
}

func (repo *txRepository) Rollback(tx *sqlx.Tx) error {
	if tx == nil {
		return nil
	}
	return tx.Rollback()
}
