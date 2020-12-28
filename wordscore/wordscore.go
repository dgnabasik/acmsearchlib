package wordscore

//  wordscore database interface

import (
	"time"

	dbase "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	"github.com/lib/pq"
	_ "github.com/lib/pq"
)

func Version() string {
	return "1.0.10"
}

// GetWordScore func
func GetWordScore(word string) hd.WordScore {
	db, err := dbase.GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT id,word,timeframetype,startDate,endDate,density,linkage,growth,score FROM WordScore WHERE Word='" + word + "';"
	rows, err := db.Query(SELECT)
	dbase.CheckErr(err)
	defer rows.Close()

	// fields to read
	var id uint64
	var timeframetype int
	var wordA string
	var dt1, dt2 time.Time
	var density, linkage, growth, score float32

	err = db.QueryRow(SELECT).Scan(&id, &wordA, &timeframetype, &dt1, &dt2, &density, &linkage, &growth, &score)
	dbase.CheckErr(err)

	startDate := nt.New_NullTime2(dt1)
	endDate := nt.New_NullTime2(dt2)
	tfType := nt.TimeFrameType(timeframetype)

	wordScore := hd.New_WordScore(id, wordA, tfType, startDate, endDate, density, linkage, growth, score)

	return wordScore
}

// BulkInsertWordScores func populates [WordScore] table. Uses CopyIn. Assumes explicit schema path (search_path=public) in connection string.
func BulkInsertWordScores(wordScoreList []hd.WordScore) error {
	db, err := dbase.GetDatabaseReference()
	defer db.Close()

	txn, err := db.Begin()
	dbase.CheckErr(err)

	//original: stmt, err := db.Prepare("INSERT INTO WordScore (word, timeframetype, startdate, enddate, density, linkage, growth, score) VALUES ($1, $2, $3, $4, $5, $6, $7, $8);")
	// [wordscore] in schema wordscore. Must use lowercase column names!
	stmt, err := txn.Prepare(pq.CopyIn("wordscore", "word", "timeframetype", "startdate", "enddate", "density", "linkage", "growth", "score"))
	dbase.CheckErr(err)

	for _, v := range wordScoreList {
		_, err = stmt.Exec(v.Word, int(v.Timeinterval.Timeframetype), v.Timeinterval.StartDate.DT, v.Timeinterval.EndDate.DT, v.Density, v.Linkage, v.Growth, v.Score)
		dbase.CheckErr(err)
	}

	_, err = stmt.Exec()
	dbase.CheckErr(err)

	err = stmt.Close()
	dbase.CheckErr(err)

	err = txn.Commit()
	dbase.CheckErr(err)

	return err
}