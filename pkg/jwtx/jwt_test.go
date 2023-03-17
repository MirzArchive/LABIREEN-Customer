package jwtx

import (
	"fmt"
	"labireen-customer/pkg/response"
	"log"
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func TestGenerateToken(t *testing.T) {
	err := godotenv.Load("here.env")
	if err != nil {
		log.Fatalln(".env load failed")
	}

	fmt.Println("HERE IS THE SECRET", response.Highlight(os.Getenv("SECRET")))

	var id uuid.UUID
	temp := "895f87ab-f448-4b28-99ed-bb743aa65de7"
	id, _ = uuid.FromBytes([]byte(temp))
	exp, _ := time.ParseDuration("24h")

	tempToken := jwt.NewWithClaims(jwt.SigningMethodHS256, NewCustomerClaims(id, exp))

	token, err := tempToken.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Println("Failed to Sign token")
	}

	fmt.Println("Here is your token", token)
}
