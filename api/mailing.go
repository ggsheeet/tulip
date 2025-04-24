package api

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ggsheet/tulip/app"
	"github.com/ggsheet/tulip/internal/database"
	"github.com/resend/resend-go/v2"
)

func loadTemplate(filePath string) (string, error) {
	isDev := os.Getenv("ENVIRONMENT") == "development"

	if isDev {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", err
		}
		return string(content), nil
	} else {
		file, err := app.FS.Open(filePath)
		if err != nil {
			return "", err
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			return "", err
		}

		return string(content), nil
	}
}

func generateOrderSummary(cart []database.OrderBook) string {
	var sb strings.Builder
	for _, book := range cart {
		sb.WriteString(fmt.Sprintf(`
			<tr>
				<td class="left"><img src="%s" alt="Portada" style="width: 50px; height: auto;"></td>
				<td class="left">%s</td>
				<td class="left">%s</td>
				<td class="left">$%s</td>
			</tr>
		`, book.PictureURL, book.Title, book.Quantity, book.UnitPrice))
	}
	return sb.String()
}

func (s *ResendServer) HandlePurchaseConfirmation(emailData app.EmailData) (string, error) {
	templatePath := "public/html/order_confirmation.html"
	emailTemplate, err := loadTemplate(templatePath)
	if err != nil {
		return "", fmt.Errorf("error loading email template: %v", err)
	}

	// Debugging
	// log.Printf("Loaded customer template (len %d): %s\n", len(emailTemplate), emailTemplate)

	orderSummary := generateOrderSummary(emailData.Cart)

	emailContent := strings.NewReplacer(
		"{{first_name}}", emailData.FirstName,
		"{{order_number}}", emailData.OrderNumber,
		"{{sub_total}}", fmt.Sprintf("%.2f", emailData.SubTotal),
		"{{shipping_cost}}", fmt.Sprintf("%.2f", emailData.Shipping),
		"{{total}}", fmt.Sprintf("%.2f", emailData.Total),
		"{{order_summary}}", orderSummary,
	).Replace(emailTemplate)

	params := &resend.SendEmailRequest{
		From:    "Publicaciones Tulip <contacto@publicacionestulip.org>",
		To:      []string{emailData.Email},
		Subject: "Gracias por tu orden en TULIP!",
		Html:    emailContent,
	}

	sentCustomer, err := s.msg.Emails.Send(params)
	if err != nil {
		return "", fmt.Errorf("error sending confirmation email to customer: %s", err)
	}

	sentAdmin, err := s.sendAdminEmail(emailData)
	if err != nil {
		return "", err
	}

	sentEmailIds := fmt.Sprintf("customer: %s, admin: %s", sentCustomer.Id, sentAdmin)

	return sentEmailIds, nil
}

func (s *ResendServer) sendAdminEmail(emailData app.EmailData) (string, error) {
	templatePath := "public/html/admin_order_notification.html"
	emailTemplate, err := loadTemplate(templatePath)
	if err != nil {
		return "", fmt.Errorf("error loading email template: %v", err)
	}
	// Debugging
	// log.Printf("Loaded admin template (len %d): %s\n", len(emailTemplate), emailTemplate)

	orderSummary := generateOrderSummary(emailData.Cart)

	emailContent := strings.NewReplacer(
		"{{first_name}}", emailData.FirstName,
		"{{last_name}}", emailData.LastName,
		"{{customer_email}}", emailData.Email,
		"{{order_number}}", emailData.OrderNumber,
		"{{sub_total}}", fmt.Sprintf("%.2f", emailData.SubTotal),
		"{{shipping_cost}}", fmt.Sprintf("%.2f", emailData.Shipping),
		"{{total}}", fmt.Sprintf("%.2f", emailData.Total),
		"{{order_summary}}", orderSummary,
	).Replace(emailTemplate)

	params := &resend.SendEmailRequest{
		From:    "contacto@publicacionestulip.org",
		To:      []string{"contacto@publicacionestulip.org"},
		Subject: "¡Tienes una compra nueva!",
		Html:    emailContent,
	}

	// Debugging
	// fmt.Println("Sending email to customer with: ", params)

	sent, err := s.msg.Emails.Send(params)
	if err != nil {
		return "", fmt.Errorf("error sending confirmation email to admin: %s", err)
	}

	// Debugging
	// fmt.Println("Email response: ", sent)
	// fmt.Println("Email error: ", err)

	return sent.Id, nil
}

// Debugging
// func (s *ResendServer) handleTestEmail(c echo.Context) error {
// 	params := &resend.SendEmailRequest{
// 		From:    "Publicaciones Tulip <contacto@publicacionestulip.org>",
// 		To:      []string{"gigisheet@gmail.com"},
// 		Subject: "Test Email",
// 		Html:    "<p>Hello from Resend!</p>",
// 	}

// 	sent, err := s.msg.Emails.Send(params)
// 	if err != nil {
// 		fmt.Println(err)
// 		return err
// 	}
// 	fmt.Println(sent.Id)
// 	return nil
// }
