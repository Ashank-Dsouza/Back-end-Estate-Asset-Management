package models

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"html/template"
	"os"
	"strconv"

	"github.com/badoux/checkmail"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	gomail "gopkg.in/mail.v2"
)

type SendMail struct {
	Email string
}

func (sm *SendMail) SendEmail(link string, email_type string) error {

	if len(sm.Email) < 1 {
		return errors.New("Required Email")
	}
	if err := checkmail.ValidateFormat(sm.Email); err != nil {
		return errors.New("Invalid Email")
	}

	templateData := struct {
		Password string
	}{
		Password: link,
	}

	m := gomail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", os.Getenv("SYSTEM_EMAIL"))

	m.SetHeader("Subject", "Default Subject")

	// Set E-Mail receivers
	fmt.Print("the email address of the user who wants to change pswd is: " + sm.Email)
	m.SetHeader("To", sm.Email)

	path := "./html/email_template.html"

	if email_type == "Welcome" {
		m.SetHeader("Subject", "Welcome to TruVest")
		//templateData.Password = models.GetFrontEndUrl() + "/" + sm.Link.String()
		path = "./html/welcome_email_template.html"
		//fmt.Print(" " + path)
	}
	if email_type == "ResetPassword" {
		m.SetHeader("Subject", "Reset TruVest Password")

		path = "./html/password_reset_email.html"
	}

	var err error
	t, err := template.ParseFiles(path)

	if err != nil {
		return err
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, templateData); err != nil {
		return err
	}
	result := tpl.String()

	// Set E-Mail body. You can set plain text or html with text/html
	m.SetBody("text/html", result)

	port, err := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if err != nil {
		fmt.Print("did not get post variable from env")
		return err
	}

	// Settings for SMTP server
	//d := gomail.NewDialer("smtp.gmail.com", 587, "serealestate97@gmail.com", "Peterparker996")

	d := gomail.NewDialer(os.Getenv("SMTP_HOST"), port, os.Getenv("SYSTEM_EMAIL"), os.Getenv("SYSTEM_EMAIL_PASSWORD"))

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		fmt.Println("The email was not sent successfully!!")
		fmt.Println(err)
	}

	return nil
}

func (sm *SendMail) SendGridMail() error {
	from := mail.NewEmail(os.Getenv("SYSTEM_ADMIN_USERNAME"), os.Getenv("SYSTEM_EMAIL"))
	subject := "Reset your credentials"
	to := mail.NewEmail("Sender User", sm.Email)
	plainTextContent := "Hi, Please set your credentials by following steps below:"

	var err error
	t, err := template.ParseFiles("./html/email_template.html")
	if err != nil {
		return err
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, nil); err != nil {
		return err
	}
	htmlContent := tpl.String()

	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(response.StatusCode)
	fmt.Println(response.Body)
	fmt.Println(response.Headers)
	return nil
}
