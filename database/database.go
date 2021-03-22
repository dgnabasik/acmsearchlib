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

   Ddatabase driver: go get -u github.com/lib/pq	(_) include this package even though the package is not explicitly referenced in code.
   pq driver: NullTime implements the sql.Scanner interface so it can be used as a scan destination, similar to sql.NullString.
   s.p. inserts into table Occurrence. The defer statement should come after you check for an error from DB.Query.
*/

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

// mapset https://github.com/deckarep/golang-set/blob/master/README.md & https://godoc.org/github.com/deckarep/golang-set

// Version func
func Version() string {
	return "1.16.2"
}

// DB struct
// dbRef, err := dbase.GetDatabaseReference()
// dbObj := &ArticleDatastore{db: dbRef}
type DB struct {
	*sql.DB
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
// defer db.Close() must follow a call to this function in the calling function. sslmode is set to 'required' with lib/pq by default.
func GetDatabaseReference() (*sql.DB, error) {
	const (
		dbDriver      = "postgres"
		dbConnections = 20
	)

	dbConn := GetDatabaseConnectionString()
	db, err := sql.Open(dbDriver, dbConn)
	CheckErr(err)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(dbConnections)
	db.SetConnMaxLifetime(0)
	err = db.Ping() // connects
	if err != nil {
		fmt.Println(err)
		fmt.Print("Press Enter to continue...")
		os.Stdin.Read([]byte{0})
	}
	return db, err
}

// TestDbConnection returns a new connection after 1 attempt if db connection is dead else user prompt.
func TestDbConnection(db *sql.DB) (*sql.DB, error) {
	err := db.Ping()
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

// CallTruncateTables truncates tables with sequences.
func CallTruncateTables() error {
	db, err := GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec("call TruncateTables()")
	CheckErr(err)

	fmt.Println("CallTruncateTables() done.")
	return nil
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
	return " (" + strings.Join(wordlist, ", ") + ") "
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
func CompileDateClause(timeInterval nt.TimeInterval) string {
	return "timeframetype=" + strconv.Itoa(int(timeInterval.Timeframetype)) + " AND startDate >= '" + timeInterval.StartDate.StandardDate() + "' AND endDate <= '" + timeInterval.EndDate.StandardDate() + "' "
}

/*************************************************************************************************/

// BulkUpdateVocabularySpeechpart concatentates parts into output. Unknown list is returned.
// Change to root word. See https://www.datamuse.com/api/ & https://www.wordsapi.com/
// curl "https://wordsapiv1.p.mashape.com/words/soliloquy" -H "X-Mashape-Key: <APIkey>"
/*
func BulkUpdateVocabularySpeechpart() []string {
	fmt.Print("BulkInsertVocabulary_Speechpart: ")
	var wordNetSpeechParts = hd.New_WordNetSpeechParts()
	var word string
	var wordSet []string

	start := time.Now()
	db, err := GetDatabaseReference()
	defer db.Close()

	// part 1: WHERE SpeechPart not assigned.
	SELECT := "SELECT word FROM Vocabulary WHERE SpeechPart='';"
	rows, err := db.Query(SELECT)
	CheckErr(err)
	for rows.Next() {
		err = rows.Scan(&word)
		CheckErr(err)
		wordSet = append(wordSet, word)
	}

	// part 2:
	txn, err := db.Begin()
	CheckErr(err)

	stmt, err := db.Prepare("UPDATE vocabulary SET SpeechPart= $1 WHERE Word= $2;")
	CheckErr(err)

	for _, w := range wordSet {
		speechPart := wordNetSpeechParts.GetSpeechpart(w)
		if speechPart == "" {
			wordNetSpeechParts.Unknown = append(wordNetSpeechParts.Unknown, w)
		}
		_, err = stmt.Exec(speechPart, w)
		CheckErr(err)
	}

	err = stmt.Close()
	CheckErr(err)

	err = txn.Commit()
	CheckErr(err)

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())

	return wordNetSpeechParts.Unknown
}
*/
/*************************************************************************************************/
