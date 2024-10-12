package model

type User struct {
	Username      string
	Email         string
	ID            string
	Differentials []*Differential
}

type Differential struct {
	ID    string
	Title string
}
