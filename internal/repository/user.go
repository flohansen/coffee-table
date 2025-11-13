package repository

import (
	"context"

	"github.com/flohansen/coffee-table/internal/domain"
	"github.com/flohansen/coffee-table/sql/generated/database"
)

type UserPostgres struct {
	q *database.Queries
}

func NewUserPostgres(db database.DBTX) *UserPostgres {
	return &UserPostgres{
		q: database.New(db),
	}
}

func (r *UserPostgres) GetAll(ctx context.Context) ([]domain.User, error) {
	dbUsers, err := r.q.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	var users []domain.User
	for _, dbUser := range dbUsers {
		users = append(users, domain.User{
			ID:        dbUser.ID,
			FirstName: dbUser.FirstName,
			LastName:  dbUser.FirstName,
			Email:     dbUser.Email,
		})
	}

	return users, nil
}

func (r *UserPostgres) Create(ctx context.Context, user domain.User) error {
	return nil
}
