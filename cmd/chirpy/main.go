package main

import (
	"database/sql"
	"io"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/goinginblind/chirpy/internal/config"
	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/server"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func main() {
	// setup logging
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Fail to open the log file: %s", err)
	}
	defer logFile.Close()

	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// load enviromental variables
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("No .env file found, proceeding with system environment")
	}
	dbURL := os.Getenv("DB_URL")
	filepathRoot := os.Getenv("FILEPATH_ROOT")
	port := os.Getenv("PORT")
	platform := os.Getenv("PLATFORM")

	// connect to db
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening db: %s", err)
	}
	dbQueries := database.New(db)
	log.Println("Connected to database.")

	// setup config
	apiServ := server.Server{Cfg: config.APIConfig{
		FileserverHits: atomic.Int32{},
		DB:             dbQueries,
		Platform:       platform,
	}}

	// setup multiplexer and handles
	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app/", apiServ.MiddlewareMetricsInc(http.FileServer(http.Dir(filepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /api/healthz", server.HandlerReadiness)
	mux.HandleFunc("GET /admin/metrics", apiServ.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiServ.HandlerReset)
	mux.HandleFunc("POST /api/validate_chirp", server.HandlerValidate)
	mux.HandleFunc("POST /api/users", apiServ.HandlerCreateUser)

	// setup server
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
