# API-Criacao-usuario-Agendamento-Consultas

Uma API RESTful desenvolvida em Golang para gerenciamento de pacientes e agendamentos de consultas em uma clínica. Utiliza autenticação JWT, arquitetura limpa com Echo, e banco de dados PostgreSQL com suporte ao Docker.

---

## 🧠 Tecnologias Utilizadas

- **Golang** (com Echo Framework)
- **PostgreSQL**
- **Docker & Docker Compose**
- **JWT (Autenticação segura)**
- **Bcrypt** (Hash de senhas)
- **Adminer** (Interface web para gerenciar o banco)

---

## ⚙️ Como Rodar o Projeto

### Pré-requisitos

- [Docker](https://www.docker.com/)
- [Go](https://go.dev/) instalado (versão 1.20+)

### 1. Subir Banco de Dados + Adminer

```bash
docker-compose up -d

PostgreSQL será iniciado em localhost:5432
Adminer disponível em http://localhost:8080

O script init.sql será executado automaticamente na primeira vez, criando as tabelas necessárias.

go run cmd/main.go


🔐 Autenticação JWT
A API utiliza autenticação baseada em token JWT.

Faça o login via POST /auth/login
Copie o token da resposta
Nas rotas protegidas, envie no header:

Authorization: Bearer SEU_TOKEN_AQUI

🔓 Endpoints Públicos
POST	/auth/register -	Cadastro de novo paciente
POST	/auth/login -	Login e geração de token

🔒 Endpoints Privados (requer token JWT)
POST	/consultas	- Agendar uma nova consulta
GET     /consultas	- Listar consultas do paciente
DELETE	/consultas/:id	- Cancelar consulta por ID