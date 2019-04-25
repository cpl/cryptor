package mwords_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/mwords"
	"cpl.li/go/cryptor/tests"
)

var invalidMnemonics = []string{
	"", "a", "a b c", "a b c d", "0 1 2", "abandon", "abandon abandon",
	"abandon abandon abandon abandon", "abandon abandon abandon 0",
	"abandon 0 abandon", "abandon abandon .", "zoo abandon zooo", " ",
	"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo",
	"dog cat pig dog cat pig dog cat pig dog cat cockamouse",
}

func TestInvalidMnemonicFromString(t *testing.T) {
	t.Parallel()

	for _, str := range invalidMnemonics {
		_, err := mwords.MnemonicFromString(str)
		if err == nil {
			t.Errorf("validated invalid mnemonic, %s\n", str)
		}
	}
}

func TestInvalidMnemonicSentenceIsValid(t *testing.T) {
	t.Parallel()

	var ms mwords.MnemonicSentence
	for _, str := range invalidMnemonics {
		ms = strings.Fields(str)
		if ms.IsValid() {
			t.Errorf("validated invalid mnemonic, %s\n", str)
		}
	}
}

func TestInvalidEntropyFromString(t *testing.T) {
	t.Parallel()

	for _, str := range invalidMnemonics {
		if _, err := mwords.EntropyFromString(str); err == nil {
			t.Errorf("validated invalid mnemonic, %s\n", str)
		}
	}
}

func TestInvalidEntropyFromMnemonic(t *testing.T) {
	t.Parallel()

	var ms mwords.MnemonicSentence
	for _, str := range invalidMnemonics {
		ms = strings.Fields(str)
		if _, err := mwords.EntropyFromMnemonic(ms); err == nil {
			t.Errorf("validated invalid mnemonic, %s\n", str)
		}
	}
}

func TestInvalidEntropyToMnemonic(t *testing.T) {
	t.Parallel()

	invalidEntropy := [][]byte{
		nil,
		[]byte{},
		crypt.RandomBytes(1),
		crypt.RandomBytes(10),
		crypt.RandomBytes(17),
		crypt.RandomBytes(21),
		crypt.RandomBytes(30),
		crypt.RandomBytes(33),
		crypt.RandomBytes(100),
		crypt.RandomBytes(1000),
	}

	for _, ent := range invalidEntropy {
		if _, err := mwords.EntropyToMnemonic(ent); err == nil {
			t.Errorf("generated mnemonic from invalid entropy, %v\n", ent)
		}
	}
}

func TestEntropyToMnemonic(t *testing.T) {
	t.Parallel()

	for _, vector := range testVectors {
		entropy, err := hex.DecodeString(vector.entropy)
		tests.AssertNil(t, err)

		mnemonic, err := mwords.EntropyToMnemonic(entropy)
		tests.AssertNil(t, err)

		tests.AssertEqual(t, vector.mnemonic, mnemonic.String(),
			"mnemonic mismatch")
	}
}

func TestEntropyFromMnemonic(t *testing.T) {
	t.Parallel()

	for _, vector := range testVectors {
		// extract and validate sentence
		sentence, err := mwords.MnemonicFromString(vector.mnemonic)
		tests.AssertNil(t, err)

		// extract initial entropy
		entropy, err := mwords.EntropyFromMnemonic(sentence)
		if err != nil {
			t.Error(err)
			t.Error(sentence)
			continue
		}

		// validate expected entropy
		tests.AssertEqual(t, hex.EncodeToString(entropy), vector.entropy,
			"entropy mismatch")
	}
}

func TestEntropyFromMnemonicInvalidChecksum(t *testing.T) {
	t.Parallel()

	_, err := mwords.EntropyFromString(
		"abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon yellow")
	if err == nil {
		t.Fatal("generated entropy from mnemonic with invalid checksum")
	}
}

type vector struct {
	entropy  string
	mnemonic string
}

// taken from https://github.com/tyler-smith/go-bip39/blob/master/bip39_test.go
// as seen on https://github.com/trezor/python-mnemonic/blob/master/vectors.json
var testVectors = []vector{
	{
		entropy:  "00000000000000000000000000000000",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon about",
	},
	{
		entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		mnemonic: "legal winner thank year wave sausage worth useful legal winner thank yellow",
	},
	{
		entropy:  "80808080808080808080808080808080",
		mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage above",
	},
	{
		entropy:  "ffffffffffffffffffffffffffffffff",
		mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo wrong",
	},
	{
		entropy:  "000000000000000000000000000000000000000000000000",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon agent",
	},
	{
		entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal will",
	},
	{
		entropy:  "808080808080808080808080808080808080808080808080",
		mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter always",
	},
	{
		entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffff",
		mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo when",
	},
	{
		entropy:  "0000000000000000000000000000000000000000000000000000000000000000",
		mnemonic: "abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon abandon art",
	},
	{
		entropy:  "7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f7f",
		mnemonic: "legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth useful legal winner thank year wave sausage worth title",
	},
	{
		entropy:  "8080808080808080808080808080808080808080808080808080808080808080",
		mnemonic: "letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic avoid letter advice cage absurd amount doctor acoustic bless",
	},
	{
		entropy:  "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		mnemonic: "zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
	},
	{
		entropy:  "77c2b00716cec7213839159e404db50d",
		mnemonic: "jelly better achieve collect unaware mountain thought cargo oxygen act hood bridge",
	},
	{
		entropy:  "b63a9c59a6e641f288ebc103017f1da9f8290b3da6bdef7b",
		mnemonic: "renew stay biology evidence goat welcome casual join adapt armor shuffle fault little machine walk stumble urge swap",
	},
	{
		entropy:  "3e141609b97933b66a060dcddc71fad1d91677db872031e85f4c015c5e7e8982",
		mnemonic: "dignity pass list indicate nasty swamp pool script soccer toe leaf photo multiply desk host tomato cradle drill spread actor shine dismiss champion exotic",
	},
	{
		entropy:  "0460ef47585604c5660618db2e6a7e7f",
		mnemonic: "afford alter spike radar gate glance object seek swamp infant panel yellow",
	},
	{
		entropy:  "72f60ebac5dd8add8d2a25a797102c3ce21bc029c200076f",
		mnemonic: "indicate race push merry suffer human cruise dwarf pole review arch keep canvas theme poem divorce alter left",
	},
	{
		entropy:  "2c85efc7f24ee4573d2b81a6ec66cee209b2dcbd09d8eddc51e0215b0b68e416",
		mnemonic: "clutch control vehicle tonight unusual clog visa ice plunge glimpse recipe series open hour vintage deposit universe tip job dress radar refuse motion taste",
	},
	{
		entropy:  "eaebabb2383351fd31d703840b32e9e2",
		mnemonic: "turtle front uncle idea crush write shrug there lottery flower risk shell",
	},
	{
		entropy:  "7ac45cfe7722ee6c7ba84fbc2d5bd61b45cb2fe5eb65aa78",
		mnemonic: "kiss carry display unusual confirm curtain upgrade antique rotate hello void custom frequent obey nut hole price segment",
	},
	{
		entropy:  "4fa1a8bc3e6d80ee1316050e862c1812031493212b7ec3f3bb1b08f168cabeef",
		mnemonic: "exile ask congress lamp submit jacket era scheme attend cousin alcohol catch course end lucky hurt sentence oven short ball bird grab wing top",
	},
	{
		entropy:  "18ab19a9f54a9274f03e5209a2ac8a91",
		mnemonic: "board flee heavy tunnel powder denial science ski answer betray cargo cat",
	},
	{
		entropy:  "18a2e1d81b8ecfb2a333adcb0c17a5b9eb76cc5d05db91a4",
		mnemonic: "board blade invite damage undo sun mimic interest slam gaze truly inherit resist great inject rocket museum chief",
	},
	{
		entropy:  "15da872c95a13dd738fbf50e427583ad61f18fd99f628c417a61cf8343c90419",
		mnemonic: "beyond stage sleep clip because twist token leaf atom beauty genius food business side grid unable middle armed observe pair crouch tonight away coconut",
	},
}
