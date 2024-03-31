package eventHandlers

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"vip-service/database"
	"vip-service/github.com/SocialMinecraft/protos/vip"
)

// vip.get
func Get(db *database.Db, msg *nats.Msg) error {
	e := &vip.GetRequest{}
	if err := proto.Unmarshal(msg.Data, e); err != nil {
		return err
	}

	if e.UserId == nil && e.MinecraftUuid == nil {
		return sendGetResp(msg, nil)
	}

	if e.UserId != nil {
		membership, err := db.GetMembership(*e.UserId)
		if err != nil {
			return err
		}
		return sendGetResp(msg, membership)
	}

	if e.MinecraftUuid != nil {
		membership, err := db.GetMembershipByMinecraft(*e.MinecraftUuid)
		if err != nil {
			return err
		}
		return sendGetResp(msg, membership)
	}

	return nil
}

func sendGetResp(msg *nats.Msg, membership *database.Membership) error {
	var m *vip.Membership = nil
	if membership != nil {
		t := membership.ToProto()
		m = &t
	}

	buf, err := proto.Marshal(&vip.GetResponse{
		HasMembership: membership != nil,
		Membership:    m,
	})
	if err != nil {
		return err
	}
	return msg.Respond(buf)
}
