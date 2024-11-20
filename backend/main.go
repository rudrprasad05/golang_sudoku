package main

import (
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"rudrprasad.com/backend/database"
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

	// Initialize the database
	db, err := database.InitDB(config)
	routes := &routes.Routes{DB: db}
	if err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer db.Close()


	// err = database.CreateTableOnce(db, "users", database.CreateUserTable())
	// if err != nil {
	// 	log.Fatalf("Error initializing table: %v", err)
	// }

	router.HandleFunc("/", routes.GetHome).Methods("GET")
	router.HandleFunc("/api/auth/register", routes.PostRegisterUser).Methods("POST")
	// router.HandleFunc("/users", createUserHandler).Methods("POST")
	// router.HandleFunc("/users/{id:[0-9]+}", getUserByIDHandler).Methods("GET")

	corsRouter := routes.CorsMiddleware(router)

	log.Fatal(http.ListenAndServe(":8080", corsRouter))
}