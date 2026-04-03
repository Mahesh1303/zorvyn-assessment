package services

import "finance-processing/internal/repository"

type Services struct {
	Transaction *TransactionService
	Dashboard   *DashboardService
	User        *UserService
	Auth        *AuthService
}

func NewServices(repos *repository.Repositories) *Services {
	return &Services{
		Transaction: NewTransactionService(repos.Tx),
		Dashboard:   NewDashboardService(repos.Dashboard),
		User:        NewUserService(repos.User),
		Auth:        NewAuthService(repos.User),
	}
}
