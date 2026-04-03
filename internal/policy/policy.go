package policy

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Role string
}

func CanCreateTransaction(u User) bool {
	return u.Role == "admin"
}

func CanViewDashboard(u User) bool {
	return true // all roles can view
}
