package main

import (
    "LeForum/internal/api"
    "LeForum/internal/auth"
    "LeForum/internal/storage"
    "html/template"
    "log"
    "net/http"
    "time"

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
        log.Printf("Erreur lors de l'exécution du template: %v", err)
        http.Error(w, "Erreur Interne du Serveur", http.StatusInternalServerError)
    }
}

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

    // Créer le ServeMux
    mux := http.NewServeMux()

    // Serveur de fichiers statiques
    fs := http.FileServer(http.Dir("web/static"))
    mux.Handle("/static/", http.StripPrefix("/static/", fs))

    // Routes principales
    mux.HandleFunc("/", api.Handler)
    mux.HandleFunc("/categories", api.CategoriesHandler)
    mux.HandleFunc("/post", api.PostHandler)
    mux.HandleFunc("/toggle-theme", api.ToggleThemeHandler)
    mux.HandleFunc("/createPost", api.CreatePostHandler)

    // Authentification
    authHandler := auth.NewHandler()
    authHandler.RegisterRoutes(mux)

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