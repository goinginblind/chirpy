package app

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/goinginblind/chirpy/internal/database"
	"github.com/goinginblind/chirpy/internal/server"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

func Run() error {
	// setup logging
	logFile, err := os.OpenFile("server.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fail to open log file: %v\n", err)
		return err
	}
	defer logFile.Close()
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))

	// load environmental variables and parse config from env
	_ = godotenv.Load()
	apiConfig, err := loadConfigFromEnv()
	if err != nil {
		log.Printf("error loading environment variables or parsing config: %v\n", err)
		return err
	}

	// connect to db
	db, err := sql.Open("postgres", apiConfig.DBUrl)
	if err != nil {
		log.Printf("Error opening db: %s\n", err)
		return err
	}
	defer db.Close()
	dbQueries := database.New(db)
	if err := db.Ping(); err != nil {
		log.Printf("Error connecting to the db: %s\n", err)
		return err
	}
	log.Println("Connected to database.")

	// setup db into config with a server wrapper used by handlers
	apiConfig.DB = dbQueries
	apiServ := server.Server{Cfg: apiConfig}

	// setup multiplexer and handles
	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app/", apiServ.MiddlewareMetricsInc(http.FileServer(http.Dir(apiConfig.FilepathRoot))))
	mux.Handle("/app/", fsHandler)

	mux.HandleFunc("GET /admin/metrics", apiServ.HandlerMetrics)
	mux.HandleFunc("POST /admin/reset", apiServ.HandlerReset)

	mux.HandleFunc("GET /api/healthz", server.HandlerReadiness)
	mux.HandleFunc("POST /api/users", apiServ.HandlerCreateUser)
	mux.HandleFunc("POST /api/login", apiServ.HandlerLogin)

	mux.HandleFunc("POST /api/chirps", apiServ.HandlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiServ.HandlerGetAllChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiServ.HandlerGetChirpByID)

	// setup server
	srv := &http.Server{
		Addr:    ":" + apiConfig.Port,
		Handler: mux,
	}

	// run it with a proper shutdown just in case
	go func() {
		log.Printf("Serving files from %s on port: %s\n", apiConfig.FilepathRoot, apiConfig.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("ListenAndServe(): %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown fail: %v\n", err)
		return err
	}

	log.Printf("Server shut down as expected")
	return nil
}
