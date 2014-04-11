package onememberoauth2_test

import (
	"net/http/httptest"

	"github.com/aminjam/onemember"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var as onemember.AccountService
var server *httptest.Server

type service struct{}

func TestOneMemberOAuth2(t *testing.T) {
	RegisterFailHandler(Fail)
	as = &service{}
	RunSpecs(t, "OneMemberOAuth2 Suite")
}
func (s *service) AuthAccount(username string, password string) (out *onemember.Account, err error) {
	out = &onemember.Account{}
	err = nil
	return
}
func (s *service) AddClaim(in *onemember.Account, claim onemember.Claim) (err error) {
	err = nil
	return
}
func (s *service) Create(username string, password string, email string) (out *onemember.Account, err error) {
	out = &onemember.Account{}
	err = nil
	return
}
func (s *service) GetByUsername(username string) (out *onemember.Account, err error) {
	out = &onemember.Account{}
	err = nil
	return
}
func (s *service) RemoveClaim(in *onemember.Account, claimType string, provider string) (err error) {
	err = nil
	return
}
func (s *service) SetLinkedAccount(in *onemember.Account, provider string, claims []onemember.Claim) (err error) {
	err = nil
	return
}
func (s *service) Update(in *onemember.Account) (err error) {
	err = nil
	return
}
