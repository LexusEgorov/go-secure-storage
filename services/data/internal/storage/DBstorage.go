package storage

import (
	"data/internal/models"
	"database/sql"
	"fmt"
	"strconv"

	"github.com/Masterminds/squirrel"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	d    *sql.DB
	psql squirrel.StatementBuilderType
}

// AddBinary implements data.Storager.
func (d *DB) AddBinary(binary []byte, userId int) (string, error) {
	sql, args, err := d.psql.Insert("binaries").
		Columns("user_id", "binary").
		Values(userId, binary).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return "", err
	}

	var id int

	err = d.d.QueryRow(sql, args...).Scan(&id)

	if err != nil {
		return "", err
	}

	filename := strconv.Itoa(id)
	return filename, nil
}

// AddCard implements data.Storager.
func (d *DB) AddCard(card models.Card, userId int) (string, error) {
	sql, args, err := d.psql.Insert("cards").
		Columns("user_id", "cvv", "date", "number", "holder").
		Values(userId, card.Cvv, card.Exp, card.Number, card.Holder).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return "", err
	}

	var id int

	err = d.d.QueryRow(sql, args...).Scan(&id)

	if err != nil {
		return "", err
	}

	filename := strconv.Itoa(id)
	return filename, nil
}

// AddPassword implements data.Storager.
func (d *DB) AddPassword(password models.Password, userId int) (string, error) {
	sql, args, err := d.psql.Insert("passwords").
		Columns("user_id", "login", "password").
		Values(userId, password.Login, password.Password).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return "", err
	}

	var id int

	err = d.d.QueryRow(sql, args...).Scan(&id)

	if err != nil {
		return "", err
	}

	filename := strconv.Itoa(id)
	return filename, nil
}

// AddText implements data.Storager.
func (d *DB) AddText(text string, userId int) (string, error) {
	sql, args, err := d.psql.Insert("binaries").
		Columns("user_id", "text").
		Values(userId, text).
		Suffix("RETURNING id").
		ToSql()

	if err != nil {
		return "", err
	}

	var id int

	err = d.d.QueryRow(sql, args...).Scan(&id)

	if err != nil {
		return "", err
	}

	filename := strconv.Itoa(id)
	return filename, nil
}

// GetBinary implements data.Storager.
func (d *DB) GetBinary(filename string) ([]byte, error) {
	sql, args, err := d.psql.Select("binary").
		From("binaries").
		Where("id = ?", filename).
		ToSql()

	var result []byte

	if err != nil {
		return result, err
	}

	err = d.d.QueryRow(sql, args...).Scan(&result)

	if err != nil {
		return result, err
	}

	return result, nil
}

// GetCard implements data.Storager.
func (d *DB) GetCard(filename string) (*models.Card, error) {
	sql, args, err := d.psql.Select("cvv", "exp", "number", "holder").
		From("cards").
		Where("id = ?", filename).
		ToSql()

	var card models.Card

	if err != nil {
		return nil, err
	}

	err = d.d.QueryRow(sql, args...).Scan(&card.Cvv, &card.Exp, &card.Number, &card.Holder)

	if err != nil {
		return nil, err
	}

	return &card, nil
}

// GetPassword implements data.Storager.
func (d *DB) GetPassword(filename string) (*models.Password, error) {
	sql, args, err := d.psql.Select("login", "password").
		From("cards").
		Where("id = ?", filename).
		ToSql()

	var password models.Password

	if err != nil {
		return nil, err
	}

	err = d.d.QueryRow(sql, args...).Scan(&password.Login, &password.Password)

	if err != nil {
		return nil, err
	}

	return &password, nil
}

// GetText implements data.Storager.
func (d *DB) GetText(filename string) (string, error) {
	sql, args, err := d.psql.Select("text").
		From("texts").
		Where("id = ?", filename).
		ToSql()

	var result string

	if err != nil {
		return result, err
	}

	err = d.d.QueryRow(sql, args...).Scan(&result)

	if err != nil {
		return result, err
	}

	return result, nil
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
