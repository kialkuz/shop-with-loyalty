package service

import (
	"context"
	authDto "kialkuz/shop-with-loyalty/internal/dto/auth"
	"kialkuz/shop-with-loyalty/internal/model"
	"kialkuz/shop-with-loyalty/internal/token"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 12

type AuthService struct {
	tokenRepository    TokenRepository
	userRepository     UserRepository
	transactionManager TransactionManager
}

func NewAuthService(
	tokenRepository TokenRepository,
	userRepository UserRepository,
	transactionManager TransactionManager,
) *AuthService {
	return &AuthService{
		tokenRepository:    tokenRepository,
		userRepository:     userRepository,
		transactionManager: transactionManager,
	}
}

func (s *AuthService) HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func (s *AuthService) RegisterAndLogin(ctx context.Context, user model.User, token *model.Token) error {
	return s.transactionManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		err := s.userRepository.AddNewUserTx(ctx, tx, user)
		if err != nil {
			return err
		}

		err = s.tokenRepository.AddNewTokenTx(ctx, tx, token)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) MakeToken(userId uuid.UUID, tokenExp time.Duration) *model.Token {
	return &model.Token{
		Jti:       uuid.NewString(),
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(tokenExp),
	}
}

func (s *AuthService) GenerateTokenString(ctx context.Context, mToken *model.Token, secretKey string) (string, error) {
	jwtToken := s.buildJWTString(mToken)
	tokenString, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) buildJWTString(mToken *model.Token) *jwt.Token {
	claims := token.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        mToken.Jti,
			ExpiresAt: jwt.NewNumericDate(mToken.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(mToken.IssuedAt),
		},
		UserID: mToken.UserId,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (s *AuthService) AddNewToken(ctx context.Context, token *model.Token) error {
	return s.tokenRepository.AddNewToken(ctx, token)
}

func (s *AuthService) ParseTokenString(tokenString, secretKey string) (*authDto.Token, error) {
	claims := &token.Claims{}
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

	return &authDto.Token{
		Jti:    jti,
		UserId: claims.UserID,
	}, nil
}

func (s *AuthService) GetTokenByJti(ctx context.Context, jti uuid.UUID) (*model.Token, error) {
	return s.tokenRepository.GetTokenByJti(ctx, jti)
}
