package firebase

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type Auth struct {
	Client *auth.Client
}

func NewAuth(ctx context.Context, credFile string) *Auth {
	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credFile))
	if err != nil {
		log.Fatalf("error initializing firebase app: %v", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		log.Fatalf("error getting auth client: %v", err)
	}

	return &Auth{Client: client}
}

func (a *Auth) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	return a.Client.VerifyIDToken(ctx, idToken)
}
