package repository

import "gorm.io/gorm"

type Repositories struct {
	User      *UserRepository
	Tx        *TransactionRepository
	Dashboard *DashboardRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		User:      NewUserRepository(db),
		Tx:        NewTransactionRepository(db),
		Dashboard: NewDashboardRepository(db),
	}
}
