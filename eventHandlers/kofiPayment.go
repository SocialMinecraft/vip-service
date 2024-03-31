package eventHandlers

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"time"
	"vip-service/database"
	"vip-service/github.com/SocialMinecraft/protos/kofi"
)

// comes in on kofi.payment
func KofiPayment(db *database.Db, msg *nats.Msg) error {
	e := &kofi.Payment{}
	if err := proto.Unmarshal(msg.Data, e); err != nil {
		return err
	}

	err := db.AddMembership(database.Membership{
		Email: e.Email,
		Start: time.Now(),
		End:   time.Now().AddDate(0, 1, 0),
	})

	return err
}
