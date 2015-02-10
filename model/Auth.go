package model

type Auth struct {
	Addr, User    string
	Password, Key []byte
}
