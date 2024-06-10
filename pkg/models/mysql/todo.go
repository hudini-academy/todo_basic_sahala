package mysql

import (
	"database/sql"
	"todo/pkg/models"
)

// Define a SnippetModel type which wraps a sql.DB connection pool.
type TodoModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.
func (m *TodoModel) Insert(title string) (int, error) {
	// Write the SQL statement we want to execute. I've split it over two lines
	// for readability (which is why it's surrounded with backquotes instead
	// of normal double quotes).
	stmt := `INSERT INTO todos (title, created, expires) 
VALUES(?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`
	// Use the Exec() method on the embedded connection pool to execute the
	// statement. The first parameter is the SQL statement, followed by the
	// title, content and expiry values for the placeholder parameters. This
	// method returns a sql.Result object, which contains some basic
	// information about what happened when the statement was executed.
	result, err := m.DB.Exec(stmt, title, 7)
	if err != nil {
		return 0, err
	}
	// Use the LastInsertId() method on the result object to get the ID of our
	// newly inserted record in the snippets table.
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	// The ID returned has the type int64, so we convert it to an int type
	// before returning.
	return int(id), nil
}

// This will return a specific todo based on its id.
func (m *TodoModel) Get(id int) (*models.Todo, error) {
	stmt := `SELECT * FROM todos 
	WHERE id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Todo{}
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Todo struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement. If the query returns no rows, then
	// row.Scan() will return a sql.ErrNoRows error. We check for that and retu
	// our own models.ErrNoRecord error instead of a Todo object.
	err := row.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord
	} else if err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippet object.
	return s, nil
}

// This will return the 10 most recently created snippets.
func (m *TodoModel) Latest() ([]*models.Todo, error) {
	stmt := `SELECT id, title, created, expires FROM todos 
    WHERE expires > UTC_TIMESTAMP() `

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	// We defer rows.Close() to ensure the sql.Rows resultset is
	// always properly closed before the Latest() method returns. This defer
	// statement should come *after* you check for an error from the Query()
	// method. Otherwise, if Query() returns an error, you'll get a panic
	// trying to close a nil resultset.
	defer rows.Close()

	// Initialize an empty slice to hold the models.Snippets objects.
	todos := []*models.Todo{}

	// Use rows.Next to iterate through the rows in the resultset. This
	// prepares the first (and then each subsequent) row to be acted on by the
	// rows.Scan() method. If iteration over all the rows completes then the
	// resultset automatically closes itself and frees-up the underlying
	// database connection.
	for rows.Next() {
		// Create a pointer to a new zeroed Snippet struct.
		s := &models.Todo{}
		// Use rows.Scan() to copy the values from each field in the row to the
		// new Snippet object that we created. Again, the arguments to row.Scan
		// must be pointers to the place you want to copy the data into, and the
		// number of arguments must be exactly the same as the number of
		// columns returned by your statement.
		err = rows.Scan(&s.ID, &s.Title, &s.Created, &s.Expires)
		if err != nil {
			return nil, err
		}
		// Append it to the slice of snippets.
		todos = append(todos, s)
	}
	// When the rows.Next() loop has finished we call rows.Err() to retrieve any
	// error that was encountered during the iteration. It's important to
	// call this - don't assume that a successful iteration was completed
	// over the whole resultset.
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// If everything went OK then return the Snippets slice.
	return todos, nil
}

func (m *TodoModel) Delete(id string) (*models.Todo, error) {
	stmt := `DELETE FROM todos 
	WHERE title = ?`

	_, err := m.DB.Exec(stmt, id)
	if err != nil {
		return nil, err
	}

	return nil, nil
}

func (m *TodoModel) Update(id int, title string) error {
	stmt := `UPDATE todos
	SET title=?
	WHERE id=?`

	_, err := m.DB.Exec(stmt, title, id)
	if err != nil {
		return err
	}

	return nil
}
