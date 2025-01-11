package activation

import (
	"errors"
)

// errors -------------------------------------------------------------------------------------------------------------

var (
	ErrAlreadyExtended         = errors.New("already extended")
	ErrTrialExtensionNotNeeded = errors.New("trial extension not needed")
	ErrAccountNotFound         = errors.New("account not found")
	ErrInternal                = errors.New("internal error")
)

// dependencies -------------------------------------------------------------------------------------------------------

type Store interface {
	Find(accountID int) (Account, error)
	Save(account Account) error
}

type Clock interface {
	Now() int64 // timestamp in seconds
}

type Monitoring interface {
	Monitor(name string, accountID int) func(errp *error)
}

// data ---------------------------------------------------------------------------------------------------------------

type Account struct {
	ID            int
	TrialExtended bool
	ExpireAt      int64
}

// write operations ---------------------------------------------------------------------------------------------------

type App struct {
	Store      Store
	Clock      Clock
	Monitoring Monitoring
}

func (a *App) ExtendTrial(accountID int) (acc Account, err error) {
	defer a.Monitoring.Monitor("extend_trial", accountID)(&err)

	acc, err = a.Store.Find(accountID)
	if err != nil {
		return Account{}, err
	}

	now := a.Clock.Now()

	acc, err = extendTrial(acc, now)
	if err != nil {
		return Account{}, err
	}

	if err := a.Store.Save(acc); err != nil {
		return Account{}, err
	}

	return acc, nil
}

func extendTrial(acc Account, now int64) (Account, error) {
	if acc.TrialExtended {
		return Account{}, ErrAlreadyExtended
	}

	if now+60*60*24 <= acc.ExpireAt {
		return Account{}, ErrTrialExtensionNotNeeded
	}

	acc = Account{ID: acc.ID, TrialExtended: true, ExpireAt: acc.ExpireAt + 60*60*24*7}
	return acc, nil
}
