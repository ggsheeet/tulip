package app

import (
	"github.com/ggsheet/tulip/internal/database"
)

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
	StreetNumber string     `json:"streetNumber"`
	City         string     `json:"city"`
	ZipCode      string     `json:"zipcode"`
	State        string     `json:"state"`
	Country      string     `json:"country"`
	Cart         []CartItem `json:"cart"`
}

type PaymentData struct {
	ID                 int                `json:"id"`
	Status             string             `json:"status"`
	AdditionalInfo     AdditionalInfo     `json:"additional_info"`
	ExternalReference  *string            `json:"external_reference"`
	ShippingAmount     float64            `json:"shipping_amount"`
	TransactionAmount  float64            `json:"transaction_amount"`
	TransactionDetails TransactionDetails `json:"transaction_details"`
}

type AdditionalInfo struct {
	ItemsInfo     []database.OrderBook `json:"items"`
	ShipmentsInfo ShipmentsInfo        `json:"shipments"`
	Payer         Payer                `json:"payer"`
}

type ShipmentsInfo struct {
	ReceiverAddress ReceiverAddress `json:"receiver_address"`
}

type ReceiverAddress struct {
	CityName     string `json:"city_name"`
	StateName    string `json:"state_name"`
	StreetName   string `json:"street_name"`
	StreetNumber string `json:"street_number"`
	ZipCode      string `json:"zip_code"`
}

type Payer struct {
	FirstName string `json:"first_name"`
}

type TransactionDetails struct {
	TotalPaidAmount float64 `json:"total_paid_amount"`
}

type EmailData struct {
	FirstName   string               `json:"firstName"`
	LastName    string               `json:"lastName"`
	Email       string               `json:"email"`
	OrderNumber string               `json:"orderNumber"`
	SubTotal    float64              `json:"subTotal"`
	Shipping    float64              `json:"shipping"`
	Total       float64              `json:"total"`
	Cart        []database.OrderBook `json:"cart"`
}

type Notification struct {
	Data NotificationData
}

type NotificationData struct {
	ID string `json:"id"`
}
