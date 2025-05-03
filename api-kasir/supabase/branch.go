package supabase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"api-kasir/models"
	"api-kasir/utils"
)

func GetBranches() ([]models.Branch, error) {
	req, err := utils.NewRequest("GET", "branches?select=*",nil)
	if err != nil {
		return nil, err
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	fmt.Println("Status Code:", res.StatusCode)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	fmt.Println("Response Body:", string(body))

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Gagal fetch branches: %s", body)
	}

	var branches []models.Branch
	err = json.Unmarshal(body, &branches)
	return branches, err
}

func CreateBranch(branch models.Branch) error {
	body, err := json.Marshal(branch)
	if err != nil {
		return err
	}

	req, err := utils.NewRequest("POST", "branches", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := utils.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated && res.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Gagal menambah cabang: %s", string(respBody))
	}

	return nil
}



func UpdateBranch(id string, input models.CreateBranchInput) error {
	body, err := json.Marshal(map[string]interface{}{
		"name":    input.Name,
		"address": input.Address,
	})
	if err != nil {
		return err
	}

	req, err := utils.NewRequest("PATCH", "branches?id=eq."+id, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := utils.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusNoContent {
		respBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Gagal update cabang: %s", string(respBody))
	}

	return nil
}


func DeleteBranch(id string) error {
	req, err := utils.NewRequest("DELETE", "branches?id=eq."+id, nil)
	if err != nil {
		return err
	}

	res, err := utils.Client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent && res.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(res.Body)
		return fmt.Errorf("Gagal menghapus cabang: %s", string(respBody))
	}

	return nil
}


