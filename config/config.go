package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type APIConfig struct {
	Secret string `yaml:"secret"`
	Key    string `yaml:"key"`
}

type TokenConfig struct {
	OAuthToken       string `yaml:"oauth_token"`
	OAuthTokenSecret string `yaml:"oauth_token_secret"`
}

type ConfigManager interface {
	GetAPIConfig() (*APIConfig, error)
	GetTokenConfig() (*TokenConfig, error)
	SetAPIConfig(secret, key string) error
	SetTokenConfig(oauthToken, oauthTokenSecret string) error
}

type FileConfig struct {
	apiFile   string
	tokenFile string
}

func GetUserProfileConfig() (ConfigManager, error) {
	apiFile, err := getConfigFilePath("api")
	if err != nil {
		return nil, fmt.Errorf("failed to get api config file: %w", err)
	}
	tokenFile, err := getConfigFilePath("token")
	if err != nil {
		return nil, fmt.Errorf("failed to get token config file: %w", err)
	}
	config := &FileConfig{
		apiFile:   apiFile,
		tokenFile: tokenFile,
	}
	return config, nil
}

// SetApiConfig sets the API key and secret in the config file
func (c *FileConfig) SetAPIConfig(secret, key string) error {
	apiConfig := APIConfig{
		Secret: secret,
		Key:    key,
	}

	if err := writeConfigFile(apiConfig, c.apiFile); err != nil {
		return fmt.Errorf("failed to save API config: %w", err)
	}

	return nil
}

// SetTokenConfig sets the OAuth token and secret in the config file
func (c *FileConfig) SetTokenConfig(oauthToken, oauthTokenSecret string) error {
	tokenConfig := TokenConfig{
		OAuthToken:       oauthToken,
		OAuthTokenSecret: oauthTokenSecret,
	}

	if err := writeConfigFile(tokenConfig, c.tokenFile); err != nil {
		return fmt.Errorf("failed to save Token config: %w", err)
	}

	return nil
}

// GetAPIConfig returns the API config from the config file
func (c *FileConfig) GetAPIConfig() (*APIConfig, error) {
	data, err := readConfigFile(c.apiFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read api config file: %w", err)
	}

	var apiConfig APIConfig
	if err := yaml.Unmarshal(data, &apiConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal API config from YAML: %w", err)
	}

	return &apiConfig, nil
}

// GetTokenConfig returns the Token config from the config file
func (c *FileConfig) GetTokenConfig() (*TokenConfig, error) {
	data, err := readConfigFile(c.tokenFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read token config file: %w", err)
	}

	var tokenConfig TokenConfig
	if err := yaml.Unmarshal(data, &tokenConfig); err != nil {
		return nil, fmt.Errorf("failed to unmarshal Token config from YAML: %w", err)
	}

	return &tokenConfig, nil
}

// readConfigFile reads the config file and returns the data
func readConfigFile(filePath string) ([]byte, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	return data, nil
}

// writeConfigFile writes the data to the config file
func writeConfigFile(data any, filePath string) error {
	stringData, err := yaml.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal Token config to YAML: %w", err)
	}

	if err := os.WriteFile(filePath, stringData, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	return nil
}

// getConfigFilePath returns the path to the config file
// It creates the config directory if it does not exist
func getConfigFilePath(file string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("cannot get user home directory")
	}

	configDir := filepath.Join(homeDir, ".flkcli")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", errors.New("cannot create config directory")
	}

	configFilePath := filepath.Join(configDir, file+".yaml")
	return configFilePath, nil
}
