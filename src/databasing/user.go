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
var UsersById map[int64]*User

type User struct {
	ID   int64
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

type DBUserIdResponse struct {
	chl       chan int64
	assembler func(*sql.Rows) int64
}

func (mr *DBUserIdResponse) send(result *sql.Rows) {
	mr.chl <- mr.assembler(result)
}
func (mr *DBUserIdResponse) close() {
	close(mr.chl)
}

func NewUserFull(id int64, name string) *User {
	member := &User{
		ID:   id,
		Name: name}
	events.FuncEvent("databasing.members.AddUserToMap:"+name, func() { AddUserToMaps(member) })
	return member
}
func AddUserToMaps(member *User) {
	UsersById[member.ID] = member
	Users[member.Name] = member
}

func SetupUsers(db *sql.DB) {
	Users = make(map[string]*User)
	UsersById = make(map[int64]*User)

	defineQuery(db, "Users_All", `SELECT name FROM user ;`)

	defineQuery(db, "UsersId_ByName", `SELECT id FROM user WHERE name=? ;`)

	defineQuery(db, "Users_ById", `SELECT id,name FROM user WHERE id=? ;`)
	defineQuery(db, "Users_ByName", `SELECT id,name FROM user WHERE name=? ;`)
	defineQuery(db, "Users_ByPwd", `SELECT user.id,user.name FROM pwds INNER JOIN user ON pwds.id=user.id WHERE pwd=?;`)

	defineQuery(db, "Users_AddUser", `INSERT INTO user VALUES (NULL,?);`)
	defineQuery(db, "Users_AddPwd", `INSERT INTO pwds VALUES(?,?);`)
	defineQuery(db, "Users_Remove", `DELETE FROM user WHERE name = ?;`)

	defineQuery(db, "Pwds_Remove", `DELETE FROM pwds WHERE id=?;`)
	defineQuery(db, "UserActivity_Remove", `DELETE FROM user_activity WHERE user_id=?;`)
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
func InsertUser(name string, pwd string) <-chan *User {
	response := make(chan *User, 1)
	go func() {
		responseUser := make(chan bool, 1)
		actions <- &DBActionResponse{
			exec: "Users_AddUser",
			args: []interface{}{name},
			chl:  responseUser,
		}
		responseInterm := make(chan *User, 1)
		<-responseUser
		queries <- &DBQueryResponse{
			query: "Users_ByName",
			args:  []interface{}{name},
			sender: &DBUserResponse{
				chl:       responseInterm,
				assembler: parseUser,
			},
		}
		responsePwd := make(chan bool, 1)
		user := <-responseInterm
		actions <- &DBActionResponse{
			exec: "Users_AddPwd",
			args: []interface{}{pwd, user.ID},
			chl:  responsePwd,
		}
		<-responsePwd
		response <- user
	}()

	return response
}
func DeleteUser(name string) <-chan bool {
	response := make(chan bool, 1)
	go func() {
		var id int64
		if user := Users[name]; user != nil {
			id = user.ID
		} else {
			responseUser := make(chan int64, 1)
			queries <- &DBQueryResponse{
				query: "UsersId_ByName",
				args:  []interface{}{name},
				sender: &DBUserIdResponse{
					chl:       responseUser,
					assembler: parseUserId,
				},
			}
			id = <-responseUser
		}
		responseRemoveUser := make(chan bool, 1)
		responseRemoveUserActivty := make(chan bool, 1)
		responseRemovePwd := make(chan bool, 1)

		actions <- &DBActionResponse{
			exec: "UserActivity_Remove",
			args: []interface{}{id},
			chl:  responseRemoveUserActivty,
		}
		actions <- &DBActionResponse{
			exec: "Pwds_Remove",
			args: []interface{}{id},
			chl:  responseRemovePwd,
		}
		<-responseRemoveUserActivty
		<-responseRemovePwd
		actions <- &DBActionResponse{
			exec: "Users_Remove",
			args: []interface{}{name},
			chl:  responseRemoveUser,
		}
		<-responseRemoveUser
		delete(Users, name)
		delete(UsersById, id)
		response <- true
	}()

	return response
}
func parseUser(rows *sql.Rows) *User {
	var (
		name string
		id   int64
	)
	if err := rows.Scan(&id, &name); err != nil {
		log.Fatalf(" databasing.members.Parse: Error: %s", err)
	}
	return NewUserFull(id, name)
}
func parseUserByName(rows *sql.Rows) *User {
	var name string
	if err := rows.Scan(&name); err != nil {

		log.Fatalf(" databasing.members.ParseNames: Error: %s", err)
	}

	return Users[name]
}
func parseUserById(rows *sql.Rows) *User {
	var id int64
	if err := rows.Scan(&id); err != nil {
		log.Fatalf(" databasing.members.ParseNames: Error: %s", err)
	}

	return UsersById[id]
}
func parseUserId(rows *sql.Rows) int64 {
	var id int64
	if err := rows.Scan(&id); err != nil {
		log.Fatalf(" databasing.members.ParseNames: Error: %s", err)
	}

	return id
}
