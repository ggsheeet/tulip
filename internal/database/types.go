package database

import (
	"database/sql"
	"time"
)

var loc, _ = time.LoadLocation("America/Monterrey")

type PostgresDB struct {
	db *sql.DB
}

type AccountInterface interface {
	DeleteAccount(string) error
	UpdateAccount(string, *Account) error
	GetAccountById(string) (*Account, error)
	GetAccounts() (*[]*Account, error)
	CreateAccount(*Account) error
}

type BookInterface interface {
	DeleteBook(string) error
	UpdateBook(string, *Book) error
	GetBookById(string) (*Book, error)
	GetBooks() (*[]*Book, error)
	CreateBook(*Book) error
	DeleteLetter(string) error
	UpdateLetter(string, *Letter) error
	GetLetterById(string) (*Letter, error)
	GetLetters() (*[]*Letter, error)
	CreateLetter(*Letter) error
	DeleteVersion(string) error
	UpdateVersion(string, *Version) error
	GetVersionById(string) (*Version, error)
	GetVersions() (*[]*Version, error)
	CreateVersion(*Version) error
	DeleteCover(string) error
	UpdateCover(string, *Cover) error
	GetCoverById(string) (*Cover, error)
	GetCovers() (*[]*Cover, error)
	CreateCover(*Cover) error
	DeletePublisher(string) error
	UpdatePublisher(string, *Publisher) error
	GetPublisherById(string) (*Publisher, error)
	GetPublishers() (*[]*Publisher, error)
	CreatePublisher(*Publisher) error
	DeleteBCategory(string) error
	UpdateBCategory(string, *BCategory) error
	GetBCategoryById(string) (*BCategory, error)
	GetBCategories() (*[]*BCategory, error)
	CreateBCategory(*BCategory) error
}

type ArticleInterface interface {
	DeleteArticle(string) error
	UpdateArticle(string, *Article) error
	GetArticleById(string) (*Article, error)
	GetArticles() (*[]*Article, error)
	CreateArticle(*Article) error
	DeleteACategory(string) error
	UpdateACategory(string, *ACategory) error
	GetACategoryById(string) (*ACategory, error)
	GetACategories() (*[]*ACategory, error)
	CreateACategory(*ACategory) error
}

type ResourceInterface interface {
	DeleteResource(string) error
	UpdateResource(string, *Resource) error
	GetResourceById(string) (*Resource, error)
	GetResources() (*[]*Resource, error)
	CreateResource(*Resource) error
	DeleteRCategory(string) error
	UpdateRCategory(string, *RCategory) error
	GetRCategoryById(string) (*RCategory, error)
	GetRCategories() (*[]*RCategory, error)
	CreateRCategory(*RCategory) error
}

type OrderInterface interface {
	DeleteOrder(string) error
	UpdateOrder(string, *Order) error
	GetOrderById(string) (*Order, error)
	GetOrders() (*[]*Order, error)
	CreateOrder(*Order) error
}

type AccountRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

type BookRequest struct {
	Title       string  `json:"title"`
	Author      string  `json:"author"`
	Description string  `json:"description"`
	ISBN        string  `json:"isbn"`
	CoverURL    string  `json:"coverUrl"`
	Price       float32 `json:"price"`
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
}

type ArticleRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	CoverURL    string `json:"coverUrl"`
	CategoryID  int    `json:"categoryId"`
}

type ACategoryRequest struct {
	ArticleCategory string `json:"articleCategory"`
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
}

type OrderRequest struct {
	FirstName string  `json:"firstName"`
	LastName  string  `json:"lastName"`
	Address   string  `json:"address"`
	Quantity  int     `json:"quantity"`
	Total     float32 `json:"total"`
	BookID    int     `json:"bookId"`
	AccountID string  `json:"accountId"`
}

type Account struct {
	ID        string    `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Book struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
	ISBN        string    `json:"isbn"`
	Price       float32   `json:"price"`
	Stock       int       `json:"stock"`
	SalesCount  int       `json:"salesCount"`
	IsActive    bool      `json:"isActive"`
	LetterID    int       `json:"letterId"`
	VersionID   int       `json:"versionId"`
	CoverID     int       `json:"coverId"`
	PublisherID int       `json:"publisherId"`
	CategoryID  int       `json:"categoryId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
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
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type Article struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
	CategoryID  int       `json:"categoryId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ACategory struct {
	ID              int       `json:"id"`
	ArticleCategory string    `json:"articleCategory"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type Resource struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	CoverURL    string    `json:"coverUrl"`
	ResourceURL string    `json:"resourceUrl"`
	CategoryID  int       `json:"categoryId"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type RCategory struct {
	ID               int       `json:"id"`
	ResourceCategory string    `json:"resourceCategory"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type Order struct {
	ID        int       `json:"id"`
	FirstName string    `json:"firstName"`
	LastName  string    `json:"lastName"`
	Address   string    `json:"address"`
	Quantity  int       `json:"quantity"`
	Total     float32   `json:"total"`
	BookID    int       `json:"bookId"`
	AccountID string    `json:"accountId"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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

func NewBook(title string, author string, description string, coverUrl string, isbn string, price float32, stock int, salesCount int, isActive bool, letterId int, versionId int, coverId int, publisherId int, categoryId int) *Book {
	return &Book{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		ISBN:        isbn,
		Price:       price,
		Stock:       stock,
		SalesCount:  salesCount,
		IsActive:    isActive,
		LetterID:    letterId,
		VersionID:   versionId,
		CoverID:     coverId,
		PublisherID: publisherId,
		CategoryID:  categoryId,
		CreatedAt:   time.Now().In(loc),
		UpdatedAt:   time.Now().In(loc),
	}
}

func NewLetter(letterType string) *Letter {
	return &Letter{
		LetterType: letterType,
		CreatedAt:  time.Now().In(loc),
		UpdatedAt:  time.Now().In(loc),
	}
}

func NewVersion(bibleVersion string) *Version {
	return &Version{
		BibleVersion: bibleVersion,
		CreatedAt:    time.Now().In(loc),
		UpdatedAt:    time.Now().In(loc),
	}
}

func NewCover(coverType string) *Cover {
	return &Cover{
		CoverType: coverType,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
	}
}

func NewPublisher(publisherName string) *Publisher {
	return &Publisher{
		PublisherName: publisherName,
		CreatedAt:     time.Now().In(loc),
		UpdatedAt:     time.Now().In(loc),
	}
}

func NewBCategory(bookCategory string) *BCategory {
	return &BCategory{
		BookCategory: bookCategory,
		CreatedAt:    time.Now().In(loc),
		UpdatedAt:    time.Now().In(loc),
	}
}

func NewArticle(title string, author string, description string, coverUrl string, categoryId int) *Article {
	return &Article{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		CategoryID:  categoryId,
		CreatedAt:   time.Now().In(loc),
		UpdatedAt:   time.Now().In(loc),
	}
}

func NewACategory(articleCategory string) *ACategory {
	return &ACategory{
		ArticleCategory: articleCategory,
		CreatedAt:       time.Now().In(loc),
		UpdatedAt:       time.Now().In(loc),
	}
}

func NewResource(title string, author string, description string, coverUrl string, resourceUrl string, categoryId int) *Resource {
	return &Resource{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		ResourceURL: resourceUrl,
		CategoryID:  categoryId,
		CreatedAt:   time.Now().In(loc),
		UpdatedAt:   time.Now().In(loc),
	}
}

func NewRCategory(resourceCategory string) *RCategory {
	return &RCategory{
		ResourceCategory: resourceCategory,
		CreatedAt:        time.Now().In(loc),
		UpdatedAt:        time.Now().In(loc),
	}
}

func NewOrder(firstName string, lastName string, address string, quantity int, total float32, bookId int, accountId string) *Order {
	return &Order{
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
		Quantity:  quantity,
		Total:     total,
		BookID:    bookId,
		AccountID: accountId,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
	}
}

func UpdateAccount(firstName, lastName, email string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		UpdatedAt: time.Now().In(loc),
	}
}

func UpdateBook(title string, author string, description string, coverUrl string, isbn string, price float32, stock int, salesCount int, isActive bool, letterId int, versionId int, coverId int, publisherId int, categoryId int) *Book {
	return &Book{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		ISBN:        isbn,
		Price:       price,
		Stock:       stock,
		SalesCount:  salesCount,
		IsActive:    isActive,
		LetterID:    letterId,
		VersionID:   versionId,
		CoverID:     coverId,
		PublisherID: publisherId,
		CategoryID:  categoryId,
		UpdatedAt:   time.Now().In(loc),
	}
}

func UpdateLetter(letterType string) *Letter {
	return &Letter{
		LetterType: letterType,
		UpdatedAt:  time.Now().In(loc),
	}
}

func UpdateVersion(bibleVersion string) *Version {
	return &Version{
		BibleVersion: bibleVersion,
		UpdatedAt:    time.Now().In(loc),
	}
}

func UpdateCover(coverType string) *Cover {
	return &Cover{
		CoverType: coverType,
		UpdatedAt: time.Now().In(loc),
	}
}

func UpdatePublisher(publisherName string) *Publisher {
	return &Publisher{
		PublisherName: publisherName,
		UpdatedAt:     time.Now().In(loc),
	}
}

func UpdateBCategory(bookCategory string) *BCategory {
	return &BCategory{
		BookCategory: bookCategory,
		UpdatedAt:    time.Now().In(loc),
	}
}

func UpdateArticle(title string, author string, description string, coverUrl string, categoryId int) *Article {
	return &Article{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		CategoryID:  categoryId,
		UpdatedAt:   time.Now().In(loc),
	}
}

func UpdateACategory(articleCategory string) *ACategory {
	return &ACategory{
		ArticleCategory: articleCategory,
		UpdatedAt:       time.Now().In(loc),
	}
}

func UpdateResource(title string, author string, description string, coverUrl string, resourceUrl string, categoryId int) *Resource {
	return &Resource{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		ResourceURL: resourceUrl,
		CategoryID:  categoryId,
		UpdatedAt:   time.Now().In(loc),
	}
}

func UpdateRCategory(resourceCategory string) *RCategory {
	return &RCategory{
		ResourceCategory: resourceCategory,
		UpdatedAt:        time.Now().In(loc),
	}
}

func UpdateOrder(firstName string, lastName string, address string, quantity int, total float32, bookId int, accountId string) *Order {
	return &Order{
		FirstName: firstName,
		LastName:  lastName,
		Address:   address,
		Total:     total,
		BookID:    bookId,
		AccountID: accountId,
		UpdatedAt: time.Now().In(loc),
	}
}
