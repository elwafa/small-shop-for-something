package entities

type Order struct {
	ID         int          `json:"id"`
	UserID     int          `json:"user_id"`
	Status     string       `json:"status"`
	OrderItems []*OrderItem `json:"order_items"`
}
