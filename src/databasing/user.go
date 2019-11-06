package databasing

import (
	"database/sql"
	"log"

	"../events"
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
| member_name  | varchar(255)     | YES  |     | NULL    |                |
| id           | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
+--------------+------------------+------+-----+---------+----------------+

**/

var Users map[string]*User

type User struct {
	Name string
}

type DBUserResponse struct {
	chl       chan *User
	assembler func(*sql.Rows) *User
}

func (mr *DBUserResponse) send(result *sql.Rows) {
	mr.chl <- mr.assembler(result)
}
func (mr *DBUserResponse) close() {
	close(mr.chl)
}

func NewUserFull(name string) *User {
	member := &User{
		Name: name}
	events.FuncEvent("databasing.members.AddUserToMap:"+name, func() { AddUserToMaps(member) })
	return member
}
func AddUserToMaps(member *User) {
	Users[member.Name] = member
}

func SetupUsers(db *sql.DB) {
	defineQuery(db, "Users_All", `SELECT name FROM user ;`)

	defineQuery(db, "Users_ByName", `SELECT name FROM user WHERE name=? ;`)
	defineQuery(db, "Users_ByPwd", `SELECT name FROM user WHERE pwd=? ;`)

	defineQuery(db, "Users_Add", `INSERT INTO user VALUES (?,?);`)
	defineQuery(db, "Users_Remove", `DELETE FROM user WHERE name = ?;`)
}

func RequestUser(name string, args ...interface{}) <-chan *User {
	response := make(chan *User, 1)
	queries <- &DBQueryResponse{
		query: "Users_" + name,
		args:  args,
		sender: &DBUserResponse{
			chl:       response,
			assembler: parseUser,
		},
	}
	return response
}
func RequestUsersByName(name string, args ...interface{}) <-chan *User {
	response := make(chan *User, 1)
	queries <- &DBQueryResponse{
		query: "Users_" + name,
		args:  args,
		sender: &DBUserResponse{
			chl:       response,
			assembler: parseUserByName,
		},
	}
	return response
}
func parseUser(rows *sql.Rows) *User {
	var name string
	if err := rows.Scan(&name); err != nil {
		log.Fatalf(" databasing.members.Parse: Error: %s", err)
	}
	return NewUserFull(name)
}
func parseUserByName(rows *sql.Rows) *User {
	var name string
	if err := rows.Scan(&name); err != nil {

		log.Fatalf(" databasing.members.ParseNames: Error: %s", err)
	}

	return Users[name]
}
