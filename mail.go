package main

import (
	"gopkg.in/mail.v2"
	"strings"
)

func sendEmail(subject, body string) error {
	m := mail.NewMessage()

	// Set E-Mail sender
	m.SetHeader("From", env.EmailFrom)

	// Set E-Mail receivers
	m.SetHeader("To", strings.Join(env.EmailTo, ","))

	// Set E-Mail subject
	m.SetHeader("Subject", subject)

	m.SetBody("text/plain", body)

	// Settings for SMTP server
	d := mail.NewDialer(env.SMTPHost, env.SMTPPort, env.SMTPUser, env.SMTPPassword)

	// This is only needed when SSL/TLS certificate is not valid on server.
	// In production this should be set to false.
	//d.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	// Now send E-Mail
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}

func SendEmail(subject, body string) {
	go func() {
		if err := sendEmail(subject, body); err != nil {
			logErrLn("Failed sending email:", err)

		}
	}()
}
