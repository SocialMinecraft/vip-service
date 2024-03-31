package eventHandlers

import (
	"github.com/nats-io/nats.go"
	"google.golang.org/protobuf/proto"
	"vip-service/database"
	"vip-service/github.com/SocialMinecraft/protos/vip"
)

// comes in on vip.claim
func Claim(db *database.Db, msg *nats.Msg) error {
	e := &vip.ClaimRequest{}
	if err := proto.Unmarshal(msg.Data, e); err != nil {
		return err
	}

	// Does the user already have a memebrishp?
	membership, err := db.GetMembership(e.UserId)
	if err != nil {
		return err
	}
	if membership != nil {
		return sendClaimResp(msg, "Membership already associated with account.")
	}

	// Does a membership exist for the email address?
	membership, err = db.GetMembershipByEmail(e.Email)
	if err != nil {
		return err
	}
	if membership == nil {
		return sendClaimResp(msg, "No membership found for the email.")
	}

	// Is the membership already claimed by someone else?
	if membership.UserId != nil {
		return sendClaimResp(msg, "Membership already claimed.")
	}

	// Save the membership
	err = db.ClaimMembership(membership.Id, e.UserId)
	if err != nil {
		return err
	}

	return sendClaimResp(msg, "")
}

func sendClaimResp(msg *nats.Msg, error_message string) error {
	buf, err := proto.Marshal(&vip.ClaimResponse{
		Success:      len(error_message) <= 0,
		ErrorMessage: &error_message,
	})
	if err != nil {
		return err
	}
	return msg.Respond(buf)
}
