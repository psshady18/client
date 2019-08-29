package contacts

import (
	"context"

	"github.com/keybase/client/go/encrypteddb"
	"github.com/keybase/client/go/libkb"
	"github.com/keybase/client/go/protocol/keybase1"
)

type ContactResolution struct {
	Description  string
	ResolvedUser keybase1.User
}

func getKeyFn(mctx libkb.MetaContext) func(ctx context.Context) ([32]byte, error) {
	keyFn := func(ctx context.Context) ([32]byte, error) {
		return encrypteddb.GetSecretBoxKey(ctx, mctx.G(),
			encrypteddb.DefaultSecretUI,
			libkb.EncryptionReasonContactsResolvedServer,
			"encrypting contact data for server")
	}
	return keyFn
}

func encryptContactBlob(mctx libkb.MetaContext, res ContactResolution) ([]byte,
	error) {
	return encrypteddb.EncodeBox(mctx.Ctx(), res, getKeyFn(mctx))
}
func DecryptContactBlob(mctx libkb.MetaContext,
	contactResBlob []byte) (res ContactResolution, err error) {
	err = encrypteddb.DecodeBox(mctx.Ctx(), contactResBlob,
		getKeyFn(mctx), &res)
	return res, err
}

// TODO: actually call this
func SendContactResolutionToServer(mctx libkb.MetaContext,
	resolutions []ContactResolution) error {

	type resolvedArg struct {
		ResolvedContactBlob []byte `json:"b"`
	}

	type resolvedRes struct {
		libkb.AppStatusEmbed
		Success bool
	}

	args := make([]resolvedArg, 0, len(resolutions))
	for _, res := range resolutions {
		blob, err := encryptContactBlob(mctx, res)
		if err != nil {
			return err
		}
		args = append(args, resolvedArg{ResolvedContactBlob: blob})
	}
	payload := make(libkb.JSONPayload)
	payload["resolved_contact_blobs"] = args

	arg := libkb.APIArg{
		Endpoint:    "contacts/resolved",
		JSONPayload: payload,
		SessionType: libkb.APISessionTypeREQUIRED,
	}
	// TODO: determine if we care about this result at all
	var resp resolvedRes
	return mctx.G().API.PostDecode(mctx, arg, &resp)
}
