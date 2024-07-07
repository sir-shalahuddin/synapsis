package service

import (
	"context"

	"github.com/sir-shalahuddin/synapsis/dto"
	"github.com/sir-shalahuddin/synapsis/model"
	"github.com/sir-shalahuddin/synapsis/pkg/auth"
	"github.com/sir-shalahuddin/synapsis/pkg/hash"
	error_helper "github.com/sir-shalahuddin/synapsis/pkg/helper"
)

type userRepository interface {
	Register(ctx context.Context, email string, password string, username string) (string, error)
	FindByEmail(stx context.Context, email string) (*model.User, error)
}

type userService struct {
	repo  userRepository
	hash  hash.Hashing
	token auth.TokenManager
}

func NewUserService(repo userRepository, token auth.TokenManager) *userService {
	return &userService{repo: repo, token: token}
}

func (svc *userService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	res, err := svc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}

	if res != nil {
		return nil, error_helper.ErrorDuplicateEmail
	}

	passwordHash, err := svc.hash.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := model.User{
		Email:    req.Email,
		Password: passwordHash,
		Username: req.Username,
	}

	_, err = svc.repo.Register(ctx, user.Email, user.Password, user.Username)
	if err != nil {
		return nil, err
	}

	response := dto.RegisterResponse{
		Email:    user.Email,
		Username: user.Username,
	}

	return &response, nil
}

func (svc *userService) Login(ctx context.Context, req *dto.LoginRequest) (string, error) {

	user, err := svc.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", error_helper.ErrorAuthentication
	}

	err = svc.hash.ComparePassword(user.Password, req.Password)
	if err != nil {
		return "", error_helper.ErrorAuthentication
	}
	return svc.createJwt(user.ID)
}

func (svc *userService) createJwt(id string) (string, error) {

	jwt, err := svc.token.NewJWT(id)

	if err != nil {
		return "", err
	}

	return jwt, err
}
