package di

import (
	"go-api/domain/port"
	"go-api/domain/repository"
	httphandler "go-api/handler/http"
	"go-api/handler/http/middleware"
	infraClerk "go-api/infrastructure/clerk"
	"go-api/infrastructure/config"
	repoGorm "go-api/repository/gorm"
	"log"

	"gorm.io/gorm"
)

type Container struct {
	AuthenticateMiddleware *middleware.AuthenticateMiddleware
	UserWebhookMiddleware  *middleware.UserWebhookMiddleware
	UserWebhookHandler     *httphandler.UserWebhookHandler
	UserHandler            *httphandler.UserHandler
}

type apiDeps struct {
	env          *config.Config
	db           *gorm.DB
	userRepo     repository.UserRepository
	jwksProvider port.TokenKeyProvider
}

func NewContainer(db *gorm.DB, env *config.Config) *Container {
	jwksProvider, err := infraClerk.NewJWKSProvider(env)
	if err != nil {
		log.Fatalf("failed to create JWKS provider: %v", err)
	}

	d := &apiDeps{
		env:          env,
		db:           db,
		userRepo:     repoGorm.NewUserRepository(db),
		jwksProvider: jwksProvider,
	}

	authBundle := wireAuth(d)

	return &Container{
		AuthenticateMiddleware: authBundle.authenticateMiddleware,
		UserWebhookMiddleware:  middleware.NewUserWebhookMiddleware(env.ClerkWebhookSecret),
		UserWebhookHandler:     authBundle.userWebhookHandler,
		UserHandler:            httphandler.NewUserHandler(),
	}
}
