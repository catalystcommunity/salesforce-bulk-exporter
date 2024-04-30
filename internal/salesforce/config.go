package salesforce

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type SalesforceConfig struct {
	BaseUrl      string `yaml:"base-url"`
	ApiVersion   string `yaml:"api-version"`
	ClientId     string `yaml:"client-id"`
	ClientSecret string `yaml:"client-secret"`
	Username     string `yaml:"username"`
	Password     string `yaml:"password"`
	GrantType    string `yaml:"grant-type"`
}

var ConfigFilePath string
var Config SalesforceConfig

// initConfig reads the config file and sets the Config struct with values if
// they are not already set by cli flags or environment variables.
func initConfig() error {
	readConfig, err := readConfigFile(ConfigFilePath)
	if err != nil {
		return err
	}

	fmt.Printf("original config: %+v\n", Config)
	fmt.Printf("read config: %+v\n", readConfig)

	if readConfig != nil {
		overrideUnsetConfigValuesFromReadConfig(readConfig)
	}

	return validateSalesforceArgs()
}

func readConfigFile(path string) (*SalesforceConfig, error) {
	// if no config file is given, default to $HOME/.salesforce-bulk-exporter.yaml
	var isDefault bool
	if path == "" {
		path = os.ExpandEnv("$HOME/.salesforce-bulk-exporter.yaml")
		isDefault = true
	}

	// if the file does not exist, return
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if !isDefault {
			fmt.Println(isDefault)
			return nil, fmt.Errorf("config file does not exist: %s", path)
		}
		return nil, nil
	}

	// read the file
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// decode the file
	var readConfig SalesforceConfig
	err = yaml.NewDecoder(file).Decode(&readConfig)
	if err != nil {
		return nil, err
	}

	return &readConfig, nil
}

func overrideUnsetConfigValuesFromReadConfig(readConfig *SalesforceConfig) {
	if Config.BaseUrl == "" {
		Config.BaseUrl = readConfig.BaseUrl
	}
	if Config.ApiVersion == "" {
		Config.ApiVersion = readConfig.ApiVersion
	}
	if Config.ClientId == "" {
		Config.ClientId = readConfig.ClientId
	}
	if Config.ClientSecret == "" {
		Config.ClientSecret = readConfig.ClientSecret
	}
	if Config.Username == "" {
		Config.Username = readConfig.Username
	}
	if Config.Password == "" {
		Config.Password = readConfig.Password
	}
	if Config.GrantType == "" {
		Config.GrantType = readConfig.GrantType
	}
}

// validateSalesforceArgs checks if all required arguments are set. we don't
// use the cli library to mark these as required, because we allow for using a
// config file. so we have to check this manually.
func validateSalesforceArgs() error {
	var missingArgs []string

	if Config.BaseUrl == "" {
		missingArgs = append(missingArgs, "base-url")
	}
	if Config.ClientId == "" {
		missingArgs = append(missingArgs, "client-id")
	}
	if Config.ClientSecret == "" {
		missingArgs = append(missingArgs, "client-secret")
	}
	if Config.Username == "" {
		missingArgs = append(missingArgs, "username")
	}
	if Config.Password == "" {
		missingArgs = append(missingArgs, "password")
	}

	if len(missingArgs) > 0 {
		return fmt.Errorf("missing salesforce arguments: %s", strings.Join(missingArgs, ", "))
	}

	return nil
}
