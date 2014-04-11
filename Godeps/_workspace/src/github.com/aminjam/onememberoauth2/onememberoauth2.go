package onememberoauth2

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

	"code.google.com/p/goauth2/oauth"
	"github.com/aminjam/onemember"
)

type ClaimsBuilderfunc func(string) (map[string]string, error)
type Client struct {
	ClientId     string
	ClientSecret string
	RedirectURL  string
	ScopeDivider string
	Scopes       []string

	AuthUrl       string
	TokenUrl      string
	ClaimsBuilder ClaimsBuilderfunc
	Transport     *oauth.Transport
}

type Consumer interface {
	Callback(w http.ResponseWriter, r *http.Request, as onemember.AccountService) (int, string)
	Request(w http.ResponseWriter, r *http.Request) (int, string)
}

type consumer struct {
	clients map[string]*Client
}

func (c *consumer) parseProvider(r *http.Request) (provider string, err error) {
	regex := regexp.MustCompile(`/`)
	path := regex.Split(r.URL.Path, -1)
	provider = path[len(path)-1]
	if c.clients[provider] == nil {
		provider = r.URL.Query().Get("provider")
		if c.clients[provider] == nil {
			provider = ""
			err = errors.New("OAuth2 provider is invalid")
		}
	}
	return
}

func (c *consumer) Callback(w http.ResponseWriter, r *http.Request, as onemember.AccountService) (status int, msg string) {
	status = http.StatusBadRequest
	provider, err := c.parseProvider(r)
	if err != nil {
		msg = err.Error()
		return
	}
	code := r.URL.Query().Get("code")
	tk, err := c.clients[provider].Transport.Exchange(code)
	if err != nil {
		msg = err.Error()
		return
	}
	decodedState, err := base64.StdEncoding.DecodeString(r.URL.Query().Get("state"))
	if err != nil {
		msg = err.Error()
		return
	}
	type StateInfo struct{ Id string }
	var stateInfo StateInfo
	json.Unmarshal(decodedState, &stateInfo)
	rawClaims, err := c.clients[provider].ClaimsBuilder(tk.AccessToken)
	if err != nil {
		msg = err.Error()
		return
	}
	account, err := as.GetByUsername(stateInfo.Id)
	if err != nil {
		msg = err.Error()
		return
	}
	var claims = make([]onemember.Claim, len(rawClaims))
	var i = 0
	for v := range rawClaims {
		claims[i] = onemember.Claim{
			Provider: provider,
			Type:     v,
			Value:    rawClaims[v],
		}
		i++
	}
	as.SetLinkedAccount(account, provider, claims)
	return 200, tk.AccessToken
}

func (c *consumer) Request(w http.ResponseWriter, r *http.Request) (int, string) {
	provider, err := c.parseProvider(r)
	id := r.URL.Query().Get("id")
	if err != nil {
		return 500, err.Error()
	}
	state := []byte(`{"id":"` + id + `"}`)
	url := c.clients[provider].Transport.Config.AuthCodeURL(base64.StdEncoding.EncodeToString(state))
	http.Redirect(w, r, url, 302)
	return 302, ""
}

func New(clients map[string]*Client) Consumer {
	return &consumer{
		clients: clients,
	}
}
