package main

import (
	"os"
	"time"

	"github.com/RockyLogic/NDAX-Websocket-Client/src/ndaxClient"
	"github.com/joho/godotenv"
)

func loadEnv() (string, string) {
	godotenv.Load()
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")
	return apiKey, secretKey
}

func main() {
	apiKey, secretKey := loadEnv()
	ndaxClient.Start(apiKey, secretKey)
	// Sleep for 5 seconds
	time.Sleep(5 * time.Second)
}
