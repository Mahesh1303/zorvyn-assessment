package policy

import "github.com/google/uuid"

type User struct {
	ID   uuid.UUID
	Role string
}

// for user Management
func CanCreateUser(u User) bool {
	return u.Role == "admin"
}

func CanManageUsers(u User) bool {
	return u.Role == "admin"
}

// for doing transactions
func CanManageTransaction(u User) bool {
	return u.Role == "admin"
}

func CanViewTransaction(u User) bool {
	return u.Role == "admin" || u.Role == "analyst" || u.Role == "viewer"
}

// for dashvoard view
func CanViewDashboard(u User) bool {
	return u.Role == "admin" || u.Role == "analyst" || u.Role == "viewer"
}

func CanViewAnalytics(u User) bool {
	return u.Role == "admin" || u.Role == "analyst"
}
