package vocabulary

//  vocabulary.go manages hd.Vocabulary in database.

import (
	"log"
	"strings"

	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

// Version func
func Version() string {
	return "1.0.10"
}

// GetVocabularyByWord func
func GetVocabularyByWord(wordX string) hd.Vocabulary {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	var word, speechPart string
	var id uint32
	var rowCount, frequency, wordRank int
	var probability float32
	SELECT := "SELECT id, word, rowcount, frequency, wordrank, probability, speechpart FROM Vocabulary WHERE Word='" + wordX + "'"
	err = db.QueryRow(SELECT).Scan(&id, &word, &rowCount, &frequency, &wordRank, &probability, &speechPart)
	dbx.CheckErr(err)
	return hd.Vocabulary{Id: id, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: wordRank, Probability: probability, SpeechPart: speechPart}
}

// GetVocabularyList method does NOT apply filtering to imported []words. Func places single quotes around each words element.
func GetVocabularyList(words []string) ([]hd.Vocabulary, error) {
	DB, err := dbx.GetDatabaseReference()
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

// GetWordListMap method returns all words if prefix is blank.
func GetWordListMap(prefix string) ([]hd.LookupMap, error) {
	DB, err := dbx.GetDatabaseReference()
	defer DB.Close()

	query := "SELECT id, word FROM vocabulary ORDER BY word"
	if len(prefix) > 0 {
		query = "SELECT id, word FROM vocabulary WHERE word LIKE '" + strings.ToLower(prefix) + "%' ORDER BY word"
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

// GetTitleWordsBigramInterval func queries [title] table which has same structure as [occurrence].
// This does NOT perform the word intersection by acmId!
func GetTitleWordsBigramInterval(bigrams []string, timeInterval nt.TimeInterval) ([]hd.Occurrence, error) {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	inPhrase := dbx.CompileInClause(bigrams)
	if len(inPhrase) < 5 {
		inPhrase = "('')"
	}

	SELECT := "SELECT acmId, archiveDate, word, nentry FROM Title WHERE word IN " + inPhrase + " AND " + dbx.GetSingleDateWhereClause("archiveDate", timeInterval)

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
