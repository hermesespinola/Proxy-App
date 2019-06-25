package utils

import env "github.com/joho/godotenv"

// LoadEnv loads .env file
func LoadEnv() {
	env.Load()
}
