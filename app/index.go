package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/ggsheet/kerigma/internal/database"
	"github.com/ggsheet/kerigma/template/layout"
	"github.com/labstack/echo/v4"
)

func fetchData(url string, origin string, token string, page int, limit int, result interface{}, wg *sync.WaitGroup, errChan chan<- error, paginate bool) {
	defer wg.Done()

	if paginate {
		url = fmt.Sprintf("%s?page=%d&limit=%d", url, page, limit)
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errChan <- fmt.Errorf("failed to create request for %s: %w", url, err)
		return
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		errChan <- fmt.Errorf("failed to fetch data from %s: %w", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		errChan <- fmt.Errorf("non-200 response from %s: %s", url, resp.Status)
		return
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		errChan <- fmt.Errorf("failed to read response body from %s: %w", url, err)
		return
	}

	if err := json.Unmarshal(bodyBytes, result); err != nil {
		errChan <- fmt.Errorf("failed to unmarshal JSON from %s: %w", url, err)
		return
	}
}

func handleIndexPage(c echo.Context) error {
	env := os.Getenv("ENVIRONMENT")
	token := os.Getenv("AUTH_TOKEN")
	var origin string
	if env == "development" || env == "docker" {
		origin = os.Getenv("AUTH_ORIGIN_DEV")
	} else if env == "production" {
		origin = os.Getenv("AUTH_ORIGIN_PROD")
	} else {
		return c.String(http.StatusInternalServerError, "Invalid environment configuration")
	}

	page := 1
	limit := 10
	if pageParam := c.QueryParam("page"); pageParam != "" {
		page, _ = strconv.Atoi(pageParam)
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		limit, _ = strconv.Atoi(limitParam)
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var articles []database.Article
	var acategories []database.ACategory
	var resources []database.Resource
	var rcategories []database.RCategory
	var books []database.Book
	var letterTypes []database.Letter
	var bibleVersions []database.Version
	var coverTypes []database.Cover
	var publisherNames []database.Publisher
	var bcategories []database.BCategory

	wg.Add(1)
	go fetchData(origin+"/api/article", origin, token, page, limit, &articles, &wg, errChan, true)

	wg.Add(1)
	go fetchData(origin+"/api/article/acategory", origin, token, 0, 0, &acategories, &wg, errChan, false)

	wg.Add(1)
	go fetchData(origin+"/api/resource", origin, token, page, limit, &resources, &wg, errChan, true)

	wg.Add(1)
	go fetchData(origin+"/api/resource/rcategory", origin, token, 0, 0, &rcategories, &wg, errChan, false)

	wg.Add(1)
	go fetchData(origin+"/api/book", origin, token, page, limit, &books, &wg, errChan, true)

	wg.Add(1)
	go fetchData(origin+"/api/book/letter", origin, token, 0, 0, &letterTypes, &wg, errChan, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/version", origin, token, 0, 0, &bibleVersions, &wg, errChan, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/cover", origin, token, 0, 0, &coverTypes, &wg, errChan, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/publisher", origin, token, 0, 0, &publisherNames, &wg, errChan, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/bcategory", origin, token, 0, 0, &bcategories, &wg, errChan, false)

	wg.Wait()
	close(errChan)

	// Check if there were any errors
	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		// Combine all error messages and return them
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	aCategoryMap := make(map[int]string)
	for _, acategory := range acategories {
		aCategoryMap[acategory.ID] = acategory.ArticleCategory
	}

	rCategoryMap := make(map[int]string)
	for _, rcategory := range rcategories {
		rCategoryMap[rcategory.ID] = rcategory.ResourceCategory
	}

	letterTypeMap := make(map[int]string)
	for _, letterType := range letterTypes {
		letterTypeMap[letterType.ID] = letterType.LetterType
	}

	bVersionsMap := make(map[int]string)
	for _, bibleVersion := range bibleVersions {
		bVersionsMap[bibleVersion.ID] = bibleVersion.BibleVersion
	}

	coverTypeMap := make(map[int]string)
	for _, coverType := range coverTypes {
		coverTypeMap[coverType.ID] = coverType.CoverType
	}

	publishNameMap := make(map[int]string)
	for _, publisherName := range publisherNames {
		publishNameMap[publisherName.ID] = publisherName.PublisherName
	}

	bCategoryMap := make(map[int]string)
	for _, bcategory := range bcategories {
		bCategoryMap[bcategory.ID] = bcategory.BookCategory
	}

	for i := range articles {
		articles[i].ArticleCategory = aCategoryMap[articles[i].CategoryID]
	}

	for i := range resources {
		resources[i].ResourceCategory = rCategoryMap[resources[i].CategoryID]
	}

	for i := range books {
		books[i].LetterType = letterTypeMap[books[i].LetterID]
		books[i].BibleVersion = bVersionsMap[books[i].VersionID]
		books[i].CoverType = coverTypeMap[books[i].CoverID]
		books[i].PublisherName = publishNameMap[books[i].PublisherID]
		books[i].BookCategory = bCategoryMap[books[i].CategoryID]
	}

	// Render the page with the fetched data
	return Render(c, layout.Index(
		articles,
		resources,
		books,
	))
}
