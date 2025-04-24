package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
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
	environment := os.Getenv("ENVIRONMENT")
	origin := os.Getenv("AUTH_ORIGIN")

	var form app.Form
	if err := c.Bind(&form); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request format"})
	}

	items := parseForm(form)
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

	var payer preference.PayerRequest
	if environment == "development" || environment == "docker" {
		payer = preference.PayerRequest{
			Name:  fmt.Sprintf("%v %v", form.FirstName, form.LastName),
			Email: "test_user_123456@testuser.com",
			Phone: phone,
		}
	} else {
		payer = preference.PayerRequest{
			Name:  fmt.Sprintf("%v %v", form.FirstName, form.LastName),
			Email: form.Email,
			Phone: phone,
		}
	}

	receiverAddress := preference.ReceiverAddressRequest{
		ZipCode:      form.ZipCode,
		StreetName:   form.Street,
		StreetNumber: form.StreetNumber,
		CountryName:  form.Country,
		StateName:    form.State,
		CityName:     form.City,
	}

	shipments := preference.ShipmentsRequest{
		Mode:            "not_specified",
		Cost:            shipping,
		ReceiverAddress: &receiverAddress,
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

	var notificationUrl string
	if environment == "development" || environment == "docker" {
		notificationUrl = "https://4vj2j6tv-8080.usw3.devtunnels.ms/notification"
	} else {
		notificationUrl = origin + "/notification"
	}

	client := preference.NewClient(s.cfg)
	request := preference.Request{
		Items:             items,
		Payer:             &payer,
		Shipments:         &shipments,
		BackURLs:          backurls,
		PaymentMethods:    &paymentMethods,
		BinaryMode:        true,
		ExternalReference: accId.String(),
		NotificationURL:   notificationUrl,
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
		bodyBytes, _ := io.ReadAll(resp.Body)
		return c.JSON(resp.StatusCode, map[string]string{
			"error": fmt.Sprintf("Failed to fetch payment: %s, details: %s", resp.Status, string(bodyBytes)),
		})
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.String(http.StatusInternalServerError, fmt.Sprintf("Failed to read response body: %v", err))
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

	if payment.ExternalReference == nil || *payment.ExternalReference == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Missing or null external reference"})
	}
	accId, err := uuid.Parse(*payment.ExternalReference)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid external reference UUID"})
	}

	orderReq := database.OrderRequest{
		Address:     fmt.Sprintf("%v, %v, %v, %v, %v, %v", payment.AdditionalInfo.ShipmentsInfo.ReceiverAddress.StreetName, payment.AdditionalInfo.ShipmentsInfo.ReceiverAddress.StreetNumber, payment.AdditionalInfo.ShipmentsInfo.ReceiverAddress.CityName, payment.AdditionalInfo.ShipmentsInfo.ReceiverAddress.ZipCode, payment.AdditionalInfo.ShipmentsInfo.ReceiverAddress.StateName, "México"),
		Total:       float64(payment.TransactionDetails.TotalPaidAmount),
		PaymentID:   payment.ID,
		IsFulfilled: false,
		Status:      "processing",
		AccountID:   accId,
		OrderBooks:  payment.AdditionalInfo.ItemsInfo,
	}

	err = s.createOrder(orderReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to create order: %v", err)})
	}

	acc, err := s.a.GetAccountById(accId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Failed to get account email: %v", err)})
	}

	emailData := app.EmailData{
		FirstName:   acc.FirstName,
		LastName:    acc.LastName,
		Email:       acc.Email,
		OrderNumber: paymentIdString,
		SubTotal:    float64(payment.TransactionAmount),
		Shipping:    float64(payment.ShippingAmount),
		Total:       float64(payment.TransactionDetails.TotalPaidAmount),
		Cart:        payment.AdditionalInfo.ItemsInfo,
	}

	emailId, err := s.m.HandlePurchaseConfirmation(emailData)
	if err != nil {
		log.Printf("Email error: %v", err)
		c.JSON(http.StatusInternalServerError, fmt.Errorf("failed to send confirmation email: %v", err))
	}

	log.Printf("email sent: ", emailId)
	os.Stdout.Sync()

	return c.JSON(http.StatusOK, echo.Map{"message": "Order created successfully"})
}

func parseForm(form app.Form) []preference.ItemRequest {
	var items []preference.ItemRequest

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
	}

	return items
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
		bookId, err := strconv.Atoi(bookOrder.ID)
		if err != nil || bookId <= 0 {
			return err
		}

		bookQuantity, err := strconv.Atoi(bookOrder.Quantity)
		if err != nil || bookQuantity <= 0 {
			return err
		}

		if err := s.o.CreateBookOrder(bookQuantity, bookId, orderId); err != nil {
			return err
		}

		if err := s.b.UpdateBookStock(bookId, bookQuantity); err != nil {
			return err
		}
	}

	return nil
}
