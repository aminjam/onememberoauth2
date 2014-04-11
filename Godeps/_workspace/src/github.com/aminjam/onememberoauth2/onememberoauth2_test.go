package onememberoauth2_test

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	. "github.com/aminjam/onememberoauth2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OneMemberOAuth2", func() {
	var (
		options  map[string]*Client
		consumer Consumer
		err      error
		res      *http.Response
		handler  func(http.ResponseWriter, *http.Request)
	)
	BeforeEach(func() {
		handler = func(w http.ResponseWriter, r *http.Request) {
			switch {
			case strings.Contains(r.URL.Path, "/token"):
				w.Header().Set("Content-Type", "application/json")
				body := `{"access_token":"token1","refresh_token":"refreshtoken1","id_token":"idtoken1","expires_in":3600}`
				io.WriteString(w, body)
			case strings.Contains(r.URL.Path, "/auth"):
				q := url.Values{
					"code":  {"c0d3"},
					"state": {r.URL.Query().Get("state")},
				}.Encode()
				url := strings.Replace(r.URL.Path, "/auth", "/callback/myprovider", -1)
				url += "?" + q
				http.Redirect(w, r, url, 302)
			case strings.Contains(r.URL.Path, "/callback"):
				status, message := consumer.Callback(w, r, as)
				w.WriteHeader(status)
				io.WriteString(w, message)
			case strings.Contains(r.URL.Path, "/request"):
				_, message := consumer.Request(w, r)
				io.WriteString(w, message)
			}
		}
		options = make(map[string]*Client)
		server = httptest.NewServer(http.HandlerFunc(handler))
		options["myprovider"] = NewOAuth2Provider(Client{
			AuthUrl:      server.URL + "/auth",
			TokenUrl:     server.URL + "/token",
			ClientId:     "ClientID",
			ClientSecret: "ClientSecret",
			ClaimsBuilder: func(accessToken string) (claims map[string]string, err error) {
				claims = make(map[string]string)
				claims["name"] = "MyName"
				return
			},
		})
		consumer = New(options)
	})
	AfterEach(func() {
		defer server.Close()
	})
	It("request should have ok (200) status code with accessToken", func() {
		res, err = http.Get(server.URL + "/request/myprovider?id=aminjam")
		Ω(err).Should(BeNil())
		Ω(res.StatusCode).Should(Equal(200))
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		Ω(string(body)).Should(Equal("token1"))
	})
	It("callback should have ok (200) status code with accessToken", func() {
		q := url.Values{
			"code":  {"c0d3"},
			"state": {"eyJpZCI6ImFtaW5qYW0ifQ=="},
		}.Encode()
		url := server.URL + "/callback/myprovider"
		url += "?" + q
		res, err = http.Get(url)
		Ω(err).Should(BeNil())
		Ω(res.StatusCode).Should(Equal(200))
		defer res.Body.Close()
		body, _ := ioutil.ReadAll(res.Body)
		Ω(string(body)).Should(Equal("token1"))
	})
})
