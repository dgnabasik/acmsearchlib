package wordscore

//  wordscore database interface

import (
	"context"
	"errors"
	"log"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	pgx "github.com/jackc/pgx/v4"
)

const wordscoreSelect = "SELECT id,word,timeframetype,startDate,endDate,density,linkage,growth,score FROM WordScore"

// GetWordScores func returns all wordscores.
func GetWordScores(word string) ([]hd.WordScore, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := wordscoreSelect + " WHERE Word='" + word + "' ORDER BY startDate"
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
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
	dbx.CheckErr(err)

	return wordscoreList, nil
}

// GetWordScoreListByTimeInterval func
func GetWordScoreListByTimeInterval(words []string, timeInterval nt.TimeInterval) ([]hd.WordScore, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := wordscoreSelect + " WHERE word IN" + dbx.CompileInClause(words) + "AND " + dbx.CompileDateClause(timeInterval, true)
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetWordScoreListByTimeInterval(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

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

	err = rows.Err()
	dbx.CheckErr(err)

	return wordscoreList, nil
}

// BulkInsertWordScores func populates [WordScore] table. Assumes explicit schema path (search_path=public) in connection string.
func BulkInsertWordScores(wordScoreList []hd.WordScore) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)
	defer txn.Rollback(context.Background())

	// Must use lowercase column names!
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"wordscore"}, // tablename
		[]string{"word", "timeframetype", "startdate", "enddate", "density", "linkage", "growth", "score"},
		pgx.CopyFromSlice(len(wordScoreList), func(i int) ([]interface{}, error) {
			return []interface{}{wordScoreList[i].Word, int(wordScoreList[i].Timeinterval.Timeframetype), wordScoreList[i].Timeinterval.StartDate.DT, wordScoreList[i].Timeinterval.EndDate.DT, wordScoreList[i].Density, wordScoreList[i].Linkage, wordScoreList[i].Growth, wordScoreList[i].Score}, nil
		}),
	)
	dbx.CheckErr(err)
	if copyCount == 0 {
		log.Printf("BulkInsertWordScores: no rows inserted")
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	return err
}

// DeleteWordscoreByPeriod func
func DeleteWordscoreByPeriod(timeInterval nt.TimeInterval) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	DELETE := " DELETE FROM Wordscore WHERE " + dbx.CompileDateClause(timeInterval, true)
	commandTag, err := db.Exec(context.Background(), DELETE)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() != 1 {
		return errors.New("no rows were deleted")
	}
	return nil
}

// GetWordscorePeriodGroup func
func GetWordscorePeriodGroup() ([]nt.TimeInterval, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT timeframetype, MAX(startdate) AS startdate, MAX(enddate) AS enddate FROM Wordscore GROUP BY timeframetype ORDER BY timeframetype"
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	// fields to read
	var startdate, enddate time.Time
	var timeframetype int
	var dateList []nt.TimeInterval

	for rows.Next() {
		err = rows.Scan(&timeframetype, &startdate, &enddate)
		dbx.CheckErr(err)
		ti := nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startdate), nt.New_NullTime2(enddate))
		dateList = append(dateList, ti)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return dateList, nil
}
