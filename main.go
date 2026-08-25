package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/JavierP/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	// Constants
	const filepathRoot = "."
	const port = "8080"

	// virtual env and db connection
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	platf := os.Getenv("PLATFORM")
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL must be set")
	}

	dbConn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening database: %s", err)
	}

	dbQueries := database.New(dbConn)

	// starting the config
	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             dbQueries,
		platform:       platf,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	mux.HandleFunc("GET /api/healthz", readinessHandler)
	mux.HandleFunc("GET /api/chirps", apiCfg.getChirpsHandler)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.getChirpsByIDHandler)
	mux.HandleFunc("POST /api/chirps", apiCfg.createChirpHandler)
	mux.HandleFunc("POST /api/users", apiCfg.createUserHandler)

	mux.HandleFunc("GET /admin/metrics", apiCfg.logHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetHandler)

	httpServer := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)

	err = httpServer.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

// check for auth
