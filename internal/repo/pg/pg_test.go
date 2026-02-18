package pg

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"tracker-server/internal/repo/suite"
)

// setupTestDB creates a new test database or cleans up the existing one.
// For simplicity in this environment, we'll try to use the existing DB_URL
// but perhaps prefix tables or truncate them.
// WARN: Truncating tables in a shared DB is risky. ideally use a separate test DB.
func setupTestDB(t *testing.T) *sql.DB {
	// Try loading .env from project root relative path
	// Assuming test is run from internal/repo/pg
	_ = godotenv.Load("../../../.env")

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		t.Skip("DATABASE_URL not set, skipping integration tests")
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		t.Fatalf("Failed to open DB: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping DB: %v", err)
	}

	// Clean up tables to ensure fresh state
	// Order matters due to FKs
	tables := []string{"transaction_tags", "transactions", "tags", "categories", "accounts", "user_secrets", "users"}
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			// Ignore error if table doesn't exist or other issue, but log it
			t.Logf("Failed to truncate %s: %v", table, err)
		}
	}

	return db
}

func TestUserRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	repo := NewPostgresUserRepo(db)
	suite := &suite.UserRepoSuite{Repo: repo}
	suite.TestAll(t)
}

func TestAccountRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewPostgresUserRepo(db)
	repo := NewPostgresAccountRepo(db)
	suite := &suite.AccountRepoSuite{Repo: repo, UserRepo: userRepo}
	suite.TestAll(t)
}

func TestCategoryRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewPostgresUserRepo(db)
	repo := NewPostgresCategoryRepo(db)
	suite := &suite.CategoryRepoSuite{Repo: repo, UserRepo: userRepo}
	suite.TestAll(t)
}

func TestTagRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewPostgresUserRepo(db)
	repo := NewPostgresTagRepo(db)
	suite := &suite.TagRepoSuite{Repo: repo, UserRepo: userRepo}
	suite.TestAll(t)
}

func TestTransactionRepo(t *testing.T) {
	db := setupTestDB(t)
	defer db.Close()

	userRepo := NewPostgresUserRepo(db)
	accountRepo := NewPostgresAccountRepo(db)
	categoryRepo := NewPostgresCategoryRepo(db)
	repo := NewPostgresTransactionRepo(db)

	suite := &suite.TransactionRepoSuite{
		Repo:         repo,
		UserRepo:     userRepo,
		AccountRepo:  accountRepo,
		CategoryRepo: categoryRepo,
	}
	suite.TestAll(t)
}
