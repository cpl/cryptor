package noise

import (
	"testing"

	"cpl.li/go/cryptor/internal/crypt/ppk"

	"github.com/stretchr/testify/assert"
)

func TestHandshakeFlow(t *testing.T) {
	t.Parallel()

	var (
		sHandshake Handshake
		rHandshake Handshake
		sSec       ppk.PrivateKey
		sPub       ppk.PublicKey
		rSec       ppk.PrivateKey
		rPub       ppk.PublicKey
		sPubOut    ppk.PublicKey
		sPubEnc    EncryptedKey
		enc        EncryptedNothing
		sSend      [ppk.KeySize]byte
		sRecv      [ppk.KeySize]byte
		rSend      [ppk.KeySize]byte
		rRecv      [ppk.KeySize]byte
	)

	// generate keys
	assert.NoError(t, ppk.NewPrivateKey(&sSec))
	assert.NoError(t, sSec.PublicKey(&sPub))
	assert.NoError(t, ppk.NewPrivateKey(&rSec))
	assert.NoError(t, rSec.PublicKey(&rPub))

	// init sender
	assert.NoError(t, sHandshake.InitializeSender(&rPub))
	assert.NoError(t, sHandshake.Exchange(&sPub, &sPubEnc))

	// generate sender temp pub
	sPubTmp := sHandshake.PublicKey()

	// abstract away how sender public temp key and sender pub encrypted
	// are sent to recipient

	// receive sender keys and init recipient
	assert.NoError(t, rHandshake.InitializeRecipient(&rSec, &sPubTmp))
	assert.NoError(t, rHandshake.Exchange(&sPubOut, &sPubEnc))

	// check pub matches after encryption and decryption
	assert.Equal(t, sPub, sPubOut)

	// recipient prepares response
	rHandshake.PrepareRecipientResponse(&sPubTmp, &sPub, &enc)
	rPubTmp := rHandshake.PublicKey()

	// abstract away how the recipient sends the pub temp and the enc nothing

	// receive recipient response
	assert.NoError(t, sHandshake.ConsumeRecipientResponse(&sSec, &rPubTmp, &enc))

	// finalize on each side
	assert.NoError(t, rHandshake.Finalize(&rSend, &rRecv))
	assert.NoError(t, sHandshake.Finalize(&sSend, &sRecv))

	assert.Equal(t, rSend, sRecv)
	assert.Equal(t, rRecv, sSend)
}
