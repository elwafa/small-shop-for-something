package entities

type Card struct {
	ID     int        `json:"id"`
	UserID int        `json:"user_id"`
	Status string     `json:"status"`
	Items  []CardItem `json:"items"`
}
