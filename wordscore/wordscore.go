package wordscore

//  wordscore database interface

import (
	"log"
	"time"

	dbase "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	"github.com/lib/pq"
)

func Version() string {
	return "1.16.2"
}

// GetWordScores func returns all wordscores.
func GetWordScores(word string) ([]hd.WordScore, error) {
	db, err := dbase.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT id,word,timeframetype,startDate,endDate,density,linkage,growth,score FROM WordScore WHERE Word='" + word + "' ORDER BY startDate"
	rows, err := db.Query(SELECT)
	dbase.CheckErr(err)
	defer rows.Close()

	// fields to read
	var id uint64
	var timeframetype int
	var dt1, dt2 time.Time
	var density, linkage, growth, score float32
	var wordscore hd.WordScore
	wordscoreList := []hd.WordScore{}

	for rows.Next() {
		err := rows.Scan(
			&id,
			&word,
			&timeframetype,
			&dt1,
			&dt2,
			&density,
			&linkage,
			&growth,
			&score)
		if err != nil {
			log.Printf("GetWordScores: %+v\n", err)
			return wordscoreList, err
		}

		timeinterval := nt.TimeInterval{Timeframetype: nt.TimeFrameType(timeframetype), StartDate: nt.New_NullTime2(dt1), EndDate: nt.New_NullTime2(dt2)}
		wordscore = hd.WordScore{Id: id, Word: word, Timeinterval: timeinterval, Density: density, Linkage: linkage, Growth: growth, Score: score}
		wordscoreList = append(wordscoreList, wordscore)
	}

	// get any iteration errors
	err = rows.Err()
	dbase.CheckErr(err)

	return wordscoreList, nil
}

// GetWordScoreListByTimeInterval func
func GetWordScoreListByTimeInterval(words []string, timeInterval nt.TimeInterval) ([]hd.WordScore, error) {
	db, err := dbase.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT id,word,timeframetype,startDate,endDate,density,linkage,growth,score FROM WordScore WHERE word IN" + dbase.CompileInClause(words) +
		"AND " + dbase.CompileDateClause(timeInterval, true)

	rows, err := db.Query(SELECT)
	dbase.CheckErr(err)
	if err != nil {
		log.Printf("GetWordScoreListByTimeInterval(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	// fields to read
	var id uint64
	var word string
	var timeframetype int
	var dt1, dt2 time.Time
	var density, linkage, growth, score float32
	var wordscore hd.WordScore
	wordscoreList := []hd.WordScore{}

	for rows.Next() {
		err := rows.Scan(
			&id,
			&word,
			&timeframetype,
			&dt1,
			&dt2,
			&density,
			&linkage,
			&growth,
			&score)
		if err != nil {
			log.Printf("GetWordScoreListByTimeInterval(2): %+v\n", err)
			return wordscoreList, err
		}

		timeinterval := nt.TimeInterval{Timeframetype: nt.TimeFrameType(timeframetype), StartDate: nt.New_NullTime2(dt1), EndDate: nt.New_NullTime2(dt2)}
		wordscore = hd.WordScore{Id: id, Word: word, Timeinterval: timeinterval, Density: density, Linkage: linkage, Growth: growth, Score: score}
		wordscoreList = append(wordscoreList, wordscore)
	}

	// get any iteration errors
	err = rows.Err()
	dbase.CheckErr(err)

	return wordscoreList, nil
}

// BulkInsertWordScores func populates [WordScore] table. Uses CopyIn. Assumes explicit schema path (search_path=public) in connection string.
func BulkInsertWordScores(wordScoreList []hd.WordScore) error {
	db, err := dbase.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin()
	dbase.CheckErr(err)

	// Must use lowercase column names! First param is table name.
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
