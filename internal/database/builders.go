package database

import (
	"time"

	"github.com/google/uuid"
)

func NewAccount(firstName, lastName, email, phone string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
	}
}

func NewBook(title string, author string, description string, coverUrl string, isbn string, price float64, stock int, salesCount int, isActive bool, letterId int, versionId int, coverId int, publisherId int, categoryId int) *Book {
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

func NewBCategory(bookCategory string, isActive bool) *BCategory {
	return &BCategory{
		BookCategory: bookCategory,
		IsActive:     isActive,
		CreatedAt:    time.Now().In(loc),
		UpdatedAt:    time.Now().In(loc),
	}
}

func NewArticle(title string, author string, excerpt string, description string, coverUrl string, categoryId int) *Article {
	return &Article{
		Title:       title,
		Author:      author,
		Excerpt:     excerpt,
		Description: description,
		CoverURL:    coverUrl,
		CategoryID:  categoryId,
		CreatedAt:   time.Now().In(loc),
		UpdatedAt:   time.Now().In(loc),
	}
}

func NewACategory(articleCategory string, isActive bool) *ACategory {
	return &ACategory{
		ArticleCategory: articleCategory,
		IsActive:        isActive,
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

func NewRCategory(resourceCategory string, isActive bool) *RCategory {
	return &RCategory{
		ResourceCategory: resourceCategory,
		IsActive:         isActive,
		CreatedAt:        time.Now().In(loc),
		UpdatedAt:        time.Now().In(loc),
	}
}

func NewOrder(address string, total float64, paymentId int, isFulfilled bool, status string, accountId uuid.UUID) *Order {
	return &Order{
		Address:     address,
		Total:       total,
		PaymentID:   paymentId,
		IsFulfilled: isFulfilled,
		Status:      status,
		AccountID:   accountId,
		CreatedAt:   time.Now().In(loc),
		UpdatedAt:   time.Now().In(loc),
	}
}

func NewBookOrder(quantity int, bookId int, orderId int) *BookOrder {
	return &BookOrder{
		Quantity:  quantity,
		BookID:    bookId,
		OrderID:   orderId,
		CreatedAt: time.Now().In(loc),
		UpdatedAt: time.Now().In(loc),
	}
}

func UpdateAccount(firstName, lastName, email, phone string) *Account {
	return &Account{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
		UpdatedAt: time.Now().In(loc),
	}
}

func UpdateBook(title string, author string, description string, coverUrl string, isbn string, price float64, stock int, salesCount int, isActive bool, letterId int, versionId int, coverId int, publisherId int, categoryId int) *Book {
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

func UpdateBCategory(bookCategory string, isActive bool) *BCategory {
	return &BCategory{
		BookCategory: bookCategory,
		IsActive:     isActive,
		UpdatedAt:    time.Now().In(loc),
	}
}

func UpdateArticle(title string, author string, excerpt, description string, coverUrl string, categoryId int) *Article {
	return &Article{
		Title:       title,
		Author:      author,
		Excerpt:     excerpt,
		Description: description,
		CoverURL:    coverUrl,
		CategoryID:  categoryId,
		UpdatedAt:   time.Now().In(loc),
	}
}

func UpdateACategory(articleCategory string, isActive bool) *ACategory {
	return &ACategory{
		ArticleCategory: articleCategory,
		IsActive:        isActive,
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

func UpdateRCategory(resourceCategory string, isActive bool) *RCategory {
	return &RCategory{
		ResourceCategory: resourceCategory,
		IsActive:         isActive,
		UpdatedAt:        time.Now().In(loc),
	}
}

func UpdateOrder(address string, total float64, paymentId int, isFulfilled bool, status string, accountId uuid.UUID) *Order {
	return &Order{
		Address:     address,
		Total:       total,
		PaymentID:   paymentId,
		IsFulfilled: isFulfilled,
		Status:      status,
		AccountID:   accountId,
		CreatedAt:   time.Now().In(loc),
		UpdatedAt:   time.Now().In(loc),
	}
}

func UpdateBookOrder(quantity int, bookId int, orderId int) *BookOrder {
	return &BookOrder{
		Quantity:  quantity,
		BookID:    bookId,
		OrderID:   orderId,
		UpdatedAt: time.Now().In(loc),
	}
}
