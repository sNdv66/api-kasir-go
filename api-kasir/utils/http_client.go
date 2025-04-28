package utils

import (
	"io"
	"net/http"
	"os"
	"time"
)

var Client *http.Client
var SupabaseUrl string
var SupabaseKey string

func Init() {
	SupabaseUrl = os.Getenv("SUPABASE_URL")
	SupabaseKey = os.Getenv("SUPABASE_API_KEY")

	Client = &http.Client{
		Timeout: 10 * time.Second,
	}
}

// Update fungsi ini untuk menerima body
func NewRequest(method, path string, body io.Reader) (*http.Request, error) {
	url := SupabaseUrl + "/rest/v1/" + path
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	
	req.Header.Set("apikey", SupabaseKey)
	req.Header.Set("Authorization", "Bearer "+SupabaseKey)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}