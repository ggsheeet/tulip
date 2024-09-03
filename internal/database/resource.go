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

func (s *PostgresDB) GetResources(page int, limit int) (*[]*Resource, error) {
	offset := (page - 1) * limit

	rows, err := s.db.Query(getResourcesQ, limit, offset)
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
			&resource.CreatedAt,
			&resource.UpdatedAt,
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
