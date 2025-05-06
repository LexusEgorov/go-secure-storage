package data

import "data/internal/models"

type Storager interface {
	AddBinary(binary []byte, userId int) (string, error)
	AddCard(card models.Card, userId int) (string, error)
	AddPassword(password models.Password, userId int) (string, error)
	AddText(text string, userId int) (string, error)

	GetBinary(filename string) ([]byte, error)
	GetCard(filename string) (*models.Card, error)
	GetPassword(filename string) (*models.Password, error)
	GetText(filename string) (string, error)
}

type DataProvider struct{}

// AddBinary implements grpcserv.DataProvider.
func (d *DataProvider) AddBinary(binary []byte, userId int) (string, error) {
	return d.AddBinary(binary, userId)
}

// AddCard implements grpcserv.DataProvider.
func (d *DataProvider) AddCard(card models.Card, userId int) (string, error) {
	return d.AddCard(card, userId)
}

// AddPassword implements grpcserv.DataProvider.
func (d *DataProvider) AddPassword(password models.Password, userId int) (string, error) {
	return d.AddPassword(password, userId)
}

// AddText implements grpcserv.DataProvider.
func (d *DataProvider) AddText(text string, userId int) (string, error) {
	return d.AddText(text, userId)
}

// GetBinary implements grpcserv.DataProvider.
func (d *DataProvider) GetBinary(filename string) ([]byte, error) {
	return d.GetBinary(filename)
}

// GetCard implements grpcserv.DataProvider.
func (d *DataProvider) GetCard(filename string) (*models.Card, error) {
	return d.GetCard(filename)
}

// GetPassword implements grpcserv.DataProvider.
func (d *DataProvider) GetPassword(filename string) (*models.Password, error) {
	return d.GetPassword(filename)
}

// GetText implements grpcserv.DataProvider.
func (d *DataProvider) GetText(filename string) (string, error) {
	return d.GetText(filename)
}

func NewData(storage Storager) *DataProvider {
	return &DataProvider{}
}
