package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"pulverizacao-api/config"
	"pulverizacao-api/database"
	"pulverizacao-api/graphql"

	"github.com/gorilla/mux"
	"github.com/graphql-go/handler"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Configurar aplicação
	cfg := config.Load()

	// Conectar ao MongoDB
	client, err := database.Connect(cfg.MongoURI)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer client.Disconnect(context.TODO())

	// Configurar banco de dados
	db := client.Database(cfg.DatabaseName)

	// Criar schema GraphQL
	schema, err := graphql.CreateSchema(db)
	if err != nil {
		log.Fatal("Failed to create GraphQL schema:", err)
	}

	// Configurar handler GraphQL
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	// Configurar rotas
	router := mux.NewRouter()
	router.Handle("/graphql", h).Methods("GET", "POST")
	router.Handle("/", http.RedirectHandler("/graphql", http.StatusFound)).Methods("GET")

	// CORS middleware
	corsHandler := func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}

			h.ServeHTTP(w, r)
		})
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on port %s", port)
	log.Printf("GraphiQL available at http://localhost:%s/graphql", port)

	log.Fatal(http.ListenAndServe(":"+port, corsHandler(router)))
}
