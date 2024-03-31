package eventHandlers

import (
	"github.com/nats-io/nats.go"
	"vip-service/db"
)

func addAccount(db *db.Db, msg *nats.Msg) error {

	return nil
}

func removeAccount(db *db.Db, msg *nats.Msg) error {

	return nil
}

func syncAccounts(db *db.Db, msg *nats.Msg) error {

	return nil
}
