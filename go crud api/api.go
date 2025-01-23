package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
	_ "modernc.org/sqlite"
)

// Student represents a student record
type Student struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

var db *sql.DB

func main() {
	fmt.Print("hello")

	var err error
	db, err = sql.Open("sqlite3", "./students.db")

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	// Create table if it doesn't exist
	createTable()

	// Define API routes
	http.HandleFunc("/students", studentsHandler)   // GET and POST
	http.HandleFunc("/students/id", studentHandler) // GET by ID, PUT, DELETE

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}

// Create the "students" table
func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS students1 (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		age INTEGER NOT NULL,
		email TEXT NOT NULL UNIQUE
	);`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating table: %v", err)
	}
}

// Handler for "/students" (GET and POST)
func studentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		// Get all students
		readStudents(w, r)
	case http.MethodPost:
		// Add a new student
		createStudent(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handler for "/students/{id}" (GET by ID, PUT, DELETE)
func studentHandler(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid student ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		// Get a student by ID
		readStudentByID(w, id)
	case http.MethodPut:
		// Update a student by ID
		updateStudent(w, r, id)
	case http.MethodDelete:
		// Delete a student by ID
		deleteStudent(w, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Create a new student
func createStudent(w http.ResponseWriter, r *http.Request) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO students1 (name, age, email) VALUES (?, ?, ?)`
	_, err = db.Exec(query, student.Name, student.Age, student.Email)
	if err != nil {
		http.Error(w, "Failed to create student", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintln(w, "Student created successfully")
}

// Get all students
func readStudents(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, name, age, email FROM students1`
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch students", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []Student
	for rows.Next() {
		var student Student
		err := rows.Scan(&student.ID, &student.Name, &student.Age, &student.Email)
		if err != nil {
			http.Error(w, "Error scanning student", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	json.NewEncoder(w).Encode(students)
}

// Get a student by ID
func readStudentByID(w http.ResponseWriter, id int) {
	query := `SELECT id, name, age, email FROM students1 WHERE id = ?`
	row := db.QueryRow(query, id)

	var student Student
	err := row.Scan(&student.ID, &student.Name, &student.Age, &student.Email)
	if err == sql.ErrNoRows {
		http.Error(w, "Student not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching student", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(student)
}

// Update a student by ID
func updateStudent(w http.ResponseWriter, r *http.Request, id int) {
	var student Student
	err := json.NewDecoder(r.Body).Decode(&student)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := `UPDATE students1 SET name = ?, age = ?, email = ? WHERE id = ?`
	_, err = db.Exec(query, student.Name, student.Age, student.Email, id)
	if err != nil {
		http.Error(w, "Failed to update student", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Student updated successfully")
}

// Delete a student by ID
func deleteStudent(w http.ResponseWriter, id int) {
	query := `DELETE FROM students1 WHERE id = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete student", http.StatusInternalServerError)
		return
	}
	fmt.Fprintln(w, "Student deleted successfully")
}
