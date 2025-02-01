// main.go
package main

import (
    "log"
    "net/http"
    "os"

    "github.com/gorilla/mux"
    "github.com/joho/godotenv"
    "whatsapp-manager/handlers"
)

func init() {
	// Carrega as variáveis de ambiente de um arquivo .env, se existir.
	err := godotenv.Load()
	if err != nil {
		log.Println("Não foi possível carregar o arquivo .env. As variáveis de ambiente devem estar configuradas no sistema.")
	}
}

func main() {
    // Carregar variáveis de ambiente
    err := godotenv.Load()
    if err != nil {
        log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente existentes")
    }

    // Configurar roteador
    r := mux.NewRouter()

    // Servir arquivos estáticos
    fs := http.FileServer(http.Dir("./static/"))
    r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

    // Rotas
    r.HandleFunc("/", handlers.Login).Methods("GET")
    r.HandleFunc("/login", handlers.Login).Methods("GET", "POST")
    r.HandleFunc("/dashboard", handlers.Dashboard).Methods("GET")
    r.HandleFunc("/logout", handlers.Logout).Methods("GET")
    r.HandleFunc("/createInstance", handlers.CreateInstance).Methods("POST") // Adicionado Methods("POST")
    r.HandleFunc("/fetchInstances", handlers.FetchInstances).Methods("GET")  // Adicionado Methods("GET")
    r.HandleFunc("/deleteInstance", handlers.DeleteInstance).Methods("DELETE") // Adicionado Methods("DELETE")
    r.HandleFunc("/instance/connect", handlers.ConnectInstance).Methods("GET")

    // Iniciar servidor
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    log.Printf("Servidor iniciado na porta %s", port)
    log.Fatal(http.ListenAndServe(":"+port, r))
}