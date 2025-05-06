package models

type DataItem struct {
	data interface{}
}

type Password struct {
	Login    string
	Password string
}

type Text struct {
	Text string
}

type Binary struct {
	Binary []byte
}

type Card struct {
	Number string
	Holder string
	Exp    string
	Cvv    string
}

type DataList []DataItem
