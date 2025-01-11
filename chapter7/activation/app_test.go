package activation_test

import (
	"errors"
	"reflect"
	"testing"

	"chapter7/activation"
)

func TestExtendTrial(t *testing.T) {
	tests := []struct {
		name      string
		accountID int
		store     *StoreMock
		clock     activation.Clock
		want      activation.Account
		wantErr   error
		wantSave  []activation.Account
	}{
		{
			name:      "happy path",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{ID: 1, TrialExtended: false, ExpireAt: 946681200}, nil
				},
				saveFn: func(_ activation.Account) error {
					return nil
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 },
			},
			want:     activation.Account{ID: 1, TrialExtended: true, ExpireAt: 947286000},
			wantSave: []activation.Account{{ID: 1, TrialExtended: true, ExpireAt: 947286000}},
		},
		{
			name:      "already extended",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{ID: 1, TrialExtended: true, ExpireAt: 946681200}, nil
				},
				saveFn: func(_ activation.Account) error {
					return nil
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 },
			},
			wantErr: activation.ErrAlreadyExtended,
		},
		{
			name:      "expire within 24h",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{ID: 1, TrialExtended: false, ExpireAt: 946681200}, nil
				},
				saveFn: func(_ activation.Account) error {
					return nil
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 + 60*60*23 },
			},
			want:     activation.Account{ID: 1, TrialExtended: true, ExpireAt: 947286000},
			wantSave: []activation.Account{{ID: 1, TrialExtended: true, ExpireAt: 947286000}},
		},
		{
			name:      "expire after 24h",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{ID: 1, TrialExtended: false, ExpireAt: 946681200}, nil
				},
				saveFn: func(_ activation.Account) error {
					return nil
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 - 60*60*25 },
			},
			wantErr: activation.ErrTrialExtensionNotNeeded,
		},
		{
			name:      "account not found",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{}, activation.ErrAccountNotFound
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 },
			},
			wantErr: activation.ErrAccountNotFound,
		},
		{
			name:      "save error",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{ID: 1, TrialExtended: false, ExpireAt: 946681200}, nil
				},
				saveFn: func(_ activation.Account) error {
					return activation.ErrInternal
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 },
			},
			wantErr:  activation.ErrInternal,
			wantSave: []activation.Account{{ID: 1, TrialExtended: true, ExpireAt: 947286000}},
		},
		{
			name:      "save is called properly",
			accountID: 1,
			store: &StoreMock{
				findFn: func(_ int) (activation.Account, error) {
					return activation.Account{ID: 1, TrialExtended: false, ExpireAt: 946681200}, nil
				},
				saveFn: func(_ activation.Account) error {
					return nil
				},
			},
			clock: ClockMock{
				nowFn: func() int64 { return 946681200 },
			},
			want:     activation.Account{ID: 1, TrialExtended: true, ExpireAt: 947286000},
			wantSave: []activation.Account{{ID: 1, TrialExtended: true, ExpireAt: 947286000}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := activation.App{Store: tt.store, Clock: tt.clock, Monitoring: MonitoringNoop{}}
			got, err := app.ExtendTrial(tt.accountID)
			if (err != nil) && !errors.Is(tt.wantErr, err) {
				t.Errorf("ExtendTrial() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExtendTrial() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(tt.store.saveCalled, tt.wantSave) {
				t.Errorf("ExtendTrial() got = %v, want %v", tt.store.saveCalled, tt.wantSave)
			}
		})
	}
}

type StoreMock struct {
	saveCalled []activation.Account
	findFn     func(accountID int) (activation.Account, error)
	saveFn     func(account activation.Account) error
}

func (s *StoreMock) Find(accountID int) (activation.Account, error) {
	return s.findFn(accountID)
}

func (s *StoreMock) Save(account activation.Account) error {
	s.saveCalled = append(s.saveCalled, account)
	return s.saveFn(account)
}

type ClockMock struct {
	nowFn func() int64
}

func (c ClockMock) Now() int64 {
	return c.nowFn()
}

type MonitoringNoop struct{}

func (m MonitoringNoop) Monitor(_ string, _ int) func(_ *error) {
	return func(_ *error) {}
}
