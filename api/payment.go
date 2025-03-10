package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/ggsheet/tulip/app"
	"github.com/ggsheet/tulip/internal/database"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/mercadopago/sdk-go/pkg/preference"
)

func (s *MPServer) handleGeneratePreference(c echo.Context) error {
	origin := os.Getenv("AUTH_ORIGIN")

	var form app.Form
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	items, _, _, err := parseForm(form)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	shipping, err := strconv.ParseFloat(os.Getenv("SHIPPING_COST"), 64)
	if err != nil {
		return err
	}

	backurls := &preference.BackURLsRequest{
		Success: origin + "/processed",
		Failure: origin + "/processed",
	}

	phone := &preference.PhoneRequest{
		AreaCode: "52",
		Number:   form.Phone,
	}

	payer := preference.PayerRequest{
		Name:  fmt.Sprintf("%v %v", form.FirstName, form.LastName),
		Email: fmt.Sprintf("%v", form.Email),
		Phone: phone,
	}

	shipments := preference.ShipmentsRequest{
		Mode: "not_specified",
		Cost: shipping,
	}

	excludedPaymentMethods := []preference.ExcludedPaymentMethodRequest{
		{ID: "redcompra"},
	}

	excludedPaymentTypes := []preference.ExcludedPaymentTypeRequest{
		{ID: "bank_transfer"},
		{ID: "ticket"},
		{ID: "atm"},
	}

	paymentMethods := preference.PaymentMethodsRequest{
		ExcludedPaymentMethods: excludedPaymentMethods,
		ExcludedPaymentTypes:   excludedPaymentTypes,
		Installments:           3,
	}

	client := preference.NewClient(s.cfg)
	request := preference.Request{
		Items:          items,
		Payer:          &payer,
		Shipments:      &shipments,
		BackURLs:       backurls,
		PaymentMethods: &paymentMethods,
		BinaryMode:     true,
	}
	resource, err := client.Create(context.Background(), request)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create preference"})
	}

	return c.JSON(http.StatusOK, map[string]string{"redirect_url": resource.InitPoint})
}

func (s *MPServer) handleConfirmedTransaction(c echo.Context) error {
	var paymentId int
	var paymentIdString string
	if paymentIdParam := c.QueryParam("payment_id"); paymentIdParam != "" {
		var err error
		paymentIdString = paymentIdParam
		paymentId, err = strconv.Atoi(paymentIdParam)
		if err != nil || paymentId <= 0 {
			return c.JSON(http.StatusBadRequest, "Invalid paymentId")
		}
	}

	mpAccessToken := os.Getenv("MP_ACCESS_TOKEN")
	mpURL := fmt.Sprintf("https://api.mercadopago.com/v1/payments/%v", paymentId)
	req, err := http.NewRequest("GET", mpURL, nil)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to create request: %v", err))
	}

	req.Header.Add("Authorization", "Bearer "+mpAccessToken)
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to fetch payment: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.String(resp.StatusCode, fmt.Sprintf("Failed to fetch payment: %s", resp.Status))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Failed to read response body")
	}

	var payment app.PaymentData
	if err := json.Unmarshal(bodyBytes, &payment); err != nil {
		return c.String(http.StatusInternalServerError, "Failed to parse payment response")
	}

	if payment.Status != "approved" {
		return c.String(http.StatusPaymentRequired, "Payment not approved")
	}

	orderNotFound := 0
	existingOrder, err := s.o.GetOrderByPaymentId(paymentId)
	if err != nil {
		return err
	}
	if existingOrder != orderNotFound {
		return c.JSON(http.StatusConflict, echo.Map{"error": "Order already exists"})
	}

	var form app.Form
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	_, orderBooks, total, err := parseForm(form)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	accReq := database.AccountRequest{
		FirstName: form.FirstName,
		LastName:  form.LastName,
		Email:     form.Email,
		Phone:     form.Phone,
	}

	var accId uuid.UUID
	accId, err = s.getAccountId(accReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to process account: %v", err)})
	}

	orderReq := database.OrderRequest{
		Address:     fmt.Sprintf("%v, %v, %v, %v, %v, %v", form.Street, form.Neighborhood, form.City, form.ZipCode, form.State, form.Country),
		Total:       total,
		PaymentID:   paymentId,
		IsFulfilled: false,
		Status:      "processing",
		AccountID:   accId,
		OrderBooks:  orderBooks,
	}

	err = s.createOrder(orderReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create order: %v", err)})
	}

	shipping, err := strconv.ParseFloat(os.Getenv("SHIPPING_COST"), 64)
	if err != nil {
		return fmt.Errorf("failed to parse shipping cost: %v", err)
	}

	emailData := app.EmailData{
		FirstName:   accReq.FirstName,
		Email:       accReq.Email,
		OrderNumber: paymentIdString,
		SubTotal:    total - shipping,
		Shipping:    shipping,
		Total:       total,
		Cart:        orderBooks,
	}

	emailId, err := s.m.HandlePurchaseConfirmation(emailData)
	if err != nil {
		return err
	}
	fmt.Println(fmt.Printf("email sent: %v", emailId))

	return c.JSON(http.StatusCreated, echo.Map{"message": "Order created successfully"})
}

func parseForm(form app.Form) ([]preference.ItemRequest, []database.OrderBook, float64, error) {
	var items []preference.ItemRequest
	var orderBooks []database.OrderBook
	total := 0.00

	shipping, err := strconv.ParseFloat(os.Getenv("SHIPPING_COST"), 64)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("failed to parse shipping cost: %v", err)
	}

	for _, item := range form.Cart {
		items = append(items, preference.ItemRequest{
			ID:          fmt.Sprintf("%d", item.Id),
			Title:       item.Title,
			Description: item.Description,
			Quantity:    item.Quantity,
			UnitPrice:   float64(item.UnitPrice),
			CategoryID:  item.CategoryID,
			CurrencyID:  "MXN",
			PictureURL:  item.PictureURL,
		})

		orderBooks = append(orderBooks, database.OrderBook{
			BookID:      item.Id,
			Title:       item.Title,
			Description: item.Description,
			Quantity:    item.Quantity,
			Price:       float64(item.UnitPrice),
			CoverURL:    item.PictureURL,
			BCategory:   item.CategoryID,
		})

		total += float64(item.Quantity) * float64(item.UnitPrice)
	}

	return items, orderBooks, total + shipping, nil
}

func (s *MPServer) getAccountId(accReq database.AccountRequest) (uuid.UUID, error) {
	existingAcc, err := s.a.GetAccountByEmail(accReq.Email)
	if err != nil {
		account := database.NewAccount(accReq.FirstName, accReq.LastName, accReq.Email, accReq.Phone)
		accId, err := s.a.CreateAccount(account)
		if err != nil {
			return uuid.Nil, err
		}
		return accId, nil
	}

	updatedAcc := existingAcc
	updatedAcc.FirstName = accReq.FirstName
	updatedAcc.LastName = accReq.LastName
	updatedAcc.Phone = accReq.Phone

	updateErr := s.a.UpdateAccount(existingAcc.ID, updatedAcc)
	if updateErr != nil {
		return uuid.Nil, updateErr
	}

	return existingAcc.ID, nil
}

func (s *MPServer) createOrder(orderReq database.OrderRequest) error {
	order := database.NewOrder(orderReq.Address, orderReq.Total, orderReq.PaymentID, orderReq.IsFulfilled, orderReq.Status, orderReq.AccountID)
	bookOrders := orderReq.OrderBooks

	orderId, err := s.o.CreateOrder(order)
	if err != nil {
		return err
	}

	for _, bookOrder := range bookOrders {
		if err := s.o.CreateBookOrder(&bookOrder, orderId); err != nil {
			return err
		}

		if err := s.b.UpdateBookStock(bookOrder.BookID, bookOrder.Quantity); err != nil {
			return err
		}
	}

	return nil
}
