// config/config.go
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	// Variáveis de configuração globais
	AppName string
	Port    string
	// Você pode adicionar outras variáveis conforme necessário
)

// init é executado automaticamente quando o pacote é importado
func init() {
	// Tenta carregar o arquivo .env, se existir
	if err := godotenv.Load(); err != nil {
		log.Println("Nenhum arquivo .env encontrado, usando variáveis de ambiente existentes")
	}

	// Recupera as variáveis de ambiente e define valores padrão se necessário
	AppName = os.Getenv("APP_NAME")
	if AppName == "" {
		AppName = "MeuApp" // Valor padrão
	}

	Port = os.Getenv("PORT")
	if Port == "" {
		Port = "8080" // Valor padrão
	}
}
