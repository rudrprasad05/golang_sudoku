package main

import (
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/rudrprasad05/go-logs/logs"
	"rudrprasad.com/backend/database"
	"rudrprasad.com/backend/routes"
)

const size = 9
var board [size][size]int;

func DrawBoard(){
	width, height := 9,9
	// var board [width][height]int;

	for i := 0; i < height; i++{
		for j := 0; j < width; j++{
			fmt.Print(".")
			if (j + 1) % 3 == 0 {
				fmt.Print("|")
			} 
			
		}
		fmt.Print("\n")
		if (i + 1) % 3 == 0 {
			for j := 0; j < width; j++{
				fmt.Print("-")
				
			}
			fmt.Print("\n")
		}
		
	}
}

func PrintBoard(board [size][size]int) {
	for i, row := range board {
		if i%3 == 0 && i != 0 {
			fmt.Println("---------------------")
		}
		for j, val := range row {
			if j%3 == 0 && j != 0 {
				fmt.Print("| ")
			}
			fmt.Printf("%d ", val)
		}
		fmt.Println()
	}
}


func main() {
	PrintBoard(board)
	
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
		logger.Error("failed to create db")
	}
	defer db.Close()

	// routes
	router.HandleFunc("/404", routes.Handle404)
	router.HandleFunc("/", routes.GetHome).Methods("GET")
	router.HandleFunc("/auth/register", routes.PostRegisterUser).Methods("POST")
	router.HandleFunc("/auth/login", routes.PostLoginUser).Methods("POST")

	protected := router.PathPrefix("/game").Subrouter()
	protected.Use(routes.AuthMiddleware)

	protected.HandleFunc("", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Welcome to the protected route!"))
	}).Methods("GET")


	// redirect handler
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/404", http.StatusTemporaryRedirect)
	})
	// middleware
	corsRouter := routes.CorsMiddleware(router)
	loggedHandler := logs.LoggingMiddleware(logger, corsRouter)

	log.Fatal(http.ListenAndServe(":8080", loggedHandler))
}