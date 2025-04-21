package login

import (
	"flkcli/config"
	"fmt"
	"testing"
)

type MockConfig struct {
	APIConfig   *config.APIConfig
	TokenConfig *config.TokenConfig
}

func (m *MockConfig) GetAPIConfig() (*config.APIConfig, error) {
	return m.APIConfig, nil
}

func (m *MockConfig) GetTokenConfig() (*config.TokenConfig, error) {
	return m.TokenConfig, nil
}

func (m *MockConfig) SetAPIConfig(secret, key string) error {
	m.APIConfig = &config.APIConfig{
		Secret: secret,
		Key:    key,
	}
	return nil
}
func (m *MockConfig) SetTokenConfig(oauthToken, oauthTokenSecret string) error {
	m.TokenConfig = &config.TokenConfig{
		OAuthToken:       oauthToken,
		OAuthTokenSecret: oauthTokenSecret,
	}
	return nil
}

type MockLoginHandler struct{}

func (m *MockLoginHandler) Login() (string, string, error) {
	return "mockToken", "mockSecret", nil
}

type MockLoginHandlerError struct{}

func (m *MockLoginHandlerError) Login() (string, string, error) {
	return "", "", fmt.Errorf("mock error")
}

func TestLoginCmd(t *testing.T) {
	cm := &MockConfig{}
	lh := &MockLoginHandler{}
	login(cm, lh)
	if cm.TokenConfig.OAuthToken != "mockToken" {
		t.Errorf("Expected token 'mockToken', got '%s'", cm.TokenConfig.OAuthToken)
	}
	if cm.TokenConfig.OAuthTokenSecret != "mockSecret" {
		t.Errorf("Expected secret 'mockSecret', got '%s'", cm.TokenConfig.OAuthTokenSecret)
	}
}

func TestLoginCmdError(t *testing.T) {
	cm := &MockConfig{}
	lh := &MockLoginHandlerError{}
	login(cm, lh)
	if cm.TokenConfig != nil {
		t.Errorf("Expected token config to be nil, got '%v'", cm.TokenConfig)
	}
}
