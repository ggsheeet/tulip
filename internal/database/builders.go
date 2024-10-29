package database

import "time"

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
