package main

// @title User Management API
// @version 1.0
// @description This is a sample server for managing users.
// @host localhost:8080
// @BasePath /

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/CRUD-Golang/docs"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// User represents the model for a user
// @Description User model
type User struct {
	ID        uint       `json:"id" gorm:"primaryKey"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty" gorm:"index"`
	Name      string     `json:"name"`
	Email     string     `json:"email"`
}

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code    int    `json:"code" example:"400"`
	Message string `json:"message" example:"status bad request"`
}

// UserCreateRequest represents the request body for creating a user
type UserCreateRequest struct {
	Name  string `json:"name" example:"Tonkhab"`
	Email string `json:"email" example:"Tonkhab@gmail.com"`
}

var db *gorm.DB

func initDB() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// Auto Migration
	err = db.AutoMigrate(&User{})
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}
}

func main() {
	initDB()
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/swagger/*", echoSwagger.EchoWrapHandler())
	e.GET("/users", getUsers)
	e.GET("/user/:id", getUserHandler)
	e.POST("/users", createUser)
	e.PUT("/users/:id", updateUser)
	e.DELETE("/users/:id", deleteUser)
	e.Logger.Fatal(e.Start(":8080"))
}

// @Summary Get all users
// @Description Get all users
// @Tags users
// @Produce json
// @Success 200 {array} User
// @Failure 500 {object} echo.HTTPError
// @Router /users [get]
func getUsers(c echo.Context) error {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, users)
}

// @Summary Get user by ID
// @Description Get user by ID
// @Tags user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /user/{id} [get]
func getUserHandler(c echo.Context) error {
	id := c.Param("id")
	var user User
	if err := db.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// @Summary Create user
// @Description Create a new user
// @Tags users
// @Accept json
// @Produce json
// @Param user body UserCreateRequest true "User data"
// @Success 201 {object} User
// @Failure 400 {object} HTTPError
// @Failure 500 {object} HTTPError
// @Router /users [post]
func createUser(c echo.Context) error {
	req := new(UserCreateRequest)
	if err := c.Bind(req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	user := &User{
		Name:  req.Name,
		Email: req.Email,
	}

	if err := db.Create(user).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, user)
}

// @Summary Update user
// @Description Update user
// @Tags user
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param user body User true "User data"
// @Success 200 {object} User
// @Failure 400 {object} echo.HTTPError
// @Failure 404 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [put]
func updateUser(c echo.Context) error {
	id := c.Param("id")
	user := new(User)
	if err := c.Bind(user); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := db.Model(&User{}).Where("id = ?", id).Updates(user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return echo.NewHTTPError(http.StatusNotFound, "User not found")
		}
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Delete user
// @Tags user
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {string} string "User deleted successfully"
// @Failure 500 {object} echo.HTTPError
// @Router /users/{id} [delete]
func deleteUser(c echo.Context) error {
	id := c.Param("id")
	if err := db.Delete(&User{}, id).Error; err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("User with ID %s deleted", id)})
}
