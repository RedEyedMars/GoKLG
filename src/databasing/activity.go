package databasing

import (
	"database/sql"
	"fmt"
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

type DBReportResponse struct {
	chl       chan string
	assembler func(*sql.Rows) string
}

func (mr *DBReportResponse) send(result *sql.Rows) {
	mr.chl <- mr.assembler(result)
}
func (mr *DBReportResponse) close() {
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

func LoadAllActivities() {
	for activity := range RequestActivity("All") {
		doNothing(activity)
	}
}

func doNothing(_ interface{}) {

}

func SetupActivities(db *sql.DB) {
	Activities = make(map[string]*Activity)
	ActivityById = make(map[int64]*Activity)

	defineQuery(db, "Activities_All", `SELECT id,name FROM activity ;`)
	defineQuery(db, "Activities_ByName", `SELECT id,name FROM activity WHERE name=? ;`)

	defineQuery(db, "Activities_AddActivity", `INSERT INTO activity VALUES (NULL,?);`)
	defineQuery(db, "Activities_LogActivity", `INSERT INTO user_activity VALUES (?,?,?);`)

	defineQuery(db, "Activities_Report", `SELECT
																					user.name as user_name,
																					activity.name as activity_name,
																					COUNT(occurrence) as amount,
																					MIN(occurrence) as first_occurrence,
																					MAX(occurrence) as last_occurrence
																			  FROM user_activity
																					INNER JOIN user ON user.id=user_id
																					INNER JOIN activity ON activity.id=activity_id
																				WHERE MONTH(occurrence)=?
																				GROUP BY user.name,activity.name;`)
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
func LogActivity(activityName string, userName string, ts time.Time) <-chan bool {
	response := make(chan bool, 1)
	go func() {
		user := Users[userName]
		if user == nil {
			user = <-RequestUser("ByName", userName)
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
		log.Printf(" log activity for:%d,%d", user.ID, activity.ID)
		actions <- &DBActionResponse{
			exec: "Activities_LogActivity",
			args: []interface{}{activity.ID, user.ID, ts},
			chl:  response,
		}
	}()
	return response
}
func RequestOctoberReport() <-chan string {
	return RequestReport("October")
}
func RequestReport(month string) <-chan string {
	var month_i int64
	switch month {
	case "January":
		month_i = 1
	case "February":
		month_i = 2
	case "March":
		month_i = 3
	case "April":
		month_i = 4
	case "May":
		month_i = 5
	case "June":
		month_i = 6
	case "July":
		month_i = 7
	case "August":
		month_i = 8
	case "September":
		month_i = 9
	case "October":
		month_i = 10
	case "November":
		month_i = 11
	case "December":
		month_i = 12
	}

	response2 := make(chan string, 1)
	go func() {
		response := make(chan string, 160)
		queries <- &DBQueryResponse{
			query: "Activities_Report",
			args:  []interface{}{month_i},
			sender: &DBReportResponse{
				chl:       response,
				assembler: parseReport,
			},
		}
		ret := `<table id="report"><tr><th>user_name</th><th>activity_name</th><th>amount</th><th>first_occurrence</th><th>last_occurrence</th></tr>`
		for row := range response {
			ret += row
		}
		ret += "</table>"
		response2 <- ret
		close(response2)
	}()
	return response2
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

func parseReport(rows *sql.Rows) string {
	var (
		username         string
		activityname     string
		amount           int64
		first_occurrence time.Time
		last_occurrence  time.Time
	)
	if err := rows.Scan(&username, &activityname, &amount, &first_occurrence, &last_occurrence); err != nil {
		log.Fatalf(" databasing.activities.Parse: Error: %s", err)
	}
	return fmt.Sprintf(`<tr><td>%s</td><td>%s</td><td>%d</td><td>%s</td><td>%s</td><td></tr>`, username, activityname, amount, first_occurrence.Format("2006-01-02 15:04:05"), last_occurrence.Format("2006-01-02 15:04:05"))
}
