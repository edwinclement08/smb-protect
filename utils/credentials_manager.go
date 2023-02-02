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
	blob, _, err := transform.Bytes(encoder, []byte("mysecret"))
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

func ListCred() {
	creds, err := wincred.List()
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range creds {
		fmt.Println(creds[i].TargetName)
	}
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
