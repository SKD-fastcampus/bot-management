package firebase

import (
	"context"
	"fmt"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"github.com/SKD-fastcampus/bot-management/pkg/config"
	"google.golang.org/api/option"
)

// TokenVerifier defines the interface for verifying tokens
type TokenVerifier interface {
	VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error)
}

type firebaseTokenVerifier struct {
	authClient *auth.Client
}

// NewTokenVerifier creates a new TokenVerifier using Firebase Admin SDK
func NewTokenVerifier(ctx context.Context, cfg config.Config) (TokenVerifier, error) {
	// 1. Initialize Firebase App
	// Assuming GOOGLE_APPLICATION_CREDENTIALS env var is set or credentials file path is in config
	// If credentials file is specified in config:
	credPath := cfg.GetString("firebase.credentials_file")

	var app *firebase.App
	var err error

	if credPath != "" {
		opt := option.WithCredentialsFile(credPath)
		app, err = firebase.NewApp(ctx, nil, opt)
	} else {
		// Use default credentials (env var GOOGLE_APPLICATION_CREDENTIALS)
		app, err = firebase.NewApp(ctx, nil)
	}

	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	// 2. Get Auth Client
	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting auth client: %v", err)
	}

	return &firebaseTokenVerifier{authClient: client}, nil
}

func (v *firebaseTokenVerifier) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	token, err := v.authClient.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}
