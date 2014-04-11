package onemember_test

import (
	"github.com/aminjam/onemember"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func init() {
	Describe("Onemember tests", func() {
		var (
			account        *onemember.Account
			err            error
			accountService onemember.AccountService
		)
		BeforeEach(func() {
			accountService = onemember.New(db)
		})
		Context("when creating an account ", func() {
			JustBeforeEach(func() {
				account, err = accountService.Create("aminjam", "none", "none@aminjam.com")
			})
			It("should not have errors", func() {
				Ω(err).Should(BeNil())
			})
			It("should have a record without error", func() {
				_, err := db.OnememberRead("aminjam")
				Ω(err).Should(BeNil())
			})
			It("should have a record", func() {
				out, _ := db.OnememberRead("aminjam")
				Ω(out).ShouldNot(BeNil())
			})
			It("should have the same ID", func() {
				out, _ := db.OnememberRead("aminjam")
				Ω(out["Username"]).Should(Equal("aminjam"))
			})

			Context("when reading back the account", func() {
				JustBeforeEach(func() {
					account, err = accountService.GetByUsername("aminjam")
				})
				It("should be without error", func() {
					Ω(err).Should(BeNil())
				})
				It("should have non-null record", func() {
					Ω(account).ShouldNot(BeNil())
				})
				It("should have the record", func() {
					Ω(account.Username).Should(Equal("aminjam"))
				})
			})

			Context("when updating an account", func() {
				JustBeforeEach(func() {
					account.Email = "changed@aminjam.com"
					err = accountService.Update(account)
				})
				It("should not have errors", func() {
					Ω(err).Should(BeNil())
				})
				It("should be updated", func() {
					na, _ := accountService.GetByUsername("aminjam")
					Ω(account.Username).Should(Equal(na.Username))
					Ω(account.Email).Should(Equal("changed@aminjam.com"))
				})
			})

			Context("when authenticating an account", func() {
				JustBeforeEach(func() {
					account, err = accountService.AuthAccount("aminjam", "none")
				})
				It("should not have errors", func() {
					Ω(err).Should(BeNil())
				})
				It("should retrieve the account", func() {
					Ω(account.Username).Should(Equal("aminjam"))
				})
			})

			Context("when authentication fails", func() {
				Describe("incorrect username", func() {
					JustBeforeEach(func() {
						account, err = accountService.AuthAccount("incorrectUsername", "none")
					})
					It("should have errors", func() {
						Ω(err).ShouldNot(BeNil())
					})
					It("should not retrieve the account", func() {
						Ω(account).Should(BeNil())
					})
				})
				Describe("incorrect password", func() {
					JustBeforeEach(func() {
						account, err = accountService.AuthAccount("aminjam", "incorrectPassword")
					})
					It("should have errors", func() {
						Ω(err).ShouldNot(BeNil())
					})
					It("should not retrieve the account", func() {
						Ω(account).Should(BeNil())
					})
				})
			})

			Context("when adding claims", func() {
				JustBeforeEach(func() {
					claim := onemember.Claim{
						Type:  "myName",
						Value: "theValue",
					}
					err = accountService.AddClaim(account, claim)
				})
				It("should not have errors", func() {
					Ω(err).Should(BeNil())
				})
				It("should have the claim", func() {
					for v := range account.Claims {
						if account.Claims[v].Type == "myName" {
							claim := account.Claims[v]
							Ω(claim).ShouldNot(BeNil())
							Ω(claim.Value).Should(Equal("theValue"))
						}
					}
				})

				Context("when removing the claim", func() {
					JustBeforeEach(func() {
						err = accountService.RemoveClaim(account, "myName", "")
					})
					It("should not have errors", func() {
						Ω(err).Should(BeNil())
					})
					It("should have removed the claim", func() {
						for v := range account.Claims {
							if account.Claims[v].Type == "myName" {
								claim := account.Claims[v]
								Ω(claim).ShouldNot(BeNil())
								Ω(claim.Value).Should(Equal(""))
							}
						}
					})
				})
			})

			Context("when adding linked account", func() {
				JustBeforeEach(func() {
					claims := make([]onemember.Claim, 1)
					claims[0] = onemember.Claim{
						Type:  "urn:facebook:name",
						Value: "aminjamFB",
					}
					err = accountService.SetLinkedAccount(account, "facebook", claims)
				})
				It("should not have errors", func() {
					Ω(err).Should(BeNil())
				})
				It("should contain the linked account", func() {
					Ω(account.LinkedAccounts).Should(HaveLen(1))
					Ω(account.LinkedAccounts[0].Provider).Should(Equal("facebook"))
				})
				It("should contain the linked account claim", func() {
					Ω(account.Claims).Should(HaveLen(1))
					Ω(account.Claims[0].Provider).Should(Equal("facebook"))
				})
			})
		})
	})
}
