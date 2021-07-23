package vocabulary

// vocabulary.go manages hd.Vocabulary in database.

import (
	"context"
	"log"
	"strings"
	"sync"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	pgx "github.com/jackc/pgx/v4"
)

const vocabularySelect = "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart, occurrencecount, stem, dateupdated FROM "

// Must be lowercase.
func getTableNameVocabulary(useVocabulary bool) string {
	if useVocabulary {
		return "vocabulary"
	}
	return "titlevocabulary"
}

// Must be lowercase.
func getTableNameOccurrence(useOccurrence bool) string {
	if useOccurrence {
		return "occurrence"
	}
	return "titleoccurrence"
}

// GetVocabularyByWord func
func GetVocabularyByWord(wordX string, useVocabulary bool) (hd.Vocabulary, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return hd.Vocabulary{}, err
	}
	defer db.Close()

	var word, speechPart, stem string
	var id uint32
	var rowCount, frequency, wordRank, occurrenceCount int
	var probability float32
	var dateupdated time.Time
	SELECT := vocabularySelect + getTableNameVocabulary(useVocabulary) + " WHERE word='" + wordX + "'"
	err = db.QueryRow(context.Background(), SELECT).Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart, &occurrenceCount, &stem, &dateupdated)
	dbx.CheckErr(err)
	if err != nil { // not found
		return hd.Vocabulary{}, err
	} else {
		return hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart, OccurrenceCount: occurrenceCount, Stem: stem, DateUpdated: dateupdated}, err
	}
}

// GetVocabularyList method does NOT apply filtering to imported []words. Func places single quotes around each words element.
func GetVocabularyList(words []string, useVocabulary bool) ([]hd.Vocabulary, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(words)
	query := vocabularySelect + getTableNameVocabulary(useVocabulary) + " WHERE word IN " + inPhrase
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
			&vocabulary.OccurrenceCount,
			&vocabulary.Stem,
			&vocabulary.DateUpdated)
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

// GetStemWords func returns all the words which have the stem of imported word, inclusive.
func GetStemWords(word string, useVocabulary bool) ([]hd.Vocabulary, error) {
	stemVoc, err := GetVocabularyByWord(strings.ToLower(word), useVocabulary)
	if err != nil {
		log.Printf("GetStemWords(0): %+v\n", err)
		return []hd.Vocabulary{}, err
	}

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return []hd.Vocabulary{}, err
	}
	defer db.Close()

	query := vocabularySelect + getTableNameVocabulary(useVocabulary) + " WHERE stem='" + stemVoc.Stem + "'"
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetStemWords(1): %+v\n", err)
		return []hd.Vocabulary{}, err
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
			&vocabulary.OccurrenceCount,
			&vocabulary.Stem,
			&vocabulary.DateUpdated)
		if err != nil {
			log.Printf("GetStemWords(2): %+v\n", err)
			return nil, err
		}
		vocabularyList = append(vocabularyList, vocabulary)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return vocabularyList, err
}

// GetStemWordList func returns all the words which have the stem of imported word, inclusive. CONCURRENT!
func GetStemWordList(queryWords []string, useVocabulary bool) ([]hd.Vocabulary, error) {
	queue := make(chan hd.Vocabulary) // Unbuffered synchronous channel
	fatalErrors := make(chan error)   // Make a channel to pass fatal errors in WaitGroup.
	var wg sync.WaitGroup

	for _, word := range queryWords {
		wg.Add(1)
		go func(word string, useVocabulary bool) {
			defer wg.Done()
			vocabulary, err := GetStemWords(word, useVocabulary)
			if err != nil {
				fatalErrors <- err
			}
			for _, voc := range vocabulary {
				queue <- voc
			}
		}(word, useVocabulary)
	}

	go func() {
		wg.Wait()
		close(queue)
	}()

	// Wait until either WaitGroup is done or fatal error is received through the channel.
	select {
	case <-queue:
		break // keep going
	case err := <-fatalErrors:
		close(fatalErrors)
		log.Printf("GetStemWordList error: %+v\n", err)
	}

	vocabularyList := []hd.Vocabulary{}
	for t := range queue {
		vocabularyList = append(vocabularyList, t)
	}

	return vocabularyList, nil
}

// GetWordListMap method returns all words.
func GetWordListMap(useVocabulary bool) ([]hd.LookupMap, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := "SELECT id, word FROM " + getTableNameVocabulary(useVocabulary) + " ORDER BY word"
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
func GetVocabularyListByDate(timeinterval nt.TimeInterval, useVocabulary bool) ([]hd.Vocabulary, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	functionCall := "GetVocabularyByDate"
	if !useVocabulary {
		functionCall = "GetTitleVocabularyByDate"
	}

	SELECT := "SELECT * FROM " + functionCall + dbx.GetFormattedDatesForProcedure(timeinterval)
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var word, speechPart, stem string
	var id uint32
	var rowCount, frequency, wordRank, occurrenceCount int
	var probability float32
	var dateUpdated time.Time
	var vocabList []hd.Vocabulary

	for rows.Next() { // this order follows the \d Vocabulary description:
		err = rows.Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart, &occurrenceCount, &stem, &dateUpdated)
		dbx.CheckErr(err)

		newWord, rule := hd.FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		vocabList = append(vocabList, hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart, OccurrenceCount: occurrenceCount, Stem: stem, DateUpdated: dateUpdated})
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return vocabList, err
}

// GetVocabularyMapProbability Read all Vocabulary.Word,Probability values if wordGrams is []string{}. Applys filtering.
// Vocabulary.probability is NOT used to calculate conditional probabilities!
func GetVocabularyMapProbability(wordGrams []string, timeInterval nt.TimeInterval, useVocabulary bool) (map[string]float32, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	wordIDMap := make(map[string]float32)
	var word string
	var floatField float32

	SELECT := "SELECT word,probability FROM " + getTableNameVocabulary(useVocabulary) + " WHERE "
	if len(wordGrams) > 0 {
		SELECT = SELECT + dbx.GetWhereClause("word", wordGrams)
	} else {
		SELECT = SELECT + "word in (SELECT DISTINCT(word) FROM " + getTableNameOccurrence(useVocabulary) + " WHERE " + dbx.GetSingleDateWhereClause("archivedate", timeInterval) + ")"
	}
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&word, &floatField)
		dbx.CheckErr(err)

		newWord, rule := hd.FilteringRules(word)
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

// GetTitleWordsBigramInterval func queries either [Occurrence] or [TitleOccurrence] tables.
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

	SELECT := "SELECT acmId, archiveDate, word, nentry FROM " + getTableNameOccurrence(useOccurrence) + " WHERE word IN " + inPhrase + " AND " + dbx.GetSingleDateWhereClause("archiveDate", timeInterval)
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
func UpdateVocabulary(recordList []hd.Vocabulary, useVocabulary bool) (int, error) {
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

	batch := &pgx.Batch{} // prepare batch updates. DB updates DateUpdated.
	for _, v := range recordList {
		sqlStatement := "UPDATE " + getTableNameVocabulary(useVocabulary) + " SET RowCount=$2, Frequency=$3, SpeechPart=$4, OccurrenceCount=$5, Stem=$6 WHERE Word = $1;"
		batch.Queue(sqlStatement, v.Word, v.RowCount, v.Frequency, v.SpeechPart, v.OccurrenceCount, v.Stem)
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
func GetVocabularyMap(fieldName string, useVocabulary bool) (map[string]int, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	wordIDmap := make(map[string]int)
	var word string
	var intField int
	SELECT := "SELECT word," + fieldName + " FROM " + getTableNameVocabulary(useVocabulary)
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&word, &intField)
		dbx.CheckErr(err)

		newWord, rule := hd.FilteringRules(word)
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
func BulkInsertVocabulary(recordList []hd.Vocabulary, useVocabulary bool) (int, error) {
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
		_, rule := hd.FilteringRules(rec.Word)
		if rule >= 0 {
			vocabList = append(vocabList, rec)
		}
	}

	tablename := getTableNameVocabulary(useVocabulary)
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{tablename},
		[]string{"word", "rowcount", "frequency", "wordrank", "probability", "speechpart", "occurrencecount", "stem"},
		pgx.CopyFromSlice(len(vocabList), func(i int) ([]interface{}, error) {
			return []interface{}{vocabList[i].Word, vocabList[i].RowCount, vocabList[i].Frequency, vocabList[i].WordRank, vocabList[i].Probability, vocabList[i].SpeechPart, vocabList[i].OccurrenceCount, vocabList[i].Stem}, nil
		}),
	)
	dbx.CheckErr(err)
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)
	if copyCount == 0 {
		log.Printf("BulkInsertVocabulary: no rows inserted")
	}

	return len(recordList), nil
}

// CallUpdateVocabulary invokes Postgresql UpdateVocabulary() which updates every Vocabulary.WordRank,Probability value.
func CallUpdateVocabulary(useVocabulary bool) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	invokeCall := "call UpdateVocabulary()"
	if !useVocabulary {
		invokeCall = "call UpdateTitleVocabulary()"
	}
	_, err = db.Exec(context.Background(), invokeCall)
	dbx.CheckErr(err)

	return err
}

/***********************************************************************************************/

// GetLookupValues func
func GetLookupValues(tableName, columnName string) ([]string, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return []string{}, err
	}
	defer db.Close()

	items := make([]string, 0)
	var item string
	SELECT := "SELECT itemValue FROM Lookup WHERE tableName='" + strings.ToLower(tableName) + "' AND columnName='" + strings.ToLower(columnName) + "' ORDER BY itemOrder"
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)

	for rows.Next() {
		err = rows.Scan(&item)
		dbx.CheckErr(err)
		items = append(items, item)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return items, nil
}
