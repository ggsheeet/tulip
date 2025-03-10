package database

import "github.com/google/uuid"

type AccountInterface interface {
	DeleteAccount(uuid.UUID) error
	UpdateAccount(uuid.UUID, *Account) error
	GetAccountById(uuid.UUID) (*Account, error)
	GetAccountByEmail(string) (*Account, error)
	GetAccounts() (*[]*Account, error)
	CreateAccount(*Account) (uuid.UUID, error)
}

type BookInterface interface {
	DeleteBook(string) error
	UpdateBook(string, *Book) error
	UpdateBookStock(int, int) error
	GetBookById(string) (*Book, error)
	GetBooks(int, int, int, string, int, string) (*[]*Book, error)
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
	GetArticles(int, int, int, string, int) (*[]*Article, error)
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
	GetResources(int, int, int, string, int) (*[]*Resource, error)
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
	GetOrderByPaymentId(int) (int, error)
	GetUnfulfilledOrders() (*[]*Order, error)
	GetFulfilledOrders() (*[]*Order, error)
	CreateOrder(*Order) (int, error)
	CreateBookOrder(*OrderBook, int) error
}
