package service

import (
	"LaoQGChat/dto"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthService interface {
	Login(ctx *gin.Context, inDto dto.AuthInDto) *dto.AuthOutDto
}

type authService struct {
	AuthTokenMap map[string]uuid.UUID
}

func NewAuthService() AuthService {
	service := new(authService)
	return service
}

func (service *authService) Login(ctx *gin.Context, inDto dto.AuthInDto) *dto.AuthOutDto {
	outDto := new(dto.AuthOutDto)
	outDto.AuthToken = uuid.New()
	return outDto
}
