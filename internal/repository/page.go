package repository

import (
	"context"
	"database/sql"
	"fmt"
	"wiki-export/pkg/database"
)

type PageRepository interface {
	GetPagesToExport(ctx context.Context) ([]PageToExport, error)
}

type PageToExport struct {
	Id          int    `db:"page_id"`
	Slug        string `db:"page_slug"`
	Name        string `db:"page_name"`
	BookName    string `db:"book_name"`
	ChapterName string `db:"chapter_name"`
}

type repository struct {
	db database.Database
}

func NewPageRepository(db database.Database) PageRepository {
	return &repository{db: db}
}

func (r *repository) GetPagesToExport(ctx context.Context) ([]PageToExport, error) {
	query := `
        SELECT 
            pages.id AS page_id,
            pages.slug AS page_slug,
            pages.name AS page_name,
            books.name AS book_name,
            chapters.name AS chapter_name
        FROM pages
        LEFT JOIN books ON pages.book_id = books.id
        LEFT JOIN chapters ON pages.chapter_id = chapters.id
        WHERE pages.draft = 0
        AND pages.deleted_at IS NULL
        AND pages.text != ''
    `

	rows, err := r.db.GetInstance().QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %w", err)
	}
	defer rows.Close()

	var pages []PageToExport

	for rows.Next() {
		var page PageToExport
		bookName := sql.NullString{}
		chapterName := sql.NullString{}

		if err := rows.Scan(&page.Id, &page.Slug, &page.Name, &bookName, &chapterName); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		if bookName.Valid {
			page.BookName = bookName.String
		} else {
			page.BookName = ""
		}

		if chapterName.Valid {
			page.ChapterName = chapterName.String
		} else {
			page.ChapterName = ""
		}

		pages = append(pages, page)
	}

	return pages, nil
}
