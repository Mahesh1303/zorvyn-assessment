package services

import (
	auth "finance-processing/internal/lib/utils"
	"finance-processing/internal/repository"
)

type Services struct {
	Transaction *TransactionService
	Dashboard   *DashboardService
	User        *UserService
	Auth        *AuthService
}

func NewServices(repos *repository.Repositories, jwtManager *auth.JWTManager) *Services {
	return &Services{
		Transaction: NewTransactionService(repos.Tx),
		Dashboard:   NewDashboardService(repos.Dashboard),
		User:        NewUserService(repos.User),
		Auth:        NewAuthService(repos.User, jwtManager),
	}
}
