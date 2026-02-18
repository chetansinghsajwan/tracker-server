package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"tracker-server/internal/controller"
	"tracker-server/internal/repo/pg"
	"tracker-server/internal/service"
)

func main() {
	// Load .env file if it exists
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found")
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
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

	// Initialize Services
	userService := service.NewUserService(userRepo)
	accountService := service.NewAccountService(accountRepo)
	categoryService := service.NewCategoryService(categoryRepo)
	tagService := service.NewTagService(tagRepo)
	transactionService := service.NewTransactionService(transactionRepo)

	// Initialize Controllers
	userController := controller.NewUserController(userService)
	accountController := controller.NewAccountController(accountService)
	categoryController := controller.NewCategoryController(categoryService)
	tagController := controller.NewTagController(tagService)
	transactionController := controller.NewTransactionController(transactionService)

	// Routes
	http.HandleFunc("/users", userController.Register)
	http.HandleFunc("/users/", userController.GetUser)

	http.HandleFunc("/accounts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			accountController.Create(w, r)
		case http.MethodGet:
			// Decide between List and Get based on path?
			// accountController.ListAndGet handles /accounts and /accounts/{id} logic
			// But ListAndGet assumes it's mounted at /accounts/??
			// Let's use ListAndGet for both for now, as it checks path length
			accountController.ListAndGet(w, r)
		default:
			accountController.ListAndGet(w, r) // Let it handle other methods
		}
	})
	// Handle /accounts/ specifically to ensure it goes to ListAndGet
	http.HandleFunc("/accounts/", accountController.ListAndGet)

	http.HandleFunc("/categories", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			categoryController.Create(w, r)
		case http.MethodGet:
			categoryController.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/tags", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			tagController.Create(w, r)
		case http.MethodGet:
			tagController.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/transactions", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			transactionController.Create(w, r)
		case http.MethodGet:
			transactionController.List(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintf(w, "Tracker Server API")
		} else {
			http.NotFound(w, r)
		}
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
