package database

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

var loc, _ = time.LoadLocation("America/Monterrey")

type PostgresDB struct {
	db *sql.DB
}

type AccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}

type BookRequest struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	ISBN        string  `json:"isbn"`
	CoverURL    string  `json:"coverUrl"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	SalesCount  int     `json:"salesCount"`
	IsActive    bool    `json:"isActive"`
	LetterID    int     `json:"letterId"`
	VersionID   int     `json:"versionId"`
	CoverID     int     `json:"coverId"`
	PublisherID int     `json:"publisherId"`
	CategoryID  int     `json:"categoryId"`
}

type LetterRequest struct {
	LetterType string `json:"letterType"`
}

type VersionRequest struct {
	BibleVersion string `json:"bibleVersion"`
}

type CoverRequest struct {
	CoverType string `json:"coverType"`
}

type PublisherRequest struct {
	PublisherName string `json:"publisherName"`
}

type BCategoryRequest struct {
	BookCategory string `json:"bookCategory"`
	IsActive     bool   `json:"isActive"`
}

type ArticleRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Excerpt     string `json:"excerpt"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	CategoryID  int    `json:"categoryId"`
}

type ACategoryRequest struct {
	ArticleCategory string `json:"articleCategory"`
	IsActive        bool   `json:"isActive"`
}

type ResourceRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	ResourceURL string `json:"resourceUrl"`
	CategoryID  int    `json:"categoryId"`
}

type RCategoryRequest struct {
	ResourceCategory string `json:"ResourceCategory"`
	IsActive         bool   `json:"isActive"`
}

type OrderRequest struct {
	Address     string    `json:"address"`
	Total       float64   `json:"total"`
	PaymentID   int       `json:"paymentId"`
	IsFulfilled bool      `json:"isFulfilled"`
	Status      string    `json:"status"`
	AccountID   uuid.UUID `json:"accountId"`
	OrderBooks  []OrderBook
}

type BookOrderRequest struct {
	Quantity int `json:"quantity"`
	BookID   int `json:"bookId"`
	OrderID  int `json:"orderId"`
}

type Account struct {
	ID        uuid.UUID `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Book struct {
	ID            int       `json:"id"`
	Title         string    `json:"title"`
	Author        string    `json:"author"`
	Description   string    `json:"description"`
	CoverURL      string    `json:"coverUrl"`
	ISBN          string    `json:"isbn"`
	Price         float64   `json:"price"`
	Stock         int       `json:"stock"`
	SalesCount    int       `json:"salesCount"`
	IsActive      bool      `json:"isActive"`
	LetterID      int       `json:"letterId"`
	LetterType    string    `json:"letterType"`
	VersionID     int       `json:"versionId"`
	BibleVersion  string    `json:"biblVersion"`
	CoverID       int       `json:"coverId"`
	CoverType     string    `json:"coverType"`
	PublisherID   int       `json:"publisherId"`
	PublisherName string    `json:"publisherName"`
	CategoryID    int       `json:"categoryId"`
	BookCategory  string    `json:"bCategory"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	RecordCount   int       `json:"recordCount,omitempty"`
}

type Letter struct {
	ID         int       `json:"id"`
	LetterType string    `json:"letterType"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type Version struct {
	ID           int       `json:"id"`
	BibleVersion string    `json:"bibleVersion"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Cover struct {
	ID        int       `json:"id"`
	CoverType string    `json:"coverType"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Publisher struct {
	ID            int       `json:"id"`
	PublisherName string    `json:"publisherName"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type BCategory struct {
	ID           int       `json:"id"`
	BookCategory string    `json:"bookCategory"`
	IsActive     bool      `json:"isActive"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Article struct {
	ID              int       `json:"id"`
	Title           string    `json:"title"`
	Author          string    `json:"author"`
	Excerpt         string    `json:"excerpt"`
	Description     string    `json:"description"`
	CoverURL        string    `json:"coverUrl"`
	CategoryID      int       `json:"categoryId"`
	ArticleCategory string    `json:"articleCategory"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
	RecordCount     int       `json:"recordCount,omitempty"`
}

type ACategory struct {
	ID              int       `json:"id"`
	ArticleCategory string    `json:"aCategory"`
	IsActive        bool      `json:"isActive"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Resource struct {
	ID               int       `json:"id"`
	Title            string    `json:"title"`
	Author           string    `json:"author"`
	Description      string    `json:"description"`
	CoverURL         string    `json:"coverUrl"`
	ResourceURL      string    `json:"resourceUrl"`
	CategoryID       int       `json:"categoryId"`
	ResourceCategory string    `json:"rCategory"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	RecordCount      int       `json:"recordCount,omitempty"`
}

type RCategory struct {
	ID               int       `json:"id"`
	ResourceCategory string    `json:"resourceCategory"`
	IsActive         bool      `json:"isActive"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type Order struct {
	ID          int       `json:"id"`
	Address     string    `json:"address"`
	Total       float64   `json:"total"`
	PaymentID   int       `json:"paymentId"`
	IsFulfilled bool      `json:"isFulfilled"`
	Status      string    `json:"status"`
	AccountID   uuid.UUID `json:"accountId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
	OrderBooks  []OrderBook
	RecordCount int `json:"recordCount,omitempty"`
}

type OrderBook struct {
	BookID      int     `json:"id"`
	Title       string  `json:"title"`
	Description string  `json:"description"`
	CoverURL    string  `json:"coverUrl"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
	BCategory   string  `json:"bCategory"`
}

type BookOrder struct {
	ID        int       `json:"id"`
	Quantity  int       `json:"quantity"`
	BookID    int       `json:"bookId"`
	OrderID   int       `json:"orderId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
