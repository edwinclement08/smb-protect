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

func Base64Decode(str string) ([]byte, bool) {
	data, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return []byte{}, true
	}
	return data, false
}

func SaveShareMapping(sharePath, roUser, roPass, rwUser, rwPass, mountLocation string) ShareMapping {
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

	roErr := CreateCred(
		fmt.Sprintf("%s:%s", shareMapping.Uuid, "ro"),
		roUser, string(roCipher),
	)

	rwErr := CreateCred(
		fmt.Sprintf("%s:%s", shareMapping.Uuid, "rw"),
		rwUser, string(rwCipher),
	)
	if roErr != nil || rwErr != nil {
		fmt.Println("Failed to save credentials")
	}

	if err := SaveConfig(); err != nil {
		fmt.Println("Error in saving Configuration")
	}
	return shareMapping
}
