package main

const (
	// Sender address must be verified with Amazon SES.
	Sender = "noreply@yourdomain.com"

	// Recipient address.
	Recipient = "receiver@domain.com"

	// Subject line for the email.
	Subject = "Yandex Cloud Postbox Test via AWS SDK for Go"

	// HtmlBody is the body for the email.
	HtmlBody = "<h1>Amazon SES Test Email (AWS SDK for Go)</h1><p>This email was sent with " +
		"<a href='https://yandex.cloud/ru/docs/postbox/quickstart'>Yandex Cloud Postbox</a> using the " +
		"<a href='https://aws.amazon.com/sdk-for-go/'>AWS SDK for Go</a>.</p>"

	// TextBody is the email body for recipients with non-HTML email clients.
	TextBody = "This email was sent with Yandex Cloud Postbox using the AWS SDK for Go."

	// CharSet The character encoding for the email.
	CharSet = "UTF-8"
)
