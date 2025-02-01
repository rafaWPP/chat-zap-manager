// handlers/instance.go
package handlers

import (
	"context"
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
	"log"
)

// Estrutura para criar uma nova instância
type CreateInstanceRequest struct {
    InstanceName string `json:"instanceName"`
}

// Estrutura padrão para respostas de erro
type ErrorResponse struct {
    Status  int      `json:"status"`
    Error   string   `json:"error"`
    Message []string `json:"message"`
}

// Função para criar uma nova instância
func CreateInstance(w http.ResponseWriter, r *http.Request) {
    // Obter credenciais dos cookies
    serverURL, apiKey, err := getCredentialsFromCookies(r)
    if err != nil || serverURL == "" || apiKey == "" {
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Obter dados do formulário
    instanceName := r.FormValue("instanceName")
    if instanceName == "" {
        http.Error(w, "Instance name is required", http.StatusBadRequest)
        return
    }

    // Preparar payload
    payload := CreateInstanceRequest{
        InstanceName: instanceName,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        http.Error(w, "Failed to encode payload", http.StatusInternalServerError)
        return
    }

    // Fazer requisição à API
    url := serverURL + "/instance/create"
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(payloadBytes))
    if err != nil {
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Accept", "application/json")
    req.Header.Add("apikey", apiKey)

    client := &http.Client{Timeout: 10 * time.Second}
    res, err := client.Do(req)
    if err != nil {
        http.Error(w, "Failed to send request to API", http.StatusInternalServerError)
        return
    }
    defer res.Body.Close()

    // Ler resposta
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        http.Error(w, "Failed to read API response", http.StatusInternalServerError)
        return
    }

    // Retornar resposta para o cliente
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(res.StatusCode)
    w.Write(body)
}

// Função para listar instâncias
func FetchInstances(w http.ResponseWriter, r *http.Request) {
    // Obter credenciais dos cookies
    serverURL, apiKey, err := getCredentialsFromCookies(r)
    if err != nil || serverURL == "" || apiKey == "" {
        log.Println("Credenciais ausentes ou inválidas:", err)
        http.Redirect(w, r, "/login", http.StatusSeeOther)
        return
    }

    // Fazer requisição à API
    url := serverURL + "/instance/fetchInstances"
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Println("Erro ao criar requisição:", err)
        http.Error(w, "Failed to create request", http.StatusInternalServerError)
        return
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("apikey", apiKey)

    client := &http.Client{Timeout: 10 * time.Second}
    res, err := client.Do(req)
    if err != nil {
        log.Println("Erro ao enviar requisição para a API:", err)
        http.Error(w, "Failed to send request to API", http.StatusInternalServerError)
        return
    }
    defer res.Body.Close()

    // Ler resposta
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Println("Erro ao ler resposta da API:", err)
        http.Error(w, "Failed to read API response", http.StatusInternalServerError)
        return
    }

    // Verificar se a resposta é JSON válida
    var jsonBody interface{}
    if err := json.Unmarshal(body, &jsonBody); err != nil {
        log.Printf("Resposta da API não é JSON válido: %s", string(body))
        http.Error(w, `{"error": "Resposta inválida da API"}`, http.StatusInternalServerError)
        return
    }

    // Retornar resposta para o cliente
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(res.StatusCode)
    w.Write(body)
}

func DeleteInstance(w http.ResponseWriter, r *http.Request) {
    log.Println("Recebida requisição para deletar instância.")

    // Obter credenciais dos cookies
    serverURL, apiKey, err := getCredentialsFromCookies(r)
    if err != nil || serverURL == "" || apiKey == "" {
        log.Printf("Erro ao obter credenciais: %v, serverURL: '%s', apiKey: '%s'\n", err, serverURL, apiKey)
        errorResponse := ErrorResponse{
            Status:  http.StatusUnauthorized,
            Error:   "Unauthorized",
            Message: []string{"Credenciais inválidas ou ausentes."},
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusUnauthorized)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }
    log.Printf("Credenciais obtidas: serverURL=%s, apiKey=***\n", serverURL)

    // Obter nome da instância da URL
    instanceName := r.URL.Query().Get("instanceName")
    if instanceName == "" {
        log.Println("Nome da instância não fornecido.")
        errorResponse := ErrorResponse{
            Status:  http.StatusBadRequest,
            Error:   "Bad Request",
            Message: []string{"Nome da instância é obrigatório."},
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusBadRequest)
        json.NewEncoder(w).Encode(errorResponse)
        return
    }

    // Realiza o logout da instância antes de deletá-la
    logoutInstance(serverURL, instanceName, apiKey)

    // Opcional: aguardar um breve intervalo para garantir que o logout foi processado
    time.Sleep(1 * time.Second)

    // Obter parâmetro 'force' da URL
    force := r.URL.Query().Get("force") == "true"
    log.Printf("Parâmetro 'force' recebido: %v\n", force)

    // Construir a URL para deletar a instância
    var apiURL string
    if force {
        apiURL = fmt.Sprintf("%s/instance/delete/%s?force=true", serverURL, instanceName)
    } else {
        apiURL = fmt.Sprintf("%s/instance/delete/%s", serverURL, instanceName)
    }
    log.Printf("Construindo requisição DELETE para URL: %s\n", apiURL)

    // Criar um contexto com timeout de 60 segundos
    ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
    defer cancel()

    // Criar requisição HTTP com o contexto
    req, err := http.NewRequestWithContext(ctx, "DELETE", apiURL, nil)
    if err != nil {
        log.Printf("Falha ao criar requisição DELETE: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Internal Server Error",
            Message: []string{"Falha ao criar a requisição."},
        })
        return
    }

    // Adicionar headers
    req.Header.Add("Accept", "application/json")
    req.Header.Add("apikey", apiKey)

    // Criar cliente HTTP e enviar a requisição
    client := &http.Client{}
    res, err := client.Do(req)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            log.Println("Erro: Tempo limite excedido na requisição à API externa.")
            w.WriteHeader(http.StatusGatewayTimeout)
            json.NewEncoder(w).Encode(ErrorResponse{
                Status:  http.StatusGatewayTimeout,
                Error:   "Gateway Timeout",
                Message: []string{"A API demorou muito para responder."},
            })
        } else {
            log.Printf("Falha ao enviar requisição para a API: %v\n", err)
            w.WriteHeader(http.StatusInternalServerError)
            json.NewEncoder(w).Encode(ErrorResponse{
                Status:  http.StatusInternalServerError,
                Error:   "Internal Server Error",
                Message: []string{"Falha ao enviar a requisição para a API externa."},
            })
        }
        return
    }
    defer res.Body.Close()
    log.Printf("Recebida resposta da API: %d %s\n", res.StatusCode, res.Status)

    // Ler resposta
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Printf("Falha ao ler resposta da API: %v\n", err)
        w.WriteHeader(http.StatusInternalServerError)
        json.NewEncoder(w).Encode(ErrorResponse{
            Status:  http.StatusInternalServerError,
            Error:   "Internal Server Error",
            Message: []string{"Falha ao ler a resposta da API externa."},
        })
        return
    }
    log.Printf("Corpo da resposta da API: %s\n", string(body))

    // Retornar resposta para o cliente
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(res.StatusCode)
    w.Write(body)
    log.Println("Resposta enviada ao cliente.", res.StatusCode)
}


// logoutInstance realiza o logout da instância antes de deletá-la.
func logoutInstance(serverURL, instanceName, apiKey string) error {
    // Constrói a URL para logout (substitua ":instanceName" pelo valor real)
    logoutURL := fmt.Sprintf("%s/instance/logout/%s", serverURL, instanceName)
    method := "DELETE"

    client := &http.Client{
        Timeout: 10 * time.Second,
    }
    req, err := http.NewRequest(method, logoutURL, nil)
    if err != nil {
        return fmt.Errorf("falha ao criar requisição de logout: %v", err)
    }
    req.Header.Add("Accept", "application/json")
    req.Header.Add("apikey", apiKey) // Se a API utiliza esse header, ou ajuste conforme necessário

    res, err := client.Do(req)
    if err != nil {
        return fmt.Errorf("falha ao enviar requisição de logout: %v", err)
    }
    defer res.Body.Close()

    // Opcional: ler a resposta e verificar se o logout foi bem-sucedido
    body, err := ioutil.ReadAll(res.Body)
    if err != nil {
        return fmt.Errorf("falha ao ler resposta do logout: %v", err)
    }
    log.Printf("Resposta do logout: %s\n", string(body))

    // Verifica se o status indica sucesso (por exemplo, 200 OK)
    if res.StatusCode != http.StatusOK {
        return fmt.Errorf("logout falhou com status %d", res.StatusCode)
    }

    return nil
}


type ConnectResponse struct {
	QRCodeURL string `json:"qrcodeUrl"`
	// Adicione outros campos se necessário.
}

// ConnectInstance faz a conexão com a API externa para obter o QR Code.
func ConnectInstance(w http.ResponseWriter, r *http.Request) {
	// Obter o nome da instância via query string
	instanceName := r.URL.Query().Get("instanceName")
	if instanceName == "" {
		http.Error(w, "instanceName is required", http.StatusBadRequest)
		return
	}

	// Obter credenciais dos cookies (se necessário)
	serverURL, apiKey, err := getCredentialsFromCookies(r)
	if err != nil || serverURL == "" || apiKey == "" {
		http.Error(w, "Credenciais inválidas ou ausentes", http.StatusUnauthorized)
		return
	}

	// Construir a URL da API de conexão. Supondo que o endpoint seja:
	// {serverURL}/instance/connect/{instanceName}
	apiURL := fmt.Sprintf("%s/instance/connect/%s", serverURL, instanceName)

	// Criar a requisição GET para a API
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Falha ao criar a requisição", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("apikey", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		http.Error(w, "Falha ao conectar à API", http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		http.Error(w, "Falha ao ler a resposta da API", http.StatusInternalServerError)
		return
	}

	// Encaminhar a resposta da API para o cliente
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.StatusCode)
	w.Write(body)
}