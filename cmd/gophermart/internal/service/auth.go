package service

import (
	"context"
	"kialkuz/shop-with-loyalty/internal/contracts/repository"
	"kialkuz/shop-with-loyalty/internal/model"
	"kialkuz/shop-with-loyalty/internal/token"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 12

type AuthService struct {
	authRepository repository.AuthRepository
	userRepository repository.UserRepository
	pool           *pgxpool.Pool
}

func NewAuthService(
	authRepository repository.AuthRepository,
	userRepository repository.UserRepository,
	pool *pgxpool.Pool,
) *AuthService {
	return &AuthService{authRepository: authRepository, userRepository: userRepository, pool: pool}
}

func (s *AuthService) HashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)
	if err != nil {
		return nil, err
	}
	return hashedPassword, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, mToken model.Token, secretKey string) (string, error) {
	jwtToken := s.buildJWTString(mToken)
	tokenString, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) buildJWTString(mToken model.Token) *jwt.Token {
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

func (s *AuthService) RegisterAndLogin(ctx context.Context, user model.User, token model.Token) error {
	tx, err := s.pool.Begin(ctx)

	if err != nil {
		return err
	}

	err = s.userRepository.AddNewUserTx(ctx, tx, user)
	if err != nil {
		return err
	}

	err = s.authRepository.AddNewTokenTx(ctx, tx, token)
	if err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (s *AuthService) AddNewToken(ctx context.Context, token model.Token) error {
	return s.authRepository.AddNewToken(ctx, token)
}
