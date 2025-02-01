// handlers/auth.go
package handlers

import (
    "html/template"
    "net/http"
    "time"
    "os"
)

var (
    tpl = template.Must(template.ParseFiles("templates/login.html", "templates/dashboard.html"))
)

type Credentials struct {
    ServerURL string `json:"server_url"`
    APIKey    string `json:"api_key"`
}


// Renderizar Templates
func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
    err := tpl.ExecuteTemplate(w, tmpl+".html", data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

// Handler para Login (GET e POST)
func Login(w http.ResponseWriter, r *http.Request) {
    if r.Method == http.MethodGet {
        // Verificar se o usuário já está logado
        serverURL, apiKey, err := getCredentialsFromCookies(r)
        if err == nil && serverURL != "" && apiKey != "" {
            http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
            return
        }

        appName := os.Getenv("APP_NAME")

        // Criar a estrutura com o AppName
        data := PageData{
            AppName: appName,
        }

        
        renderTemplate(w, "login", data)
    } else if r.Method == http.MethodPost {
        // Obter dados do formulário
        serverURL := r.FormValue("server_url")
        apiKey := r.FormValue("api_key")

        // Validar entradas (simples)
        if serverURL == "" || apiKey == "" {
            renderTemplate(w, "login", "Por favor, preencha todos os campos.")
            return
        }

        // Criar credenciais
        creds := Credentials{
            ServerURL: serverURL,
            APIKey:    apiKey,
        }

        // Definir cookies para ServerURL e APIKey
        setCookie(w, "server_url", creds.ServerURL, 24*time.Hour)
        setCookie(w, "api_key", creds.APIKey, 24*time.Hour)

        // Redirecionar para dashboard
        http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
    }
}

type PageData struct {
    AppName string
}

// Handler para Dashboard
func Dashboard(w http.ResponseWriter, r *http.Request) {
    // Obter credenciais dos cookies
    serverURL, apiKey, err := getCredentialsFromCookies(r)
    if err != nil || serverURL == "" || apiKey == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    appName := os.Getenv("APP_NAME")

    data := struct {
        Credentials
        PageData
    }{
        Credentials: Credentials{
            ServerURL: serverURL,
            APIKey:    apiKey,
        },
        PageData: PageData{
            AppName: appName,
        },
    }

    // Renderizar dashboard com dados das credenciais
    renderTemplate(w, "dashboard", data)
}

// Handler para Logout
func Logout(w http.ResponseWriter, r *http.Request) {
    // Remover cookies de ServerURL e APIKey
    removeCookie(w, "server_url")
    removeCookie(w, "api_key")

    // Redirecionar para login
    http.Redirect(w, r, "/login", http.StatusSeeOther)
}

// Função para definir um cookie
func setCookie(w http.ResponseWriter, name, value string, maxAge time.Duration) {
    http.SetCookie(w, &http.Cookie{
        Name:     name,
        Value:    value,
        Path:     "/",
        Expires:  time.Now().Add(maxAge),
        HttpOnly: true, // Recomendado para segurança
        // Secure:   true, // Descomente se usar HTTPS
    })
}

// Função para remover um cookie
func removeCookie(w http.ResponseWriter, name string) {
    http.SetCookie(w, &http.Cookie{
        Name:     name,
        Value:    "",
        Path:     "/",
        Expires:  time.Unix(0, 0),
        MaxAge:   -1,
        HttpOnly: true,
        // Secure:   true, // Descomente se usar HTTPS
    })
}

// Função para obter credenciais dos cookies
func getCredentialsFromCookies(r *http.Request) (serverURL, apiKey string, err error) {
    serverCookie, err := r.Cookie("server_url")
    if err != nil {
        return
    }
    apiCookie, err := r.Cookie("api_key")
    if err != nil {
        return
    }
    serverURL = serverCookie.Value
    apiKey = apiCookie.Value
    return
}
