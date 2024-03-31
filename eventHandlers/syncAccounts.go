package eventHandlers

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"time"
	"vip-service/database"
	"vip-service/github.com/SocialMinecraft/protos"
	"vip-service/github.com/SocialMinecraft/protos/vip"
)

func changeAcconut(db *database.Db, msg *nats.Msg) error {
	e := &protos.MinecraftAccountChanged{}
	if err := proto.Unmarshal(msg.Data, e); err != nil {
		return err
	}

	var err error
	switch e.Change {
	case protos.MinecraftAccountChangeType_ADDED:
		err = db.AddAccount(database.Account{
			UserId:        e.UserId,
			MinecraftUuid: e.Account.MinecraftUuid,
		})
		break
	case protos.MinecraftAccountChangeType_REMOVED:
		err = db.RemoveAccountByUuid(e.Account.MinecraftUuid)
		break
	}

	return err
}

func syncAccounts(nc *nats.Conn, db *database.Db, msg *nats.Msg) error {
	var buf []byte
	var err error

	e := &vip.SyncRequest{}
	if err = proto.Unmarshal(msg.Data, e); err != nil {
		return err
	}

	buf, err = proto.Marshal(&protos.ListMinecraftAccountsRequest{UserId: e.UserId})
	if err != nil {
		return err
	}
	resp, err := nc.Request("accounts.minecraft.list", buf, time.Second*1)
	if err != nil {
		return err
	}
	list := &protos.ListMinecraftAccountsResponse{}
	if err = proto.Unmarshal(resp.Data); err != nil {
		return err
	}

	for _, account := range list.Accounts {
		if err = db.AddAccount(database.Account{UserId: list.UserId, MinecraftUuid: account.MinecraftUuid}); err != nil {
			return err
		}
	}

	// did we get zero accounts? error
	if len(list.Accounts) <= 0 {
		str := "No minecraft accounts found."
		buf, err = proto.Marshal(&vip.SyncResponse{
			Success:      false,
			ErrorMessage: &str,
		})
		if err != nil {
			return err
		}
		if err = msg.Respond(buf); err != nil {
			return err
		}
	}

	// Send back a success.
	buf, err = proto.Marshal(&vip.SyncResponse{
		Success: true,
	})
	if err != nil {
		return err
	}
	err = msg.Respond(buf)
	return nil
}
