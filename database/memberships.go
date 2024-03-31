package database

import (
	"database/sql"
	"errors"
	"time"
)

type Membership struct {
	Id     int
	UserId *string
	Email  string
	Start  time.Time
	End    time.Time
}

func (r *Db) AddMembership(membership Membership) error {

	// https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-upsert/

	_, err := r.db.Query(
		`
INSERT INTO 
    memberships 
    (email, user_id, start, expire) 
VALUES 
    ($1, $2, $3, $4) 
ON CONFLICT (email) DO UPDATE 
    SET
        expire = $4
`,
		membership.Email,
		membership.UserId,
		membership.Start,
		membership.End,
	)
	return err
}

func (r *Db) ClaimMembership(membershipId int, userId string) error {
	_, err := r.db.Query(
		`
UPDATE 
    memberships 
SET 
    user_id = $2,
WHERE
	id = $1
`,
		membershipId,
		userId,
	)
	return err
}

func (r *Db) GetMembership(userId string) (*Membership, error) {

	var re Membership

	err := r.db.QueryRow(
		`
SELECT
    id, email, user_id, start, expire
FROM 
    memberships
WHERE 
    user_id = $1;
`,
		userId,
	).Scan(
		&re.Id,
		&re.Email,
		&re.UserId,
		&re.Start,
		&re.End,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &re, err
}

func (r *Db) GetMembershipByEmail(email string) (*Membership, error) {

	var re Membership

	err := r.db.QueryRow(
		`
SELECT
    id, email, user_id, start, expire
FROM 
    memberships
WHERE 
    email = $1;
`,
		email,
	).Scan(
		&re.Id,
		&re.Email,
		&re.UserId,
		&re.Start,
		&re.End,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &re, err
}
