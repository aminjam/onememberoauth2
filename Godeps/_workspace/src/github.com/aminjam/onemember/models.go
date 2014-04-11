package onemember

type Claim struct {
	Provider string `json:"provider"`
	Type     string `json:"type"`
	Value    string `json:"value"`
}

type LinkedAccount struct {
	LastLogin string `json:"lastLogin"`
	Provider  string `json:"provider"`
}

type Account struct {
	Created        string          `json:"created"`
	Claims         []Claim         `json:"claims"`
	Email          string          `json:"email"`
	HashedPassword string          `json:"hashedPassword,omitempty" out:"false"`
	LinkedAccounts []LinkedAccount `json:"linkedAccount"`
	Salt           string          `json:"salt,omitempty" out:"false"`
	Tenant         string          `json:"tenant"`
	Username       string          `json:"username"`
}

func (a *Account) AddClaim(c Claim) {
	a.Claims = append(a.Claims, c)
}
func (a *Account) RemoveClaim(claimType string, provider string) {
	for v := range a.Claims {
		if a.Claims[v].Type == claimType && (provider == "local" || a.Claims[v].Provider == provider) {
			a.Claims = append(a.Claims[:v], a.Claims[v+1:]...)
			break
		}
	}
}
