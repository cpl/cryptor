package aes

import "io/ioutil"

// EncryptFiles ...
func EncryptFiles(key Key, files ...string) error {
	for _, file := range files {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		eData, err := Encrypt(key, data)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(file, eData, 0666); err != nil {
			return err
		}
	}

	return nil
}

// DecryptFiles ...
func DecryptFiles(key Key, files ...string) error {
	for _, file := range files {
		eData, err := ioutil.ReadFile(file)
		if err != nil {
			return err
		}
		data, err := Decrypt(key, eData)
		if err != nil {
			return err
		}
		if err := ioutil.WriteFile(file, data, 0666); err != nil {
			return err
		}
	}

	return nil
}
