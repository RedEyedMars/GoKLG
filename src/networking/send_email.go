package networking

import (
	"log"
	"net/smtp"
)

func send_email(user, email, body string) {
	from := "redeyedmars@gmail.com"
	pass := "ZX!@cv34qwas"
	to := "redeyedmars@gmail.com"

	msg := "From: " + user + "\n" +
		"To: " + to + "\n" +
		"Subject: " + user + "\n\n" +
		"Email:" + email +
		body

	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
