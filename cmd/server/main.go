package main

import (
	"LeForum/internal/api"
	"LeForum/internal/config"
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("Pas de fichier .env chargé")
	}
}

func main() {
	// Chargement de la configuration
	appConfig, err := config.NewAppConfig()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la configuration: %v", err)
	}

	// Configuration du routeur avec la configuration de l'application
	mux := api.SetupRouter(appConfig)

	// Nettoyage des sessions expirées périodiquement
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			appConfig.SessionService.CleanExpiredSessions()
		}
	}()

	// Démarrage du serveur
	log.Println("Serveur démarré sur http://localhost:3002")
	if err := http.ListenAndServe(":3002", mux); err != nil {
		log.Fatal(err)
	}
}
