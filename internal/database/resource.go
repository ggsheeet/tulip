package database

import (
	"database/sql"
	"strconv"
	"time"
)

func (s *PostgresDB) createResourceTable() error {
	_, err := s.db.Exec(createResourceTabQ)

	return err
}

func (s *PostgresDB) CreateResource(resource *Resource) error {
	_, err := s.db.Query(
		createResourceQ,
		&resource.Title,
		&resource.Author,
		&resource.Description,
		&resource.CoverURL,
		&resource.ResourceURL,
		&resource.CategoryID,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteResource(id string) error {
	_, err := s.db.Exec(deleteResourceQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateResource(id string, resource *Resource) error {
	resId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateResourceQ,
		resId,
		&resource.Title,
		&resource.Author,
		&resource.Description,
		&resource.CoverURL,
		&resource.ResourceURL,
		&resource.CategoryID,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetResourceById(id string) (*Resource, error) {
	row := s.db.QueryRow(getResourceQ, id)

	var resource Resource

	err := row.Scan(
		&resource.ID,
		&resource.Title,
		&resource.Author,
		&resource.Description,
		&resource.CoverURL,
		&resource.ResourceURL,
		&resource.CategoryID,
		&resource.ResourceCategory,
		&resource.CreatedAt,
		&resource.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &resource, nil
}

func (s *PostgresDB) GetResources(page int, limit int, category int, order string, resourceId int) (*[]*Resource, error) {
	offset := (page - 1) * limit

	query := getResourcesQ

	whereClause := ""

	if category != 0 && resourceId != 0 {
		whereClause = " WHERE r.category_id = $3 AND r.id != $4"
	} else if category != 0 && resourceId == 0 {
		whereClause += " WHERE r.category_id = $3"
	} else if category == 0 && resourceId != 0 {
		whereClause += " WHERE r.id != $4"
	}

	query += whereClause

	switch order {
	case "newer":
		query += " ORDER BY r.created_at DESC"
	case "older":
		query += " ORDER BY r.created_at ASC"
	default:
		query += " ORDER BY r.created_at DESC, r.id DESC"
	}

	query += " LIMIT $1 OFFSET $2"

	args := []interface{}{limit, offset}
	if category != 0 {
		args = append(args, category)
	}
	if resourceId != 0 {
		args = append(args, resourceId)
	}

	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	resources := []*Resource{}
	for rows.Next() {
		resource := new(Resource)
		err := rows.Scan(
			&resource.ID,
			&resource.Title,
			&resource.Author,
			&resource.Description,
			&resource.CoverURL,
			&resource.ResourceURL,
			&resource.CategoryID,
			&resource.ResourceCategory,
			&resource.CreatedAt,
			&resource.UpdatedAt,
			&resource.RecordCount,
		)

		if err != nil {
			return nil, err
		}

		resources = append(resources, resource)
	}

	return &resources, nil
}

func (s *PostgresDB) createRCategoryTable() error {
	_, err := s.db.Exec(createRCategoryTabQ)

	return err
}

func (s *PostgresDB) CreateRCategory(rCategory *RCategory) error {
	_, err := s.db.Query(
		createRCategoryQ,
		&rCategory.ResourceCategory,
		&rCategory.CreatedAt,
		&rCategory.UpdatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) DeleteRCategory(id string) error {
	_, err := s.db.Exec(deleteRCategoryQ, id)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) UpdateRCategory(id string, rCategory *RCategory) error {
	rCategoryId, idErr := strconv.Atoi(id)
	if idErr != nil {
		return idErr
	}

	_, err := s.db.Query(
		updateRCategoryQ,
		rCategoryId,
		&rCategory.ResourceCategory,
		time.Now().In(loc),
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresDB) GetRCategoryById(id string) (*RCategory, error) {
	row := s.db.QueryRow(getRCategoryQ, id)

	var rCategory RCategory

	err := row.Scan(
		&rCategory.ID,
		&rCategory.ResourceCategory,
		&rCategory.CreatedAt,
		&rCategory.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	return &rCategory, nil
}

func (s *PostgresDB) GetRCategories() (*[]*RCategory, error) {
	rows, err := s.db.Query(getRCategoriesQ)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	rCategories := []*RCategory{}
	for rows.Next() {
		rCategory := new(RCategory)
		err := rows.Scan(
			&rCategory.ID,
			&rCategory.ResourceCategory,
			&rCategory.CreatedAt,
			&rCategory.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		rCategories = append(rCategories, rCategory)
	}
	return &rCategories, nil
}
