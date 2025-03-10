package app

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/ggsheet/tulip/internal/database"
	"github.com/ggsheet/tulip/template/component"
	"github.com/ggsheet/tulip/template/layout"
	"github.com/labstack/echo/v4"
)

func fetchData(url string, origin string, token string, page int, limit int, category int, order string, itemId int, result interface{}, wg *sync.WaitGroup, errChan chan<- error, paginate bool, filter bool) {
	defer wg.Done()

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

func handleIndexPage(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")
	token := os.Getenv("AUTH_TOKEN")

	log.Printf("origin: %v", origin)
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

	wg.Add(1)
	go fetchData(origin+"/api/article", origin, token, page, limit, 0, "", 0, &articles, &wg, errChan, true, false)

	wg.Add(1)
	go fetchData(origin+"/api/resource", origin, token, page, limit, 0, "", 0, &resources, &wg, errChan, true, false)

	wg.Add(1)
	go fetchData(origin+"/api/book", origin, token, page, limit, 0, "", 0, &books, &wg, errChan, true, false)

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

func handleProcesedTransaction(c echo.Context) error {
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

func handleLoginPage(c echo.Context) error {
	html := `
		<form hx-post="/login" hx-target="#error" hx-swap="innerHTML">
			<input type="email" name="email" placeholder="Email" required>
			<input type="password" name="password" placeholder="Password" required>
			<button type="submit">Login</button>
		</form>
		<div id="error"></div>
	`
	return c.HTML(http.StatusOK, html)
}

func handleLoginAuth(c echo.Context) error {
	adminEmail := os.Getenv("ADMIN_EMAIL_DEV")
	adminPassword := os.Getenv("ADMIN_PASSWORD_DEV")

	inputEmail := c.FormValue("adminEmail")
	inputPassword := c.FormValue("adminPassword")

	if inputEmail == adminEmail && inputPassword == adminPassword {
		cookie := new(http.Cookie)
		cookie.Name = "authSession"
		cookie.Value = "TulipAdminAuthenticated"
		cookie.Path = "/"
		cookie.Expires = time.Now().Add(24 * time.Hour)
		c.SetCookie(cookie)
		return c.Redirect(http.StatusSeeOther, "/admin")
	}
	return c.HTML(http.StatusUnauthorized, "<p style='color:red;'>Invalid credentials</p>")
}

func handleLogoutAuth(c echo.Context) error {
	cookie := new(http.Cookie)
	cookie.Name = "authSession"
	cookie.Value = ""
	cookie.Path = "/"
	cookie.MaxAge = -1
	c.SetCookie(cookie)
	return c.Redirect(http.StatusSeeOther, "/login")
}

func handleAdminPage(c echo.Context) error {
	cookie, err := c.Cookie("authSession")
	if err != nil || cookie.Value != "TulipAdminAuthenticated" {
		return c.Redirect(http.StatusSeeOther, "/login")
	}

	return c.HTML(http.StatusOK, "<h1>Welcome to Admin Dashboard</h1><a href='/logout'>Logout</a>")
}
