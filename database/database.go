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

	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

// mapset https://github.com/deckarep/golang-set/blob/master/README.md & https://godoc.org/github.com/deckarep/golang-set

// DB struct
// dbRef, err := dbase.GetDatabaseReference()
// dbObj := &ArticleDatastore{db: dbRef}
type DB struct {
	*sql.DB
}

// CheckErr database error handler.
func CheckErr(err error) {
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
}

// GetDatabaseReference opens a database specified by its database driver name and a driver-specific data source name: db,err := GetDatabaseReference()
// defer db.Close() must follow a call to this function in the calling function. sslmode is set to 'required' with lib/pq by default.
func GetDatabaseReference() (*sql.DB, error) {
	const (
		dbHost        = "localhost"
		dbPort        = 5432
		dbUser        = "postgres"
		dbPassword    = "Ski7Vail!"
		dbName        = "postgres"
		dbDriver      = "postgres"
		dbSchema      = "acmsearch,public"
		dbConnections = 10
	)

	dbConn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s search_path=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName, dbSchema)
	db, err := sql.Open(dbDriver, dbConn)
	CheckErr(err)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(dbConnections)
	db.SetConnMaxLifetime(0)
	err = db.Ping() // connects
	CheckErr(err)
	return db, err
}

// CallTruncateTables truncates tables with sequences.
func CallTruncateTables() error {
	db, err := GetDatabaseReference()
	defer db.Close()

	_, err = db.Exec("call TruncateTables()")
	CheckErr(err)

	fmt.Println("CallTruncateTables() done.")
	return nil
}

// GetArticleCount func
func GetArticleCount() int {
	db, err := GetDatabaseReference()
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM AcmData;").Scan(&count)
	CheckErr(err)
	return count
}

// GetLastDateSavedFromDb returns the earliest and latest AcmData.ArchiveDate values else default time.
func GetLastDateSavedFromDb() (nt.NullTime, nt.NullTime, error) {
	articleCount := GetArticleCount()
	if articleCount == 0 {
		return nt.New_NullTime(""), nt.New_NullTime(""), nil // default time.
	}

	db, err := GetDatabaseReference()
	defer db.Close()

	var archiveDate1, archiveDate2 nt.NullTime // NullTime supports Scan() interface.

	err = db.QueryRow("SELECT MIN(ArchiveDate) FROM AcmData;").Scan(&archiveDate1)
	CheckErr(err)

	err = db.QueryRow("SELECT MAX(ArchiveDate) FROM AcmData;").Scan(&archiveDate2)
	CheckErr(err)

	return archiveDate1, archiveDate2, nil
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
