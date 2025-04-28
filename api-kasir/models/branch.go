package models

import (
    "time"
    "encoding/json"
)
// structur branch mentode get
type Branch struct {
    ID        string    `json:"id"`       // ID akan di-generate otomatis oleh Supabase
    Name      string    `json:"name"`
    Address   string    `json:"address"`
    CreatedAt time.Time `json:"created_at"`
}

// structur branch metode get
type CreateBranchInput struct {
    Name    string `json:"name"`
    Address string `json:"address"`
}

 // parse waktu dengan format dari Supabase
func (b *Branch) UnmarshalJSON(data []byte) error {
    type Alias Branch
    aux := &struct {
        CreatedAt string `json:"created_at"`
        *Alias
    }{
        Alias: (*Alias)(b),
    }

    if err := json.Unmarshal(data, &aux); err != nil {
        return err
    }

    parsedTime, err := time.Parse("2006-01-02T15:04:05.999999", aux.CreatedAt)
    if err != nil {
        return err
    }

    b.CreatedAt = parsedTime
    return nil
}