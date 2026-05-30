package auth

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/luponetn/noitrex/internal/db"
	"github.com/luponetn/noitrex/utils"
)

type Service interface {
	Login(ctx context.Context, email, password string) (string, string, error)
	CreateOperator(ctx context.Context, args RegisterRequest) (db.Operator, string, string, error)
	Refresh(ctx context.Context, refreshToken string) (string, string, error)
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
		if errors.Is(err, pgx.ErrNoRows) {
			return "", "", ErrInvalidCredentials // normal — no logging needed
		}
		slog.Error("unexpected error fetching operator", "error", err)
		return "", "", err
	}

	isPasswordCorrect, err := utils.ComparePasswordHash(password, existingUser.PasswordHash.String)
	if err != nil {
		slog.Error("Something went wrong when comparing password hash for user:", "user", existingUser.ID)
		return "", "", ErrInvalidCredentials
	}

	if !isPasswordCorrect {
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

	// Note: We removed the GetOperatorByEmail check to prevent a Time-of-Check to Time-of-Use (TOCTOU) race condition.
	// We now rely on the database's UNIQUE constraint on the email column to enforce uniqueness.

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

	params := db.CreateOperatorParams{
		Name:          args.Name,
		Description:   utils.ToPgtypeText(args.Description),
		LogoUrl:       utils.ToPgtypeText(args.LogoUrl),
		SupportEmail:  pgtype.Text{String: args.SupportEmail, Valid: true},
		WebsiteUrl:    utils.ToPgtypeText(args.WebsiteUrl),
		PasswordHash:  pgtype.Text{String: hashedPassword, Valid: true},
		ApiKeyHash:    apiKeyHash,
		WebhookSecret: webhookSecret,
		Status:        db.OperatingStatusActive,
	}

	//create operator
	operator, err := s.db.CreateOperator(ctx, params)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is PostgreSQL's unique_violation code
			slog.Warn("user registration failed: email already exists", "email", args.SupportEmail)
			return db.Operator{}, "", "", ErrUserAlreadyExists
		}
		slog.Error("something went wrong when creating operator", "error", err)
		return db.Operator{}, "", "", err
	}

	return operator, apiKey, webhookSecret, nil
}

func (s *Svc) Refresh(ctx context.Context, refreshToken string) (string, string, error) {
	claims, err := utils.VerifyJwt(refreshToken, s.JWTRefreshSecret)
	if err != nil {
		slog.Error("error verifying refresh token", "error", err)
		return "", "", ErrInvalidToken
	}

	// Issue a new pair of access and refresh tokens
	accessToken, newRefreshToken, err := utils.GenerateJwtPair(s.JWTAccessSecret, s.JWTRefreshSecret, claims.Name, claims.OperatorID, claims.Role)
	if err != nil {
		slog.Error("something went wrong when generating token pair during refresh", "error", err)
		return "", "", err
	}
	return accessToken, newRefreshToken, nil
}
