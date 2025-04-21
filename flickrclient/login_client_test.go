package flickrclient

import (
	"testing"

	"gopkg.in/masci/flickr.v3"
)

type FakeOauthFlickrHandler struct {
	Client                       any
	GetRequestTokenCalled        bool
	GetAuthorizeUrlCalled        bool
	GetAuthorizeUrlTokenReceived *flickr.RequestToken
	confirmationCodeReceived     string
	GetAccessTokenCalled         bool
}

func (b *FakeOauthFlickrHandler) GetRequestToken() (*flickr.RequestToken, error) {
	b.GetRequestTokenCalled = true
	return &flickr.RequestToken{
		OauthToken: "fake_request_token",
	}, nil
}

func (b *FakeOauthFlickrHandler) GetAuthorizeUrl(requestTok *flickr.RequestToken) (string, error) {
	b.GetAuthorizeUrlCalled = true
	b.GetAuthorizeUrlTokenReceived = requestTok
	return "fake_authorize_url", nil
}

func (b *FakeOauthFlickrHandler) GetAccessToken(requestTok *flickr.RequestToken, confirmationCode string) (*flickr.OAuthToken, error) {
	b.GetAccessTokenCalled = true
	b.confirmationCodeReceived = confirmationCode
	return &flickr.OAuthToken{}, nil
}

type FakeCodeReader struct {
	confirmationCode string
}

func (r *FakeCodeReader) ReadCode() (string, error) {
	return r.confirmationCode, nil
}

func TestLogin(t *testing.T) {
	confirmationCode := "fake_confirmation_code"
	fakeCodeReader := &FakeCodeReader{
		confirmationCode: confirmationCode,
	}
	fakeClient := &FakeOauthFlickrHandler{}
	oauthLogin(fakeClient, fakeCodeReader)
	if !fakeClient.GetRequestTokenCalled {
		t.Errorf("GetRequestToken was not called")
	}
	if !fakeClient.GetAuthorizeUrlCalled {
		t.Errorf("GetAuthorizeUrl was not called")
	}
	if fakeClient.GetAuthorizeUrlTokenReceived.OauthToken != "fake_request_token" {
		t.Errorf("GetAuthorizeUrl was called with wrong token: %s", fakeClient.GetAuthorizeUrlTokenReceived.OauthToken)
	}
	if !fakeClient.GetAccessTokenCalled {
		t.Errorf("GetAccessToken was not called")
	}
	if fakeClient.confirmationCodeReceived != confirmationCode {
		t.Errorf("GetAccessToken was called with wrong confirmation code: %s", fakeClient.confirmationCodeReceived)
	}

}
