package storage

import (
	"auth/internal/models"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	d    *sql.DB
	psql squirrel.StatementBuilderType
}

// AddToken implements auth.Storager.
func (d *DB) AddToken(userID int, token string) error {
	sql, args, err := d.psql.Insert("tokens").
		Columns("user_id", "token").
		Values(userID, token).
		ToSql()

	if err != nil {
		return err
	}

	_, err = d.d.Exec(sql, args...)

	if err != nil {
		return err
	}

	return nil
}

// AddUser implements auth.Storager.
func (d *DB) AddUser(email, password string) (int, error) {
	sql, args, err := d.psql.Insert("users").
		Columns("login", "password").
		Values(email, password).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return 0, err
	}

	var uID int
	err = d.d.QueryRow(sql, args...).Scan(&uID)

	if err != nil {
		return 0, err
	}

	return uID, nil
}

// GetToken implements auth.Storager.
func (d *DB) CheckRefresh(token string) (bool, error) {
	sql, args, err := d.psql.Select("id", "closed").
		From("tokens").
		Where("token = ?", token).
		ToSql()

	if err != nil {
		return false, err
	}

	tId := -1
	isUseful := false

	d.d.QueryRow(sql, args...).Scan(&tId, &isUseful)

	if tId == -1 {
		return false, models.ErrNotFound
	}

	if isUseful {
		sql, args, err = d.psql.Update("tokens").
			Set("closed", true).
			Where("id = ?", tId).
			ToSql()

		if err != nil {
			return false, err
		}

		_, err = d.d.Exec(sql, args...)

		if err != nil {
			return false, err
		}
	}

	return isUseful, nil
}

// GetUser implements auth.Storager.
func (d *DB) GetUser(email string) (*models.User, error) {
	sql, args, err := d.psql.Select("id", "email", "password").
		From("users").
		Where("email = ?", email).
		ToSql()

	if err != nil {
		return nil, err
	}

	user := models.User{}

	err = d.d.QueryRow(sql, args...).Scan(&user)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

var db *DB = nil

func NewDB(connection string) *DB {
	if db != nil {
		return db
	}

	db = &DB{
		d:    nil,
		psql: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
	db.connect(connection)

	return db
}

func (d *DB) connect(connection string) error {
	var err error
	d.d, err = sql.Open("pgx", connection)

	if err != nil {
		return err
	}

	err = d.d.Ping()

	if err != nil {
		fmt.Printf("connection error!\n %v\n", err)
		return err
	}

	fmt.Printf("connected to DB with %s successfully\n", connection)

	return nil
}
