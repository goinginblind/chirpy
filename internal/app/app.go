// Provides the Run() function, which is then run in main and handles all the app logic.
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
	"github.com/goinginblind/chirpy/internal/handlers/admin"
	"github.com/goinginblind/chirpy/internal/handlers/chirps"
	"github.com/goinginblind/chirpy/internal/handlers/hooks"
	"github.com/goinginblind/chirpy/internal/handlers/tokens"
	"github.com/goinginblind/chirpy/internal/handlers/users"
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
	cfg, err := loadConfigFromEnv()
	if err != nil {
		log.Printf("error loading environment variables or parsing config: %v\n", err)
		return err
	}

	// connect to db
	db, err := sql.Open("postgres", cfg.DBUrl)
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

	// setup db into config
	cfg.DB = dbQueries

	// setup multiplexer and handles
	mux := http.NewServeMux()
	fsHandler := http.StripPrefix("/app/", cfg.MiddlewareMetricsInc(http.FileServer(http.Dir(cfg.FilepathRoot))))
	mux.Handle("/app/", fsHandler)

	// Admin handles
	mux.HandleFunc("GET /admin/metrics", cfg.InjectConfig(admin.Metrics))
	mux.HandleFunc("GET /api/healthz", admin.HandlerReadiness)
	mux.HandleFunc("POST /admin/reset", cfg.InjectConfig(admin.Reset))

	// Users handles
	mux.HandleFunc("POST /api/users", cfg.InjectConfig(users.Create))
	mux.HandleFunc("POST /api/login", cfg.InjectConfig(users.Login))
	mux.HandleFunc("PUT /api/users", cfg.InjectConfig(users.ChangeLoginInfo))

	// Tokens handles
	mux.HandleFunc("POST /api/refresh", cfg.InjectConfig(tokens.RefreshAccessToken))
	mux.HandleFunc("POST /api/revoke", cfg.InjectConfig(tokens.RevokeRefreshToken))

	// Chirps handles
	mux.HandleFunc("POST /api/chirps", cfg.InjectConfig(chirps.Create))
	mux.HandleFunc("GET /api/chirps", cfg.InjectConfig(chirps.GetChirps))
	mux.HandleFunc("GET /api/chirps/{chirpID}", cfg.InjectConfig(chirps.GetOneByID))
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", cfg.InjectConfig(chirps.DeleteOneByID))

	// Webhooks
	mux.HandleFunc("POST /api/polka/webhooks", cfg.InjectConfig(hooks.UpgradeToChirpyRed))

	// setup server
	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: mux,
	}

	// run it with a proper shutdown just in case
	go func() {
		log.Printf("Serving files from %s on port: %s\n", cfg.FilepathRoot, cfg.Port)
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
