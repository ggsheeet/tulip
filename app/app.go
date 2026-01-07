package app

import (
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ggsheet/tulip/internal/database"
	"github.com/ggsheet/tulip/template/component"
	"github.com/ggsheet/tulip/template/layout"
	"github.com/labstack/echo/v4"
)

var FS fs.FS

// New fetchData function that doesn't require WaitGroup parameter
func fetchDataSimple(url string, origin string, token string, page int, limit int, category int, order string, itemId int, result interface{}, errChan chan<- error, paginate bool, filter bool) {

	if paginate {
		if filter {
			if category != 0 {
				if order != "" {
					url = fmt.Sprintf("%s?page=%d&limit=%d&category=%d&order=%s", url, page, limit, category, order)
				} else if itemId != 0 {
					url = fmt.Sprintf("%s?page=%d&limit=%d&category=%d&itemId=%d", url, page, limit, category, itemId)
				} else {
					url = fmt.Sprintf("%s?page=%d&limit=%d&category=%d", url, page, limit, category)
				}
			} else if order != "" {
				url = fmt.Sprintf("%s?page=%d&limit=%d&order=%s", url, page, limit, order)
			} else {
				url = fmt.Sprintf("%s?page=%d&limit=%d", url, page, limit)
			}
		} else {
			url = fmt.Sprintf("%s?page=%d&limit=%d", url, page, limit)
		}
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

	log.Printf("Fetched data successfully from %s", url)
}

// Legacy fetchData function for backward compatibility
func fetchData(url string, origin string, token string, page int, limit int, category int, order string, itemId int, result interface{}, wg *sync.WaitGroup, errChan chan<- error, paginate bool, filter bool) {
	defer wg.Done()
	fetchDataSimple(url, origin, token, page, limit, category, order, itemId, result, errChan, paginate, filter)
}

func handleIndexPage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	page := 1
	limit := 10
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var articles []database.Article
	var resources []database.Resource
	var books []database.Book

	wg.Go(func() {
		fetchDataSimple(origin+"/api/article", origin, token, page, limit, 0, "", 0, &articles, errChan, true, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/resource", origin, token, page, limit, 0, "", 0, &resources, errChan, true, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book", origin, token, page, limit, 0, "", 0, &books, errChan, true, false)
	})

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	return Render(c, layout.Index(
		articles,
		resources,
		books,
	))
}

func handleStorePage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	page := 1
	limit := 10
	filter := false
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		} else {
			filter = true
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}
	category := 0
	order := ""
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil {
			log.Printf("Error parsing category parameter: %v", err)
			category = 0
		} else {
			filter = true
		}
	}
	if orderParam := c.QueryParam("order"); orderParam != "" {
		order = orderParam
		filter = true
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var books []database.Book
	var bcategories []database.BCategory

	wg.Add(1)
	go fetchData(origin+"/api/book", origin, token, page, limit, category, order, 0, &books, &wg, errChan, true, filter)

	wg.Add(1)
	go fetchData(origin+"/api/book/bcategory", origin, token, page, limit, 0, "", 0, &bcategories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		log.Printf("Data fetch errors: %v", errMessages)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	var totalPages int
	if len(books) > 0 {
		totalPages = int(math.Ceil(float64(books[0].RecordCount) / float64(limit)))
	} else {
		totalPages = 1
	}

	if clearParam := c.QueryParam("clear"); clearParam == "true" || filter {
		return Render(c, component.BookGrid(books, page, totalPages))
	}

	return Render(c, layout.Store(books, bcategories, page, totalPages))
}

func handleArticlesPage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	page := 1
	limit := 10
	filter := false
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		} else {
			filter = true
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}
	category := 0
	order := ""
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil {
			log.Printf("Error parsing category parameter: %v", err)
			category = 0
		} else {
			filter = true
		}
	}
	if orderParam := c.QueryParam("order"); orderParam != "" {
		order = orderParam
		filter = true
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var articles []database.Article
	var acategories []database.ACategory

	wg.Add(1)
	go fetchData(origin+"/api/article", origin, token, page, limit, category, order, 0, &articles, &wg, errChan, true, filter)

	wg.Add(1)
	go fetchData(origin+"/api/article/acategory", origin, token, page, limit, 0, "", 0, &acategories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		log.Printf("Data fetch errors: %v", errMessages)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	var totalPages int
	if len(articles) > 0 {
		totalPages = int(math.Ceil(float64(articles[0].RecordCount) / float64(limit)))
	} else {
		totalPages = 1
	}

	if clearParam := c.QueryParam("clear"); clearParam == "true" || filter {
		return Render(c, component.ArticleGrid(articles, page, totalPages))
	}

	return Render(c, layout.Articles(articles, acategories, page, totalPages))
}

func handleResourcesPage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	page := 1
	limit := 10
	filter := false
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		} else {
			filter = true
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}
	category := 0
	order := ""
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil {
			log.Printf("Error parsing category parameter: %v", err)
			category = 0
		} else {
			filter = true
		}
	}
	if orderParam := c.QueryParam("order"); orderParam != "" {
		order = orderParam
		filter = true
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var resources []database.Resource
	var rcategories []database.RCategory

	wg.Add(1)
	go fetchData(origin+"/api/resource", origin, token, page, limit, category, order, 0, &resources, &wg, errChan, true, filter)

	wg.Add(1)
	go fetchData(origin+"/api/resource/rcategory", origin, token, page, limit, 0, "", 0, &rcategories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		log.Printf("Data fetch errors: %v", errMessages)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	var totalPages int
	if len(resources) > 0 {
		totalPages = int(math.Ceil(float64(resources[0].RecordCount) / float64(limit)))
	} else {
		totalPages = 1
	}

	if clearParam := c.QueryParam("clear"); clearParam == "true" || filter {
		return Render(c, component.ResourceGrid(resources, page, totalPages))
	}

	return Render(c, layout.Resources(resources, rcategories, page, totalPages))
}

func handleBookPage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	id := ""
	bookId := 0
	if idParam := c.QueryParam("id"); idParam != "" {
		id = idParam
		var err error
		bookId, err = strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			bookId = 0
		}
	}
	page := 1
	limit := 10
	category := 0
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil {
			log.Printf("Error parsing category parameter: %v", err)
			category = 0
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var book database.Book
	var books []database.Book

	wg.Add(1)
	go fetchData(origin+"/api/book/"+id, origin, token, 0, 0, 0, "", 0, &book, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book", origin, token, page, limit, category, "", bookId, &books, &wg, errChan, true, true)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		log.Printf("Data fetch errors: %v", errMessages)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	return Render(c, layout.Book(book, books))
}

func handleArticlePage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	id := ""
	articleId := 0
	if idParam := c.QueryParam("id"); idParam != "" {
		id = idParam
		var err error
		articleId, err = strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			articleId = 0
		}
	}
	page := 1
	limit := 10
	category := 0
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil {
			log.Printf("Error parsing category parameter: %v", err)
			category = 0
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var article database.Article
	var articles []database.Article

	wg.Add(1)
	go fetchData(origin+"/api/article/"+id, origin, token, 0, 0, 0, "", 0, &article, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/article", origin, token, page, limit, category, "", articleId, &articles, &wg, errChan, true, true)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		log.Printf("Data fetch errors: %v", errMessages)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	return Render(c, layout.Article(article, articles))
}

func handleResourcePage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	id := ""
	resourceId := 0
	if idParam := c.QueryParam("id"); idParam != "" {
		id = idParam
		var err error
		resourceId, err = strconv.Atoi(idParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			resourceId = 0
		}
	}
	page := 1
	limit := 10
	category := 0
	if pageParam := c.QueryParam("page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing page parameter: %v", err)
			page = 1
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}
	if categoryParam := c.QueryParam("category"); categoryParam != "" {
		var err error
		category, err = strconv.Atoi(categoryParam)
		if err != nil {
			log.Printf("Error parsing category parameter: %v", err)
			category = 0
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	var resource database.Resource
	var resources []database.Resource

	wg.Add(1)
	go fetchData(origin+"/api/resource/"+id, origin, token, 0, 0, 0, "", 0, &resource, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/resource", origin, token, page, limit, category, "", resourceId, &resources, &wg, errChan, true, true)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		log.Printf("Data fetch errors: %v", errMessages)
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	return Render(c, layout.Resource(resource, resources))
}

func handleResourceDownload(c echo.Context) error {
	resourceUrl := c.QueryParam("rUrl")
	if resourceUrl == "" {
		return c.String(http.StatusBadRequest, "Missing resource URL")
	}

	req, err := http.NewRequest("GET", resourceUrl, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create request: %v", err))
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch resource: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.String(resp.StatusCode, fmt.Sprintf("Failed to fetch resource: %s", resp.Status))
	}

	c.Response().Header().Set("Content-Type", "application/pdf")

	if _, err := io.Copy(c.Response().Writer, resp.Body); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to stream file: %v", err))
	}

	return nil
}

func handleCartPage(c echo.Context) error {
	return Render(c, layout.Cart())
}

func handleProcesedPage(c echo.Context) error {
	paymentId := ""
	status := ""
	if paymentIdParam := c.QueryParam("payment_id"); paymentIdParam != "null" {
		paymentId = paymentIdParam
	}
	if statusParam := c.QueryParam("status"); statusParam != "" {
		if statusParam == "approved" {
			status = "Exitosa"
		} else {
			status = "Fallda"
		}
	}

	return Render(c, layout.Processed(paymentId, status))
}

func handlePaymentNotification(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var notification Notification
	if err := c.Bind(&notification); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}
	log.Printf("Received webhook: %+v", notification)

	if notification.Data.ID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Order ID cannot be empty"})
	}

	url := fmt.Sprintf("%s/api/payment/confirmed?payment_id=%s", origin, notification.Data.ID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to create order: %v", err)})
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	log.Println("Response from order confirmation:", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": fmt.Sprintf("Failed to confirm order: %s", resp.Status),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Order successfully created"})
}

func handleAuthCheck(c echo.Context) error {
	cookieName := os.Getenv("COOKIE_NAME")
	cookieValue := os.Getenv("COOKIE_VALUE")

	cookie, err := c.Cookie(cookieName)
	if err != nil || cookie.Value != cookieValue {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User must login"})
	}

	if cookie.Value == cookieValue {
		return c.JSON(http.StatusOK, map[string]string{"error": "User authenticated"})
	}
	return c.JSON(http.StatusForbidden, map[string]string{"error": "Unknown error occured"})
}

func handleLoginPage(c echo.Context) error {
	return Render(c, layout.Login())
}

func handleLoginAuth(c echo.Context) error {
	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	cookieName := os.Getenv("COOKIE_NAME")
	cookieValue := os.Getenv("COOKIE_VALUE")

	inputEmail := c.FormValue("adminEmail")
	inputPassword := c.FormValue("adminPassword")

	if inputEmail == adminEmail && inputPassword == adminPassword {
		cookie := new(http.Cookie)
		cookie.Name = cookieName
		cookie.Value = cookieValue
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)

		return c.JSON(http.StatusOK, map[string]string{
			"redirect": "/admin",
		})
	}

	emailError := ""
	passwordError := ""

	if inputEmail != adminEmail {
		emailError = "Correo incorrecto"
	}
	if inputPassword != adminPassword {
		passwordError = "Contraseña incorrecta"
	}

	return c.JSON(http.StatusUnauthorized, map[string]string{
		"emailError":    emailError,
		"passwordError": passwordError,
	})
}

func handleLogoutAuth(c echo.Context) error {
	cookieName := os.Getenv("COOKIE_NAME")
	cookie := new(http.Cookie)
	cookie.Name = cookieName
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}

func renderAdminTableContent(c echo.Context, tableData component.CollectiveTableData, activeTab string) error {
	switch activeTab {
	case "book":
		return Render(c, component.AdminTableSection(tableData.Books, tableData.BooksPage, tableData.BooksTotalPages, "book"))
	case "article":
		return Render(c, component.AdminTableSection(tableData.Articles, tableData.ArticlesPage, tableData.ArticlesTotalPages, "article"))
	case "resource":
		return Render(c, component.AdminTableSection(tableData.Resources, tableData.ResourcesPage, tableData.ResourcesTotalPages, "resource"))
	case "bcategory":
		return Render(c, component.AdminTableSection(tableData.BookCategories, tableData.BookCategoriesPage, tableData.BookCategoriesTotalPages, "bcategory"))
	case "acategory":
		return Render(c, component.AdminTableSection(tableData.ArticleCategories, tableData.ArticleCategoriesPage, tableData.ArticleCategoriesTotalPages, "acategory"))
	case "rcategory":
		return Render(c, component.AdminTableSection(tableData.ResourceCategories, tableData.ResourceCategoriesPage, tableData.ResourceCategoriesTotalPages, "rcategory"))
	case "publisher":
		return Render(c, component.AdminTableSection(tableData.Publishers, tableData.PublishersPage, tableData.PublishersTotalPages, "publisher"))
	case "version":
		return Render(c, component.AdminTableSection(tableData.Versions, tableData.VersionsPage, tableData.VersionsTotalPages, "version"))
	case "letter":
		return Render(c, component.AdminTableSection(tableData.Letters, tableData.LettersPage, tableData.LettersTotalPages, "letter"))
	case "cover":
		return Render(c, component.AdminTableSection(tableData.Covers, tableData.CoversPage, tableData.CoversTotalPages, "cover"))
	case "order":
		return Render(c, component.AdminTableSection(tableData.Orders, tableData.OrdersPage, tableData.OrdersTotalPages, "order"))
	default:
		return Render(c, component.AdminTableSection(tableData.Books, tableData.BooksPage, tableData.BooksTotalPages, "book"))
	}
}

func renderUpdatedTableSection(c echo.Context, tabType string, page int) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")
	limit := 10

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	switch tabType {
	case "book":
		var books []database.Book
		wg.Add(1)
		go fetchData(origin+"/api/book/admin", origin, token, page, limit, 0, "", 0, &books, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(books) > 0 {
			totalPages = int(math.Ceil(float64(books[0].RecordCount) / float64(limit)))
		} else {
			totalPages = 1
		}

		bookData := getBookTableData(books)
		return Render(c, component.AdminTableSection(bookData, page, totalPages, "book"))

	case "article":
		var articles []database.Article
		wg.Add(1)
		go fetchData(origin+"/api/article", origin, token, page, limit, 0, "", 0, &articles, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(articles) > 0 {
			totalPages = int(math.Ceil(float64(articles[0].RecordCount) / float64(limit)))
		} else {
			totalPages = 1
		}

		articleData := getArticleTableData(articles)
		return Render(c, component.AdminTableSection(articleData, page, totalPages, "article"))

	case "resource":
		var resources []database.Resource
		wg.Add(1)
		go fetchData(origin+"/api/resource", origin, token, page, limit, 0, "", 0, &resources, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(resources) > 0 {
			totalPages = int(math.Ceil(float64(resources[0].RecordCount) / float64(limit)))
		} else {
			totalPages = 1
		}

		resourceData := getResourceTableData(resources)
		return Render(c, component.AdminTableSection(resourceData, page, totalPages, "resource"))

	case "bcategory":
		var bookCategories []database.BCategory
		wg.Add(1)
		go fetchData(origin+"/api/book/bcategory", origin, token, page, limit, 0, "", 0, &bookCategories, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(bookCategories) > 0 {
			totalPages = int(math.Ceil(float64(len(bookCategories)) / float64(limit)))
		} else {
			totalPages = 1
		}

		bookCategoryData := getBookCategoryTableData(bookCategories)
		return Render(c, component.AdminTableSection(bookCategoryData, page, totalPages, "bcategory"))

	case "acategory":
		var articleCategories []database.ACategory
		wg.Add(1)
		go fetchData(origin+"/api/article/acategory", origin, token, page, limit, 0, "", 0, &articleCategories, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(articleCategories) > 0 {
			totalPages = int(math.Ceil(float64(len(articleCategories)) / float64(limit)))
		} else {
			totalPages = 1
		}

		articleCategoryData := getArticleCategoryTableData(articleCategories)
		return Render(c, component.AdminTableSection(articleCategoryData, page, totalPages, "acategory"))

	case "rcategory":
		var resourceCategories []database.RCategory
		wg.Add(1)
		go fetchData(origin+"/api/resource/rcategory", origin, token, page, limit, 0, "", 0, &resourceCategories, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(resourceCategories) > 0 {
			totalPages = int(math.Ceil(float64(len(resourceCategories)) / float64(limit)))
		} else {
			totalPages = 1
		}

		resourceCategoryData := getResourceCategoryTableData(resourceCategories)
		return Render(c, component.AdminTableSection(resourceCategoryData, page, totalPages, "rcategory"))

	case "publisher":
		var publishers []database.Publisher
		wg.Add(1)
		go fetchData(origin+"/api/book/publisher", origin, token, page, limit, 0, "", 0, &publishers, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(publishers) > 0 {
			totalPages = int(math.Ceil(float64(len(publishers)) / float64(limit)))
		} else {
			totalPages = 1
		}

		publisherData := getPublisherTableData(publishers)
		return Render(c, component.AdminTableSection(publisherData, page, totalPages, "publisher"))

	case "version":
		var versions []database.Version
		wg.Add(1)
		go fetchData(origin+"/api/book/version", origin, token, page, limit, 0, "", 0, &versions, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(versions) > 0 {
			totalPages = int(math.Ceil(float64(len(versions)) / float64(limit)))
		} else {
			totalPages = 1
		}

		versionData := getVersionTableData(versions)
		return Render(c, component.AdminTableSection(versionData, page, totalPages, "version"))

	case "letter":
		var letters []database.Letter
		wg.Add(1)
		go fetchData(origin+"/api/book/letter", origin, token, page, limit, 0, "", 0, &letters, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(letters) > 0 {
			totalPages = int(math.Ceil(float64(len(letters)) / float64(limit)))
		} else {
			totalPages = 1
		}

		letterData := getLetterTableData(letters)
		return Render(c, component.AdminTableSection(letterData, page, totalPages, "letter"))

	case "cover":
		var covers []database.Cover
		wg.Add(1)
		go fetchData(origin+"/api/book/cover", origin, token, page, limit, 0, "", 0, &covers, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(covers) > 0 {
			totalPages = int(math.Ceil(float64(len(covers)) / float64(limit)))
		} else {
			totalPages = 1
		}

		coverData := getCoverTableData(covers)
		return Render(c, component.AdminTableSection(coverData, page, totalPages, "cover"))

	case "order":
		var orders []database.Order
		wg.Add(1)
		go fetchData(origin+"/api/order", origin, token, page, limit, 0, "", 0, &orders, &wg, errChan, true, false)
		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		}

		var totalPages int
		if len(orders) > 0 {
			totalPages = int(math.Ceil(float64(orders[0].RecordCount) / float64(limit)))
		} else {
			totalPages = 1
		}

		orderData := getOrderTableData(orders)
		return Render(c, component.AdminTableSection(orderData, page, totalPages, "order"))

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid tab type"})
	}
}

func handleAdminPage(c echo.Context) error {
	cookieName := os.Getenv("COOKIE_NAME")
	cookieValue := os.Getenv("COOKIE_VALUE")
	cookie, err := c.Cookie(cookieName)
	if err != nil || cookie.Value != cookieValue {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	c.Response().Header().Set("X-Robots-Tag", "noindex, nofollow")

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	// Get current active tab (default to books)
	activeTab := c.QueryParam("tab")
	if activeTab == "" {
		activeTab = "book"
	}

	// Get pagination parameters for each tab
	booksPage := 1
	articlesPage := 1
	resourcesPage := 1
	bookCategoriesPage := 1
	articleCategoriesPage := 1
	resourceCategoriesPage := 1
	publishersPage := 1
	versionsPage := 1
	lettersPage := 1
	coversPage := 1
	ordersPage := 1
	limit := 10

	if pageParam := c.QueryParam("book_page"); pageParam != "" {
		var err error
		booksPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing books page parameter: %v", err)
			booksPage = 1
		}
	}
	if pageParam := c.QueryParam("article_page"); pageParam != "" {
		var err error
		articlesPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing articles page parameter: %v", err)
			articlesPage = 1
		}
	}
	if pageParam := c.QueryParam("resource_page"); pageParam != "" {
		var err error
		resourcesPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing resources page parameter: %v", err)
			resourcesPage = 1
		}
	}
	if pageParam := c.QueryParam("book_category_page"); pageParam != "" {
		var err error
		bookCategoriesPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing book categories page parameter: %v", err)
			bookCategoriesPage = 1
		}
	}
	if pageParam := c.QueryParam("article_category_page"); pageParam != "" {
		var err error
		articleCategoriesPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing article categories page parameter: %v", err)
			articleCategoriesPage = 1
		}
	}
	if pageParam := c.QueryParam("resource_category_page"); pageParam != "" {
		var err error
		resourceCategoriesPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing resource categories page parameter: %v", err)
			resourceCategoriesPage = 1
		}
	}
	if pageParam := c.QueryParam("publisher_page"); pageParam != "" {
		var err error
		publishersPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing publishers page parameter: %v", err)
			publishersPage = 1
		}
	}
	if pageParam := c.QueryParam("version_page"); pageParam != "" {
		var err error
		versionsPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing versions page parameter: %v", err)
			versionsPage = 1
		}
	}
	if pageParam := c.QueryParam("letter_page"); pageParam != "" {
		var err error
		lettersPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing letters page parameter: %v", err)
			lettersPage = 1
		}
	}
	if pageParam := c.QueryParam("cover_page"); pageParam != "" {
		var err error
		coversPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing covers page parameter: %v", err)
			coversPage = 1
		}
	}
	if pageParam := c.QueryParam("order_page"); pageParam != "" {
		var err error
		ordersPage, err = strconv.Atoi(pageParam)
		if err != nil {
			log.Printf("Error parsing orders page parameter: %v", err)
			ordersPage = 1
		}
	}
	if limitParam := c.QueryParam("limit"); limitParam != "" {
		var err error
		limit, err = strconv.Atoi(limitParam)
		if err != nil {
			log.Printf("Error parsing limit parameter: %v", err)
			limit = 10
		}
	}

	var wg sync.WaitGroup
	errChan := make(chan error, 11)

	var articles []database.Article
	var resources []database.Resource
	var books []database.Book
	var bookCategories []database.BCategory
	var articleCategories []database.ACategory
	var resourceCategories []database.RCategory
	var publishers []database.Publisher
	var versions []database.Version
	var letters []database.Letter
	var covers []database.Cover
	var orders []database.Order

	wg.Go(func() {
		fetchDataSimple(origin+"/api/article", origin, token, articlesPage, limit, 0, "", 0, &articles, errChan, true, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/resource", origin, token, resourcesPage, limit, 0, "", 0, &resources, errChan, true, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book/admin", origin, token, booksPage, limit, 0, "", 0, &books, errChan, true, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book/bcategory", origin, token, bookCategoriesPage, limit, 0, "", 0, &bookCategories, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/article/acategory", origin, token, articleCategoriesPage, limit, 0, "", 0, &articleCategories, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/resource/rcategory", origin, token, resourceCategoriesPage, limit, 0, "", 0, &resourceCategories, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book/publisher", origin, token, publishersPage, limit, 0, "", 0, &publishers, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book/version", origin, token, versionsPage, limit, 0, "", 0, &versions, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book/letter", origin, token, lettersPage, limit, 0, "", 0, &letters, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/book/cover", origin, token, coversPage, limit, 0, "", 0, &covers, errChan, false, false)
	})

	wg.Go(func() {
		fetchDataSimple(origin+"/api/order", origin, token, ordersPage, limit, 0, "", 0, &orders, errChan, true, false)
	})

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	// Calculate pagination data
	var booksTotalPages, articlesTotalPages, resourcesTotalPages int
	var bookCategoriesTotalPages, articleCategoriesTotalPages, resourceCategoriesTotalPages int
	var publishersTotalPages, versionsTotalPages, lettersTotalPages, coversTotalPages, ordersTotalPages int

	if len(books) > 0 {
		booksTotalPages = int(math.Ceil(float64(books[0].RecordCount) / float64(limit)))
	} else {
		booksTotalPages = 1
	}

	if len(articles) > 0 {
		articlesTotalPages = int(math.Ceil(float64(articles[0].RecordCount) / float64(limit)))
	} else {
		articlesTotalPages = 1
	}

	if len(resources) > 0 {
		resourcesTotalPages = int(math.Ceil(float64(resources[0].RecordCount) / float64(limit)))
	} else {
		resourcesTotalPages = 1
	}

	if len(bookCategories) > 0 {
		bookCategoriesTotalPages = int(math.Ceil(float64(len(bookCategories)) / float64(limit)))
	} else {
		bookCategoriesTotalPages = 1
	}

	if len(articleCategories) > 0 {
		articleCategoriesTotalPages = int(math.Ceil(float64(len(articleCategories)) / float64(limit)))
	} else {
		articleCategoriesTotalPages = 1
	}

	if len(resourceCategories) > 0 {
		resourceCategoriesTotalPages = int(math.Ceil(float64(len(resourceCategories)) / float64(limit)))
	} else {
		resourceCategoriesTotalPages = 1
	}

	if len(publishers) > 0 {
		publishersTotalPages = int(math.Ceil(float64(len(publishers)) / float64(limit)))
	} else {
		publishersTotalPages = 1
	}

	if len(versions) > 0 {
		versionsTotalPages = int(math.Ceil(float64(len(versions)) / float64(limit)))
	} else {
		versionsTotalPages = 1
	}

	if len(letters) > 0 {
		lettersTotalPages = int(math.Ceil(float64(len(letters)) / float64(limit)))
	} else {
		lettersTotalPages = 1
	}

	if len(covers) > 0 {
		coversTotalPages = int(math.Ceil(float64(len(covers)) / float64(limit)))
	} else {
		coversTotalPages = 1
	}

	if len(orders) > 0 {
		ordersTotalPages = int(math.Ceil(float64(orders[0].RecordCount) / float64(limit)))
	} else {
		ordersTotalPages = 1
	}

	var tableData component.CollectiveTableData
	tableData.Books = getBookTableData(books)
	tableData.Articles = getArticleTableData(articles)
	tableData.Resources = getResourceTableData(resources)
	tableData.BookCategories = getBookCategoryTableData(bookCategories)
	tableData.ArticleCategories = getArticleCategoryTableData(articleCategories)
	tableData.ResourceCategories = getResourceCategoryTableData(resourceCategories)
	tableData.Publishers = getPublisherTableData(publishers)
	tableData.Versions = getVersionTableData(versions)
	tableData.Letters = getLetterTableData(letters)
	tableData.Covers = getCoverTableData(covers)
	tableData.Orders = getOrderTableData(orders)
	tableData.BooksPage = booksPage
	tableData.ArticlesPage = articlesPage
	tableData.ResourcesPage = resourcesPage
	tableData.BookCategoriesPage = bookCategoriesPage
	tableData.ArticleCategoriesPage = articleCategoriesPage
	tableData.ResourceCategoriesPage = resourceCategoriesPage
	tableData.PublishersPage = publishersPage
	tableData.VersionsPage = versionsPage
	tableData.LettersPage = lettersPage
	tableData.CoversPage = coversPage
	tableData.OrdersPage = ordersPage
	tableData.BooksTotalPages = booksTotalPages
	tableData.ArticlesTotalPages = articlesTotalPages
	tableData.ResourcesTotalPages = resourcesTotalPages
	tableData.BookCategoriesTotalPages = bookCategoriesTotalPages
	tableData.ArticleCategoriesTotalPages = articleCategoriesTotalPages
	tableData.ResourceCategoriesTotalPages = resourceCategoriesTotalPages
	tableData.PublishersTotalPages = publishersTotalPages
	tableData.VersionsTotalPages = versionsTotalPages
	tableData.LettersTotalPages = lettersTotalPages
	tableData.CoversTotalPages = coversTotalPages
	tableData.OrdersTotalPages = ordersTotalPages
	tableData.ActiveTab = activeTab

	if c.Request().Header.Get("HX-Request") == "true" {
		return renderAdminTableContent(c, tableData, activeTab)
	}

	return Render(c, layout.Admin(tableData))
}

func handleSitemap(c echo.Context) error {
	if os.Getenv("ENVIRONMENT") == "development" {
		file, err := os.Open("public/sitemap.xml")
		if err != nil {
			return echo.NewHTTPError(http.StatusNotFound, "Sitemap not found")
		}
		defer file.Close()
		return c.Stream(http.StatusOK, "application/xml", file)
	}

	file, err := fs.ReadFile(FS, "public/sitemap.xml")
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Sitemap not found")
	}

	return c.Blob(http.StatusOK, "application/xml", file)
}

// Debugging
func handleTestEmail(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	url := fmt.Sprintf("%s/api/test-email", origin)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("failed to send email: %v", err)})
	}
	defer resp.Body.Close()

	return c.JSON(http.StatusOK, map[string]string{"message": "Email sent successfully"})
}

func getBookTableData(books []database.Book) component.TableData {
	bookData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Título"},
			{Value: "Autor"},
			{Value: "Precio"},
			{Value: "Stock"},
			{Value: "Estado"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, book := range books {
		var statusAction string
		if book.IsActive {
			statusAction = "deactivate"
		} else {
			statusAction = "activate"
		}

		bookData.Items = append(bookData.Items, component.TableItem{
			Key:    strconv.Itoa(book.ID),
			Values: []any{book.ID, book.Title, book.Author, book.Price, book.Stock, book.IsActive, "edit", statusAction},
		})
	}

	return bookData
}

func getArticleTableData(articles []database.Article) component.TableData {
	articleData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Título"},
			{Value: "Autor"},
			{Value: "Categoría"},
			{Value: "Fecha"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, article := range articles {
		articleData.Items = append(articleData.Items, component.TableItem{
			Key:    strconv.Itoa(article.ID),
			Values: []any{article.ID, article.Title, article.Author, article.ArticleCategory, article.CreatedAt, "edit", "delete"},
		})
	}

	return articleData
}

func getResourceTableData(resources []database.Resource) component.TableData {
	resourceData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Título"},
			{Value: "Autor"},
			{Value: "Categoría"},
			{Value: "Fecha"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, resource := range resources {
		resourceData.Items = append(resourceData.Items, component.TableItem{
			Key:    strconv.Itoa(resource.ID),
			Values: []any{resource.ID, resource.Title, resource.Author, resource.ResourceCategory, resource.CreatedAt, "edit", "delete"},
		})
	}

	return resourceData
}

func getBookCategoryTableData(categories []database.BCategory) component.TableData {
	categoryData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Categoría"},
			{Value: "Estado"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, category := range categories {
		var statusAction string
		if category.IsActive {
			statusAction = "deactivate"
		} else {
			statusAction = "activate"
		}

		categoryData.Items = append(categoryData.Items, component.TableItem{
			Key:    strconv.Itoa(category.ID),
			Values: []any{category.ID, category.BookCategory, category.IsActive, category.CreatedAt, "edit", statusAction, "delete"},
		})
	}

	return categoryData
}

func getArticleCategoryTableData(categories []database.ACategory) component.TableData {
	categoryData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Categoría"},
			{Value: "Estado"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, category := range categories {
		var statusAction string
		if category.IsActive {
			statusAction = "deactivate"
		} else {
			statusAction = "activate"
		}

		categoryData.Items = append(categoryData.Items, component.TableItem{
			Key:    strconv.Itoa(category.ID),
			Values: []any{category.ID, category.ArticleCategory, category.IsActive, category.CreatedAt, "edit", statusAction, "delete"},
		})
	}

	return categoryData
}

func getResourceCategoryTableData(categories []database.RCategory) component.TableData {
	categoryData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Categoría"},
			{Value: "Estado"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, category := range categories {
		var statusAction string
		if category.IsActive {
			statusAction = "deactivate"
		} else {
			statusAction = "activate"
		}

		categoryData.Items = append(categoryData.Items, component.TableItem{
			Key:    strconv.Itoa(category.ID),
			Values: []any{category.ID, category.ResourceCategory, category.IsActive, category.CreatedAt, "edit", statusAction, "delete"},
		})
	}

	return categoryData
}

func getPublisherTableData(publishers []database.Publisher) component.TableData {
	publisherData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Editorial"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, publisher := range publishers {
		publisherData.Items = append(publisherData.Items, component.TableItem{
			Key:    strconv.Itoa(publisher.ID),
			Values: []any{publisher.ID, publisher.PublisherName, publisher.CreatedAt, "edit", "delete"},
		})
	}

	return publisherData
}

func getVersionTableData(versions []database.Version) component.TableData {
	versionData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Versión Bíblica"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, version := range versions {
		versionData.Items = append(versionData.Items, component.TableItem{
			Key:    strconv.Itoa(version.ID),
			Values: []any{version.ID, version.BibleVersion, version.CreatedAt, "edit", "delete"},
		})
	}

	return versionData
}

func getLetterTableData(letters []database.Letter) component.TableData {
	letterData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Tipo de Letra"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, letter := range letters {
		letterData.Items = append(letterData.Items, component.TableItem{
			Key:    strconv.Itoa(letter.ID),
			Values: []any{letter.ID, letter.LetterType, letter.CreatedAt, "edit", "delete"},
		})
	}

	return letterData
}

func getCoverTableData(covers []database.Cover) component.TableData {
	coverData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Tipo de Cubierta"},
			{Value: "Fecha Creación"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, cover := range covers {
		coverData.Items = append(coverData.Items, component.TableItem{
			Key:    strconv.Itoa(cover.ID),
			Values: []any{cover.ID, cover.CoverType, cover.CreatedAt, "edit", "delete"},
		})
	}

	return coverData
}

func getOrderTableData(orders []database.Order) component.TableData {
	orderData := component.TableData{
		Fields: []component.TableField{
			{Value: "ID"},
			{Value: "Total"},
			{Value: "Estado"},
			{Value: "Dirección"},
			{Value: "Fecha"},
			{Value: "Acciones"},
		},
		Items: []component.TableItem{},
	}

	for _, order := range orders {
		orderData.Items = append(orderData.Items, component.TableItem{
			Key:    strconv.Itoa(order.ID),
			Values: []any{order.ID, order.Total, order.Status, order.Address, order.CreatedAt, "view"},
		})
	}

	return orderData
}

func adminAuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		cookieName := os.Getenv("COOKIE_NAME")
		cookieValue := os.Getenv("COOKIE_VALUE")

		cookie, err := c.Cookie(cookieName)
		if err != nil || cookie.Value != cookieValue {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Unauthorized"})
		}

		return next(c)
	}
}

// Admin action handlers
func handleBookAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "edit":
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/edit/book?id=%s", id))
	case "activate", "deactivate":
		return handleToggleActiveRequest(c, fmt.Sprintf("%s/api/book/%s", origin, id), origin, token, action == "activate")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
	// Delete action
	// return handleDeleteRequest(c, fmt.Sprintf("%s/api/book/%s", origin, id), origin, token)
}

func handleArticleAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "edit":
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/edit/article?id=%s", id))
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/article/%s", origin, id), origin, token)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleResourceAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "edit":
		return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/admin/edit/resource?id=%s", id))
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/resource/%s", origin, id), origin, token)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleDeleteRequest(c echo.Context, url, origin, token string) error {
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete item"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Delete operation failed"})
	}

	// For HTMX requests, return updated table content
	if c.Request().Header.Get("HX-Request") == "true" {
		// Get the current tab from the request path to determine what table to refresh
		path := c.Request().URL.Path
		var tabType string
		if strings.Contains(path, "/book/") {
			tabType = "book"
		} else if strings.Contains(path, "/article/") {
			tabType = "article"
		} else if strings.Contains(path, "/resource/") {
			tabType = "resource"
		}

		// Fetch updated data and render the table section
		return renderUpdatedTableSection(c, tabType, 1) // Default to page 1
	}

	// For non-HTMX requests, return JSON response
	return c.JSON(http.StatusOK, map[string]string{"message": "Item deleted successfully"})
}

func handleToggleActiveRequest(c echo.Context, url, origin, token string, activate bool) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get item"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get item data"})
	}

	var book database.Book
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
	}

	if err := json.Unmarshal(bodyBytes, &book); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse book data"})
	}

	// Update the book with new active status
	bookReq := database.BookRequest{
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description,
		ISBN:        book.ISBN,
		CoverURL:    book.CoverURL,
		Price:       book.Price,
		Stock:       book.Stock,
		SalesCount:  book.SalesCount,
		IsActive:    activate,
		LetterID:    book.LetterID,
		VersionID:   book.VersionID,
		CoverID:     book.CoverID,
		PublisherID: book.PublisherID,
		CategoryID:  book.CategoryID,
	}

	// Send PUT request to update
	bookJSON, err := json.Marshal(bookReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to marshal book data"})
	}

	putReq, err := http.NewRequest("PUT", url, strings.NewReader(string(bookJSON)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create update request"})
	}

	putReq.Header.Add("Authorization", "Bearer "+token)
	putReq.Header.Add("Origin", origin)
	putReq.Header.Add("Content-Type", "application/json")

	putResp, err := client.Do(putReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update item"})
	}
	defer putResp.Body.Close()

	if putResp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Update operation failed"})
	}

	// For HTMX requests, return updated table content
	if c.Request().Header.Get("HX-Request") == "true" {
		// Get the current tab from the request path to determine what table to refresh
		path := c.Request().URL.Path
		var tabType string
		if strings.Contains(path, "/book/") {
			tabType = "book"
		} else if strings.Contains(path, "/article/") {
			tabType = "article"
		} else if strings.Contains(path, "/resource/") {
			tabType = "resource"
		}

		// Fetch updated data and render the table section
		return renderUpdatedTableSection(c, tabType, 1) // Default to page 1
	}

	// For non-HTMX requests, return JSON response
	action := "deactivated"
	if activate {
		action = "activated"
	}
	return c.JSON(http.StatusOK, map[string]string{"message": fmt.Sprintf("Item %s successfully", action)})
}

// Edit page handlers
func handleBookEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 6)

	var book database.Book
	var letters []database.Letter
	var versions []database.Version
	var covers []database.Cover
	var publishers []database.Publisher
	var categories []database.BCategory

	// Fetch book data
	wg.Add(1)
	go fetchData(origin+"/api/book/"+id, origin, token, 0, 0, 0, "", 0, &book, &wg, errChan, false, false)

	// Fetch dropdown data
	wg.Add(1)
	go fetchData(origin+"/api/book/letter", origin, token, 0, 0, 0, "", 0, &letters, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/version", origin, token, 0, 0, 0, "", 0, &versions, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/cover", origin, token, 0, 0, 0, "", 0, &covers, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/publisher", origin, token, 0, 0, 0, "", 0, &publishers, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/bcategory", origin, token, 0, 0, 0, "", 0, &categories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	editData := layout.EditFormData{
		Type: "book",
		Book: &component.BookEditData{
			Book:       book,
			Letters:    letters,
			Versions:   versions,
			Covers:     covers,
			Publishers: publishers,
			Categories: categories,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleArticleEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	var article database.Article
	var categories []database.ACategory

	// Fetch article data
	wg.Add(1)
	go fetchData(origin+"/api/article/"+id, origin, token, 0, 0, 0, "", 0, &article, &wg, errChan, false, false)

	// Fetch categories
	wg.Add(1)
	go fetchData(origin+"/api/article/acategory", origin, token, 0, 0, 0, "", 0, &categories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	editData := layout.EditFormData{
		Type: "article",
		Article: &component.ArticleEditData{
			Article:    article,
			Categories: categories,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleResourceEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 2)

	var resource database.Resource
	var categories []database.RCategory

	// Fetch resource data
	wg.Add(1)
	go fetchData(origin+"/api/resource/"+id, origin, token, 0, 0, 0, "", 0, &resource, &wg, errChan, false, false)

	// Fetch categories
	wg.Add(1)
	go fetchData(origin+"/api/resource/rcategory", origin, token, 0, 0, 0, "", 0, &categories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	editData := layout.EditFormData{
		Type: "resource",
		Resource: &component.ResourceEditData{
			Resource:   resource,
			Categories: categories,
		},
	}

	return Render(c, layout.Edit(editData))
}

// Create page handlers
func handleBookCreatePage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 5)

	var letters []database.Letter
	var versions []database.Version
	var covers []database.Cover
	var publishers []database.Publisher
	var categories []database.BCategory

	// Fetch dropdown data
	wg.Add(1)
	go fetchData(origin+"/api/book/letter", origin, token, 0, 0, 0, "", 0, &letters, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/version", origin, token, 0, 0, 0, "", 0, &versions, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/cover", origin, token, 0, 0, 0, "", 0, &covers, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/publisher", origin, token, 0, 0, 0, "", 0, &publishers, &wg, errChan, false, false)

	wg.Add(1)
	go fetchData(origin+"/api/book/bcategory", origin, token, 0, 0, 0, "", 0, &categories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	createData := layout.CreateFormData{
		Type: "book",
		Book: &component.BookCreateData{
			Letters:    letters,
			Versions:   versions,
			Covers:     covers,
			Publishers: publishers,
			Categories: categories,
		},
	}

	return Render(c, layout.Create(createData))
}

func handleArticleCreatePage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var categories []database.ACategory

	// Fetch categories
	wg.Add(1)
	go fetchData(origin+"/api/article/acategory", origin, token, 0, 0, 0, "", 0, &categories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	createData := layout.CreateFormData{
		Type: "article",
		Article: &component.ArticleCreateData{
			Categories: categories,
		},
	}

	return Render(c, layout.Create(createData))
}

func handleResourceCreatePage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var categories []database.RCategory

	// Fetch categories
	wg.Add(1)
	go fetchData(origin+"/api/resource/rcategory", origin, token, 0, 0, 0, "", 0, &categories, &wg, errChan, false, false)

	wg.Wait()
	close(errChan)

	var errMessages []string
	for err := range errChan {
		if err != nil {
			errMessages = append(errMessages, err.Error())
		}
	}

	if len(errMessages) > 0 {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %s", fmt.Sprintf("%v", errMessages)))
	}

	createData := layout.CreateFormData{
		Type: "resource",
		Resource: &component.ResourceCreateData{
			Categories: categories,
		},
	}

	return Render(c, layout.Create(createData))
}

// Update handlers
func handleBookUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Parse and validate form data
	title := strings.TrimSpace(c.FormValue("title"))
	author := strings.TrimSpace(c.FormValue("author"))
	description := strings.TrimSpace(c.FormValue("description"))
	isbn := strings.TrimSpace(c.FormValue("isbn"))
	coverUrl := strings.TrimSpace(c.FormValue("coverUrl"))

	// Server-side validation
	if title == "" || len(title) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Título es requerido y debe tener máximo 255 caracteres"})
	}
	if author == "" || len(author) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Autor es requerido y debe tener máximo 255 caracteres"})
	}
	if description == "" || len(description) > 300 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción es requerida y debe tener máximo 300 caracteres"})
	}
	if len(isbn) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ISBN debe tener máximo 255 caracteres"})
	}
	if coverUrl == "" || len(coverUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL de portada es requerida y debe tener máximo 255 caracteres"})
	}

	bookRequest := database.BookRequest{
		Title:       title,
		Author:      author,
		Description: description,
		ISBN:        isbn,
		CoverURL:    coverUrl,
		IsActive:    c.FormValue("isActive") == "true",
	}

	// Parse and validate numeric fields
	if price := c.FormValue("price"); price != "" {
		if p, err := strconv.ParseFloat(price, 64); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Precio debe ser un número válido"})
		} else if p < 0 || p > 99999999.99 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Precio debe estar entre 0 y 99,999,999.99"})
		} else {
			bookRequest.Price = p
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Precio es requerido"})
	}

	if stock := c.FormValue("stock"); stock != "" {
		if s, err := strconv.Atoi(stock); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock debe ser un número entero válido"})
		} else if s < 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock no puede ser negativo"})
		} else {
			bookRequest.Stock = s
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock es requerido"})
	}

	if categoryId := c.FormValue("categoryId"); categoryId != "" {
		if cId, err := strconv.Atoi(categoryId); err == nil {
			bookRequest.CategoryID = cId
		}
	}

	if letterId := c.FormValue("letterId"); letterId != "" {
		if lId, err := strconv.Atoi(letterId); err == nil {
			bookRequest.LetterID = lId
		}
	}

	if versionId := c.FormValue("versionId"); versionId != "" {
		if vId, err := strconv.Atoi(versionId); err == nil {
			bookRequest.VersionID = vId
		}
	}

	if coverId := c.FormValue("coverId"); coverId != "" {
		if cId, err := strconv.Atoi(coverId); err == nil {
			bookRequest.CoverID = cId
		}
	}

	if publisherId := c.FormValue("publisherId"); publisherId != "" {
		if pId, err := strconv.Atoi(publisherId); err == nil {
			bookRequest.PublisherID = pId
		}
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/book/%s", origin, id), origin, token, bookRequest)
}

func handleArticleUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Parse and validate form data
	title := strings.TrimSpace(c.FormValue("title"))
	author := strings.TrimSpace(c.FormValue("author"))
	excerpt := strings.TrimSpace(c.FormValue("excerpt"))
	description := strings.TrimSpace(c.FormValue("description"))
	coverUrl := strings.TrimSpace(c.FormValue("coverUrl"))

	// Server-side validation
	if title == "" || len(title) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Título es requerido y debe tener máximo 255 caracteres"})
	}
	if author == "" || len(author) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Autor es requerido y debe tener máximo 255 caracteres"})
	}
	if excerpt == "" || len(excerpt) > 200 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Extracto es requerido y debe tener máximo 200 caracteres"})
	}
	if description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción es requerida"})
	}
	if coverUrl == "" || len(coverUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL de portada es requerida y debe tener máximo 255 caracteres"})
	}

	articleRequest := database.ArticleRequest{
		Title:       title,
		Author:      author,
		Excerpt:     excerpt,
		Description: description,
		CoverURL:    coverUrl,
	}

	if categoryId := c.FormValue("categoryId"); categoryId != "" {
		if cId, err := strconv.Atoi(categoryId); err == nil {
			articleRequest.CategoryID = cId
		}
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/article/%s", origin, id), origin, token, articleRequest)
}

func handleResourceUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	// Parse and validate form data
	title := strings.TrimSpace(c.FormValue("title"))
	author := strings.TrimSpace(c.FormValue("author"))
	description := strings.TrimSpace(c.FormValue("description"))
	coverUrl := strings.TrimSpace(c.FormValue("coverUrl"))
	resourceUrl := strings.TrimSpace(c.FormValue("resourceUrl"))

	// Server-side validation
	if title == "" || len(title) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Título es requerido y debe tener máximo 255 caracteres"})
	}
	if author == "" || len(author) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Autor es requerido y debe tener máximo 255 caracteres"})
	}
	if description == "" || len(description) > 300 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción es requerida y debe tener máximo 300 caracteres"})
	}
	if coverUrl == "" || len(coverUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL de portada es requerida y debe tener máximo 255 caracteres"})
	}
	if resourceUrl == "" || len(resourceUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL del recurso es requerida y debe tener máximo 255 caracteres"})
	}

	resourceRequest := database.ResourceRequest{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		ResourceURL: resourceUrl,
	}

	if categoryId := c.FormValue("categoryId"); categoryId != "" {
		if cId, err := strconv.Atoi(categoryId); err == nil {
			resourceRequest.CategoryID = cId
		}
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/resource/%s", origin, id), origin, token, resourceRequest)
}

func handleUpdateRequest(c echo.Context, url, origin, token string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to marshal data"})
	}

	req, err := http.NewRequest("PUT", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to update: %v", err)})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "Update failed"})
	}

	// Redirect back to admin panel with success message
	return c.Redirect(http.StatusSeeOther, "/admin")
}

// Create handlers
func handleBookCreate(c echo.Context) error {
	// Parse and validate form data (same validation as update)
	title := strings.TrimSpace(c.FormValue("title"))
	author := strings.TrimSpace(c.FormValue("author"))
	description := strings.TrimSpace(c.FormValue("description"))
	isbn := strings.TrimSpace(c.FormValue("isbn"))
	coverUrl := strings.TrimSpace(c.FormValue("coverUrl"))

	// Server-side validation
	if title == "" || len(title) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Título es requerido y debe tener máximo 255 caracteres"})
	}
	if author == "" || len(author) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Autor es requerido y debe tener máximo 255 caracteres"})
	}
	if description == "" || len(description) > 300 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción es requerida y debe tener máximo 300 caracteres"})
	}
	if len(isbn) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ISBN debe tener máximo 255 caracteres"})
	}
	if coverUrl == "" || len(coverUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL de portada es requerida y debe tener máximo 255 caracteres"})
	}

	bookRequest := database.BookRequest{
		Title:       title,
		Author:      author,
		Description: description,
		ISBN:        isbn,
		CoverURL:    coverUrl,
		IsActive:    c.FormValue("isActive") == "true",
	}

	// Parse and validate numeric fields
	if price := c.FormValue("price"); price != "" {
		if p, err := strconv.ParseFloat(price, 64); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Precio debe ser un número válido"})
		} else if p < 0 || p > 99999999.99 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Precio debe estar entre 0 y 99,999,999.99"})
		} else {
			bookRequest.Price = p
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Precio es requerido"})
	}

	if stock := c.FormValue("stock"); stock != "" {
		if s, err := strconv.Atoi(stock); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock debe ser un número entero válido"})
		} else if s < 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock no puede ser negativo"})
		} else {
			bookRequest.Stock = s
		}
	} else {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Stock es requerido"})
	}

	// Parse optional fields
	if categoryId := c.FormValue("categoryId"); categoryId != "" {
		if cId, err := strconv.Atoi(categoryId); err == nil {
			bookRequest.CategoryID = cId
		}
	}

	if letterId := c.FormValue("letterId"); letterId != "" {
		if lId, err := strconv.Atoi(letterId); err == nil {
			bookRequest.LetterID = lId
		}
	}

	if versionId := c.FormValue("versionId"); versionId != "" {
		if vId, err := strconv.Atoi(versionId); err == nil {
			bookRequest.VersionID = vId
		}
	}

	if coverId := c.FormValue("coverId"); coverId != "" {
		if cId, err := strconv.Atoi(coverId); err == nil {
			bookRequest.CoverID = cId
		}
	}

	if publisherId := c.FormValue("publisherId"); publisherId != "" {
		if pId, err := strconv.Atoi(publisherId); err == nil {
			bookRequest.PublisherID = pId
		}
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/book", origin), origin, token, bookRequest)
}

func handleArticleCreate(c echo.Context) error {
	// Parse and validate form data
	title := strings.TrimSpace(c.FormValue("title"))
	author := strings.TrimSpace(c.FormValue("author"))
	excerpt := strings.TrimSpace(c.FormValue("excerpt"))
	description := strings.TrimSpace(c.FormValue("description"))
	coverUrl := strings.TrimSpace(c.FormValue("coverUrl"))

	// Server-side validation
	if title == "" || len(title) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Título es requerido y debe tener máximo 255 caracteres"})
	}
	if author == "" || len(author) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Autor es requerido y debe tener máximo 255 caracteres"})
	}
	if excerpt == "" || len(excerpt) > 200 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Extracto es requerido y debe tener máximo 200 caracteres"})
	}
	if description == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción es requerida"})
	}
	if coverUrl == "" || len(coverUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL de portada es requerida y debe tener máximo 255 caracteres"})
	}

	articleRequest := database.ArticleRequest{
		Title:       title,
		Author:      author,
		Excerpt:     excerpt,
		Description: description,
		CoverURL:    coverUrl,
	}

	if categoryId := c.FormValue("categoryId"); categoryId != "" {
		if cId, err := strconv.Atoi(categoryId); err == nil {
			articleRequest.CategoryID = cId
		}
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/article", origin), origin, token, articleRequest)
}

func handleResourceCreate(c echo.Context) error {
	// Parse and validate form data
	title := strings.TrimSpace(c.FormValue("title"))
	author := strings.TrimSpace(c.FormValue("author"))
	description := strings.TrimSpace(c.FormValue("description"))
	coverUrl := strings.TrimSpace(c.FormValue("coverUrl"))
	resourceUrl := strings.TrimSpace(c.FormValue("resourceUrl"))

	// Server-side validation
	if title == "" || len(title) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Título es requerido y debe tener máximo 255 caracteres"})
	}
	if author == "" || len(author) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Autor es requerido y debe tener máximo 255 caracteres"})
	}
	if description == "" || len(description) > 300 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Descripción es requerida y debe tener máximo 300 caracteres"})
	}
	if coverUrl == "" || len(coverUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL de portada es requerida y debe tener máximo 255 caracteres"})
	}
	if resourceUrl == "" || len(resourceUrl) > 255 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "URL del recurso es requerida y debe tener máximo 255 caracteres"})
	}

	resourceRequest := database.ResourceRequest{
		Title:       title,
		Author:      author,
		Description: description,
		CoverURL:    coverUrl,
		ResourceURL: resourceUrl,
	}

	if categoryId := c.FormValue("categoryId"); categoryId != "" {
		if cId, err := strconv.Atoi(categoryId); err == nil {
			resourceRequest.CategoryID = cId
		}
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/resource", origin), origin, token, resourceRequest)
}

func handleCreateRequest(c echo.Context, url, origin, token string, data interface{}) error {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to marshal data"})
	}

	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create: %v", err)})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(resp.StatusCode, map[string]string{"error": "Create failed"})
	}

	// Redirect back to admin panel with success message
	return c.Redirect(http.StatusSeeOther, "/admin")
}

// Missing handler for admin table content
func handleAdminTableContent(c echo.Context) error {
	cookieName := os.Getenv("COOKIE_NAME")
	cookieValue := os.Getenv("COOKIE_VALUE")
	cookie, err := c.Cookie(cookieName)
	if err != nil || cookie.Value != cookieValue {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	tabType := c.QueryParam("tab")
	if tabType == "" {
		tabType = "book"
	}

	page := 1
	if pageParam := c.QueryParam(tabType + "_page"); pageParam != "" {
		var err error
		page, err = strconv.Atoi(pageParam)
		if err != nil {
			page = 1
		}
	}

	return renderUpdatedTableSection(c, tabType, page)
}

// Book Category Handlers
func handleBookCategoryAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/book/bcategory/%s", origin, id), origin, token)
	case "activate", "deactivate":
		return handleToggleCategoryActiveRequest(c, fmt.Sprintf("%s/api/book/bcategory/%s", origin, id), origin, token, action == "activate", "book")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleBookCategoryEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var category database.BCategory

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/book/bcategory/"+id, origin, token, &category); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "bcategory",
		BookCategory: &component.BookCategoryEditData{
			Category: category,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleBookCategoryCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:         "bcategory",
		BookCategory: &component.BookCategoryCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handleBookCategoryUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	bookCategory := strings.TrimSpace(c.FormValue("book_category"))
	if bookCategory == "" || len(bookCategory) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de categoría es requerido y debe tener máximo 100 caracteres"})
	}

	categoryRequest := database.BCategoryRequest{
		BookCategory: bookCategory,
		IsActive:     c.FormValue("is_active") == "on",
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/book/bcategory/%s", origin, id), origin, token, categoryRequest)
}

func handleBookCategoryCreate(c echo.Context) error {
	bookCategory := strings.TrimSpace(c.FormValue("book_category"))
	if bookCategory == "" || len(bookCategory) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de categoría es requerido y debe tener máximo 100 caracteres"})
	}

	categoryRequest := database.BCategoryRequest{
		BookCategory: bookCategory,
		IsActive:     c.FormValue("is_active") == "on",
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/book/bcategory", origin), origin, token, categoryRequest)
}

// Article Category Handlers
func handleArticleCategoryAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/article/acategory/%s", origin, id), origin, token)
	case "activate", "deactivate":
		return handleToggleCategoryActiveRequest(c, fmt.Sprintf("%s/api/article/acategory/%s", origin, id), origin, token, action == "activate", "article")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleArticleCategoryEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var category database.ACategory

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/article/acategory/"+id, origin, token, &category); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "acategory",
		ArticleCategory: &component.ArticleCategoryEditData{
			Category: category,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleArticleCategoryCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:            "acategory",
		ArticleCategory: &component.ArticleCategoryCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handleArticleCategoryUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	articleCategory := strings.TrimSpace(c.FormValue("article_category"))
	if articleCategory == "" || len(articleCategory) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de categoría es requerido y debe tener máximo 100 caracteres"})
	}

	categoryRequest := database.ACategoryRequest{
		ArticleCategory: articleCategory,
		IsActive:        c.FormValue("is_active") == "on",
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/article/acategory/%s", origin, id), origin, token, categoryRequest)
}

func handleArticleCategoryCreate(c echo.Context) error {
	articleCategory := strings.TrimSpace(c.FormValue("article_category"))
	if articleCategory == "" || len(articleCategory) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de categoría es requerido y debe tener máximo 100 caracteres"})
	}

	categoryRequest := database.ACategoryRequest{
		ArticleCategory: articleCategory,
		IsActive:        c.FormValue("is_active") == "on",
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/article/acategory", origin), origin, token, categoryRequest)
}

// Resource Category Handlers
func handleResourceCategoryAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/resource/rcategory/%s", origin, id), origin, token)
	case "activate", "deactivate":
		return handleToggleCategoryActiveRequest(c, fmt.Sprintf("%s/api/resource/rcategory/%s", origin, id), origin, token, action == "activate", "resource")
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleResourceCategoryEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var category database.RCategory

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/resource/rcategory/"+id, origin, token, &category); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "rcategory",
		ResourceCategory: &component.ResourceCategoryEditData{
			Category: category,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleResourceCategoryCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:             "rcategory",
		ResourceCategory: &component.ResourceCategoryCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handleResourceCategoryUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	resourceCategory := strings.TrimSpace(c.FormValue("resource_category"))
	if resourceCategory == "" || len(resourceCategory) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de categoría es requerido y debe tener máximo 100 caracteres"})
	}

	categoryRequest := database.RCategoryRequest{
		ResourceCategory: resourceCategory,
		IsActive:         c.FormValue("is_active") == "on",
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/resource/rcategory/%s", origin, id), origin, token, categoryRequest)
}

func handleResourceCategoryCreate(c echo.Context) error {
	resourceCategory := strings.TrimSpace(c.FormValue("resource_category"))
	if resourceCategory == "" || len(resourceCategory) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de categoría es requerido y debe tener máximo 100 caracteres"})
	}

	categoryRequest := database.RCategoryRequest{
		ResourceCategory: resourceCategory,
		IsActive:         c.FormValue("is_active") == "on",
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/resource/rcategory", origin), origin, token, categoryRequest)
}

// Publisher Handlers
func handlePublisherAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/book/publisher/%s", origin, id), origin, token)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handlePublisherEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var publisher database.Publisher

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/book/publisher/"+id, origin, token, &publisher); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "publisher",
		Publisher: &component.PublisherEditData{
			Publisher: publisher,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handlePublisherCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:      "publisher",
		Publisher: &component.PublisherCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handlePublisherUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	publisherName := strings.TrimSpace(c.FormValue("publisher_name"))
	if publisherName == "" || len(publisherName) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de editorial es requerido y debe tener máximo 100 caracteres"})
	}

	publisherRequest := database.PublisherRequest{
		PublisherName: publisherName,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/book/publisher/%s", origin, id), origin, token, publisherRequest)
}

func handlePublisherCreate(c echo.Context) error {
	publisherName := strings.TrimSpace(c.FormValue("publisher_name"))
	if publisherName == "" || len(publisherName) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de editorial es requerido y debe tener máximo 100 caracteres"})
	}

	publisherRequest := database.PublisherRequest{
		PublisherName: publisherName,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/book/publisher", origin), origin, token, publisherRequest)
}

// Version Handlers
func handleVersionAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/book/version/%s", origin, id), origin, token)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleVersionEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var version database.Version

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/book/version/"+id, origin, token, &version); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "version",
		Version: &component.VersionEditData{
			Version: version,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleVersionCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:    "version",
		Version: &component.VersionCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handleVersionUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	bibleVersion := strings.TrimSpace(c.FormValue("bible_version"))
	if bibleVersion == "" || len(bibleVersion) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de versión es requerido y debe tener máximo 100 caracteres"})
	}

	versionRequest := database.VersionRequest{
		BibleVersion: bibleVersion,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/book/version/%s", origin, id), origin, token, versionRequest)
}

func handleVersionCreate(c echo.Context) error {
	bibleVersion := strings.TrimSpace(c.FormValue("bible_version"))
	if bibleVersion == "" || len(bibleVersion) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Nombre de versión es requerido y debe tener máximo 100 caracteres"})
	}

	versionRequest := database.VersionRequest{
		BibleVersion: bibleVersion,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/book/version", origin), origin, token, versionRequest)
}

// Letter Handlers
func handleLetterAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/book/letter/%s", origin, id), origin, token)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleLetterEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var letter database.Letter

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/book/letter/"+id, origin, token, &letter); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "letter",
		Letter: &component.LetterEditData{
			Letter: letter,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleLetterCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:   "letter",
		Letter: &component.LetterCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handleLetterUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	letterType := strings.TrimSpace(c.FormValue("letter_type"))
	if letterType == "" || len(letterType) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tipo de letra es requerido y debe tener máximo 100 caracteres"})
	}

	letterRequest := database.LetterRequest{
		LetterType: letterType,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/book/letter/%s", origin, id), origin, token, letterRequest)
}

func handleLetterCreate(c echo.Context) error {
	letterType := strings.TrimSpace(c.FormValue("letter_type"))
	if letterType == "" || len(letterType) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tipo de letra es requerido y debe tener máximo 100 caracteres"})
	}

	letterRequest := database.LetterRequest{
		LetterType: letterType,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/book/letter", origin), origin, token, letterRequest)
}

// Cover Handlers
func handleCoverAction(c echo.Context) error {
	action := c.Param("action")
	id := c.QueryParam("id")

	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	switch action {
	case "delete":
		return handleDeleteRequest(c, fmt.Sprintf("%s/api/book/cover/%s", origin, id), origin, token)
	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid action"})
	}
}

func handleCoverEditPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var wg sync.WaitGroup
	errChan := make(chan error, 1)

	var cover database.Cover

	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := fetchSingleData(origin+"/api/book/cover/"+id, origin, token, &cover); err != nil {
			errChan <- err
		}
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch data: %v", err))
		}
	}

	editData := layout.EditFormData{
		Type: "cover",
		Cover: &component.CoverEditData{
			Cover: cover,
		},
	}

	return Render(c, layout.Edit(editData))
}

func handleCoverCreatePage(c echo.Context) error {
	createData := layout.CreateFormData{
		Type:  "cover",
		Cover: &component.CoverCreateData{},
	}

	return Render(c, layout.Create(createData))
}

func handleCoverUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "ID is required"})
	}

	coverType := strings.TrimSpace(c.FormValue("cover_type"))
	if coverType == "" || len(coverType) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tipo de cubierta es requerido y debe tener máximo 100 caracteres"})
	}

	coverRequest := database.CoverRequest{
		CoverType: coverType,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleUpdateRequest(c, fmt.Sprintf("%s/api/book/cover/%s", origin, id), origin, token, coverRequest)
}

func handleCoverCreate(c echo.Context) error {
	coverType := strings.TrimSpace(c.FormValue("cover_type"))
	if coverType == "" || len(coverType) > 100 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Tipo de cubierta es requerido y debe tener máximo 100 caracteres"})
	}

	coverRequest := database.CoverRequest{
		CoverType: coverType,
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	return handleCreateRequest(c, fmt.Sprintf("%s/api/book/cover", origin), origin, token, coverRequest)
}

// Order Handler
func handleOrderPage(c echo.Context) error {
	id := c.QueryParam("id")
	if id == "" {
		return c.Redirect(http.StatusSeeOther, "/admin?tab=order")
	}

	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	var account database.Account

	// Fetch order with books from admin endpoint
	var orderResponse struct {
		Order      database.Order       `json:"order"`
		BookOrders []database.BookOrder `json:"bookOrders"`
		Books      []database.Book      `json:"books"`
	}

	if err := fetchSingleData(origin+"/api/order/admin/"+id, origin, token, &orderResponse); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch order: %v", err))
	}

	// Now fetch account details using the order.AccountID (sequential)
	if err := fetchSingleData(origin+"/api/account/"+orderResponse.Order.AccountID.String(), origin, token, &account); err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch account: %v", err))
	}

	// Create order data
	orderData := layout.OrderData{
		Order:      orderResponse.Order,
		Account:    account,
		OrderBooks: orderResponse.BookOrders,
		Books:      orderResponse.Books,
	}

	return layout.Order(orderData).Render(c.Request().Context(), c.Response().Writer)
}

func handleOrderStatusUpdate(c echo.Context) error {
	id := c.QueryParam("id")
	status := c.FormValue("status")

	if id == "" {
		return c.HTML(http.StatusBadRequest, "ID es requerido")
	}

	if status == "" {
		return c.HTML(http.StatusBadRequest, "Estado es requerido")
	}

	// Validate status value
	validStatuses := []string{"processing", "delivered", "returned"}
	if !slices.Contains(validStatuses, status) {
		return c.HTML(http.StatusBadRequest, "Valor de estado inválido")
	}

	// Call API to update status
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	// Create form data for the API call
	formData := url.Values{}
	formData.Set("status", status)

	req, err := http.NewRequest("PUT", origin+"/api/order/status/"+id, strings.NewReader(formData.Encode()))
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Error al crear la solicitud")
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, "Error al actualizar el estado del pedido")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.HTML(http.StatusInternalServerError, "Error al actualizar el estado")
	}

	return c.HTML(http.StatusOK, "Estado actualizado")
}

// Helper function for fetching single data items
func fetchSingleData(url, origin, token string, result interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(body, result)
}

// Category-specific toggle active function
func handleToggleCategoryActiveRequest(c echo.Context, url, origin, token string, activate bool, categoryType string) error {
	// First, get the current category data
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create request"})
	}

	req.Header.Add("Authorization", "Bearer "+token)
	req.Header.Add("Origin", origin)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get category"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get category data"})
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to read response"})
	}

	// Handle different category types
	switch categoryType {
	case "book":
		var category database.BCategory
		if err := json.Unmarshal(bodyBytes, &category); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse category data"})
		}

		categoryReq := database.BCategoryRequest{
			BookCategory: category.BookCategory,
			IsActive:     activate,
		}

		return handleUpdateRequest(c, url, origin, token, categoryReq)

	case "article":
		var category database.ACategory
		if err := json.Unmarshal(bodyBytes, &category); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse category data"})
		}

		categoryReq := database.ACategoryRequest{
			ArticleCategory: category.ArticleCategory,
			IsActive:        activate,
		}

		return handleUpdateRequest(c, url, origin, token, categoryReq)

	case "resource":
		var category database.RCategory
		if err := json.Unmarshal(bodyBytes, &category); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse category data"})
		}

		categoryReq := database.RCategoryRequest{
			ResourceCategory: category.ResourceCategory,
			IsActive:         activate,
		}

		return handleUpdateRequest(c, url, origin, token, categoryReq)

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid category type"})
	}
}
