package database

import (
	"database/sql"
	"errors"
	"time"
	"vip-service/github.com/SocialMinecraft/protos/vip"
)

type Membership struct {
	Id     int64
	UserId *string
	Email  string
	Start  time.Time
	End    time.Time
}

func (r *Membership) ToProto() vip.Membership {
	return vip.Membership{
		Id:     r.Id,
		Email:  r.Email,
		Start:  r.Start.Unix(),
		Expire: r.End.Unix(),
		UserId: r.UserId,
	}
}

func (r *Db) AddMembership(membership Membership) error {

	// https://www.postgresqltutorial.com/postgresql-tutorial/postgresql-upsert/

	_, err := r.db.Query(
		`
INSERT INTO 
    memberships 
    (email, user_id, start, expire) 
VALUES 
    (lower($1), $2, $3, $4) 
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

func (r *Db) ClaimMembership(membershipId int64, userId string) error {
	_, err := r.db.Query(
		`
UPDATE 
    memberships 
SET 
    user_id = $2
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
    lower(email) = lower($1);
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

func (r *Db) GetMembershipByMinecraft(uuid string) (*Membership, error) {

	var re Membership

	err := r.db.QueryRow(
		`
SELECT
    m.id, m.email, m.user_id, m.start, m.expire
FROM 
    memberships m
INNER JOIN accounts a on m.user_id = a.user_id
WHERE 
    a.minecraft_uuid = $1;
`,
		uuid,
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
