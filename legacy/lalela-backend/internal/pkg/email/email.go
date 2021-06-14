// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package email

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"html/template"
	"log"
	"sync"

	"github.com/mailgun/mailgun-go/v4"
)

type TemplateData struct {
	To         string
	Subject    string
	Name       string
	Text       string
	Link       string
	ButtonText string
	MainText   string
}

func Send(data TemplateData) (bool, error) {
	waitGroup := sync.WaitGroup{} // a WaitGroup waits for a collection of goroutines to finish, pass this by address
	// context.WithCancel returns a copy of parent with a new Done channel.
	// The returned context's Done channel is closed when the returned cancel function is called or when the parent
	// context's Done channel is closed, whichever happens first.
	ctx := context.Background()

	// Get our API key from the environment; configure.
	// apiKey := os.Getenv("EMAIL_API_KEY")
	//emailFrom := os.Getenv("EMAIL_FROM")
	emailFrom := viper.Get("emailFrom").(string)
	mailgunDomain := viper.Get("mailgunDomain").(string)
	mailgunApiKey := viper.Get("mailgunApiKey").(string)

	// Create an instance of the Mailgun Client
	mg := mailgun.NewMailgun(mailgunDomain, mailgunApiKey)

	body, err := ParseTemplate("templates/actionable_email.html", data)
	if err != nil {
		log.Print(err)
		return false, err
	}
	// The message object allows you to add attachments and Bcc recipients
	message := mg.NewMessage(emailFrom, data.Subject, "", data.To)
	message.SetHtml(body)

	resp, id, err := mg.Send(ctx, message)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	// TODO: check if this is safe?
	waitGroup.Wait() // it blocks until the WaitGroup counter is zero

	fmt.Printf("ID: %s Resp: %s\n", id, resp)
	return true, nil
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	buf := new(bytes.Buffer)
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
