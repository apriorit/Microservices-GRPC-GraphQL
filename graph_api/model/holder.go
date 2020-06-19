package model

type Holder struct {
	ID        string   `json:"id"`
	FirstName string   `json:"firstName"`
	LastName  string   `json:"lastName"`
	Phone     string   `json:"phone"`
	Email     string   `json:"email"`
	HeldBooks []string `json:"heldBooks"`
}
