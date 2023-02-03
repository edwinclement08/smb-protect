package utils

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/google/uuid"
	"github.com/mitchellh/mapstructure"
)

type MountLocationType int

const (
	DriveMount MountLocationType = 2
	PathMount
)

var LoadedConfig ConfigType

type ShareMapping struct {
	SharePath         string
	MountLocation     string
	MountLocationType MountLocationType
	Uuid              string
	ROUser            string
	RWUser            string
	ROCryptoSalt      string
	ROCryptoNonce     string
	RWCryptoSalt      string
	RWCryptoNonce     string
}

type ConfigType struct {
	Version       int
	ShareMappings []ShareMapping
}

func getDefaultConfigPath() string {
	appDataPath, present := os.LookupEnv("APPDATA")

	if present {
		return path.Join(appDataPath, "smb-protect")
	} else {
		return "."
	}
}

func CreateDefaultConfig() (ConfigType, error) {
	config := ConfigType{
		Version:       1,
		ShareMappings: []ShareMapping{},
	}
	bytes, err := json.Marshal(config)
	if err != nil {
		fmt.Println("Marshaling failed!, Unable to create default config")
		log.Fatal(err)
	}

	configPath := getDefaultConfigPath()
	fileName := "config.json"
	filePath := path.Join(configPath, fileName)

	os.MkdirAll(configPath, os.ModePerm)

	// Open or create file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open config file: %s\n", err.Error())
		return ConfigType{}, err
	}
	defer file.Close()

	// Write to file
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Printf("Failed to write config file: %s\n", err.Error())
		return ConfigType{}, err
	}
	fmt.Println("Config File initialized")
	LoadedConfig = config
	return config, nil
}

func SaveConfig() error {
	bytes, err := json.Marshal(LoadedConfig)
	if err != nil {
		fmt.Println("Marshaling failed!, Unable to create default config")
		log.Fatal(err)
	}

	configPath := getDefaultConfigPath()
	fileName := "config.json"
	filePath := path.Join(configPath, fileName)

	// Open or create file
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Printf("Failed to open config file: %s\n", err.Error())
		return err
	}
	defer file.Close()

	// Write to file
	_, err = file.Write(bytes)
	if err != nil {
		fmt.Printf("Failed to write config file: %s\n", err.Error())
		return err
	}
	fmt.Println("Current Config Saved")
	return nil
}

func LoadConfig() (ConfigType, error) {
	configPath := getDefaultConfigPath()
	fileName := "config.json"
	filePath := path.Join(configPath, fileName)

	configFileData, err := os.ReadFile(filePath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return CreateDefaultConfig()
		} else {
			fmt.Printf("Failed to open existing config file: %s\n", err.Error())
			return ConfigType{}, err
		}
	}
	// Unmarshal JSON
	var raw map[string]interface{}
	err = json.Unmarshal(configFileData, &raw)
	if err != nil {
		fmt.Printf("Unable to parse config: %s\n", err.Error())
		return ConfigType{}, err
	}

	// Convert to struct
	var config ConfigType
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Result:  &config,
		TagName: "json",
	})
	if err != nil {
		fmt.Printf("Unable to configure mapstructure: %s\n", err.Error())
		return ConfigType{}, err
	}
	if err = decoder.Decode(raw); err != nil {
		fmt.Printf("Unable to Unmarshal: %s\n", err.Error())
		return ConfigType{}, err
	}

	fmt.Println("loaded config")
	LoadedConfig = config
	return config, nil
}

func Base64Encode(byteArray []byte) string {
	return base64.StdEncoding.EncodeToString(byteArray)
}

func Base64Decode(str string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}

func SaveShareMappingAndPasswords(sharePath, roUser, roPass, rwUser, rwPass, mountLocation string) ShareMapping {
	masterPass := "pass23"
	roCipher, roNonce, roSalt := Encrypt(roPass, masterPass)
	rwCipher, rwNonce, rwSalt := Encrypt(rwPass, masterPass)

	shareMapping := ShareMapping{
		SharePath:         sharePath,
		MountLocation:     mountLocation,
		MountLocationType: DriveMount,
		Uuid:              uuid.New().String(),
		ROUser:            roUser,
		RWUser:            rwUser,
		ROCryptoSalt:      Base64Encode(roSalt),
		ROCryptoNonce:     Base64Encode(roNonce),
		RWCryptoSalt:      Base64Encode(rwSalt),
		RWCryptoNonce:     Base64Encode(rwNonce),
	}
	LoadedConfig.ShareMappings = append(LoadedConfig.ShareMappings, shareMapping)

	roCipherStr := Base64Encode(roCipher)
	rwCipherStr := Base64Encode(rwCipher)
	roErr := CreateCred(getLabel(shareMapping.Uuid, readOnlyUser), roUser, roCipherStr)
	rwErr := CreateCred(getLabel(shareMapping.Uuid, readWriteUser), rwUser, rwCipherStr)

	if roErr != nil || rwErr != nil {
		fmt.Println("Failed to save credentials")
	}

	if err := SaveConfig(); err != nil {
		fmt.Println("Error in saving Configuration")
	}
	return shareMapping
}

type UserAccessLevel int

const (
	readOnlyUser  UserAccessLevel = 1
	readWriteUser UserAccessLevel = 2
)

func getLabel(uuid string, userAccessLevel UserAccessLevel) string {
	access := "rw"
	if userAccessLevel == readOnlyUser {
		access = "ro"
	}
	label := fmt.Sprintf("%s:%s", uuid, access)
	return label
}

func LoadPasswords(shareMapping ShareMapping) (roUser, roPass, rwUser, rwPass string, err error) {
	masterPass := "pass23"
	roUser, roPassEnc, err1 := RetrieveCredential(getLabel(shareMapping.Uuid, readOnlyUser))
	rwUser, rwPassEnc, err2 := RetrieveCredential(getLabel(shareMapping.Uuid, readWriteUser))

	if err1 != nil || err2 != nil {
		fmt.Println("Failed to retrieve Credentials")
		return "", "", "", "", CoalesceError(err1, err2)
	}

	roPassEncByte, err1 := Base64Decode(roPassEnc)
	roNonce, err4 := Base64Decode(shareMapping.ROCryptoNonce)
	roSalt, err3 := Base64Decode(shareMapping.ROCryptoSalt)

	rwPassEncByte, err2 := Base64Decode(rwPassEnc)
	rwNonce, err6 := Base64Decode(shareMapping.RWCryptoNonce)
	rwSalt, err5 := Base64Decode(shareMapping.RWCryptoSalt)
	totalErr := CoalesceError(err1, err2, err3, err4, err5, err6)
	if totalErr != nil {
		fmt.Println("Failed to base64 decode config")
		return "", "", "", "", totalErr
	}

	roPassDecStr := Decrypt(roPassEncByte, masterPass, roNonce, roSalt)
	rwPassDecStr := Decrypt(rwPassEncByte, masterPass, rwNonce, rwSalt)

	return roUser, roPassDecStr, rwUser, rwPassDecStr, nil
}
