package service

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"
	"wiki-export/internal/client/http"
	"wiki-export/internal/config/clients"
	"wiki-export/internal/repository"
)

type ExportService interface {
	ExportPagesToMarkdown(ctx context.Context) error
}

type exportService struct {
	repository     repository.PageRepository
	wikiHttpClient http.WikiHttpClient
}

func NewExportService(repository repository.PageRepository, cfg *clients.Wiki) ExportService {
	return &exportService{
		repository:     repository,
		wikiHttpClient: http.NewWikiHttpClient(cfg),
	}
}

func (s *exportService) ExportPagesToMarkdown(ctx context.Context) error {
	dir := "./tmp/markdown"

	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return fmt.Errorf("failed to create tmp directory: %w", err)
		}
	}

	pages, err := s.repository.GetPagesToExport(ctx)
	if err != nil {
		return err
	}

	chunkSize := 60
	chunks := len(pages) / chunkSize
	if len(pages)%chunkSize != 0 {
		chunks++
	}

	for i := 0; i < chunks; i++ {
		start := i * chunkSize
		end := (i + 1) * chunkSize
		if end > len(pages) {
			end = len(pages)
		}

		for _, page := range pages[start:end] {
			if err := s.exportPageToMarkdown(ctx, page); err != nil {
				log.Printf("Ignoring error while exporting page: %s\n", err)
			}
		}

		time.Sleep(1 * time.Minute)
	}

	f, err := os.OpenFile("./tmp/markdown.zip", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	zw := zip.NewWriter(f)
	defer zw.Close()

	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fh, err := zw.Create(path)
		if err != nil {
			return err
		}

		ff, err := os.Open(path)
		if err != nil {
			return err
		}
		defer ff.Close()

		_, err = io.Copy(fh, ff)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (s *exportService) exportPageToMarkdown(ctx context.Context, page repository.PageToExport) error {
	fileName := fmt.Sprintf("./tmp/markdown/%s.md", page.Slug)

	if _, err := os.Stat(fileName); err == nil || !os.IsNotExist(err) {
		return nil
	}

	content, err := s.wikiHttpClient.PagesExportMarkdown(ctx, page.Id)
	if err != nil {
		return err
	}

	markdown := string(content)

	re := regexp.MustCompile(`<[^>]*>`)
	markdown = re.ReplaceAllString(markdown, "")

	imageRe := regexp.MustCompile(`data:image\/[bmp,gif,ico,jpg,png,svg,webp,x\-icon,svg+xml]+;base64,[a-zA-Z0-9,+,/]+={0,2}`)
	markdown = imageRe.ReplaceAllString(markdown, "")

	breadcrumbs := fmt.Sprintf("Книги > %s", page.BookName)

	if page.ChapterName == "" {
		breadcrumbs = fmt.Sprintf("%s > %s", breadcrumbs, page.Name)
	} else {
		breadcrumbs = fmt.Sprintf("%s > %s > %s", breadcrumbs, page.ChapterName, page.Name)
	}

	newMarkdown := fmt.Sprintf("# %s\n\n%s", breadcrumbs, markdown)

	err = os.WriteFile(fileName, []byte(newMarkdown), 0644)
	if err != nil {
		return fmt.Errorf("failed to save file %s: %w", fileName, err)
	}

	return nil
}
