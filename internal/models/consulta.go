package models

import (
	"database/sql"
	"errors"
	"saudemais-api/internal/types"
	"time"
)

func ListaConsultasPorPaciente(db *sql.DB, patientID int) ([]types.Consulta, error) {
	query := `SELECT id, patient_id, datetime 
	          FROM consultas 
	          WHERE patient_id = $1 
	          ORDER BY datetime`

	rows, err := db.Query(query, patientID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var consultas []types.Consulta
	for rows.Next() {
		var c types.Consulta
		if err := rows.Scan(&c.ID, &c.PatientID, &c.Datetime); err != nil {
			return nil, err
		}
		consultas = append(consultas, c)
	}

	return consultas, nil
}

func ConsultaExistente(db *sql.DB, patientID int, datetime time.Time) (bool, error) {
	var id int
	err := db.QueryRow(
		`SELECT id FROM consultas 
		 WHERE patient_id = $1 AND datetime = $2`,
		patientID, datetime,
	).Scan(&id)

	if err == sql.ErrNoRows {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

func CriarConsulta(db *sql.DB, patientID int, datetime time.Time) error {
	_, err := db.Exec(`
		INSERT INTO consultas (patient_id, datetime)
		VALUES ($1, $2)
	`, patientID, datetime)
	return err
}

func CancelarConsultaDoPaciente(db *sql.DB, pacienteID int, consultaID int) error {
	result, err := db.Exec(`
		DELETE FROM consultas 
		WHERE id = $1 AND patient_id = $2
	`, consultaID, pacienteID)

	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("Consulta não encontrada ou acesso não autorizado")
	}

	return nil
}
