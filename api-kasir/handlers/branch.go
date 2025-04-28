package handlers // Package handlers berisi implementasi handler untuk route API

import (
    "github.com/gofiber/fiber/v2" // Framework web Fiber untuk membuat API
    "github.com/google/uuid"      // Library untuk generate unique ID (UUID)
    "api-kasir/models"           // Package models berisi struktur data aplikasi
    "api-kasir/supabase"         // Package untuk berinteraksi dengan database Supabase
    "api-kasir/utils"
)

// GetBranches - Handler untuk mendapatkan semua data cabang
// Menerima parameter context dari Fiber
// Mengembalikan error jika terjadi masalah
func GetBranches(c *fiber.Ctx) error {
    // Memanggil fungsi GetBranches dari package supabase
    // untuk mengambil data semua cabang dari database
    branches, err := supabase.GetBranches() // Perbaikan: menggunakan := bukan ,

    // Jika terjadi error saat mengambil data
    if err != nil {
        // Mengembalikan response error 500 (Internal Server Error)
        // dengan pesan error dari supabase
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    // Jika sukses, kembalikan data branches dalam format JSON
    // Fiber akan otomatis mengkonversi struct ke JSON
    return c.JSON(branches)
}

// CreateBranch - Handler untuk membuat cabang baru
// Menerima parameter context dari Fiber
// Mengembalikan error jika terjadi masalah
func CreateBranch(c *fiber.Ctx) error {
    // Deklarasi variabel untuk menampung data input dari request
    var input models.CreateBranchInput

    // Parsing request body ke dalam struct CreateBranchInput
    if err := c.BodyParser(&input); err != nil { // Perbaikan: menggunakan := bukan =
        // Jika parsing gagal, kembalikan response error 400 (Bad Request)
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body", // Pesan error sederhana
        })
    }

    // Membuat objek Branch baru dengan:
    branch := models.Branch{ // Perbaikan: menggunakan := bukan :
        ID:      uuid.New().String(), // Generate UUID v4 sebagai ID unik
        Name:    input.Name,          // Nama cabang dari input
        Address: input.Address,       // Alamat cabang dari input
    }

    // Memanggil fungsi CreateBranch dari package supabase
    // untuk menyimpan data cabang baru ke database
    err := supabase.CreateBranch(branch) // Perbaikan: menggunakan := bukan :
    if err != nil {
        // Jika terjadi error saat menyimpan ke database
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(), // Mengembalikan pesan error dari supabase
        })
    }

    // Jika sukses, kembalikan response 201 (Created)
    // dengan pesan sukses dalam bahasa Indonesia
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "message": "Cabang berhasil ditambahkan",
    })
}


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

	err := supabase.UpdateBranch(id, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cabang berhasil diperbarui",
	})
}


func DeleteBranch(c *fiber.Ctx) error {
	id := c.Params("id")

	err := supabase.DeleteBranch(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Cabang berhasil dihapus",
	})
}


// Login handler untuk autentikasi
func Login(c *fiber.Ctx) error {
    var user map[string]string
    if err := c.BodyParser(&user); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
    }

    email := user["email"]
    password := user["password"]

    // Verifikasi email dan password dengan database (di sini hanya contoh)
    // Misalnya, kita punya data pengguna seperti berikut:
    dbEmail := "admin123@gmail.com"
    dbPassword := "password123" // Password yang di-hash sebaiknya
    dbRole := "admin"
    dbBranchID := "0e125bde-08ed-4018-86a1-3ea4435335e3"

    if email == dbEmail && password == dbPassword {
        // Generate JWT token
        token, err := utils.GenerateJWT(dbEmail, dbRole, dbBranchID)
        if err != nil {
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate token"})
        }

        return c.JSON(fiber.Map{
            "message": "Login successful",
            "token":   token,
        })
    }

    return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
}
