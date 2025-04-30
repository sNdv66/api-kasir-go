package handlers

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"api-kasir/models"
	"api-kasir/supabase"
	"api-kasir/utils"
	"io"
	"bytes"
)

// GetBranches
func GetBranches(c *fiber.Ctx) error {
	branches, err := supabase.GetBranches()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	return c.JSON(branches)
}

// CreateBranch
func CreateBranch(c *fiber.Ctx) error {
	var input models.CreateBranchInput

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	branch := models.Branch{
		ID:      uuid.New().String(),
		Name:    input.Name,
		Address: input.Address,
	}

	if err := supabase.CreateBranch(branch); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
    "message": "Cabang berhasil ditambahkan",
    "branch": branch,
         })
    
}

// UpdateBranch
func UpdateBranch(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID is required",
		})
	}

	var input models.CreateBranchInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := supabase.UpdateBranch(id, input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cabang berhasil diperbarui",
	})
}

// DeleteBranch
func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := supabase.DeleteBranch(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cabang berhasil dihapus",
	})
}

// Login
func Login(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request",
		})
	}

	email := user.Email
	password := user.Password
	
	if email == "" || password == "" {
    return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
        "error": "Email and password are required",
    })
    }

	req, err := utils.NewRequest("GET", "users?email=eq."+email+"&select=id,email,password,role,branch_id", nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	res, err := utils.Client.Do(req)
	if err != nil || res.StatusCode != 200 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}
	defer res.Body.Close()

	var users []models.User
	if err := json.NewDecoder(res.Body).Decode(&users); err != nil || len(users) == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	foundUser := users[0]

	if err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid credentials",
		})
	}

	token, err := utils.GenerateJWT(foundUser.Email, foundUser.Role, foundUser.BranchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate token",
		})
	}
    
    // melakukan perubahan ... 
    
	/*return c.JSON(fiber.Map{
	"message":   "Login successful",
	"token":     token,
	"branch_id": foundUser.BranchID,
     })
     */
     
     return c.JSON(fiber.Map{
    "token": token,
    "user": fiber.Map{
        "id":        foundUser.ID,
        "role":      foundUser.Role,
        "branch_id": foundUser.BranchID,
    },
})
     
}





// yang di kerjakan sekarang kalau tidak cocok tinggal hapus 


func GetProductsByBranch(c *fiber.Ctx) error {
    branchID := c.Params("branch_id")
    query := "products?branch_id=eq." + branchID

    req, err := utils.NewRequest("GET", query, nil)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    body, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(body)
}




func CreateProduct(c *fiber.Ctx) error {
    body := c.Body()

    req, err := utils.NewRequest("POST", "products", bytes.NewReader(body))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    resBody, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(resBody)
}




func UpdateProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    path := "products?id=eq." + id

    req, err := utils.NewRequest("PATCH", path, bytes.NewReader(c.Body()))
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    resBody, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(resBody)
}




func DeleteProduct(c *fiber.Ctx) error {
    id := c.Params("id")
    path := "products?id=eq." + id

    req, err := utils.NewRequest("DELETE", path, nil)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }

    res, err := utils.Client.Do(req)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    defer res.Body.Close()

    resBody, _ := io.ReadAll(res.Body)
    return c.Status(res.StatusCode).Send(resBody)
}



