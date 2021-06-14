package authenticator

import (
	"lalela-backend/internal/pkg/security/claims"
)

type Authenticator interface {
	Login(LoginRequest) (*LoginResponse, error)
	AuthenticateService(request AuthenticateServiceRequest) (*AuthenticateServiceResponse, error)
}

const ServiceProvider = "Authenticator"
const LoginService = ServiceProvider + ".Login"
const AuthenticateServiceService = ServiceProvider + ".AuthenticateService"

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	JWT string
}

type AuthenticateServiceRequest struct {
	Claims  claims.Claims
	Service string
	OrganizationId string
}

type AuthenticateServiceResponse struct {
}
