package utils

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

// Secret key untuk signing JWT
var jwtSecret = []byte("582889w67888")

// Struct untuk claims JWT
type Claims struct {
    UserID   string `json:"user_id"`
    Email    string `json:"email"`
    Role     string `json:"role"`
    BranchID string `json:"branch_id"`
    jwt.StandardClaims
}

func GenerateJWT(userID, email, role, branchID string) (string, error) {
    claims := Claims{
        UserID:   userID,
        Email:    email,
        Role:     role,
        BranchID: branchID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(),
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}


// Validate JWT token
func ValidateJWT(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil || !token.Valid {
        return nil, err
    }

    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, err
    }
    return claims, nil
}
