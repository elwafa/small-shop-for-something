package entities

import (
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"regexp"
)
import "golang.org/x/crypto/bcrypt"

type User struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,password"`
	Phone    string `json:"phone" validate:"required,saudi_phone"`
	Type     string `json:"type" validate:"required,oneof=admin customer seller"`
	IsActive bool   `json:"is_active"`
}

// Create a new validator instance
var validate *validator.Validate

func init() {
	validate = validator.New()
	validate.RegisterValidation("saudi_phone", validateSaudiPhoneNumber)
	validate.RegisterValidation("password", validatePassword)
}

func NewUser(name, email, password, phone, userType string, isActive bool) (*User, error) {
	// hash the password before saving to the database
	hashedPassword, _ := HashPassword(password)
	user := &User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
		Phone:    phone,
		Type:     userType,
		IsActive: isActive,
	}
	err := validate.Struct(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// HashPassword hashes the password before saving to the database
func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPasswordHash(password, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func validateSaudiPhoneNumber(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	// Regex for Saudi phone number validation (local and international formats)
	re := regexp.MustCompile(`^(?:\+966|0)(5\d{8})$`)

	return re.MatchString(phone)
}

func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()

	// Check length >= 8
	if len(password) < 8 {
		return false
	}

	// Regular expressions to check password complexity
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString
	hasUppercase := regexp.MustCompile(`[A-Z]`).MatchString
	hasSpecialChar := regexp.MustCompile(`[!@#~$%^&*(),.?":{}|<>]`).MatchString

	// Ensure password contains at least one number, one uppercase letter, and one special character
	if !hasNumber(password) || !hasUppercase(password) || !hasSpecialChar(password) {
		return false
	}

	return true
}
