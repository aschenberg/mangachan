package config

import (
	"context"
	"log"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

func NewOidcProvider(env *Env) *oidc.Provider {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	provider, err := oidc.NewProvider(ctx, env.IssuerURL)
	if err != nil {
		log.Fatalf("Failed to get provider: %v", err)
	}

	return provider

}

func NewOidcConfig(env *Env, provider *oidc.Provider) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     env.ClientID,
		ClientSecret: env.ClientSecret,
		RedirectURL:  env.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       []string{oidc.ScopeOpenID, "profile", "email", "openid"},
	}
}
