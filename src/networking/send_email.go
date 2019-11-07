package networking

import (
	"log"
	"net/smtp"
)

func send_email(user, email, body string) {
	from := "greg_estouffey@hotmail.com"
	pass := "ZX!@cv34qwas"
	to := "redeyedmars@gmail.com"

	msg := "From: " + user + "\n" +
		"To: " + to + "\n" +
		"Subject: " + user + "\n\n" +
		"Email:" + email +
		body

	err := smtp.SendMail("smtp.gmail.com:465",
		smtp.PlainAuth("", from, pass, "smtp.live.com"),
		from, []string{to}, []byte(msg))

	if err != nil {
		log.Printf("smtp error: %s", err)
		return
	}
}
