package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"rudrprasad.com/backend/database"
	"rudrprasad.com/backend/logs"
	"rudrprasad.com/backend/routes"
)


func main() {
	
	router := mux.NewRouter()
	config := database.Config{
		Username: "root",
		Password: "",
		Host:     "127.0.0.1",
		Port:     3306,
		DbName:   "golangtest",
	}

	logger, err := logs.NewLogger()
	if err != nil {
		log.Fatalf("Error initializing logger: %v", err)
	}
	defer logger.Close()

	// Initialize the database
	db, err := database.InitDB(config)
	routes := &routes.Routes{DB: db, LOG: logger}
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()

	// err = database.CreateTableOnce(db, "users", database.CreateUserTable())
	// if err != nil {
	// 	log.Fatalf("Error initializing table: %v", err)
	// }


	// routes
	router.HandleFunc("/404", routes.Handle404)
	router.HandleFunc("/", routes.GetHome).Methods("GET")
	router.HandleFunc("/api/auth/register", routes.PostRegisterUser).Methods("POST")

	// redirect handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
	})


	// middleware
	corsRouter := routes.CorsMiddleware(router)
	loggedHandler := logs.LoggingMiddleware(logger, corsRouter)

	log.Fatal(http.ListenAndServe(":8080", loggedHandler))
}