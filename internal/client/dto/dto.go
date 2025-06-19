package dto

type PWD struct {
	Name     string
	Login    string
	Password string
}

type BIN struct {
	Name string
	Ext  string
	Data []byte
}

type BankCard struct {
	Name    string
	Number  string
	ExpDate string
	CVC     string
}

type Text struct {
	Name string
	Data string
}
