// internal/handlers/handlers.go

package handlers

import "finance-processing/internal/services"

type Handlers struct {
	Transaction *TransactionHandler
	User        *UserHandler
	Dashboard   *DashboardHandler
	AuthHandler *AuthHandler
}

func NewHandlers(s *services.Services) *Handlers {
	return &Handlers{
		Transaction: NewTransactionHandler(s.Transaction),
		User:        NewUserHandler(s.User),
		Dashboard:   NewDashboardHandler(s.Dashboard),
		AuthHandler: NewAuthHandler(s.Auth),
	}
}
