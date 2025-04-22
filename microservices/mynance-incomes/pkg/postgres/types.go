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

type Category struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	CreatedAt time.Time `json:"createdAt"`
	Active    bool      `json:"active"`
}

type Income struct {
	ID          int       `json:"id"`
	UserID      int       `json:"userId"`
	Description string    `json:"description"`
	Amount      float64   `json:"amount"`
	ReceivedAt  time.Time `json:"receivedAt"`
	IsRecurring bool      `json:"isRecurring"`
	Deleted     bool      `json:"deleted"`
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

type CreditCard struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userId"`
	Name        string  `json:"name"`
	Bank        string  `json:"bank"`
	LimitAmount float64 `json:"limitAmount"`
	DueDay      int     `json:"dueDay"`
	Active      bool    `json:"active"`
}

type CreditCardExpense struct {
	ID               int       `json:"id"`
	UserID           int       `json:"userId"`
	CardID           int       `json:"cardId"`
	Description      string    `json:"description"`
	Amount           float64   `json:"amount"`
	PurchaseDate     time.Time `json:"purchaseDate"`
	InstallmentCount int       `json:"installmentCount"`
	CategoryID       *int      `json:"categoryId,omitempty"`
	CreatedAt        time.Time `json:"createdAt"`
	Deleted          bool      `json:"deleted"`
}

type Payment struct {
	ID        int       `json:"id"`
	UserID    int       `json:"userId"`
	ExpenseID int       `json:"expenseId"`
	PaidAt    time.Time `json:"paidAt"`
	Amount    float64   `json:"amount"`
	Deleted   bool      `json:"deleted"`
}
