package main

import (
	"log"

	"github.com/goinginblind/chirpy/internal/app"

	_ "github.com/lib/pq"
)

func main() {
	if err := app.Run(); err != nil {
		log.Printf("Failed to run app: %v\n", err)
	}
}
