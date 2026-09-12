package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/auth/dto"
	"kialkuz/shop-with-loyalty/internal/auth/model"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

type TokenService struct {
	repository TokenRepository
}

func NewTokenService(repository TokenRepository) *TokenService {
	return &TokenService{repository: repository}
}

func (s TokenService) AddNewToken(ctx context.Context, token model.Token) error {
	return s.repository.AddNewToken(ctx, token)
}

func (s TokenService) AddTx(ctx context.Context, tx pgx.Tx, token model.Token) error {
	return s.repository.AddTx(ctx, tx, token)
}

func (s *TokenService) MakeToken(userID uuid.UUID, tokenExp time.Duration) model.Token {
	return model.Token{
		Jti:       uuid.NewString(),
		UserID:    userID,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(tokenExp),
	}
}

func (s *TokenService) GenerateTokenString(
	ctx context.Context,
	mToken model.Token,
	secretKey string,
) (string, error) {
	jwtToken := s.buildJWTString(mToken)
	tokenString, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *TokenService) buildJWTString(mToken model.Token) *jwt.Token {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        mToken.Jti,
			ExpiresAt: jwt.NewNumericDate(mToken.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(mToken.IssuedAt),
		},
		UserID: mToken.UserID,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (s *TokenService) ParseTokenString(tokenString, secretKey string) (*dto.Token, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
	if err != nil {
		return nil, pkgErrors.ErrTokenIsNotValid
	}

	if !token.Valid {
		return nil, pkgErrors.ErrTokenIsNotValid
	}

	jti, err := uuid.Parse(claims.ID)
	if err != nil {
		return nil, err
	}

	return &dto.Token{
		Jti:    jti,
		UserID: claims.UserID,
	}, nil
}

func (s *TokenService) GetTokenByJti(ctx context.Context, jti uuid.UUID) (*model.Token, error) {
	return s.repository.GetTokenByJti(ctx, jti)
}
