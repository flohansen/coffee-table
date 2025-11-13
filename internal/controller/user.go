package controller

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/flohansen/coffee-table/internal/api"
	v1 "github.com/flohansen/coffee-table/internal/api/v1"
	"github.com/flohansen/coffee-table/internal/domain"
	"github.com/flohansen/coffee-table/pkg/logging"
)

type UserRepo interface {
	Create(ctx context.Context, user domain.User) error
}

type UserController struct {
	repo UserRepo
}

func NewUserController(repo UserRepo) *UserController {
	return &UserController{
		repo: repo,
	}
}

func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	var req v1.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Error("could not decode request", "error", err)
		api.ErrorBadRequest(w, "request body has invalid json format")
		return
	}

	if err := req.Validate(); err != nil {
		log.Error("invalid request", "error", err)
		api.ErrorBadRequest(w, "request body is invalid", "error", err)
		return
	}

	if err := c.repo.Create(r.Context(), domain.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
	}); err != nil {
		log.Error("could not create user", "error", err)
		api.ErrorInternal(w, "internal error")
		return
	}

	api.OK(w)
}
