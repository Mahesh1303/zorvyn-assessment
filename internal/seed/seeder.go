package seeder

import (
	"fmt"
	"time"

	auth "finance-processing/internal/lib/utils"
	"finance-processing/internal/models"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type seedUser struct {
	Name     string
	Email    string
	Password string
	Role     models.UserRole
}

func RunSeed(db *gorm.DB, log zerolog.Logger) error {
	seedUsers := []seedUser{
		{Name: "Admin User", Email: "admin@test.com", Password: "password123", Role: models.RoleAdmin},
		{Name: "Analyst User", Email: "analyst@test.com", Password: "password123", Role: models.RoleAnalyst},
		{Name: "Viewer User", Email: "viewer@test.com", Password: "password123", Role: models.RoleViewer},
	}

	emails := make([]string, 0, len(seedUsers))
	for _, u := range seedUsers {
		emails = append(emails, u.Email)
	}

	now := time.Now().UTC()

	err := db.Transaction(func(tx *gorm.DB) error {
		for _, su := range seedUsers {
			hashedPassword, err := auth.EncryptPassWord(su.Password)
			if err != nil {
				return fmt.Errorf("failed to hash password for %s: %w", su.Email, err)
			}

			user := models.User{
				Name:     su.Name,
				Email:    su.Email,
				Password: hashedPassword,
				Role:     su.Role,
				IsActive: true,
			}

			err = tx.Clauses(clause.OnConflict{
				Columns: []clause.Column{{Name: "email"}},
				DoUpdates: clause.Assignments(map[string]any{
					"name":       user.Name,
					"password":   user.Password,
					"role":       user.Role,
					"is_active":  true,
					"deleted_at": nil,
				}),
			}).Create(&user).Error
			if err != nil {
				return fmt.Errorf("failed to upsert user %s: %w", su.Email, err)
			}
		}

		var users []models.User
		if err := tx.Where("email IN ?", emails).Find(&users).Error; err != nil {
			return fmt.Errorf("failed to load seeded users: %w", err)
		}

		userByEmail := make(map[string]models.User, len(users))
		userIDs := make([]uuid.UUID, 0, len(users))
		for _, u := range users {
			userByEmail[u.Email] = u
			userIDs = append(userIDs, u.ID)
		}

		if len(userIDs) != len(seedUsers) {
			return fmt.Errorf("failed to resolve all seeded users")
		}

		if err := tx.Exec("DELETE FROM financial_records WHERE created_by IN ?", userIDs).Error; err != nil {
			return fmt.Errorf("failed to clear existing seeded records: %w", err)
		}

		records := []map[string]any{
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 5000.00, "type": "income", "category": "salary", "description": "Monthly salary", "date": now.AddDate(0, -2, -2)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 1200.00, "type": "expense", "category": "rent", "description": "Apartment rent", "date": now.AddDate(0, -2, -1)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 350.00, "type": "expense", "category": "groceries", "description": "Groceries", "date": now.AddDate(0, -2, -5)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 200.00, "type": "expense", "category": "transport", "description": "Fuel", "date": now.AddDate(0, -2, -7)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 5050.00, "type": "income", "category": "salary", "description": "Monthly salary", "date": now.AddDate(0, -1, -2)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 1250.00, "type": "expense", "category": "rent", "description": "Apartment rent", "date": now.AddDate(0, -1, -1)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 420.00, "type": "expense", "category": "groceries", "description": "Groceries", "date": now.AddDate(0, -1, -5)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 180.00, "type": "expense", "category": "utilities", "description": "Electricity bill", "date": now.AddDate(0, -1, -8)},
			{"created_by": userByEmail["analyst@test.com"].ID, "amount": 3200.00, "type": "income", "category": "freelance", "description": "Consulting project", "date": now.AddDate(0, -1, -10)},
			{"created_by": userByEmail["analyst@test.com"].ID, "amount": 650.00, "type": "expense", "category": "software", "description": "Tool subscriptions", "date": now.AddDate(0, -1, -11)},
			{"created_by": userByEmail["analyst@test.com"].ID, "amount": 430.00, "type": "expense", "category": "travel", "description": "Client visit", "date": now.AddDate(0, -1, -12)},
			{"created_by": userByEmail["viewer@test.com"].ID, "amount": 2800.00, "type": "income", "category": "salary", "description": "Monthly salary", "date": now.AddDate(0, 0, -15)},
			{"created_by": userByEmail["viewer@test.com"].ID, "amount": 900.00, "type": "expense", "category": "rent", "description": "Shared rent", "date": now.AddDate(0, 0, -14)},
			{"created_by": userByEmail["viewer@test.com"].ID, "amount": 280.00, "type": "expense", "category": "groceries", "description": "Weekly groceries", "date": now.AddDate(0, 0, -12)},
			{"created_by": userByEmail["viewer@test.com"].ID, "amount": 120.00, "type": "expense", "category": "transport", "description": "Bus pass", "date": now.AddDate(0, 0, -10)},
			{"created_by": userByEmail["analyst@test.com"].ID, "amount": 1500.00, "type": "income", "category": "bonus", "description": "Quarterly bonus", "date": now.AddDate(0, 0, -8)},
			{"created_by": userByEmail["analyst@test.com"].ID, "amount": 300.00, "type": "expense", "category": "education", "description": "Course fee", "date": now.AddDate(0, 0, -6)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 450.00, "type": "expense", "category": "health", "description": "Medical checkup", "date": now.AddDate(0, 0, -4)},
			{"created_by": userByEmail["admin@test.com"].ID, "amount": 600.00, "type": "income", "category": "investment", "description": "Dividend income", "date": now.AddDate(0, 0, -3)},
			{"created_by": userByEmail["viewer@test.com"].ID, "amount": 220.00, "type": "expense", "category": "entertainment", "description": "Weekend outing", "date": now.AddDate(0, 0, -2)},
		}

		if err := tx.Table("financial_records").Create(&records).Error; err != nil {
			return fmt.Errorf("failed to insert seed records: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	log.Info().Int("users", len(seedUsers)).Int("records", 20).Msg("seed data inserted")
	return nil
}
