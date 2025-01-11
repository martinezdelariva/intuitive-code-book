package infra

import (
	"database/sql"
	"errors"
	"fmt"

	"chapter7/activation"

	_ "github.com/mattn/go-sqlite3"
)

func NewStoreSQLite(path string) (*StoreSQLite, error) {
	db, err := sql.Open("sqlite3", path+"?_journal=wal")
	if err != nil {
		return nil, err
	}
	return &StoreSQLite{db: db}, nil
}

type StoreSQLite struct {
	db *sql.DB
}

func (s *StoreSQLite) Find(accountID int) (activation.Account, error) {
	var (
		trialExtended bool
		expiredAt     int64
	)
	row := s.db.QueryRow("select trial_extended,expired_at from accounts where id=?", accountID)
	if err := row.Scan(&trialExtended, &expiredAt); errors.Is(err, sql.ErrNoRows) {
		return activation.Account{}, activation.ErrAccountNotFound
	} else if err != nil {
		return activation.Account{}, fmt.Errorf("%w: %w", activation.ErrInternal, err)
	}

	return activation.Account{
		ID:            accountID,
		TrialExtended: trialExtended,
		ExpireAt:      expiredAt,
	}, nil
}

func (s *StoreSQLite) Save(account activation.Account) error {
	_, err := s.db.Exec("update accounts set trial_extended=?, expired_at=?  where id=?", account.TrialExtended, account.ExpireAt, account.ID)
	if err != nil {
		return fmt.Errorf("%w: %w", activation.ErrInternal, err)
	}
	return nil
}

func (s *StoreSQLite) Close() error {
	return s.db.Close()
}
