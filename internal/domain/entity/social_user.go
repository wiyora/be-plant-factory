package entity

import "strings"

type Provider string

const (
	ProviderGoogle Provider = "google"
)

type SocialUser struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
	Picture   string
	Provider  Provider
}

func (su SocialUser) FullName() string {
	return strings.TrimSpace(su.FirstName) + " " + strings.TrimSpace(su.LastName)
}
