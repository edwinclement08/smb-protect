package utils

import (
	"fmt"

	"github.com/danieljoos/wincred"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func getTargetName(label string) string {
	return fmt.Sprintf("smb-protect:%s", label)
}

func CreateCred(label, username, password string) error {
	cred := wincred.NewGenericCredential(getTargetName(label))
	cred.UserName = username

	encoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewEncoder()
	blob, _, err := transform.Bytes(encoder, []byte(password))
	if err != nil {
		fmt.Println(err)
		return err
	}

	cred.CredentialBlob = blob
	err = cred.Write()

	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func RetrieveCredential(label string) (user, pass string, err error) {
	cred, err := wincred.GetGenericCredential(getTargetName(label))
	if err != nil {
		fmt.Printf("Failed to Retrieve Password: %s\n", err)
		return "", "", err
	}

	decoder := unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM).NewDecoder()
	blob, _, err := transform.Bytes(decoder, cred.CredentialBlob)
	if err != nil {
		fmt.Printf("Failed to Decode Password: %s\n", err)
		return "", "", err
	}
	pass = string(blob)
	return cred.UserName, pass, nil

}

func ListCreds() ([]*wincred.Credential, error) {
	creds, err := wincred.FilteredList(getTargetName("*"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	for i := range creds {
		fmt.Println(creds[i])
	}
	return creds, nil
}

func DeleteCred(label string) error {
	cred, err := wincred.GetGenericCredential(getTargetName(label))
	if err != nil {
		fmt.Printf("Failed to retrieve Credential for: %s\n", label)
		return err
	}

	err = cred.Delete()
	if err != nil {
		fmt.Printf("Failed to delete Credential for: %s\n", label)
		return err
	}
	return nil
}
