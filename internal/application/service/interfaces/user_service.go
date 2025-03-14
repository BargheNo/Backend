package service

import userdto "github.com/BargheNo/Backend/internal/application/dto/user"

type UserService interface {
	Register(registerInfo userdto.BasicRegisterRequest)
}
