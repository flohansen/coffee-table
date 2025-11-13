package v1

import "errors"

type RegisterRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

func (r *RegisterRequest) Validate() error {
	if r.FirstName == "" {
		return errors.New("firstName must not be empty")
	}

	if r.LastName == "" {
		return errors.New("lastName must not be empty")
	}

	if r.Email == "" {
		return errors.New("email must not be empty")
	}

	return nil
}
