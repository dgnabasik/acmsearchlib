package vocabulary

//  vocabulary.go manages hd.Vocabulary in database.

import (
	"log"
	"os"
	"strings"

	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	"github.com/lib/pq"
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

	var word, speechPart string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	SELECT := "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart FROM Vocabulary WHERE Word='" + wordX + "'"
	err = db.QueryRow(SELECT).Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart)
	dbx.CheckErr(err)
	return hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart}, nil
}

// GetVocabularyList method does NOT apply filtering to imported []words. Func places single quotes around each words element.
func GetVocabularyList(words []string) ([]hd.Vocabulary, error) {
	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	inPhrase := dbx.CompileInClause(words)
	query := "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart FROM vocabulary WHERE word IN " + inPhrase
	rows, err := DB.Query(query)
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
			&vocabulary.SpeechPart)
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
	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	occurrenceCount := getAcmGraphCount()

	query := "SELECT id, word FROM vocabulary WHERE occurrenceCount > " + occurrenceCount + " ORDER BY word"
	if len(prefix) > 0 {
		query = "SELECT id, word FROM vocabulary WHERE occurrenceCount > " + occurrenceCount + " AND word LIKE '" + strings.ToLower(prefix) + "%' ORDER BY word"
	}
	rows, err := DB.Query(query)
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
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var word, speechPart string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	var vocabList []hd.Vocabulary

	for rows.Next() { // this order follows the \d Vocabulary description:
		err = rows.Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart)
		dbx.CheckErr(err)

		newWord, rule := cond.FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		vocabList = append(vocabList, hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart})
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
	rows, err := db.Query(SELECT)
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

	rows, err := db.Query(SELECT)
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

// UpdateVocabulary updates Vocabulary.RowCount, Frequency, SpeechPart for EVERY row!. AcmData table never needs updating.
func UpdateVocabulary(recordList []hd.Vocabulary) (int, error) {
	if len(recordList) == 0 {
		return 0, nil
	}

	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return -1, err
	}
	defer DB.Close()

	txn, err := DB.Begin()
	dbx.CheckErr(err)

	stmt, err := DB.Prepare("UPDATE vocabulary SET RowCount = $2, Frequency = $3, SpeechPart = $4 WHERE Word = $1;")
	dbx.CheckErr(err)

	for _, v := range recordList {
		_, err = stmt.Exec(v.Word, v.RowCount, v.Frequency, v.SpeechPart) // lastId, err := res.LastInsertId()
		dbx.CheckErr(err)
	}

	err = stmt.Close()
	dbx.CheckErr(err)

	err = txn.Commit()
	dbx.CheckErr(err)

	return len(recordList), nil
}

// GetVocabularyMap reads all Vocabulary.Word,{Id, Frequency, Wordrank} values. Applys filtering.
func GetVocabularyMap(fieldName string) (map[string]int, error) {
	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	wordIDmap := make(map[string]int)
	var word string
	var intField int
	SELECT := "SELECT Word," + fieldName + " FROM vocabulary;" // WHERE word LIKE 'tech%'
	rows, err := DB.Query(SELECT)
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

	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return -1, err
	}
	defer DB.Close()

	txn, err := DB.Begin()
	dbx.CheckErr(err)

	// tableName, field list (except Id)
	stmt, _ := txn.Prepare(pq.CopyIn("vocabulary", "word", "rowcount", "frequency", "wordrank", "probability", "speechpart"))
	for _, rec := range recordList {
		_, rule := cond.FilteringRules(rec.Word)
		if rule >= 0 {
			_, err := stmt.Exec(rec.Word, rec.RowCount, rec.Frequency, rec.WordRank, rec.Probability, rec.SpeechPart)
			dbx.CheckErr(err)
		}
	}

	_, err = stmt.Exec() // flush needed
	dbx.CheckErr(err)
	err = stmt.Close()
	dbx.CheckErr(err)
	err = txn.Commit()
	dbx.CheckErr(err)

	return len(recordList), nil
}

// CallUpdateVocabulary invokes Postgresql UpdateVocabulary() which updates every Vocabulary.WordRank,Probability value.
func CallUpdateVocabulary() error {
	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer DB.Close()

	_, err = DB.Exec("call UpdateVocabulary();")
	dbx.CheckErr(err)

	return err
}
