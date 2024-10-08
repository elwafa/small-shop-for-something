package entities

type CardItem struct {
	ID     int   `json:"id"`
	CardID int   `json:"card_id"`
	ItemID int   `json:"item_id"`
	Item   *Item `json:"item"`
}
