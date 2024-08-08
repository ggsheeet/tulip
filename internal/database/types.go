package database

import (
	"database/sql"
	"time"
)

var loc, _ = time.LoadLocation("America/Monterrey")

type CreateAccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type CreateBookRequest struct {
	Title        string  `json:"title"`
	Author       string  `json:"author"`
	Description  string  `json:"description"`
	ISBN         string  `json:"isbn"`
	CoverURL     string  `json:"coverUrl"`
	Price        float32 `json:"price"`
	Stock        int     `json:"stock"`
	SalesCount   int     `json:"salesCount"`
	IsActive     int     `json:"isActive"`
	LetterSizeID int     `json:"letterSizeId"`
	VersionID    int     `json:"versionId"`
	CoverID      int     `json:"coverId"`
	CategoryID   int     `json:"categoryId"`
	PublisherID  int     `json:"publisherId"`
}

type Account struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Book struct {
	ID           int       `json:"id"`
	Title        string    `json:"title"`
	Author       string    `json:"author"`
	Description  string    `json:"description"`
	CoverURL     string    `json:"coverUrl"`
	ISBN         string    `json:"isbn"`
	Price        float32   `json:"price"`
	Stock        int       `json:"stock"`
	SalesCount   int       `json:"salesCount"`
	IsActive     int       `json:"isActive"`
	LetterSizeID int       `json:"letterSizeId"`
	VersionID    int       `json:"versionId"`
	CoverID      int       `json:"coverId"`
	CategoryID   int       `json:"categoryId"`
	PublisherID  int       `json:"publisherId"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type AccountInterface interface {
	DeleteAccount(string) error
	UpdateaAccount(*Account) error
	GetAccountById(string) (*Account, error)
	GetAccounts() (*[]*Account, error)
	CreateAccount(acc *Account) error
}

type BookInterface interface {
	DeleteBook(string) error
	UpdateaBook(*Book) error
	GetBookById(string) (*Book, error)
	GetBooks() (*[]*Book, error)
	CreateBook(book *Book) error
}

type PostgresDB struct {
	db *sql.DB
}

func NewAccount(firstName, lastName, email string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
	}
}

func NewBook(title string, author string, description string, coverUrl string, isbn string, price float32, stock int, salesCount int, isActive int, letterSizeId int, versionId int, coverId int, categoryId int, publisherId int) *Book {
	return &Book{
		Title:        title,
		Author:       author,
		Description:  description,
		CoverURL:     coverUrl,
		ISBN:         isbn,
		Price:        price,
		Stock:        stock,
		SalesCount:   salesCount,
		IsActive:     isActive,
		LetterSizeID: letterSizeId,
		VersionID:    versionId,
		CoverID:      coverId,
		CategoryID:   categoryId,
		PublisherID:  publisherId,
		CreatedAt:    time.Now().In(loc),
		UpdatedAt:    time.Now().In(loc),
	}
}
