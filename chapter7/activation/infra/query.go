package infra

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewQuerySQLite(path string) (*QueryDB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal=wal&mode=ro")
	if err != nil {
		return nil, err
	}
	return &QueryDB{db: db}, nil
}

type QueryDB struct {
	db *sql.DB
}

func (q *QueryDB) AccountIDsAboutExpire(before int64) ([]int, error) {
	query := "select id from accounts a where a.expired_at < ?"
	rows, err := q.db.Query(query, before)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ret []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		ret = append(ret, id)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return ret, nil
}

func (q *QueryDB) TrialExtendedActivationRate() (float64, error) {
	var (
		extended  int
		activated int
	)

	query := "select count(*) from accounts where trial_extended=true"

	err := q.db.QueryRow(query).Scan(&extended)
	if err != nil {
		return 0, err
	}

	if extended == 0 {
		return 0, nil
	}

	query = `select count(*) from accounts a 
    		 inner join payment_methods pm on pm.account_id =a.id  
          	 where a.trial_extended=true`

	err = q.db.QueryRow(query).Scan(&activated)
	if err != nil {
		return 0, err
	}

	return (float64(activated) * 100) / float64(extended), nil
}

func (q *QueryDB) Close() error {
	return q.db.Close()
}
