package main

import (
	"LeForum/internal/auth"
	"html/template"
	"log"
	"net/http"
	"github.com/joho/godotenv"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	err := tmpl.Execute(w, nil)
	if err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Pas de fichier .env chargé")
	}
}

func main() {
	// Créer le multiplexeur
	mux := http.NewServeMux()

	// Serveur de fichiers statiques
	fs := http.FileServer(http.Dir("web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Route de base
	mux.HandleFunc("/", indexHandler)

	// Configuration de l'authentification
	authHandler := auth.NewHandler()
	authHandler.RegisterRoutes(mux)

	// Démarrage du serveur
	log.Println("Serveur démarré sur http://localhost:3002")
	if err := http.ListenAndServe(":3002", mux); err != nil {
		log.Fatal(err)
	}
}
