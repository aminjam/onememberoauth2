package onememberoauth2

import (
	"net/http"
	"strings"
	
	"code.google.com/p/goauth2/oauth"
)

func Google(c Client) *Client {
	c.AuthUrl = "https://accounts.google.com/o/oauth2/auth"
	c.TokenUrl = "https://accounts.google.com/o/oauth2/token"
	return NewOAuth2Provider(c)
}

func NewOAuth2Provider(c Client) *Client {
	config := &oauth.Config{
		ClientId:     c.ClientId,
		ClientSecret: c.ClientSecret,
		RedirectURL:  c.RedirectURL,
		Scope:        strings.Join(c.Scopes, " "),
		AuthURL:      c.AuthUrl,
		TokenURL:     c.TokenUrl,
	}
	transport := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
	}
	c.Transport = transport
	return &c
}


