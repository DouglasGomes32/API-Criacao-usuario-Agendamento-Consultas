package types

import "time"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Paciente struct {
	ID           int
	Nome         string
	Email        string
	PasswordHash string
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AppointmentRequest struct {
	Datetime string `json:"datetime"` // "2006-01-02T15:04:05"
}

type AppointmentResponse struct {
	ID       int    `json:"id"`
	Datetime string `json:"datetime"` // "2006-01-02 15:04:05"
}

type APIResponse struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Consulta struct {
	ID        int       `json:"id"`
	PatientID int       `json:"patient_id"`
	Datetime  time.Time `json:"datetime"`
}
