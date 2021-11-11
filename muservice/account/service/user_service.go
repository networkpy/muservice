package service

import (
	"context"
	"log"

	"github.com/NetworkPy/muserv/muservice/account/models"
	"github.com/NetworkPy/muserv/muservice/account/models/apperrors"
	"github.com/NetworkPy/muserv/muservice/account/security"

	"github.com/google/uuid"
)

// UserService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
// userService acts as a struct for injecting an implementation of UserRepository
// for use in service methods
type userService struct {
	UserRepository models.UserRepository
}

// USConfig will hold repositories that will eventually be injected into this
// this service layer
type USConfig struct {
	UserRepository models.UserRepository
}

// NewUserService is a factory function for
// initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) models.UserService {
	return &userService{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *userService) Get(ctx context.Context, uid uuid.UUID) (*models.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)
	return u, err
}

// SignUp reaches our to a UserRepository to verify the
// email address is available and signs up the user if this is the case
func (s *userService) Signup(ctx context.Context, u *models.User) error {
	pw, err := security.HashPassword(u.Passowrd)

	if err != nil {
		log.Printf("Unable to signup user for email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	// now I realize why I originally used Signup(ctx, email, password)
	// then created a user. It's somewhat un-natural to mutate the user here
	u.Passowrd = pw
	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	// ...

	return nil
}

// Signin reaches our to a UserRepository check if the user exists
// and then compares the supplied password with the provided password.
// If a valid email/password combo is provided, u will hold all
// available user fields
func (s *userService) Signin(ctx context.Context, u *models.User) error {
	uFetched, err := s.UserRepository.FindByEmail(ctx, u.Email)

	// Will return NotAuthorized to client to omit details of why
	if err != nil {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	// verify password - we previously created this method
	match, err := security.ComparePasswords(uFetched.Passowrd, u.Passowrd)

	if err != nil {
		return apperrors.NewInternal()
	}

	if !match {
		return apperrors.NewAuthorization("Invalid email and password combination")
	}

	*u = *uFetched

	return nil
}
