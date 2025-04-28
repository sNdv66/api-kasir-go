package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	SupabaseURL string
	APIKey      string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("File .env tidak ditemukan, menggunakan env dari sistem")
	}

	return &Config{
		SupabaseURL: os.Getenv("SUPABASE_URL"),
		APIKey:      os.Getenv("SUPABASE_API_KEY"),
	}
}
