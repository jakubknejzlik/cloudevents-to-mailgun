package main

import (
	"context"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go"
)

// SMTPTransport ...
type MailgunTransport struct {
	mg     mailgun.Mailgun
	sender string
}

// SMTPTransportMessage ...
type SMTPTransportMessage struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
	HTML    string   `json:"html"`
}

// NewMailgunTransport ...
func NewMailgunTransport(domain, privateAPIKey, sender string) *MailgunTransport {
	mg := mailgun.NewMailgun(domain, privateAPIKey)

	return &MailgunTransport{mg, sender}
}

// SendMessage ...
func (t *MailgunTransport) SendMessage(msg SMTPTransportMessage) error {
	// sender := msg.From
	// if sender == "" {
	// 	sender = t.sender
	// }

	// address, err := mail.ParseAddress(sender)
	// if err != nil {
	// 	return err
	// }

	sender := msg.From
	if sender == "" {
		sender = t.sender
	}

	fmt.Println("sending message", sender, msg.Subject, msg.HTML, msg.To)
	message := t.mg.NewMessage(sender, msg.Subject, msg.HTML, msg.To...)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	// Send the message	with a 10 second timeout
	resp, id, err := t.mg.Send(ctx, message)

	if err != nil {
		panic(err)
		return err
	}

	fmt.Printf("ID: %s Resp: %s\n", id, resp)

	return nil
}
