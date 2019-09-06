// Copyright 2019 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package contacts

import (
	"github.com/keybase/client/go/encrypteddb"
	"github.com/keybase/client/go/kbtest"
	"github.com/keybase/client/go/libkb"
	"testing"

	"github.com/keybase/client/go/protocol/keybase1"
	"github.com/stretchr/testify/require"
)

func TestEncryptContactResolutionForServer(t *testing.T) {
	tc := libkb.SetupTest(t, "contacts", 2)
	defer tc.Cleanup()

	u, err := kbtest.CreateAndSignupFakeUser("t", tc.G)
	require.NoError(t, err)
	err = u.Login(tc.G)
	require.NoError(t, err)

	contact := ContactResolution{
		Description: "Jakob - (216) 555-2222",
		ResolvedUser: keybase1.User{
			Uid:      keybase1.UID(34),
			Username: "jakob223",
		},
	}
	enc, err := encryptContactBlob(tc.MetaContext(), contact)
	require.NoError(t, err)

	// Provision a new device to roll forward the PUK
	tc2 := libkb.SetupTest(t, "contacts", 2)
	defer tc2.Cleanup()
	kbtest.ProvisionNewDeviceKex(&tc, &tc2, u, libkb.DeviceTypeDesktop)

	// Require the original device can decrypt the contact
	dec, err := DecryptContactBlob(tc.MetaContext(), enc)
	require.NoError(t, err)
	require.Equal(t, contact, dec)

	// Require the new device can decrypt the contact
	dec, err = DecryptContactBlob(tc2.MetaContext(), enc)
	require.NoError(t, err)
	require.Equal(t, contact, dec)

	// Log in as a different user
	u2, err := kbtest.CreateAndSignupFakeUser("t", tc.G)
	require.NoError(t, err)
	err = u2.Login(tc.G)
	require.NoError(t, err)

	// New user should not be able to decrypt contact blob
	dec, err = DecryptContactBlob(tc.MetaContext(), enc)
	require.Equal(t, encrypteddb.ErrDecryptionFailed, err)
	require.NotEqual(t, contact, dec)

	// Reset the user and check that they can no longer decrypt
	kbtest.ResetAccount(tc2, u)
	dec, err = DecryptContactBlob(tc2.MetaContext(), enc)
	require.Error(t, err)
	require.NotEqual(t, contact, dec)
}
