package databasing

import (
	"database/sql"
	"log"
)

/**
user
+-------+--------------+------+-----+---------+-------+
| Field | Type         | Null | Key | Default | Extra |
+-------+--------------+------+-----+---------+-------+
| name  | varchar(255) | NO   | PRI | NULL    |       |
| pwd   | varchar(255) | NO   | PRI | NULL    |       |
+-------+--------------+------+-----+---------+-------+

channels_names
+--------------+------------------+------+-----+---------+----------------+
| Field        | Type             | Null | Key | Default | Extra          |
+--------------+------------------+------+-----+---------+----------------+
| channel_name | varchar(255)     | NO   |     | NULL    |                |
| email_name  | varchar(255)     | YES  |     | NULL    |                |
| id           | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
+--------------+------------------+------+-----+---------+----------------+

**/

type Email struct {
	User  string
	Email string
	Body  string
}

type DBEmailResponse struct {
	chl       chan *Email
	assembler func(*sql.Rows) *Email
}

func (mr *DBEmailResponse) send(result *sql.Rows) {
	mr.chl <- mr.assembler(result)
}
func (mr *DBEmailResponse) close() {
	close(mr.chl)
}

func NewEmailFull(user, email, body string) *Email {
	return &Email{
		User:  user,
		Email: email,
		Body:  body,
	}
}
func SetupEmails(db *sql.DB) {

	defineQuery(db, "Emails_All", `SELECT user,email,body FROM emails ;`)

	defineQuery(db, "Emails_AddEmail", `INSERT INTO emails VALUES (?,?,?);`)
}

func RequestEmail(name string, args ...interface{}) <-chan *Email {
	response := make(chan *Email, 1)
	queries <- &DBQueryResponse{
		query: "Emails_" + name,
		args:  args,
		sender: &DBEmailResponse{
			chl:       response,
			assembler: parseEmail,
		},
	}
	return response
}
func SaveEmail(user, email, body string) <-chan bool {
	response := make(chan bool, 1)
	actions <- &DBActionResponse{
		exec: "Emails_AddEmail",
		args: []interface{}{user, email, body},
		chl:  response,
	}
	return response
}
func parseEmail(rows *sql.Rows) *Email {
	var user, email, body string
	if err := rows.Scan(&user, &email, &body); err != nil {
		log.Fatalf(" databasing.emails.Parse: Error: %s", err)
	}
	return NewEmailFull(user, email, body)
}
