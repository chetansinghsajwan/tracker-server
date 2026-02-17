package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"tracker-server/internal/repo/pg"
)

func main() {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://user:password@localhost:5432/money_tracker?sslmode=disable"
	}

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	fmt.Println("Connected to database")

	// Initialize Repositories
	userRepo := pg.NewPostgresUserRepo(db)
	accountRepo := pg.NewPostgresAccountRepo(db)
	categoryRepo := pg.NewPostgresCategoryRepo(db)
	tagRepo := pg.NewPostgresTagRepo(db)
	transactionRepo := pg.NewPostgresTransactionRepo(db)

	// Prevent unused variable errors
	fmt.Printf("Initialized repositories: User(%v), Account(%v), Category(%v), Tag(%v), Transaction(%v)\n",
		userRepo != nil, accountRepo != nil, categoryRepo != nil, tagRepo != nil, transactionRepo != nil)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server starting on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error starting server: %s\n", err)
	}
}
