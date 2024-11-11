package models

// User schema
type User struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
	Age      int64  `json:"age"`
	Role     string `json:"role"` // "receptionist" or "doctor"
}

// Patient schema
type Patient struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Age      int64  `json:"age"`
	Location string `json:"location"`
	DoctorID int64  `json:"doctor_id"`   // ID of the doctor assigned to the patient
	Notes    string `json:"notes"`
}
