package postgres

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	Active    bool      `json:"active"`
}

type Profile struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
	Active    bool      `json:"active"`
}

type CreditCard struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userId"`
	Name        string  `json:"name"`
	Bank        string  `json:"bank"`
	LimitAmount float64 `json:"limitAmount"`
	DueDay      int     `json:"dueDay"`
	Active      bool    `json:"active"`
}
