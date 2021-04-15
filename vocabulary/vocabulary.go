package vocabulary

//  vocabulary.go manages hd.Vocabulary in database.

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	pgx "github.com/jackc/pgx/v4"
)

// Version func
func Version() string {
	return "1.16.2"
}

// GetVocabularyByWord func
func GetVocabularyByWord(wordX string) (hd.Vocabulary, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return hd.Vocabulary{}, err
	}
	defer db.Close()

	var word, speechPart, stem string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	SELECT := "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart, stem FROM Vocabulary WHERE Word='" + wordX + "'"
	err = db.QueryRow(context.Background(), SELECT).Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart, &stem)
	dbx.CheckErr(err)
	return hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart, Stem: stem}, nil
}

// GetVocabularyList method does NOT apply filtering to imported []words. Func places single quotes around each words element.
func GetVocabularyList(words []string) ([]hd.Vocabulary, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(words)
	query := "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart, stem FROM vocabulary WHERE word IN " + inPhrase
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetVocabularyList(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var vocabulary hd.Vocabulary
	vocabularyList := []hd.Vocabulary{}
	for rows.Next() {
		err := rows.Scan(
			&vocabulary.Id,
			&vocabulary.Word,
			&vocabulary.RowCount,
			&vocabulary.Frequency,
			&vocabulary.WordRank,
			&vocabulary.Probability,
			&vocabulary.SpeechPart,
			&vocabulary.Stem)
		if err != nil {
			log.Printf("GetVocabularyList(2): %+v\n", err)
			return nil, err
		}
		vocabularyList = append(vocabularyList, vocabulary)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return vocabularyList, err
}

func getAcmGraphCount() string {
	count := os.Getenv("REACT_ACM_GRAPH_COUNT")
	if count == "" {
		count = "1"
	}
	return count
}

// GetWordListMap method returns all words if prefix is blank. Also filters by REACT_ACM_GRAPH_COUNT.
func GetWordListMap(prefix string) ([]hd.LookupMap, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	occurrenceCount := getAcmGraphCount()

	query := "SELECT id, word FROM vocabulary WHERE occurrenceCount > " + occurrenceCount + " ORDER BY word"
	if len(prefix) > 0 {
		query = "SELECT id, word FROM vocabulary WHERE occurrenceCount > " + occurrenceCount + " AND word LIKE '" + strings.ToLower(prefix) + "%' ORDER BY word"
	}
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetWordListMap(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var word string
	var id int
	lookupMap := make([]hd.LookupMap, 0)
	for rows.Next() {
		err := rows.Scan(&id, &word)
		if err != nil {
			log.Printf("GetWordListMap(2): %+v\n", err)
			return nil, err
		}
		lookupMap = append(lookupMap, hd.LookupMap{Value: id, Label: word})
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return lookupMap, err
}

// GetVocabularyListByDate reads Vocabulary table filtered by articleList.
func GetVocabularyListByDate(timeinterval nt.TimeInterval) ([]hd.Vocabulary, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT * FROM GetVocabularyByDate" + dbx.GetFormattedDatesForProcedure(timeinterval)
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var word, speechPart, stem string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	var vocabList []hd.Vocabulary

	for rows.Next() { // this order follows the \d Vocabulary description:
		err = rows.Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart, &stem)
		dbx.CheckErr(err)

		newWord, rule := cond.FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		vocabList = append(vocabList, hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart, Stem: stem})
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return vocabList, err
}

// GetVocabularyMapProbability Read all Vocabulary.Word,Probability values if wordGrams is []string{}. Applys filtering.
// Vocabulary.probability is NOT used to calculate conditional probabilities!
func GetVocabularyMapProbability(wordGrams []string, timeInterval nt.TimeInterval) (map[string]float32, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	wordIDMap := make(map[string]float32)
	var word string
	var floatField float32

	SELECT := "SELECT word,probability FROM Vocabulary WHERE "
	if len(wordGrams) > 0 {
		SELECT = SELECT + dbx.GetWhereClause("word", wordGrams)
	} else {
		SELECT = SELECT + "word in (SELECT DISTINCT(word) FROM Occurrence WHERE " + dbx.GetSingleDateWhereClause("archivedate", timeInterval) + ")"
	}
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&word, &floatField)
		dbx.CheckErr(err)

		newWord, rule := cond.FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		wordIDMap[word] = floatField
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return wordIDMap, err
}

// GetTitleWordsBigramInterval func queries either [Occurrence] or [title] tables.
// This does NOT perform the word intersection by acmId!
func GetTitleWordsBigramInterval(bigrams []string, timeInterval nt.TimeInterval, useOccurrence bool) ([]hd.Occurrence, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(bigrams)
	if len(inPhrase) < 5 {
		inPhrase = "('')"
	}
	tableName := "Occurrence"
	if !useOccurrence {
		tableName = "Title"
	}
	SELECT := "SELECT acmId, archiveDate, word, nentry FROM " + tableName + " WHERE word IN " + inPhrase + " AND " + dbx.GetSingleDateWhereClause("archiveDate", timeInterval)

	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var acmID uint32
	var archiveDate nt.NullTime
	var word string
	var nentry int
	titleList := make([]hd.Occurrence, 0)

	for rows.Next() {
		err = rows.Scan(&acmID, &archiveDate, &word, &nentry)
		dbx.CheckErr(err)
		titleList = append(titleList, hd.Occurrence{AcmId: acmID, ArchiveDate: archiveDate, Word: word, Nentry: nentry})
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return titleList, nil
}

// UpdateVocabulary updates Vocabulary.RowCount, Frequency, SpeechPart, Stem for EVERY row!. AcmData table never needs updating.
func UpdateVocabulary(recordList []hd.Vocabulary) (int, error) {
	if len(recordList) == 0 {
		return 0, nil
	}

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	batch := &pgx.Batch{} // prepare batch updates
	for _, v := range recordList {
		sqlStatement := "UPDATE vocabulary SET RowCount=$2, Frequency=$3, SpeechPart=$4, Stem=$5 WHERE Word = $1;"
		batch.Queue(sqlStatement, v.Word, v.RowCount, v.Frequency, v.SpeechPart, v.Stem)
	}

	batchResults := txn.SendBatch(context.Background(), batch)

	var qerr error
	var rows pgx.Rows
	for qerr == nil {
		rows, qerr = batchResults.Query()
		rows.Close()
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	return len(recordList), nil
}

// GetVocabularyMap reads all Vocabulary.Word,{Id, Frequency, Wordrank} values. Applys filtering.
func GetVocabularyMap(fieldName string) (map[string]int, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	wordIDmap := make(map[string]int)
	var word string
	var intField int
	SELECT := "SELECT Word," + fieldName + " FROM vocabulary;" // WHERE word LIKE 'tech%'
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&word, &intField)
		dbx.CheckErr(err)

		newWord, rule := cond.FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		wordIDmap[word] = intField
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return wordIDmap, err
}

// BulkInsertVocabulary gets []Vocabulary from wordFrequencyList(). Does not insert Scores!
func BulkInsertVocabulary(recordList []hd.Vocabulary) (int, error) {
	if len(recordList) == 0 {
		return 0, nil
	}

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	vocabList := make([]hd.Vocabulary, 0)
	for _, rec := range recordList {
		_, rule := cond.FilteringRules(rec.Word)
		if rule >= 0 {
			vocabList = append(vocabList, rec)
		}
	}

	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"vocabulary"}, // tablename
		[]string{"word", "rowcount", "frequency", "wordrank", "probability", "speechpart", "stem"},
		pgx.CopyFromSlice(len(vocabList), func(i int) ([]interface{}, error) {
			return []interface{}{vocabList[i].Word, vocabList[i].RowCount, vocabList[i].Frequency, vocabList[i].WordRank, vocabList[i].Probability, vocabList[i].SpeechPart, vocabList[i].Stem}, nil
		}),
	)
	dbx.CheckErr(err)
	if copyCount == 0 {
		fmt.Println("BulkInsertVocabulary: no rows inserted")
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	return len(recordList), nil
}

// CallUpdateVocabulary invokes Postgresql UpdateVocabulary() which updates every Vocabulary.WordRank,Probability value.
func CallUpdateVocabulary() error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	_, err = db.Exec(context.Background(), "call UpdateVocabulary();")
	dbx.CheckErr(err)

	return err
}
