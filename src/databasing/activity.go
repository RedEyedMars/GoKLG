package databasing

import (
	"database/sql"
	"log"
	"time"
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
| activity_name  | varchar(255)     | YES  |     | NULL    |                |
| id           | int(10) unsigned | NO   | PRI | NULL    | auto_increment |
+--------------+------------------+------+-----+---------+----------------+

**/

var Activities map[string]*Activity
var ActivityById map[int64]*Activity

type Activity struct {
	ID   int64
	Name string
}

type DBActivityResponse struct {
	chl       chan *Activity
	assembler func(*sql.Rows) *Activity
}

func (mr *DBActivityResponse) send(result *sql.Rows) {
	mr.chl <- mr.assembler(result)
}
func (mr *DBActivityResponse) close() {
	close(mr.chl)
}

func NewActivity(id int64, name string) *Activity {
	newActivity := &Activity{
		ID:   id,
		Name: name,
	}
	AddActivityToMaps(newActivity)
	return newActivity
}

func AddActivityToMaps(activity *Activity) {
	ActivityById[activity.ID] = activity
	Activities[activity.Name] = activity
}

func SetupActivities(db *sql.DB) {
	Activities = make(map[string]*Activity)
	ActivityById = make(map[int64]*Activity)

	defineQuery(db, "Activities_All", `SELECT id,name FROM activity ;`)
	defineQuery(db, "Activities_ByName", `SELECT id,name FROM activity WHERE name=? ;`)

	defineQuery(db, "Activities_AddActivity", `INSERT INTO activity VALUES (NULL,?);`)
	defineQuery(db, "Activities_LogActivity", `INSERT INTO user_activity VALUES (?,?,?);`)
}

func RequestActivity(name string, args ...interface{}) <-chan *Activity {
	response := make(chan *Activity, 1)
	queries <- &DBQueryResponse{
		query: "Activities_" + name,
		args:  args,
		sender: &DBActivityResponse{
			chl:       response,
			assembler: parseActivity,
		},
	}
	return response
}
func InsertActivity(name string) chan bool {
	response := make(chan bool, 1)
	actions <- &DBActionResponse{
		exec: "Activities_AddActivity",
		args: []interface{}{name},
		chl:  response,
	}
	return response
}
func LogActivity(userId int64, activityName string, ts time.Time) <-chan bool {
	response := make(chan bool, 1)
	go func() {

		if user := UsersById[userId]; user == nil {
			user = <-RequestUser("ById", userId)
			if user == nil {
				response <- false
				return
			}
		}
		activity := Activities[activityName]
		if activity == nil {
			<-InsertActivity(activityName)
			activity = <-RequestActivity("ByName", activityName)
		}
		actions <- &DBActionResponse{
			exec: "Activities_LogActivity",
			args: []interface{}{userId, activity.ID, ts},
			chl:  response,
		}
	}()
	return response
}
func parseActivity(rows *sql.Rows) *Activity {
	var (
		id   int64
		name string
	)
	if err := rows.Scan(&id, &name); err != nil {
		log.Fatalf(" databasing.activities.Parse: Error: %s", err)
	}
	return NewActivity(id, name)
}
