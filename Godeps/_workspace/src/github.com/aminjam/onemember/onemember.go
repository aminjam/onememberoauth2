package onemember 
import (
	"errors"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/jameskeane/bcrypt"
)
const (
	PACKAGE string = "oneMember"
)

type DataConnector interface {
	OnememberCreate(in *Account) (error)
	OnememberRead(uid string) (map[string]interface{},error)
	OnememberUpdate(in *Account)(error)
}

type AccountService interface {
	AuthAccount(username string,password string) (*Account,error)
	AddClaim(in *Account,claim Claim) (error)
	Create(username string, password string, email string) (*Account,error)
	GetByUsername(username string) (*Account,error)
	RemoveClaim(in *Account,claimType string, provider string)(error)
	SetLinkedAccount(in *Account, provider string, claims []Claim) (error)
	Update(in * Account) (error)
}

type service struct {
	db DataConnector
}

func New(db DataConnector) AccountService{
	return &service{
		db:db,
	}
}

func (as *service) AuthAccount(username string, password string) (out *Account,err error){
	out, err = as.GetByUsername(username)
	if err != nil { 
		return 
	}
	password += out.Salt
	if !bcrypt.Match(password,out.HashedPassword) {
		err = errors.New("password is a mismatch")
		out = nil
	}
	return 
}

func (as *service) AddClaim(in *Account,claim Claim) (err error){
	in.AddClaim(claim)
	err = as.Update(in)
	return 
}

func (as *service) Create(username string, password string, email string)(out *Account,err error){
	salt, _ := bcrypt.Salt()
	password, _ = bcrypt.Hash(password+salt)
	out = &Account{
			Created: time.Now().Format("20060102150405"),
			Claims: make([]Claim,0),
			Email: email,
			HashedPassword: password,
			LinkedAccounts: make([]LinkedAccount,0),
			Salt: salt,
			Username: username,
	}
	err = as.db.OnememberCreate(out)
	if err != nil {
		out = nil
	}
	return
}

func (as *service) GetByUsername(username string) (out *Account,err error){
	item,err := as.db.OnememberRead(username)
	if (err == nil && item == nil){
		err = errors.New("username is invalid")
		return
	}
	out, err = asAccount(item)
	if (err == nil && out == nil){
		err = errors.New("account is not parsable")
	}
	return 
}

func (as *service) RemoveClaim(in *Account,claimType string, provider string) (err error){
	if provider == "" { provider = "local" }
	in.RemoveClaim(claimType,provider)
	err = as.Update(in)
	return 
}

func (as *service) SetLinkedAccount(in *Account, provider string, claims []Claim) (err error) {
	la := LinkedAccount{
		LastLogin: time.Now().Format("20060102150405"),
		Provider: provider,
	}
	for v := range in.LinkedAccounts {
		if in.LinkedAccounts[v].Provider == provider {
			in.LinkedAccounts = append(in.LinkedAccounts[:v], in.LinkedAccounts[v+1:]...)
			break;
		}
	}
	in.LinkedAccounts = append(in.LinkedAccounts,la)
	for v := range claims {
		claims[v].Provider = provider
		as.RemoveClaim(in,claims[v].Type,provider)
		as.AddClaim(in,claims[v])
	}
	//in.Claims = append(in.Claims,claims...)
	as.Update(in)
	return 
}

func (as *service) Update(in *Account)(err error){
	err = as.db.OnememberUpdate(in)
	return 
}

func asAccount(aMap map[string]interface{}) (out *Account,err error){
	out = &Account{}
	err = mapstructure.Decode(aMap,out)
	return
}
