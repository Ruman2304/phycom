package main

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strings"

	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
)

var db *sql.DB

type USER struct {
	ID                     int       `json:"id"`
	Age                    int       `json:"age"`
	Gender                 string    `json:"gender"`
	Location               string    `json:"location"`
	Education              string    `json:"education"`
	Occupation             string    `json:"occupation"`
	PrimaryApp             string    `json:"primary_app"`
	UsageFrequency         string    `json:"usage_frequency"`
	DailyUsageTime         string    `json:"daily_usage_time"`
	ReasonForUsing         string    `json:"reason_for_using"`
	Satisfaction           int       `json:"satisfaction"`
	DesiredFeatures        string    `json:"desired_features"`
	PreferredCommunication string    `json:"preferred_communication"`
	PartnerPriorities      string    `json:"partner_priorities"`
	TimeUpdated            time.Time `json:"timeupdated"`
	TimeCreated            time.Time `json:"timecreated"`
}

func main() {

	var err error
	db, err = sql.Open("postgres", "user=postgres password=root dbname=tinder host=localhost port=5433 sslmode=disable")
	if err != nil {
		log.Fatalf("Cannot connect to the database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	createTables()
	datainsert()

	currentTime := time.Now()
	fmt.Print("currenttime", currentTime)

	app := fiber.New()
	app.Put("/api/update/:id", updateUserByID)
	app.Get("api/get", getuser(logrus.New()))
	app.Get("api/get/:id", getuserID(logrus.New()))
	app.Delete("api/delete/:id", deleteUserByID(logrus.New()))

	log.Fatal(app.Listen(":8080"))

	defer db.Close()
}

func createTables() {
	query := `
	CREATE TABLE IF NOT EXISTS userdata (
		User_id SERIAL PRIMARY KEY,
		Age INT,
		Gender VARCHAR(50),
		Location VARCHAR(100),
		Education VARCHAR(100),
		Occupation VARCHAR(100),
		Primary_App VARCHAR(100),
		Usage_Frequency VARCHAR(50),
		Daily_Usage_Time VARCHAR(50),
		Reason_for_Using VARCHAR(255),
		Satisfaction INT,
		Desired_Features VARCHAR(255),
		Preferred_Communication VARCHAR(100),
		Partner_Priorities VARCHAR(255),
		TimeCreated TIMESTAMP,
		TimeUpdated TIMESTAMP,
		UNIQUE (
			Age, Gender, Location, Occupation, Primary_App, Usage_Frequency,
			Daily_Usage_Time, Reason_for_Using, Satisfaction, Desired_Features,
			Preferred_Communication, Partner_Priorities
		)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	}
	fmt.Println("Table created successfully")
}

func datainsert() {
	// Open the CSV file
	file, err := os.Open("userdatatinder.csv")
	if err != nil {
		log.Fatalf("Error while opening the file: %v", err)
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the file
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error while reading records from the file: %v", err)
	}

	// Iterate through each record and process it
	for i, eachRecord := range records {
		// Skip the header row (if your CSV has a header)
		if i == 0 {
			continue
		}

		// Define the query to check for the record
		checkQuery := `
		SELECT age, gender, location, education, occupation, primary_app, usage_frequency, 
		       daily_usage_time, reason_for_using, satisfaction, desired_features, 
		       preferred_communication, partner_priorities
		FROM userdata
		WHERE age=$1 AND gender=$2 AND location=$3 AND education=$4 AND occupation=$5
		  AND primary_app=$6 AND usage_frequency=$7 AND daily_usage_time=$8
		  AND reason_for_using=$9 AND satisfaction=$10 AND desired_features=$11
		  AND preferred_communication=$12 AND partner_priorities=$13`

		var existingRecord USER

		// Execute the query
		err = db.QueryRow(checkQuery,
			eachRecord[1], eachRecord[2], eachRecord[3], eachRecord[4], eachRecord[5],
			eachRecord[6], eachRecord[7], eachRecord[8], eachRecord[9], eachRecord[10],
			eachRecord[11], eachRecord[12], eachRecord[13],
		).Scan(
			&existingRecord.Age, &existingRecord.Gender, &existingRecord.Location,
			&existingRecord.Education, &existingRecord.Occupation, &existingRecord.PrimaryApp,
			&existingRecord.UsageFrequency, &existingRecord.DailyUsageTime, &existingRecord.ReasonForUsing,
			&existingRecord.Satisfaction, &existingRecord.DesiredFeatures,
			&existingRecord.PreferredCommunication, &existingRecord.PartnerPriorities,
		)

		if err != nil && err != sql.ErrNoRows {
			log.Printf("Error checking for existing record at row %d: %v", i+1, err)
			continue
		}

		if err == sql.ErrNoRows {
			// No matching record found, so insert a new one
			insertQuery := `
			INSERT INTO userdata (
				age, gender, location, education, occupation, primary_app,
				usage_frequency, daily_usage_time, reason_for_using, satisfaction,
				desired_features, preferred_communication, partner_priorities, timecreated, timeupdated
			) VALUES (
				$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
			)`
			_, err = db.Exec(
				insertQuery,
				eachRecord[1], eachRecord[2], eachRecord[3], eachRecord[4], eachRecord[5],
				eachRecord[6], eachRecord[7], eachRecord[8], eachRecord[9], eachRecord[10],
				eachRecord[11], eachRecord[12], eachRecord[13], time.Now(), time.Now(),
			)
			if err != nil {
				log.Printf("Error inserting new record at row %d: %v", i+1, err)
				continue
			}
			log.Printf("New record inserted successfully at row %d", i+1)
		} else {
			// Record exists, check for changes
			if existingRecord.Age != atoi(eachRecord[1]) ||
				existingRecord.Gender != eachRecord[2] ||
				existingRecord.Location != eachRecord[3] ||
				existingRecord.Education != eachRecord[4] ||
				existingRecord.Occupation != eachRecord[5] ||
				existingRecord.PrimaryApp != eachRecord[6] ||
				existingRecord.UsageFrequency != eachRecord[7] ||
				existingRecord.DailyUsageTime != eachRecord[8] ||
				existingRecord.ReasonForUsing != eachRecord[9] ||
				existingRecord.Satisfaction != atoi(eachRecord[10]) ||
				existingRecord.DesiredFeatures != eachRecord[11] ||
				existingRecord.PreferredCommunication != eachRecord[12] ||
				existingRecord.PartnerPriorities != eachRecord[13] {
				// Update the record
				updateQuery := `
				UPDATE userdata SET
					age=$1, gender=$2, location=$3, education=$4, occupation=$5,
					primary_app=$6, usage_frequency=$7, daily_usage_time=$8,
					reason_for_using=$9, satisfaction=$10, desired_features=$11,
					preferred_communication=$12, partner_priorities=$13, timeupdated=$14
				WHERE age=$1 AND gender=$2 AND location=$3 AND education=$4 AND occupation=$5
				  AND primary_app=$6 AND usage_frequency=$7 AND daily_usage_time=$8
				  AND reason_for_using=$9 AND satisfaction=$10 AND desired_features=$11
				  AND preferred_communication=$12 AND partner_priorities=$13`
				_, err = db.Exec(
					updateQuery,
					eachRecord[1], eachRecord[2], eachRecord[3], eachRecord[4], eachRecord[5],
					eachRecord[6], eachRecord[7], eachRecord[8], eachRecord[9], eachRecord[10],
					eachRecord[11], eachRecord[12], eachRecord[13], time.Now(),
				)
				if err != nil {
					log.Printf("Error updating record at row %d: %v", i+1, err)
					continue
				}
				log.Printf("Record updated successfully at row %d", i+1)
			}
		}
	}
}

func atoi(str string) int {
	value, _ := strconv.Atoi(str)
	return value
}
func updateUserByID(c *fiber.Ctx) error {

	var userData map[string]interface{}

	if err := c.BodyParser(&userData); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid input data")
	}

	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
	}

	var existingUser USER

	fmt.Println("Fetching user with ID:", id)

	err = db.QueryRow("SELECT age, gender, location, education, occupation FROM userdata WHERE user_id = $1", id).Scan(
		&existingUser.Age, &existingUser.Gender, &existingUser.Location, &existingUser.Education, &existingUser.Occupation,
	)

	if err != nil {
		fmt.Println("Error retrieving user:", err)
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).SendString("User not found")
		}
		return c.Status(fiber.StatusInternalServerError).SendString("Error retrieving user: " + err.Error())
	}

	updatedFields := []string{}
	args := []interface{}{}
	i := 1
	if age, ok := userData["age"]; ok && age != existingUser.Age {
		updatedFields = append(updatedFields, fmt.Sprintf("age = $%d", i))
		args = append(args, age)
		i++
	}

	if gender, ok := userData["gender"]; ok && gender != existingUser.Gender {
		updatedFields = append(updatedFields, fmt.Sprintf("gender = $%d", i))
		args = append(args, gender)
		i++
	}

	if location, ok := userData["location"]; ok && location != existingUser.Location {
		updatedFields = append(updatedFields, fmt.Sprintf("location = $%d", i))
		args = append(args, location)
		i++
	}

	if education, ok := userData["education"]; ok && education != existingUser.Education {
		updatedFields = append(updatedFields, fmt.Sprintf("education = $%d", i))
		args = append(args, education)
		i++
	}

	if occupation, ok := userData["occupation"]; ok && occupation != existingUser.Occupation {
		updatedFields = append(updatedFields, fmt.Sprintf("occupation = $%d", i))
		args = append(args, occupation)
		i++
	}

	if len(updatedFields) == 0 {
		return c.Status(fiber.StatusBadRequest).SendString("No fields to update")
	}

	args = append(args, id)
	fmt.Print(len(args))

	fmt.Println("Executing query:", fmt.Sprintf("UPDATE userdata SET %s WHERE user_id = $%d", strings.Join(updatedFields, ", "), len(args)))

	query := fmt.Sprintf("UPDATE userdata SET %s WHERE user_id = $%d", strings.Join(updatedFields, ", "), len(args))
	_, err = db.Exec(query, args...)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to update user: " + err.Error())
	}

	return c.SendString("User updated successfully")
}

func getuser(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		var users []USER

		// Query to fetch all user data
		rows, err := db.Query("SELECT * FROM userdata")
		if err != nil {
			logger.Println("Error fetching user data:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user data")
		}
		defer rows.Close()

		// Iterate through each row and populate the users slice
		for rows.Next() {
			var user USER
			err := rows.Scan(&user.ID, &user.Age, &user.Gender, &user.Location, &user.Education, &user.Occupation, &user.PrimaryApp, &user.UsageFrequency, &user.DailyUsageTime, &user.ReasonForUsing, &user.Satisfaction, &user.DesiredFeatures, &user.PreferredCommunication, &user.PartnerPriorities, &user.TimeUpdated, &user.TimeCreated)
			if err != nil {
				logger.Println("Error scanning user data:", err)
				continue
			}
			users = append(users, user)
		}

		// Check for errors after iterating through rows
		if err = rows.Err(); err != nil {
			logger.Println("Error iterating rows:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to process user data")
		}

		logger.Println("Users retrieved successfully")
		return c.Status(fiber.StatusOK).JSON(users)
	}
}
func getuserID(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {

		// Get the user ID from URL parameters
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {
			// Invalid user ID
			return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
		}

		var user USER

		// Query to fetch user data by ID
		query := "SELECT User_id, Age, Gender, Location, Education, Occupation, Primary_App, Usage_Frequency, Daily_Usage_Time, Reason_for_Using, Satisfaction, Desired_Features, Preferred_Communication, Partner_Priorities, TimeCreated, TimeUpdated FROM userdata WHERE User_id=$1"
		err = db.QueryRow(query, id).Scan(
			&user.ID,
			&user.Age,
			&user.Gender,
			&user.Location,
			&user.Education,
			&user.Occupation,
			&user.PrimaryApp,
			&user.UsageFrequency,
			&user.DailyUsageTime,
			&user.ReasonForUsing,
			&user.Satisfaction,
			&user.DesiredFeatures,
			&user.PreferredCommunication,
			&user.PartnerPriorities,
			&user.TimeCreated,
			&user.TimeUpdated,
		)

		// Check for errors in the query execution
		if err != nil {
			if err == sql.ErrNoRows {
				// User not found
				logger.Println("User not found for ID:", id)
				return c.Status(fiber.StatusNotFound).SendString("User not found")
			}
			// General query error
			logger.Println("Error fetching user data:", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Failed to fetch user data")
		}

		// User retrieved successfully
		logger.Println("User retrieved successfully for ID:", id)
		return c.Status(fiber.StatusOK).JSON(user)
	}
}

func deleteUserByID(logger *logrus.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")
		id, err := strconv.Atoi(idParam)
		if err != nil {

			logger.Error("Invalid user ID")
			return c.Status(fiber.StatusBadRequest).SendString("Invalid user ID")
		}

		// Single query to delete from both tables
		query := `
            DELETE FROM userdata WHERE user_id = $1
    `

		// Execute the query
		_, err = db.Exec(query, id)
		if err != nil {
			logger.Error("Error deleting user and address")
			return c.Status(fiber.StatusInternalServerError).SendString("Error deleting user and address")

		}
		logger.Info("User and associated address deleted successfully with ID: ", id)

		return c.SendString("User and associated address deleted successfully with ID: " + idParam)
	}
}
