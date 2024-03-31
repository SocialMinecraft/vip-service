package eventHandlers

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"vip-service/database"
	"vip-service/github.com/SocialMinecraft/protos/vip"
)

func Get(db *database.Db, msg *nats.Msg) error {
	e := &vip.GetRequest{}
	if err := proto.Unmarshal(msg.Data, e); err != nil {
		return err
	}

	if e.UserId == nil && e.MinecraftUuid == nil {
		return sendGetResp(msg, nil)
	}

	if e.UserId != nil {
		panic("todo")
	}

	if e.MinecraftUuid != nil {
		panic("todo")
	}

	return nil
}

func sendGetResp(msg *nats.Msg, membership *vip.Membership) error {
	buf, err := proto.Marshal(&vip.GetResponse{
		HasMembership: membership != nil,
		Membership:    membership,
	})
	if err != nil {
		return err
	}
	return msg.Respond(buf)
}
