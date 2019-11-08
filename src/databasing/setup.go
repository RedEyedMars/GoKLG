package databasing

import (
	"database/sql"
	"fmt"
	"log"
	"regexp"

	"../events"

	_ "github.com/go-sql-driver/mysql"
)

/**
+---------------+
| Tables_in_klp |
+---------------+
| activity      |
| user          |
| user_activity |
+---------------+
**/

var dbQueries map[string]*sql.Stmt

var queries chan dbQuery
var actions chan dbQuery

var reSanatizeDatabase *regexp.Regexp
var reIsName *regexp.Regexp

type sendable interface{}
type dbQuery interface {
	execute()
}
type dbSender interface {
	send(*sql.Rows)
	close()
}
type DBActionResponse struct {
	exec string
	args []interface{}
	chl  chan bool
}
type DBQueryResponse struct {
	query  string
	args   []interface{}
	sender dbSender
}

func (r *DBActionResponse) execute() {
	if result, err := dbQueries[r.exec].Exec(r.args...); err != nil {
		log.Fatalf("databasing.query.action.%s : caused %s", r.exec, err)
	} else {
		if _, err := result.RowsAffected(); err != nil {
			log.Fatalf("databasing.query.action.%s : caused %s", r.exec, err)
		} else {
			events.FuncEvent("databasing.query.action."+r.exec, func() {
				r.chl <- true
				close(r.chl)
			})
		}

	}
}
func (r *DBQueryResponse) execute() {
	if rows, err := dbQueries[r.query].Query(r.args...); err != nil {
		log.Fatalf("databasing.query.query.%s : caused %s", r.query, err)
	} else {
		events.FuncEvent("databasing.query.request."+r.query, func() {
			for rows.Next() {
				r.sender.send(rows)
			}
			r.sender.close()
		})
	}
}
func RequestAction(mode string, name string, args ...interface{}) <-chan bool {
	response := make(chan bool, 1)
	actions <- &DBActionResponse{
		exec: mode + "_" + name,
		args: args,
		chl:  response,
	}
	return response
}

func Setup() {
	dbQueries = make(map[string]*sql.Stmt)
	SetupAdminCommands()

	queries = make(chan dbQuery, 16)
	actions = make(chan dbQuery, 16)

	reSanatizeDatabase = regexp.MustCompile(`(\n, \r, \, ', ")`)
	reIsName = regexp.MustCompile(`[a-zA-Z][a-zA-Z0-9_-]*`)
}
func defineQuery(db *sql.DB, name string, query string) {
	if stmt, err := db.Prepare(query); err != nil {

		log.Fatalf("databasing.defineQuery: Failed to define:%s because: %s", name, err)
	} else {
		dbQueries[name] = stmt
	}
}
func Run(Shutdown chan bool) {
	log.Println(" databasing.Run: Setting up database...")
	events.GoFuncEvent("databasing.Run", func() {
		events.FuncEvent("databasing.Setup", Setup)
		events.FuncEvent("databasing.StartDatabase", func() { StartDatabase(Shutdown) })
	})
}

func StartDatabase(Shutdown chan bool) {
	dbUser := "chat_root"
	dbName := "klg"
	dbEndpoint := "chat-service.c84g8cm4el5a.us-west-2.rds.amazonaws.com"

	// Create the MySQL DNS string for the DB connection
	// user:password@protocol(endpoint)/dbname?<params>

	dnsStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", dbUser, dbPassword, dbEndpoint, dbName)

	// Use db to perform SQL operations on database
	if db, err := sql.Open("mysql", dnsStr); err != nil {
		log.Fatal(err)
	} else {
		if err = db.Ping(); err != nil {
			log.Fatal(err)
		}

		onClose = func() {
			log.Printf("Closing database...")
			db.Close()
		}

		events.FuncEvent("databasing.SetupUsers", func() { SetupUsers(db) })
		events.FuncEvent("databasing.SetupEmails", func() { SetupEmails(db) })
		events.FuncEvent("databasing.SetupActivities", func() { SetupActivities(db) })

		events.FuncEvent("databasing.StartMessageListening", func() { StartMessageListening(db) })

	}

}

var onClose func()

func End() {
	events.FuncEvent("Databasing.End", func() {
		if onClose != nil {
			onClose()
		}
		close(queries)
		close(actions)
	})
}

func StartMessageListening(db *sql.DB) {
	for {
		select {
		case request := <-queries:
			if request == nil {
				return
			}
			go request.execute()
		case request := <-actions:
			if request == nil {
				return
			}
			request.execute()
		}
	}
}

func IsName(input string) bool {
	return reIsName.FindString(input) == input
}

func SanatizeDatabaseInput(input string) string {
	return reSanatizeDatabase.ReplaceAllStringFunc(input, func(match string) string {
		return "\\" + match
	})
}
