package entities

type Item struct {
	ID          int     `json:"id"`
	Name        string  `json:"name" validate:"required"`
	Price       float64 `json:"price" validate:"required,numeric"`
	Picture     string  `json:"picture" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Status      string  `json:"status" validate:"required,oneof=New Used Sold"`
	UserId      int     `json:"user_id" validate:"required,numeric"`
	Receive     string  `json:"receive" validate:"required,oneof=PickUp Deliver BothOptionsArePossible"`
	Color       string  `json:"colour" validate:"required"`
	Category    string  `json:"category" validate:"required"`
}

func NewItem(name, picture, description, status, receive, color, category string, price float64, userId int) (*Item, error) {
	item := &Item{
		Name:        name,
		Price:       price,
		Picture:     picture,
		Description: description,
		Status:      status,
		UserId:      userId,
		Receive:     receive,
		Color:       color,
		Category:    category,
	}
	err := validate.Struct(item)
	if err != nil {
		return nil, err
	}
	return item, nil
}
