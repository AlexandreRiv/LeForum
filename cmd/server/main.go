package main

import (
	"LeForum/internal/api"
	"LeForum/internal/auth"
	"LeForum/internal/storage"
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
	// Initialisation de la base de données
	if err := storage.InitDB(); err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données: %v", err)
	}

	// Configuration du routeur
	mux := api.SetupRouter()

	// Nettoyage des sessions expirées périodiquement
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		for range ticker.C {
			auth.CleanExpiredSessions()
		}
	}()

	// Démarrage du serveur
	log.Println("Serveur démarré sur http://localhost:3002")
	if err := http.ListenAndServe(":3002", mux); err != nil {
		log.Fatal(err)
	}
}
