package middleware

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// response format
type response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

// Patient struct to map with patients table
type Patient struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Health string `json:"health"`
}

// Doctor struct to map with doctors table
type Doctor struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Specialty string `json:"specialty"`
	Experience int  `json:"experience"`
}

// create connection with postgres db
func createConnection() *sql.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	db, err := sql.Open("postgres", os.Getenv("POSTGRES_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}

// ------------------ Patient Handlers ------------------

// CreatePatient godoc
// @Summary      Creates a new patient
// @Description  Adds a new patient record to the database
// @Tags         Patients
// @Accept       json
// @Produce      json
// @Param        patient  body      Patient  true  "Patient data"
// @Success      201      {object}  response
// @Failure      400      {object}  string
// @Router       /api/newpatient [post]
func CreatePatient(w http.ResponseWriter, r *http.Request) {
	var patient Patient
	err := json.NewDecoder(r.Body).Decode(&patient)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertPatient(patient)

	res := response{
		ID:      insertID,
		Message: "Patient created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// GetPatient godoc
// @Summary      Retrieves a single patient by ID
// @Description  Gets a patient's details from the database using their ID
// @Tags         Patients
// @Produce      json
// @Param        id     path      int     true  "Patient ID"
// @Success      200    {object}  Patient
// @Failure      404    {object}  string
// @Router       /api/patient/{id} [get]
func GetPatient(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	patient, err := getPatient(int64(id))
	if err != nil {
		log.Fatalf("Unable to get patient. %v", err)
	}

	json.NewEncoder(w).Encode(patient)
}

// ------------------ Doctor Handlers ------------------


// CreateDoctor godoc
// @Summary      Creates a new doctor
// @Description  Adds a new doctor record to the database
// @Tags         Doctors
// @Accept       json
// @Produce      json
// @Param        doctor  body      Doctor  true  "Doctor data"
// @Success      201     {object}  response
// @Failure      400     {object}  string
// @Router       /api/newdoctor [post]
func CreateDoctor(w http.ResponseWriter, r *http.Request) {
	var doctor Doctor
	err := json.NewDecoder(r.Body).Decode(&doctor)
	if err != nil {
		log.Fatalf("Unable to decode the request body. %v", err)
	}

	insertID := insertDoctor(doctor)

	res := response{
		ID:      insertID,
		Message: "Doctor created successfully",
	}

	json.NewEncoder(w).Encode(res)
}

// GetDoctor godoc
// @Summary      Retrieves a single doctor by ID
// @Description  Gets a doctor's details from the database using their ID
// @Tags         Doctors
// @Produce      json
// @Param        id      path      int     true  "Doctor ID"
// @Success      200     {object}  Doctor
// @Failure      404     {object}  string
// @Router       /api/doctor/{id} [get]
func GetDoctor(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["id"])
	if err != nil {
		log.Fatalf("Unable to convert the string into int. %v", err)
	}

	doctor, err := getDoctor(int64(id))
	if err != nil {
		log.Fatalf("Unable to get doctor. %v", err)
	}

	json.NewEncoder(w).Encode(doctor)
}

// ------------------ Database Helper Functions ------------------

// insertPatient inserts a new patient into the database
func insertPatient(patient Patient) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO patients (name, age, health) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := db.QueryRow(sqlStatement, patient.Name, patient.Age, patient.Health).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single patient %v", id)
	return id
}

// getPatient retrieves a single patient from the database by ID
func getPatient(id int64) (Patient, error) {
	db := createConnection()
	defer db.Close()

	var patient Patient
	sqlStatement := `SELECT * FROM patients WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&patient.ID, &patient.Name, &patient.Age, &patient.Health)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return patient, nil
	case nil:
		return patient, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return patient, err
}

// insertDoctor inserts a new doctor into the database
func insertDoctor(doctor Doctor) int64 {
	db := createConnection()
	defer db.Close()

	sqlStatement := `INSERT INTO doctors (name, specialty, experience) VALUES ($1, $2, $3) RETURNING id`
	var id int64
	err := db.QueryRow(sqlStatement, doctor.Name, doctor.Specialty, doctor.Experience).Scan(&id)
	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	fmt.Printf("Inserted a single doctor %v", id)
	return id
}

// getDoctor retrieves a single doctor from the database by ID
func getDoctor(id int64) (Doctor, error) {
	db := createConnection()
	defer db.Close()

	var doctor Doctor
	sqlStatement := `SELECT * FROM doctors WHERE id=$1`
	row := db.QueryRow(sqlStatement, id)

	err := row.Scan(&doctor.ID, &doctor.Name, &doctor.Specialty, &doctor.Experience)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return doctor, nil
	case nil:
		return doctor, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	return doctor, err
}
