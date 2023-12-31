package postgres

import (
	"context"
	"fmt"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	twitter "github.com/mobamoh/twitter-go-graphql"
)

type UserRepo struct {
	db *DB
}

func NewUserRepo(db *DB) *UserRepo {
	return &UserRepo{db: db}
}
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE username= $1 LIMIT 1;`
	user := twitter.User{}
	if err := pgxscan.Get(ctx, r.db.Pool, &user, query, username); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error select: %v", err)
	}
	return user, nil
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE email= $1 LIMIT 1;`
	user := twitter.User{}
	if err := pgxscan.Get(ctx, r.db.Pool, &user, query, email); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error select: %v", err)
	}
	return user, nil
}

func (r *UserRepo) Create(ctx context.Context, user twitter.User) (twitter.User, error) {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return twitter.User{}, fmt.Errorf("error starting transaction: %v", err)
	}
	defer tx.Rollback(ctx)
	user, err = createUser(ctx, tx, user)
	if err != nil {
		return twitter.User{}, fmt.Errorf("error commiting: %v", err)
	}
	return user, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id string) (twitter.User, error) {
	query := `SELECT * FROM users WHERE id= $1 LIMIT 1;`
	user := twitter.User{}
	if err := pgxscan.Get(ctx, r.db.Pool, &user, query, id); err != nil {
		if pgxscan.NotFound(err) {
			return twitter.User{}, twitter.ErrNotFound
		}
		return twitter.User{}, fmt.Errorf("error select: %v", err)
	}
	return user, nil
}

func (r *UserRepo) GetByIDs(ctx context.Context, id []string) ([]twitter.User, error) {
	query := `SELECT * FROM users WHERE id= ANY($1) LIMIT 1;`
	var users []twitter.User
	if err := pgxscan.Select(ctx, r.db.Pool, &users, query, id); err != nil {
		return []twitter.User{}, fmt.Errorf("error select by ids: %v", err)
	}
	return users, nil
}

func createUser(ctx context.Context, tx pgx.Tx, user twitter.User) (twitter.User, error) {
	query := `INSERT INTO users(username, email, password) VALUES ($1,$2,$3) RETURNING *;`
	newUser := twitter.User{}
	if err := pgxscan.Get(ctx, tx, &newUser, query, user.Username, user.Email, user.Password); err != nil {
		return twitter.User{}, fmt.Errorf("error insert: %v", err)
	}
	tx.Commit(ctx)
	return newUser, nil
}
