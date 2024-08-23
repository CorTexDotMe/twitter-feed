package internal

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"twitter-feed/internal/database"
	"twitter-feed/internal/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func Run() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal(err)
	}
	port := os.Getenv("SERVER_PORT")

	DB := database.NewPostgresConnection()
	messageHandler := handler.MessageHandler{DB: DB}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", messageHandler.GetMessage)
	router.Post("/", messageHandler.CreateMessage)

	fmt.Printf("Server started at port %s", port)
	http.ListenAndServe(":"+port, router)
}
