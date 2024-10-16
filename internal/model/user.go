package model

type User struct {
	Username      string
	Email         string
	ID            string
	Differentials []*Differential
	Reviews       []*Differential
}

type Differential struct {
	ID        string
	Title     string
	LineCount string
}
