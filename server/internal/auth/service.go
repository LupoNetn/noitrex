package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/utils"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	CreateOperator(ctx context.Context, args RegisterRequest) (db.Operator, string, string, error)
}

type Svc struct {
	db               db.Querier
	JWTAccessSecret  string
	JWTRefreshSecret string
}

func NewService(db db.Querier, JWTAccessSecret, JWTRefreshSecret string) Service {
	return &Svc{
		db:               db,
		JWTAccessSecret:  JWTAccessSecret,
		JWTRefreshSecret: JWTRefreshSecret,
	}
}

// implement service interface
func (s *Svc) Login(ctx context.Context, email string, password string) (string, string, error) {
	emailTxt := pgtype.Text{
		String: email,
		Valid:  true,
	}
	existingUser, err := s.db.GetOperatorByEmail(ctx, emailTxt)
	if err != nil {
		slog.Error("an error occured when trying to fetch operator", "error", err)
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrInvalidCredentials
		}
		return "", "", err
	}

	IsPasswordCorrect, err := utils.ComparePasswordHash(password, existingUser.PasswordHash.String)
	if err != nil {
		slog.Error("Something went wrong when comparing password hash for user:", "user", existingUser.ID)
		return "", "", ErrInvalidCredentials
	}

	if !IsPasswordCorrect {
		slog.Error("Password is incorrect for user:", "user", existingUser.ID)
		return "", "", ErrInvalidCredentials
	}

	accessToken, refreshToken, err := utils.GenerateJwtPair(s.JWTAccessSecret, s.JWTRefreshSecret, existingUser.Name, existingUser.ID.String(), "operator")
	if err != nil {
		slog.Error("something went wrong when generating token pair", "error", err)
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

func (s *Svc) CreateOperator(ctx context.Context, args RegisterRequest) (db.Operator, string, string, error) {

	apiKey, err := utils.GenerateApiKey()
	if err != nil {
		slog.Error("something went wrong when generating api key", "error", err)
		return db.Operator{}, "", "", err
	}

	webhookSecret, err := utils.GenerateWebhookSecret()
	if err != nil {
		slog.Error("something went wrong when generating webhook secret", "error", err)
		return db.Operator{}, "", "", err
	}

	hashedPassword, err := utils.GeneratePasswordHash(args.Password)
	if err != nil {
		slog.Error("something went wrong when generating password hash", "error", err)
		return db.Operator{}, "", "", err
	}

	// Hash api key and webhook secret for storage.
	apiKeyHash := utils.HashSecrets(apiKey)
	webhookSecretHash := utils.HashSecrets(webhookSecret)

	params := db.CreateOperatorParams{
		Name:          args.Name,
		Description:   pgtype.Text{String: *args.Description, Valid: true},
		LogoUrl:       pgtype.Text{String: *args.LogoUrl, Valid: true},
		SupportEmail:  pgtype.Text{String: args.SupportEmail, Valid: true},
		WebsiteUrl:    pgtype.Text{String: *args.WebsiteUrl, Valid: true},
		PasswordHash:  pgtype.Text{String: hashedPassword, Valid: true},
		ApiKeyHash:    apiKeyHash,
		WebhookSecret: webhookSecretHash,
		Status:        db.OperatingStatusActive,
	}

	//create operator
	operator, err := s.db.CreateOperator(ctx, params)
	if err != nil {
		slog.Error("something went wrong when creating operator", "error", err)
		return db.Operator{}, "", "", err
	}

	return operator, apiKey, webhookSecret, nil
}
