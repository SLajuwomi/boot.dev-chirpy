package auth

import (
	"fmt"
	"testing"
)

func TestHashPassword(t *testing.T) {
	password := "MarkHanck7868"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Hashing password failed: %v\n", err)
	}
	validHash, err := CheckPasswordHash(password, hashedPassword)
	if err != nil {
		t.Errorf("Hash is not valid: %v\n", err)
	}
	fmt.Printf("The validated hash: %v\nand it's password: %v\nValid hash: %v\n", hashedPassword, password, validHash)
}

func TestInvalidHash(t *testing.T) {
	password := "MarkHanck7868"
	password2 := "JoeTheDuck"
	hashedPassword, err := HashPassword(password)
	if err != nil {
		t.Errorf("Hashing password failed: %v\n", err)
	}
	hashedPassword2, err := HashPassword(password2)
	if err != nil {
		t.Errorf("Hashing password 2 failed: %v\n", err)
	}
	validHash, err := CheckPasswordHash(password, hashedPassword)
	if err != nil || !validHash {
		t.Errorf("Hash is not valid: %v\n", err)
	}
	validHash2, err := CheckPasswordHash(password2, hashedPassword2)
	if err != nil || !validHash2 {
		t.Errorf("Hash 2 is not valid: %v\n", err)
	}
	invalidHash, err := CheckPasswordHash(password, hashedPassword2)
	if err != nil || invalidHash {
		t.Errorf("Hash is valid (not expected): %v\n", err)
	}
	invalidHash2, err := CheckPasswordHash(password2, hashedPassword)
	if err != nil || invalidHash2 {
		t.Errorf("Hash 2 is valid (not expected): %v\n", err)
	}
	fmt.Printf("The validated hash: %v\nand it's password: %v\nValid hash: %v\n", hashedPassword, password, validHash)
	fmt.Printf("The 2nd validated hash: %v\nand it's password: %v\nValid hash: %v\n", hashedPassword2, password2, validHash2)
	fmt.Printf("The invalidated hash: %v\nand it's password: %v\nValid hash: %v\n", hashedPassword, password, invalidHash)
	fmt.Printf("The 2nd invalidated hash: %v\nand it's password: %v\nValid hash: %v\n", hashedPassword2, password2, invalidHash2)
}
