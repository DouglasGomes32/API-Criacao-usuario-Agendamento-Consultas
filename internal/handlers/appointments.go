package handlers

import (
	"database/sql"
	"net/http"
	"saudemais-api/internal/models"
	"saudemais-api/internal/types"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func AgendarConsulta(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		var req types.AppointmentRequest

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Dados inválidos",
				Code:    400,
			})
		}

		if req.Datetime == "" {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Data e hora são obrigatórios",
				Code:    400,
			})
		}

		agendamento, err := time.Parse("2006-01-02T15:04:05", req.Datetime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Formato de data inválido",
				Code:    400,
			})
		}

		if agendamento.Before(time.Now()) {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "Consulta no passado não é permitida",
				Code:    400,
			})
		}

		email, ok := c.Get("userEmail").(string)
		if !ok || email == "" {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Usuário não autenticado",
				Code:    401,
			})
		}

		paciente, err := models.BuscarPacientePorEmail(db, email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Paciente não encontrado",
				Code:    401,
			})
		}

		existe, err := models.ConsultaExistente(db, paciente.ID, agendamento)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, types.APIResponse{
				Message: "Erro ao verificar agendamento",
				Code:    500,
			})
		}

		if existe {
			return c.JSON(http.StatusConflict, types.APIResponse{
				Message: "Consulta já agendada para este horário",
				Code:    409,
			})
		}

		err = models.CriarConsulta(db, paciente.ID, agendamento)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, types.APIResponse{
				Message: "Erro ao agendar consulta",
				Code:    500,
			})
		}

		return c.JSON(http.StatusCreated, types.APIResponse{
			Message: "Consulta agendada com sucesso",
			Code:    201,
		})
	}
}

func ListarConsulta(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email, ok := c.Get("userEmail").(string)
		if !ok || email == "" {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Usuário não autenticado",
				Code:    401,
			})
		}

		paciente, err := models.BuscarPacientePorEmail(db, email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Paciente não encontrado",
				Code:    401,
			})
		}

		consultas, err := models.ListaConsultasPorPaciente(db, paciente.ID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, types.APIResponse{
				Message: "Erro ao listar consultas",
				Code:    500,
			})
		}

		var resposta []types.AppointmentResponse
		for _, consulta := range consultas {
			resposta = append(resposta, types.AppointmentResponse{
				ID:       consulta.ID,
				Datetime: consulta.Datetime.Format("2006-01-02 15:04:05"),
			})
		}

		return c.JSON(http.StatusOK, resposta)
	}
}

func CancelarConsulta(db *sql.DB) echo.HandlerFunc {
	return func(c echo.Context) error {
		email, ok := c.Get("userEmail").(string)
		if !ok || email == "" {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Usuário não autenticado",
				Code:    401,
			})
		}

		paciente, err := models.BuscarPacientePorEmail(db, email)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, types.APIResponse{
				Message: "Paciente não encontrado",
				Code:    401,
			})
		}

		idParam := c.Param("id")
		idConsulta, err := strconv.Atoi(idParam)
		if err != nil {
			return c.JSON(http.StatusBadRequest, types.APIResponse{
				Message: "ID inválido",
				Code:    400,
			})
		}

		err = models.CancelarConsultaDoPaciente(db, paciente.ID, idConsulta)
		if err != nil {
			if err.Error() == "Consulta não encontrada ou acesso não autorizado" {
				return c.JSON(http.StatusForbidden, types.APIResponse{
					Message: err.Error(),
					Code:    403,
				})
			}
			return c.JSON(http.StatusInternalServerError, types.APIResponse{
				Message: "Erro ao cancelar consulta",
				Code:    500,
			})
		}

		return c.JSON(http.StatusOK, types.APIResponse{
			Message: "Consulta cancelada com sucesso",
			Code:    200,
		})
	}
}
