package sqlite

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"

	"library_importer/internal/importer/schema"
)

func Parse(path string) ([]schema.Record, error) {
	database, err := sql.Open("sqlite", path)

	if err != nil {
		return nil, fmt.Errorf("sql.Open Error: %v", err)
	}

	defer database.Close()

	rows, err := database.Query(`
		SELECT
			authors_firstName,
			authors_lastName,
			
			books_ISBN,
			books_title,
			
			genres_title,
			
			issues_issueDate,
			issues_period,
			issues_returnDate,
			
			readers_firstName,
			readers_lastName,
			readers_phoneNumber
		FROM library
	`)

	if err != nil {
		return nil, fmt.Errorf("database.Query Error: %w", err)
	}

	defer rows.Close()

	var records []schema.Record

	for rows.Next() {
		var row schema.Record

		err := rows.Scan(
			&row.AuthorFirstName,
			&row.AuthorLastName,

			&row.BookIsbn,
			&row.BookTitle,

			&row.GenreTitle,

			&row.IssueDate,
			&row.IssuePeriod,
			&row.IssueReturnDate,

			&row.ReaderFirstName,
			&row.ReaderLastName,
			&row.ReaderPhoneNumber,
		)

		if err != nil {
			return nil, fmt.Errorf("rows.Scan Error: %w", err)
		}

		records = append(records, row)
	}

	err = rows.Err()

	if err != nil {
		return nil, fmt.Errorf("rows.Err Error: %w", err)
	}

	return records, nil
}
