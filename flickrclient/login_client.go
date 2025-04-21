package flickrclient

import (
	"fmt"

	"gopkg.in/masci/flickr.v3"
)

type LoginClient struct {
	client *flickr.FlickrClient
}

type LoginHandler interface {
	Login() (string, string, error)
}

type oauthFlickrHandler interface {
	GetRequestToken() (*flickr.RequestToken, error)
	GetAuthorizeUrl(requestTok *flickr.RequestToken) (string, error)
	GetAccessToken(requestTok *flickr.RequestToken, confirmationCode string) (*flickr.OAuthToken, error)
}

func (b *LoginClient) GetRequestToken() (*flickr.RequestToken, error) {
	requestTok, err := flickr.GetRequestToken(b.client)
	if err != nil {
		return nil, fmt.Errorf("cannot get request token: %w", err)
	}
	return requestTok, nil
}

func (b *LoginClient) GetAuthorizeUrl(requestTok *flickr.RequestToken) (string, error) {
	url, err := flickr.GetAuthorizeUrl(b.client, requestTok)
	if err != nil {
		return "", fmt.Errorf("cannot get authorize url: %w", err)
	}
	return url, nil
}

func (b *LoginClient) GetAccessToken(requestTok *flickr.RequestToken, confirmationCode string) (*flickr.OAuthToken, error) {
	accessTok, err := flickr.GetAccessToken(b.client, requestTok, confirmationCode)
	if err != nil {
		return nil, fmt.Errorf("cannot get access token: %w", err)
	}
	return accessTok, nil
}

func (lc *LoginClient) Login() (string, string, error) {
	return oauthLogin(lc, &ScanCodeReader{})
}

type CodeReader interface {
	ReadCode() (string, error)
}

type ScanCodeReader struct{}

func (r *ScanCodeReader) ReadCode() (string, error) {
	var confirmationCode string
	_, err := fmt.Scanln(&confirmationCode)
	if err != nil {
		return "", fmt.Errorf("failed to read confirmation code: %w", err)
	}
	return confirmationCode, nil
}

func oauthLogin(b oauthFlickrHandler, cr CodeReader) (string, string, error) {
	requestTok, err := b.GetRequestToken()
	if err != nil {
		return "", "", fmt.Errorf("cannot get request token: %w\nplease check your api key and secret. You can set them using flkcli setup", err)
	}

	url, _ := b.GetAuthorizeUrl(requestTok)
	// Print url
	fmt.Printf("Please visit this URL to authorize the application: %s\n", url)
	// Ask user to input the code
	fmt.Print("Please enter the OAuth confirmation code: ")
	confirmationCode, err := cr.ReadCode()
	if err != nil {
		return "", "", fmt.Errorf("failed to read confirmation code: %w", err)
	}

	accessTok, err := b.GetAccessToken(requestTok, confirmationCode)
	if err != nil {
		return "", "", fmt.Errorf("cannot get access token: %w", err)
	}
	return accessTok.OAuthToken, accessTok.OAuthTokenSecret, nil
}

func GetLoginHandler() (LoginHandler, error) {
	c, err := GetFlickrClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get flickr client: %w", err)
	}
	return &LoginClient{client: c}, nil
}
