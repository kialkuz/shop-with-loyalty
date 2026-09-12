package usecase

import (
	"context"
	authModel "kialkuz/shop-with-loyalty/internal/auth/model"
	tokenServ "kialkuz/shop-with-loyalty/internal/auth/service"
	"kialkuz/shop-with-loyalty/internal/infrastructure"
	balanceModel "kialkuz/shop-with-loyalty/internal/modules/balance/model"
	balanceServ "kialkuz/shop-with-loyalty/internal/modules/balance/service"
	userModel "kialkuz/shop-with-loyalty/internal/modules/user/model"
	userServ "kialkuz/shop-with-loyalty/internal/modules/user/service"

	"github.com/jackc/pgx/v5"
)

type RegisterAndLogin struct {
	userService        *userServ.UserService
	tokenService       *tokenServ.TokenService
	balanceService     *balanceServ.BalanceService
	transactionManager infrastructure.Transaction
}

func NewRegisterAndLoginUseCase(
	userService *userServ.UserService,
	tokenService *tokenServ.TokenService,
	balanceService *balanceServ.BalanceService,
	transactionManager infrastructure.Transaction,
) *RegisterAndLogin {
	return &RegisterAndLogin{
		tokenService:       tokenService,
		userService:        userService,
		balanceService:     balanceService,
		transactionManager: transactionManager,
	}
}

func (s *RegisterAndLogin) Execute(
	ctx context.Context,
	user userModel.User,
	token authModel.Token,
	balance balanceModel.Balance,
) error {
	return s.transactionManager.RunTransaction(ctx, func(tx pgx.Tx) error {
		if err := s.userService.AddTx(ctx, tx, user); err != nil {
			return err
		}

		if err := s.tokenService.AddTx(ctx, tx, token); err != nil {
			return err
		}

		if err := s.balanceService.AddTx(ctx, tx, balance); err != nil {
			return err
		}

		return nil
	})
}
