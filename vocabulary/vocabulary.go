package vocabulary

//  vocabulary.go manages hd.Vocabulary in database.

import (
	"fmt"
	"log"
	"strings"
	"time"

	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

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

// GetWordList method returns all words if prefix is blank.
func GetWordList(prefix string) ([]string, error) {
	DB, err := dbx.GetDatabaseReference()
	defer DB.Close()

	query := "SELECT word FROM vocabulary ORDER BY word"
	if len(prefix) > 0 {
		query = "SELECT word FROM vocabulary WHERE word LIKE '" + strings.ToLower(prefix) + "%' ORDER BY word"
	}
	rows, err := DB.Query(query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetWordList(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var word string
	wordList := make([]string, 0)
	for rows.Next() {
		err := rows.Scan(&word)
		if err != nil {
			log.Printf("GetWordList(2): %+v\n", err)
			return nil, err
		}
		wordList = append(wordList, word)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return wordList, err
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

// GetVocabularyListByDate reads Vocabulary table filtered by articleList.
func GetVocabularyListByDate(timeinterval nt.TimeInterval) ([]hd.Vocabulary, error) {
	start := time.Now()
	fmt.Print("GetVocabularyListByDate() ")

	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT * FROM GetVocabularyByDate('" + timeinterval.StartDate.StandardDate() + "', '" + timeinterval.EndDate.StandardDate() + "')"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	// fields to read
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

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	elapsed := time.Since(start)
	fmt.Println(elapsed.String())

	return vocabList, err
}

// GetVocabularyMapProbability Read all Vocabulary.Word,Probability values if wordGrams is empty. Applys filtering.
func GetVocabularyMapProbability(wordGrams []string) (map[string]float32, error) {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	wordIDMap := make(map[string]float32)
	var word string
	var floatField float32

	SELECT := "SELECT Word,Probability FROM vocabulary"
	if len(wordGrams) > 0 {
		SELECT = SELECT + dbx.GetWhereClause("Word", wordGrams)
	} else {
		SELECT = SELECT + ";"
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

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return wordIDMap, err
}