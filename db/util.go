package db

import (
	"database/sql"
	"errors"
	"time"
)

func (r *Db) IsVip(uuid string) (bool, error) {

	re := false
	var expire time.Time

	err := r.db.QueryRow(
		`
SELECT
    memberships.expire 
FROM 
    memberships 
        INNER JOIN accounts a on memberships.user_id = a.user_id
        WHERE a.minecraft_uuid = $1;
`,
		uuid,
	).Scan(
		&expire,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return re, nil
	}
	if err != nil {
		return re, err
	}

	re = time.Now().After(expire.AddDate(0, 0, 1))

	return re, nil
}
