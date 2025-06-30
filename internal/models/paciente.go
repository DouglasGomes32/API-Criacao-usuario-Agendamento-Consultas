package models

import (
	"database/sql"
	"errors"
	"saudemais-api/internal/types"

	"golang.org/x/crypto/bcrypt"
)

func CriarPaciente(db *sql.DB, nome, email, senha string) error {
	var id int
	err := db.QueryRow("SELECT id FROM pacientes WHERE email = $1", email).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if id != 0 {
		return errors.New("Paciente já cadastrado")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(senha), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = db.Exec(
		"INSERT INTO pacientes (name, email, password_hash) VALUES ($1, $2, $3)",
		nome, email, string(hash),
	)
	return err
}

func BuscarPacientePorEmail(db *sql.DB, email string) (*types.Paciente, error) {
	var p types.Paciente

	err := db.QueryRow("SELECT id, name, email, password_hash FROM pacientes WHERE email = $1", email).
		Scan(&p.ID, &p.Nome, &p.Email, &p.PasswordHash)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("Paciente não encontrado")
		}
		return nil, err
	}

	return &p, nil
}
