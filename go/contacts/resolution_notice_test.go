// Copyright 2019 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package contacts

import (
	"github.com/keybase/client/go/libkb"
	"testing"

	"github.com/keybase/client/go/protocol/keybase1"
	"github.com/stretchr/testify/require"
)

func TestEncryptContactResolutionForServer(t *testing.T) {
	// TODO: how do I get this to be logged in? or do I need to move it to
	//  engine?
	tc := libkb.SetupTest(t, "contacts", 2)

	contact := ContactResolution{
		Description: "Jakob - (216) 555-2222",
		ResolvedUser: keybase1.User{
			Uid:      keybase1.UID(34),
			Username: "jakob223",
		},
	}
	enc, err := encryptContactBlob(tc.MetaContext(), contact)
	require.NoError(t, err)

	dec, err := DecryptContactBlob(tc.MetaContext(), enc)
	require.NoError(t, err)
	require.Equal(t, contact, dec)
}
