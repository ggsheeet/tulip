package app

import "github.com/ggsheet/tulip/internal/database"

type CartItem struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	CategoryID  string `json:"bCategory"`
	PictureURL  string `json:"coverUrl"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int    `json:"price"`
}

type Form struct {
	FirstName    string     `json:"firstName"`
	LastName     string     `json:"lastName"`
	Email        string     `json:"email"`
	Phone        string     `json:"phone"`
	Street       string     `json:"street"`
	Neighborhood string     `json:"neighborhood"`
	City         string     `json:"city"`
	ZipCode      string     `json:"zipcode"`
	State        string     `json:"state"`
	Country      string     `json:"country"`
	Cart         []CartItem `json:"cart"`
}

type PaymentData struct {
	ID     int    `json:"id"`
	Status string `json:"status"`
}

type EmailData struct {
	FirstName   string               `json:"firstName"`
	Email       string               `json:"email"`
	OrderNumber string               `json:"orderNumber"`
	SubTotal    float64              `json:"subTotal"`
	Shipping    float64              `json:"shipping"`
	Total       float64              `json:"total"`
	Cart        []database.OrderBook `json:"cart"`
}
