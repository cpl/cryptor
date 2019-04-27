package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"syscall"
	"unicode"

	"cpl.li/go/cryptor/crypt/mwords"

	"cpl.li/go/cryptor/crypt"
	"cpl.li/go/cryptor/crypt/ppk"
	"github.com/fatih/color"
	"golang.org/x/crypto/ssh/terminal"
)

var keyPrivate ppk.PrivateKey

func commandKey(argc int, argv []string) error {
	// expect arguments
	if argc == 0 {
		return errors.New("invalid argument count")
	}

	switch argv[0] {
	case "gen":
		return commandKeyGen()
	case "pass":
		return commandKeyPass()
	case "bip39":
		return commandKeyBip39()
	default:
		return errors.New("unexpected argument " + argv[0])
	}
}

func commandKeyPass() error {
	// read password
	pass, err := terminal.ReadPassword(syscall.Stdin)
	if err != nil {
		return err
	}
	defer crypt.ZeroBytes(pass)

	// check for minimum required len, reject if not met
	if len(pass) < 8 {
		return errors.New("password len is bellow 8")
	}

	// derive key
	key := crypt.Key(pass, nil)

	fmt.Println()
	commandKeyRatePassword(pass)

	// display key
	fmt.Println()
	fmt.Printf("%s %s\n",
		color.YellowString("key"),
		color.BlueString(hex.EncodeToString(key[:])))

	// display password strength and derived key
	fmt.Println()
	commandKeyRatePassword(pass)
	color.Red("[never share your password or private key]")

	return nil
}

func commandKeyRatePassword(pass []byte) {
	var hasDigit, hasSymbol, hasUpper, hasLower bool

	// check runes
	for _, r := range string(pass) {
		switch {
		case unicode.IsSymbol(r) || unicode.IsPunct(r) || unicode.IsMark(r):
			hasSymbol = true
		case unicode.IsDigit(r):
			hasDigit = true
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		}
	}

	// check conditions
	if len(string(pass)) < 10 {
		color.Red("[we recommend a password length of at least 10 characters]")
	}
	if !hasDigit {
		color.Red("[we recommend your password contains at least one digit]")
	}
	if !hasSymbol {
		color.Red("[we recommend your password contains at least one symbol]")
	}
	if !hasUpper {
		color.Red("[we recommend your password contains at least one upper letter]")
	}
	if !hasLower {
		color.Red("[we recommend your password contains at least one lower letter]")
	}
}

func commandKeyGen() error {
	// generate key
	keyPrivate, err := ppk.NewPrivateKey()
	if err != nil {
		return err
	}

	// remove key from memory
	defer crypt.ZeroBytes(keyPrivate[:])

	// display keys
	fmt.Println()
	fmt.Printf("%s %s\n",
		color.YellowString("private"),
		color.BlueString(keyPrivate.ToHex()))
	fmt.Printf("%s  %s\n",
		color.YellowString("public"),
		color.BlueString(keyPrivate.PublicKey().ToHex()))

	// display message
	color.Red("\n[never share your private key]")

	return nil
}

func commandKeyBip39() error {
	color.Yellow("paste your hex key to export as bip39 or paste a mnemonic to get your key")

	// get input
	text, err := reader.ReadString('\n')
	if err != nil {
		return err
	}

	// check for key
	if len(text)-1 == 64 {
		// decode key
		if err := keyPrivate.FromHex(text[:64]); err != nil {
			return err
		}
		defer crypt.ZeroBytes(keyPrivate[:])

		// display mnemonic
		color.Yellow("\nyour mnemonic are the 24 words bellow")
		color.Blue(keyPrivate.ToMnemonic().String())

		return nil
	}

	// check for mnemonic
	mnemonic, err := mwords.MnemonicFromString(text)
	if err != nil {
		return err
	}

	// decode mnemonic to key
	if err := keyPrivate.FromMnemonic(mnemonic); err != nil {
		return err
	}
	defer crypt.ZeroBytes(keyPrivate[:])

	// display key
	color.Yellow("\nkey extracted from mnemonic is bellow as hex")
	color.Blue(keyPrivate.ToHex())

	return nil
}
