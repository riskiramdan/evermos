package user

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"time"

	"github.com/riskiramdan/evermos/internal/data"
	"github.com/riskiramdan/evermos/internal/types"
	"golang.org/x/crypto/bcrypt"
)

// Errors
var (
	ErrWrongPassword      = errors.New("wrong password")
	ErrWrongEmail         = errors.New("wrong email")
	ErrEmailAlreadyExists = errors.New("Email Already Exists")
	ErrNotFound           = errors.New("not found")
	ErrNoInput            = errors.New("no input")
	ErrLimitInput         = errors.New("Name should be more than 5 char")
	ErrNameAlreadyExist   = errors.New(("Name Already Exits"))
)

// User user
type User struct {
	ID             int        `json:"id" db:"id"`
	Name           string     `json:"name" db:"name"`
	Email          string     `json:"email" db:"email"`
	Password       string     `json:"password" db:"password"`
	Token          *string    `json:"token" db:"merchantToken"`
	TokenExpiredAt *time.Time `json:"tokenExpiredAt" db:"tokenExpiredAt"`
	CreatedAt      time.Time  `json:"createdAt" db:"createdAt"`
	UpdatedAt      *time.Time `json:"updatedAt" db:"updatedAt"`
}

//FindAllUsersParams params for find all
type FindAllUsersParams struct {
	UserID int    `json:"userId"`
	Page   int    `json:"page"`
	Search string `json:"search"`
	Limit  int    `json:"limit"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	Token  string `json:"token"`
}

// CreateUserParams represent the http request data for create user
// swagger:model
type CreateUserParams struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UpdateUserParams represent the http request data for update user
// swagger:model
type UpdateUserParams struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

// LoginParams represent the http request data for login user
// swagger:model
type LoginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginResponse represents the response of login function
// swagger:model
type LoginResponse struct {
	SessionID string `json:"sessionId"`
	User      *User  `json:"user"`
}

// ChangePasswordParams represent the http request data for change password
// swagger:model
type ChangePasswordParams struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

// Storage represents the user storage interface
type Storage interface {
	FindAll(ctx context.Context, params *FindAllUsersParams) ([]*User, *types.Error)
	FindByID(ctx context.Context, userID int) (*User, *types.Error)
	FindByEmail(ctx context.Context, email string) (*User, *types.Error)
	FindByToken(ctx context.Context, token string) (*User, *types.Error)
	Insert(ctx context.Context, user *User) (*User, *types.Error)
	Update(ctx context.Context, user *User) (*User, *types.Error)
	Delete(ctx context.Context, userID int) *types.Error
}

// ServiceInterface represents the user service interface
type ServiceInterface interface {
	ListUsers(ctx context.Context, params *FindAllUsersParams) ([]*User, int, *types.Error)
	GetUser(ctx context.Context, userID int) (*User, *types.Error)
	CreateUser(ctx context.Context, params *CreateUserParams) (*User, *types.Error)
	UpdateUser(ctx context.Context, userID int, params *UpdateUserParams) (*User, *types.Error)
	DeleteUser(ctx context.Context, userID int) *types.Error
	ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) *types.Error
	Login(ctx context.Context, email string, password string) (*LoginResponse, *types.Error)
	Logout(ctx context.Context, token string) *types.Error
	GetByToken(ctx context.Context, token string) (*User, *types.Error)
}

func generateToken() (string, error) {
	buff := make([]byte, 32)
	_, err := rand.Read(buff)
	if err != nil {
		return "", err
	}
	token := fmt.Sprintf("%x", buff)
	return token, nil
}

// Service is the domain logic implementation of user Service interface
type Service struct {
	userStorage Storage
}

// ListUsers is listing users
func (s *Service) ListUsers(ctx context.Context, params *FindAllUsersParams) ([]*User, int, *types.Error) {
	users, err := s.userStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserService->ListUsers()" + err.Path
		return nil, 0, err
	}
	params.Page = 0
	params.Limit = 0
	allUsers, err := s.userStorage.FindAll(ctx, params)
	if err != nil {
		err.Path = ".UserService->ListUsers()" + err.Path
		return nil, 0, err
	}

	return users, len(allUsers), nil
}

// GetUser is get user
func (s *Service) GetUser(ctx context.Context, userID int) (*User, *types.Error) {
	user, err := s.userStorage.FindByID(ctx, userID)
	if err != nil {
		err.Path = ".UserService->GetUser()" + err.Path
		return nil, err
	}

	return user, nil
}

// CreateUser create user
func (s *Service) CreateUser(ctx context.Context, params *CreateUserParams) (*User, *types.Error) {
	users, _, errType := s.ListUsers(ctx, &FindAllUsersParams{
		Email: params.Email,
	})
	if errType != nil {
		errType.Path = ".UserService->CreateUser()" + errType.Path
		return nil, errType
	}
	if len(users) > 0 {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: ErrEmailAlreadyExists.Error(),
			Error:   ErrEmailAlreadyExists,
			Type:    "validation-error",
		}
	}

	bcryptHash, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: err.Error(),
			Error:   err,
			Type:    "golang-error",
		}
	}

	now := time.Now()

	user := &User{
		Name:           params.Name,
		Email:          params.Email,
		Password:       string(bcryptHash),
		Token:          nil,
		TokenExpiredAt: nil,
		CreatedAt:      now,
		UpdatedAt:      &now,
	}

	user, errType = s.userStorage.Insert(ctx, user)
	if errType != nil {
		errType.Path = ".UserService->CreateUser()" + errType.Path
		return nil, errType
	}

	return user, nil
}

// UpdateUser update a user
func (s *Service) UpdateUser(ctx context.Context, userID int, params *UpdateUserParams) (*User, *types.Error) {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		err.Path = ".UserService->UpdateUser()" + err.Path
		return nil, err
	}

	users, _, err := s.ListUsers(ctx, &FindAllUsersParams{
		Email: params.Email,
	})
	if err != nil {
		err.Path = ".UserService->UpdateUser()" + err.Path
		return nil, err
	}
	if len(users) > 0 {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: data.ErrAlreadyExist.Error(),
			Error:   data.ErrAlreadyExist,
			Type:    "validation-error",
		}
	}

	user.Name = params.Name
	user.Email = params.Email

	user, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->UpdateUser()" + err.Path
		return nil, err
	}

	return user, nil
}

// DeleteUser delete a user
func (s *Service) DeleteUser(ctx context.Context, userID int) *types.Error {
	err := s.userStorage.Delete(ctx, userID)
	if err != nil {
		err.Path = ".UserService->DeleteUser()" + err.Path
		return err
	}

	return nil
}

// ChangePassword change password
func (s *Service) ChangePassword(ctx context.Context, userID int, oldPassword, newPassword string) *types.Error {
	user, err := s.GetUser(ctx, userID)
	if err != nil {
		err.Path = ".UserService->ChangePassword()" + err.Path
		return err
	}

	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword))
	if errBcrypt != nil {
		return &types.Error{
			Path:    ".UserService->ChangePassword()",
			Message: ErrWrongPassword.Error(),
			Error:   ErrWrongPassword,
			Type:    "golang-error",
		}
	}

	bcryptHash, errBcrypt := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if errBcrypt != nil {
		return &types.Error{
			Path:    ".UserService->ChangePassword()",
			Message: errBcrypt.Error(),
			Error:   errBcrypt,
			Type:    "golang-error",
		}
	}

	user.Password = string(bcryptHash)
	_, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->ChangePassword()" + err.Path
		return err
	}

	return nil
}

// Login login
func (s *Service) Login(ctx context.Context, email string, password string) (*LoginResponse, *types.Error) {
	users, err := s.userStorage.FindAll(ctx, &FindAllUsersParams{
		Email: email,
	})
	if err != nil {
		err.Path = ".UserService->Login()" + err.Path
		return nil, err
	}
	if len(users) < 1 {
		return nil, &types.Error{
			Path:    ".UserService->Login()",
			Message: ErrWrongEmail.Error(),
			Error:   ErrWrongEmail,
			Type:    "validation-error",
		}
	}

	user := users[0]
	errBcrypt := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if errBcrypt != nil {
		return nil, &types.Error{
			Path:    ".UserService->ChangePassword()",
			Message: ErrWrongPassword.Error(),
			Error:   ErrWrongPassword,
			Type:    "golang-error",
		}
	}

	token, errToken := generateToken()
	if errToken != nil {
		return nil, &types.Error{
			Path:    ".UserService->CreateUser()",
			Message: errToken.Error(),
			Error:   errToken,
			Type:    "golang-error",
		}
	}

	now := time.Now()
	tokenExpiredAt := now.Add(72 * time.Hour)

	user.Token = &token
	user.TokenExpiredAt = &tokenExpiredAt
	user.UpdatedAt = &now

	user, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->CreateUser()" + err.Path
		return nil, err
	}

	return &LoginResponse{
		SessionID: token,
		User:      user,
	}, nil
}

// Logout logout
func (s *Service) Logout(ctx context.Context, token string) *types.Error {
	user, err := s.userStorage.FindByToken(ctx, token)
	if err != nil {
		err.Path = ".UserService->Logout()" + err.Path
		return err
	}

	user.Token = nil
	user.TokenExpiredAt = nil
	user, err = s.userStorage.Update(ctx, user)
	if err != nil {
		err.Path = ".UserService->Logout()" + err.Path
		return err
	}

	return nil
}

// GetByToken get user by its token
func (s *Service) GetByToken(ctx context.Context, token string) (*User, *types.Error) {
	user, err := s.userStorage.FindByToken(ctx, token)
	if err != nil {
		err.Path = ".UserService->GetByToken()" + err.Path
		return nil, err
	}

	return user, nil
}

// NewService creates a new user AppService
func NewService(
	userStorage Storage,
) *Service {
	return &Service{
		userStorage: userStorage,
	}
}
