package api

import (
	"fmt"
	"os"
	"strings"

	"github.com/ggsheet/tulip/app"
	"github.com/resend/resend-go/v2"
)

func loadTemplate(filePath string) (string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

func (s *ResendServer) HandlePurchaseConfirmation(emailData app.EmailData) (string, error) {
	templatePath := "template/email/order_confirmation.html"
	emailTemplate, err := loadTemplate(templatePath)
	if err != nil {
		return "", fmt.Errorf("error loading email template: %v", err)
	}

	orderSummary := ""
	for _, book := range emailData.Cart {
		orderSummary += fmt.Sprintf(`
			<tr>
				<td class="left"><img src="%s" alt="Portada" style="width: 50px; height: auto;"></td>
				<td class="left">%s</td>
				<td class="left">%d</td>
				<td class="left">$%.2f</td>
			</tr>
		`, book.CoverURL, book.Title, book.Quantity, book.Price)
	}

	emailContent := strings.NewReplacer(
		"{{first_name}}", emailData.FirstName,
		"{{order_number}}", emailData.OrderNumber,
		"{{sub_total}}", fmt.Sprintf("%.2f", emailData.SubTotal),
		"{{shipping_cost}}", fmt.Sprintf("%.2f", emailData.Shipping),
		"{{total}}", fmt.Sprintf("%.2f", emailData.Total),
		"{{order_summary}}", orderSummary,
	).Replace(emailTemplate)

	params := &resend.SendEmailRequest{
		From:    "contacto@publicacionestulip.org",
		To:      []string{emailData.Email},
		Subject: "Gracias por tu orden en TULIP!",
		Html:    emailContent,
	}

	sent, err := s.msg.Emails.Send(params)
	if err != nil {
		return "", fmt.Errorf("error sending confirmation email to customer: %s", err)
	}

	return sent.Id, nil
}
