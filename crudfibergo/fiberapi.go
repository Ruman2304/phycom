package main

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
)

var db *sql.DB

type USER struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Lastname string  `json:"lastname"`
	Age      int     `json:"age"`
	Sex      string  `json:"sex"`
	Address  Address `json:"address"`
}

type Address struct {
	Houseno int    `json:"houseno"`
	Street  string `json:"street"`
	City    string `json:"city"`
}

func main() {
	var err error
	db, err = sql.Open("postgres", "user=postgres password=root dbname=student host=localhost port=5433 sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	// Create tables if they don't exist
	createTables()

	defer db.Close()

	app := fiber.New()
	fmt.Println("Connected to the database")

	// Define routes
	app.Post("/api/create", createuser)
	app.Get("/api/get", getUsers)
	app.Get("/api/user/:id", getUserByID)
	app.Delete("/api/user/:id", deleteUserByID)
	app.Put("/api/user/:id", updateUserByID)

	log.Fatal(app.Listen(":8080"))
}

// createTables ensures the required tables exist
func createTables() {
	queryUser := `
		CREATE TABLE IF NOT EXISTS userinfo (
			id SERIAL PRIMARY KEY,
			name VARCHAR(10),
			lastname VARCHAR(10),
			age INT,
			sex VARCHAR(10)
		)`
	queryAddress := `
		CREATE TABLE IF NOT EXISTS addressinfo (
			user_id INT NOT NULL,
			houseno INT,
			street VARCHAR(10),
			city VARCHAR(10),
			CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES userinfo(id)
		)`
	_, err := db.Exec(queryUser)
	if err != nil {
		log.Printf("Cannot create table userinfo: %v", err)
	}

	_, err = db.Exec(queryAddress)
	if err != nil {
		log.Printf("Cannot create table addressinfo: %v", err)
	}
}

// createuser handles the insertion of a user and their address
func createuser(c *fiber.Ctx) error {
	var user USER

	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input data")
	}

	query := "INSERT INTO userinfo(name, lastname, age, sex) VALUES($1, $2, $3, $4) RETURNING id"
	var user_id int
	err := db.QueryRow(query, user.Name, user.Lastname, user.Age, user.Sex).Scan(&user_id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert user")
	}

	addQuery := "INSERT INTO addressinfo(user_id, houseno, street, city) VALUES($1, $2, $3, $4)"
	_, err = db.Exec(addQuery, user_id, user.Address.Houseno, user.Address.Street, user.Address.City)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to insert address")
	}

	return c.SendString("Data inserted successfully")
}

// getUsers fetches all users and their addresses
func getUsers(c *fiber.Ctx) error {
	var users []USER

	// Query to fetch users
	query := `
		SELECT 
			u.id, u.name, u.lastname, u.age, u.sex,
			a.houseno, a.street, a.city
		FROM 
			userinfo u
		LEFT JOIN 
			addressinfo a ON u.id = a.user_id
	`
	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching users")
	}
	defer rows.Close()

	// Populate the users slice
	for rows.Next() {
		var user USER
		err := rows.Scan(&user.ID, &user.Name, &user.Lastname, &user.Age, &user.Sex, &user.Address.Houseno, &user.Address.Street, &user.Address.City)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		users = append(users, user)
	}

	// Return users as JSON
	return c.Status(fiber.StatusOK).JSON(users)
}

// getUserByID fetches a specific user and their address by ID
func getUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	var user USER
	query := `
		SELECT 
			u.id, u.name, u.lastname, u.age, u.sex,
			a.houseno, a.street, a.city
		FROM 
			userinfo u
		LEFT JOIN 
			addressinfo a ON u.id = a.user_id
		WHERE 
			u.id = $1
	`
	err = db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Lastname,
		&user.Age,
		&user.Sex,
		&user.Address.Houseno,
		&user.Address.Street,
		&user.Address.City,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error fetching user")
	}

	return c.JSON(user)
}
func deleteUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	// Single query to delete from both tables
	query := `
        WITH deleted_address AS (
            DELETE FROM addressinfo WHERE user_id = $1
        )
        DELETE FROM userinfo WHERE id = $1
    `

	// Execute the query
	_, err = db.Exec(query, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Error deleting user and address")
	}

	return c.SendString("User and associated address deleted successfully with ID: " + idParam)
}
func updateUserByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	var user USER
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input data")
	}

	query := "UPDATE userinfo SET name=$1, lastname=$2, age=$3, sex=$4 WHERE id=$5"

	_, err = db.Exec(query, user.Name, user.Lastname, user.Age, user.Sex, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update user")
	}

	addQuery := "UPDATE addressinfo SET houseno=$1, street=$2, city=$3 WHERE user_id=$4"
	_, err = db.Exec(addQuery, user.Address.Houseno, user.Address.Street, user.Address.City, id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update address")
	}

	return c.SendString("Data updated successfully")
}
