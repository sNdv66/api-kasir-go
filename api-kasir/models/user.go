package models

import "time"

type User struct {
    ID        string    `db:"id"`
    Email     string    `db:"email"`
    Role      string    `db:"role"`
    BranchID  string    `db:"branch_id"`
    CreatedAt time.Time `db:"created_at"`
}