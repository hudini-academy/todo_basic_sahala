package mysql

import (
	"database/sql"
	"todo/pkg/models"
)

type SpecialModel struct {
	DB *sql.DB
}

func (m *SpecialModel) Insert(title string) (int, error) {

	stmt := `INSERT INTO specials (title, created, expires) 
VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY)) `

	result, err := m.DB.Exec(stmt, title, 7)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// func (m *SpecialModel) Get(id int) (*models.Special, error) {
// 	stmt := `SELECT * FROM specials
// 	WHERE id = ?`

// 	row := m.DB.QueryRow(stmt, id)
// 	s := &models.Special{}

// 	err := row.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
// 	if err == sql.ErrNoRows {
// 		return nil, models.ErrNoRecord
// 	} else if err != nil {
// 		return nil, err
// 	}

// 	return s, nil
// }

func (m *SpecialModel) Latest() ([]*models.Special, error) {
	stmt := `SELECT id, title, created, expires FROM specials 
    WHERE expires > UTC_TIMESTAMP() `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects.
	specials := []*models.Special{}
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Special{}
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		specials = append(specials, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return specials, nil
}
func (m *SpecialModel) Delete(title string) (*models.Special, error) {
	stmt := `DELETE FROM specials 
	WHERE title = ?`
	_, err := m.DB.Exec(stmt, title)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
