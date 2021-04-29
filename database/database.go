package database

/* Do not communicate by sharing memory; instead, share memory by communicating.
   for line:=0; line < len(source); line++ {.}	for k, v := range data {.}; for i := range myconfig{.}

   Channels are a typed conduit through which you can send and receive values with the channel operator, <-.
   By default, sends and receives block until the other side is ready.
   Sends to a buffered channel block only when the buffer is full.
   Closing a channel: v, ok := <-ch  Only the sender should close a channel, never the receiver.
   Closing is only necessary when the receiver must be told there are no more values coming.

   Goroutines run in the same address space, so access to shared memory must be synchronized. A goroutine is context-switched over an OS thread, not a CPU core.
   The Go scheduler (which runs in user space) is cooperative (not preemptive) and uses a work-stealing (not work-sharing) scheduling strategy.
   The select (case) statement lets a goroutine wait on multiple communication operations. A select blocks until one of its cases can run, then it executes that case.

   The compiler uses a technique called escape analysis to decide if a variable is going to be placed on the heap or the stack, but new always allocates on the heap.
   if the compiler cannot prove that the variable is not referenced after the function returns, then the compiler must allocate the variable on the garbage-collected heap to avoid dangling pointer errors. If you need to know where your variables are allocated pass the "-m" gc flag to "go build" or "go run" (e.g., go run -gcflags -m app.go).
   Most memory allocations are served from local thread caches.

   Ddatabase driver: go get -u github.com/jackc/pgx	(_) include this package even though the package is not explicitly referenced in code.
   s.p. inserts into table Occurrence. The defer statement should come after you check for an error from DB.Query.
*/

import (
	"context" // pgx driver uses context: see https://golang.org/pkg/context/
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	// https://pkg.go.dev/github.com/jackc/pgx/v4/pgxpool
	"github.com/jackc/pgx/v4/pgxpool"
)

// mapset https://github.com/deckarep/golang-set/blob/master/README.md & https://godoc.org/github.com/deckarep/golang-set

// Version func
func Version() string {
	return "1.16.2"
}

// DB struct
type DB struct {
	*pgxpool.Pool
}

// CheckErr database error handler.
func CheckErr(err error) {
	if err != nil {
		log.Printf("Database CheckErr %+v\n", err)
		fmt.Println(err)
		fmt.Print("Press Enter to continue...")
		os.Stdin.Read([]byte{0})
	}
}

// GetDatabaseConnectionString func uses environment var ACM_DATABASE_URL
func GetDatabaseConnectionString() string {
	connStr := os.Getenv("ACM_DATABASE_URL")
	if connStr == "" {
		log.Panic("ACM_DATABASE_URL not found in environment variables")
	}
	//fmt.Println(" Connected to " + connStr)
	return connStr
}

// GetDatabaseReference opens a database specified by its database driver name and a driver-specific data source name: db,err := GetDatabaseReference()
// defer db.Close() must follow a call to this function in the calling function. sslmode is set to 'required' by default.
// This is a postgres-only database drive! Background() returns a non-nil, empty Context. It is never canceled, has no values, and has no deadline.
func GetDatabaseReference() (*pgxpool.Pool, error) {
	dbConn := GetDatabaseConnectionString()
	db, err := pgxpool.Connect(context.Background(), dbConn)
	CheckErr(err)
	err = db.Ping(context.Background())
	if err != nil {
		fmt.Println(err)
		fmt.Print("Press Enter to continue...")
		os.Stdin.Read([]byte{0})
	}
	return db, err
}

// TestDbConnection returns a new connection after 1 attempt if db connection is dead else user prompt.
func TestDbConnection(db *pgxpool.Pool) (*pgxpool.Pool, error) {
	err := db.Ping(context.Background())
	if err != nil {
		db.Close()
		time.Sleep(1000)
		dbx, err := GetDatabaseReference()
		if err != nil {
			fmt.Print("There is a problem accessing the database. Press Enter to try again.")
			os.Stdin.Read([]byte{0})
			dbx, err = GetDatabaseReference()
		}
		return dbx, err
	}
	return db, nil
}

// NoRowsReturned func
func NoRowsReturned(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "no rows in result set")
}

// CompileInClause func inserts quote marks for IN db clause.
func CompileInClause(words []string) string {
	wordlist := make([]string, 0)
	for _, word := range words {
		w := strings.TrimSpace(word)
		if len(w) > 0 {
			wordlist = append(wordlist, "'"+w+"'")
		}
	}
	if len(wordlist) > 0 {
		return " (" + strings.Join(wordlist, ", ") + ") "
	} else {
		return " ('') "
	}
}

// GetFormattedDatesForProcedure func includes parentheses.
func GetFormattedDatesForProcedure(timeInterval nt.TimeInterval) string {
	return "('" + timeInterval.StartDate.StandardDate() + "', '" + timeInterval.EndDate.StandardDate() + "')"
}

// GetWhereClause func. Don't know PostgreSQL limit of IN values.
func GetWhereClause(columnName string, wordGrams []string) string {
	var sb strings.Builder
	sb.WriteString(columnName + " IN (")
	for ndx := 0; ndx < len(wordGrams); ndx++ {
		sb.WriteString("'" + wordGrams[ndx] + "'")
		if ndx < len(wordGrams)-1 {
			sb.WriteString(",")
		}
	}
	sb.WriteString(") ")
	return sb.String()
}

// GetSingleDateWhereClause func
func GetSingleDateWhereClause(columnName string, timeInterval nt.TimeInterval) string {
	return columnName + " >= '" + timeInterval.StartDate.StandardDate() + "' AND " + columnName + " <= '" + timeInterval.EndDate.StandardDate() + "' "
}

// CompileDateClause func
func CompileDateClause(timeInterval nt.TimeInterval, useTimeframetype bool) string {
	if useTimeframetype {
		return "timeframetype=" + strconv.Itoa(int(timeInterval.Timeframetype)) + " AND startDate >= '" + timeInterval.StartDate.StandardDate() + "' AND endDate <= '" + timeInterval.EndDate.StandardDate() + "' "
	}
	return "startDate >= '" + timeInterval.StartDate.StandardDate() + "' AND endDate <= '" + timeInterval.EndDate.StandardDate() + "' "
}

/* Timestamptz support for pgx driver:
func (tw *timeWrapper) Scan(in interface{}) error {
	var t pgtype.Timestamptz
	err := t.Scan(in)
	if err != nil {
		return err
	}

	tp, err := ptypes.TimestampProto(t.Time)
	if err != nil {
		return err
	}

	*tw = (timeWrapper)(*tp)
	return nil
} */

/*************************************************************************************************/

type UserProfile struct {
	ID          int       `json:"id"`
	UserName    string    `json:"username"`
	Password    string    `json:"password"`
	DateUpdated time.Time `json:"dateupdated"`
}

// GetUser func assumes unique case-insensitive userName.
func GetUser(userName string) (UserProfile, error) {
	var user UserProfile
	db, err := GetDatabaseReference()
	if err != nil {
		return user, err
	}
	defer db.Close()

	SELECT := "SELECT id, UserName, Password, DateUpdated FROM Vocabulary WHERE LOWER(UserName)='" + strings.ToLower(userName) + "'"
	err = db.QueryRow(context.Background(), SELECT).Scan(&user.ID, &user.UserName, &user.Password, &user.DateUpdated)
	CheckErr(err)
	if err != nil {
		return user, err
	}

	return user, nil
}
