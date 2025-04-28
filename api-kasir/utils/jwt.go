package utils

import (
    "time"
    "github.com/dgrijalva/jwt-go"
)

// Secret key untuk signing JWT
var jwtSecret = []byte("582889w67888")

// Struct untuk claims JWT
type Claims struct {
    Email   string `json:"email"`
    Role    string `json:"role"`
    BranchID string `json:"branch_id"`
    jwt.StandardClaims
}

// Generate JWT token
func GenerateJWT(email, role, branchID string) (string, error) {
    claims := Claims{
        Email:   email,
        Role:    role,
        BranchID: branchID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expired in 24 hours
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    signedToken, err := token.SignedString(jwtSecret)
    if err != nil {
        return "", err
    }
    return signedToken, nil
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