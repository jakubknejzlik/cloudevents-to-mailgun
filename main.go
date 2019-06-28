package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go"
)

func main() {
	ctx := context.Background()
	if err := startReceiver(ctx); err != nil {
		panic(err)
	}
}
func startReceiver(ctx context.Context) (err error) {
	_port := os.Getenv("PORT")
	if _port == "" {
		_port = "80"
	}
	port, err := strconv.Atoi(_port)
	if err != nil {
		return
	}

	sender := os.Getenv("MAILGUN_SENDER")
	mailgunDomain := os.Getenv("MAILGUN_DOMAIN")
	if mailgunDomain == "" {
		return fmt.Errorf("Missing MAILGUN_DOMAIN environment variable")
	}
	mailgunPrivateAPIKey := os.Getenv("MAILGUN_PRIVATE_API_KEY")
	if mailgunPrivateAPIKey == "" {
		return fmt.Errorf("Missing MAILGUN_PRIVATE_API_KEY environment variable")
	}

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithPort(port),
		// cloudevents.WithPath(env.Path),
	)
	if err != nil {
		return fmt.Errorf("failed to create transport, %v", err)
	}
	c, err := cloudevents.NewClient(t)
	if err != nil {
		return fmt.Errorf("failed to create client, %v", err)
	}

	templateHandler, err := NewMessageComposer()
	if err != nil {
		return fmt.Errorf("failed to create template handler, %v", err)
	}
	mailgunTransport := NewMailgunTransport(mailgunDomain, mailgunPrivateAPIKey, sender)

	log.Printf("will listen on :%d\n", port)
	log.Fatalf("failed to start receiver: %s", c.StartReceiver(ctx, gotEvent(templateHandler, mailgunTransport)))

	return nil
}

func gotEvent(templateHandler *MessageComposer, transport *MailgunTransport) func(ctx context.Context, event cloudevents.Event) error {
	return func(ctx context.Context, event cloudevents.Event) error {
		fmt.Printf("Got Event Context: %+v\n", event.Context)
		// data := &map[string]interface{}{}
		// if err := event.DataAs(data); err != nil {
		// 	fmt.Printf("Got Data Error: %s\n", err.Error())
		// }
		// fmt.Printf("Got Data: %+v\n", data)

		message, err := templateHandler.MessageFromEvent(event)
		if err != nil {
			return err
		}
		if message == nil {
			return nil
		}

		log.Println("Sending message type", event.Type(), ", to", message.To)
		if err := transport.SendMessage(*message); err != nil {
			return err
		}

		// fmt.Printf("Got Transport Context: %+v\n", cloudevents.HTTPTransportContextFrom(ctx))

		fmt.Printf("----------------------------\n")
		return nil
	}
}
