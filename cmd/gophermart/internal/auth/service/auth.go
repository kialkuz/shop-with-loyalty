package service

import (
	"context"
	authDto "kialkuz/shop-with-loyalty/internal/auth/dto"
	authModel "kialkuz/shop-with-loyalty/internal/auth/model"
	balanceModel "kialkuz/shop-with-loyalty/internal/domain/balance/model"
	balanceServ "kialkuz/shop-with-loyalty/internal/domain/balance/service"
	userModel "kialkuz/shop-with-loyalty/internal/domain/user/model"
	userServ "kialkuz/shop-with-loyalty/internal/domain/user/service"
	infrastructureInterfaces "kialkuz/shop-with-loyalty/internal/infrastructure/interfaces"
	pkgErrors "kialkuz/shop-with-loyalty/pkg/errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"golang.org/x/crypto/bcrypt"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
}

const passwordCost = 12

type AuthService struct {
	tokenService       *TokenService
	userService        *userServ.UserService
	balanceService     *balanceServ.BalanceService
	transactionManager infrastructureInterfaces.TransactionManager
}

func NewAuthService(
	tokenService *TokenService,
	userService *userServ.UserService,
	balanceService *balanceServ.BalanceService,
	transactionManager infrastructureInterfaces.TransactionManager,
) *AuthService {
	return &AuthService{
		tokenService:       tokenService,
		userService:        userService,
		balanceService:     balanceService,
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

func (s *AuthService) RegisterAndLogin(
	ctx context.Context,
	user userModel.User,
	token authModel.Token,
	balance balanceModel.Balance,
) error {
	return s.transactionManager.WithinTransaction(ctx, func(tx pgx.Tx) error {
		err := s.userService.AddTx(ctx, tx, user)
		if err != nil {
			return err
		}

		err = s.tokenService.AddTx(ctx, tx, token)
		if err != nil {
			return err
		}

		err = s.balanceService.AddTx(ctx, tx, balance)
		if err != nil {
			return err
		}

		return nil
	})
}

func (s *AuthService) MakeToken(userId uuid.UUID, tokenExp time.Duration) authModel.Token {
	return authModel.Token{
		Jti:       uuid.NewString(),
		UserId:    userId,
		IssuedAt:  time.Now(),
		ExpiresAt: time.Now().Add(tokenExp),
	}
}

func (s *AuthService) GenerateTokenString(
	ctx context.Context,
	mToken authModel.Token,
	secretKey string,
) (string, error) {
	jwtToken := s.buildJWTString(mToken)
	tokenString, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) buildJWTString(mToken authModel.Token) *jwt.Token {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        mToken.Jti,
			ExpiresAt: jwt.NewNumericDate(mToken.ExpiresAt),
			IssuedAt:  jwt.NewNumericDate(mToken.IssuedAt),
		},
		UserID: mToken.UserId,
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
}

func (s *AuthService) AddNewToken(ctx context.Context, token authModel.Token) error {
	return s.tokenService.AddNewToken(ctx, token)
}

func (s *AuthService) ParseTokenString(tokenString, secretKey string) (*authDto.Token, error) {
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

	return &authDto.Token{
		Jti:    jti,
		UserId: claims.UserID,
	}, nil
}

func (s *AuthService) GetTokenByJti(ctx context.Context, jti uuid.UUID) (*authModel.Token, error) {
	return s.tokenService.GetTokenByJti(ctx, jti)
}
