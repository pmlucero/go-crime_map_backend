package security

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
)

// GenerateAPIKey genera una nueva API Key
func GenerateAPIKey() (string, error) {
	// Generar 32 bytes aleatorios
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error generando API Key: %v", err)
	}

	// Codificar en base64
	key := base64.StdEncoding.EncodeToString(bytes)
	return key, nil
}

// TODO: Implementar funciones para:
// - Almacenar API Keys en la base de datos
// - Validar API Keys contra la base de datos
// - Revocar API Keys
// - Listar API Keys por usuario
