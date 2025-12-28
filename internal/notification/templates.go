package notification

const (
	TemplateOrderCreatedSubject = "Order Confirmation - Order #%s"
	TemplateOrderCreatedBody    = "<h1>Order Confirmation</h1><p>Thank you for your order!</p><p>Order ID: <strong>%s</strong></p><p>We will process it shortly.</p>"

	TemplatePaymentSuccessSubject = "Payment Receipt - Order #%s"
	TemplatePaymentSuccessBody    = "<h1>Payment Successful</h1><p>We received your payment for Order ID: <strong>%s</strong></p><p>Amount: %s</p>"

	TemplateOrderCancelledSubject = "Order Cancelled - Order #%s"
	TemplateOrderCancelledBody    = "<h1>Order Cancelled</h1><p>Your order #%s has been cancelled as requested.</p>"

	TemplatePaymentFailedSubject = "Payment Failed - Order #%s"
	TemplatePaymentFailedBody    = "<h1>Payment Failed</h1><p>We could not process payment for order #%s. Please try again.</p>"
)

// In a real system, these might assume a struct payload and use html/template.
// For now, we use simple Sprintf format strings as "templates".
