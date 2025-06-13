package main

import (
	"LeForum/internal/api"
	"LeForum/internal/config"
	"LeForum/internal/service"
	"LeForum/internal/storage/repositories"
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
	reportService := service.NewReportService(repositories.NewReportRepository(appConfig.DB.DB))

	mux := api.SetupRouter(appConfig, reportService)

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
