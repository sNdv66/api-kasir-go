package database

import (
    "github.com/jmoiron/sqlx"
    "time" // Import time package
)

var DB *sqlx.DB

// Struct untuk menyimpan hasil JOIN user dan branch
type UserWithBranch struct {
    ID         string    `db:"id"`
    Email      string    `db:"email"`
    Role       string    `db:"role"`
    BranchID   string    `db:"branch_id"`
    BranchName string    `db:"branch_name"`
    CreatedAt  time.Time `db:"created_at"`
}

// Mengambil data user berdasarkan email
func GetUserByEmail(email string) (UserWithBranch, error) {
    var user UserWithBranch
    err := DB.Get(&user, `
    SELECT u.id, u.email, u.role, u.branch_id, u.created_at, b.name AS branch_name
    FROM users u
    LEFT JOIN branches b ON u.branch_id = b.id
    WHERE u.email = $1
    `, email)

    return user, err
}