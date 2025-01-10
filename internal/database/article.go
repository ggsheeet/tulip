package database

import (
	"database/sql"
	"fmt"
	"strconv"
	"time"
)

func (s *PostgresDB) createArticleTable() error {
	_, err := s.db.Exec(createArticleTabQ)

	return err
}

func (s *PostgresDB) CreateArticle(article *Article) error {
	_, err := s.db.Query(
		createArticleQ,
		&article.Title,
		&article.Author,
		&article.Excerpt,
		&article.Description,
		&article.CoverURL,
		&article.CategoryID,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteArticle(id string) error {
	_, err := s.db.Exec(deleteArticleQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateArticle(id string, article *Article) error {
	artId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateArticleQ,
		artId,
		&article.Title,
		&article.Author,
		&article.Excerpt,
		&article.Description,
		&article.CoverURL,
		&article.CategoryID,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetArticleById(id string) (*Article, error) {
	row := s.db.QueryRow(getArticleQ, id)

	var article Article

	err := row.Scan(
		&article.ID,
		&article.Title,
		&article.Author,
		&article.Excerpt,
		&article.Description,
		&article.CoverURL,
		&article.CategoryID,
		&article.ArticleCategory,
		&article.CreatedAt,
		&article.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Article with id '%v' does not exist", id)
		}
		return nil, err
	}

	return &article, nil
}

func (s *PostgresDB) GetArticles(page int, limit int, category int, order string, articleId int) (*[]*Article, error) {
	offset := (page - 1) * limit

	query := getArticlesQ

	whereClause := ""

	if category != 0 && articleId != 0 {
		whereClause = " WHERE a.category_id = $3 AND a.id != $4"
	} else if category != 0 && articleId == 0 {
		whereClause = " WHERE a.category_id = $3"
	} else if category == 0 && articleId != 0 {
		whereClause = " WHERE a.id != $4"
	}

	query += whereClause

	switch order {
	case "newer":
		query += " ORDER BY a.created_at DESC"
	case "older":
		query += " ORDER BY a.created_at ASC"
	default:
		query += " ORDER BY a.created_at DESC, a.id DESC"
	}

	query += " LIMIT $1 OFFSET $2"

	args := []interface{}{limit, offset}
	if category != 0 {
		args = append(args, category)
	}
	if articleId != 0 {
		args = append(args, articleId)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	articles := []*Article{}
	for rows.Next() {
		article := new(Article)
		err := rows.Scan(
			&article.ID,
			&article.Title,
			&article.Author,
			&article.Excerpt,
			&article.Description,
			&article.CoverURL,
			&article.CategoryID,
			&article.ArticleCategory,
			&article.CreatedAt,
			&article.UpdatedAt,
			&article.RecordCount,
		)

		if err != nil {
			return nil, err
		}

		articles = append(articles, article)
	}

	return &articles, nil
}

func (s *PostgresDB) createACategoryTable() error {
	_, err := s.db.Exec(createACategoryTabQ)

	return err
}

func (s *PostgresDB) CreateACategory(aCategory *ACategory) error {
	_, err := s.db.Query(
		createACategoryQ,
		&aCategory.ArticleCategory,
		&aCategory.CreatedAt,
		&aCategory.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteACategory(id string) error {
	_, err := s.db.Exec(deleteACategoryQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateACategory(id string, aCategory *ACategory) error {
	aCategoryId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateACategoryQ,
		aCategoryId,
		&aCategory.ArticleCategory,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetACategoryById(id string) (*ACategory, error) {
	row := s.db.QueryRow(getACategoryQ, id)

	var aCategory ACategory

	err := row.Scan(
		&aCategory.ID,
		&aCategory.ArticleCategory,
		&aCategory.CreatedAt,
		&aCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &aCategory, nil
}

func (s *PostgresDB) GetACategories() (*[]*ACategory, error) {
	rows, err := s.db.Query(getACategoriesQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	aCategories := []*ACategory{}
	for rows.Next() {
		aCategory := new(ACategory)
		err := rows.Scan(
			&aCategory.ID,
			&aCategory.ArticleCategory,
			&aCategory.CreatedAt,
			&aCategory.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		aCategories = append(aCategories, aCategory)
	}
	return &aCategories, nil
}
