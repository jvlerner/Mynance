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

type Expense struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	DueDate     time.Time `json:"dueDate"`
	Paid        bool      `json:"paid"`
	CategoryID  *int      `json:"categoryId,omitempty"`
	CreatedAt   time.Time `json:"createdAt"`
	Deleted     bool      `json:"deleted"`
}
