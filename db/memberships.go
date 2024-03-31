package db

import "time"

type Membership struct {
	Id     int
	UserId *string
	Email  string
	start  time.Time
	end    time.Time
}

func (r *Db) AddMembership(membership Membership) error {
	// Insert or update...
	return nil
}

func (r *Db) GetMembership() (*Membership, error) {
	return nil, nil
}
