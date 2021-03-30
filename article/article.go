package article

//  manages articles.

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"

	cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	"github.com/lib/pq"
)

func Version() string {
	return "1.16.2"
}

// GetLastDateSavedFromDb returns the earliest and latest AcmData.ArchiveDate values else default time.
func GetLastDateSavedFromDb() (nt.NullTime, nt.NullTime, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nt.New_NullTime(""), nt.New_NullTime(""), err // default time.
	}
	defer db.Close()

	var articleCount int
	err = db.QueryRow("SELECT COUNT(*) FROM AcmData").Scan(&articleCount)
	dbx.CheckErr(err)
	if articleCount == 0 {
		return nt.New_NullTime(""), nt.New_NullTime(""), nil // default time.
	}

	var archiveDate1, archiveDate2 nt.NullTime // NullTime supports Scan() interface.

	err = db.QueryRow("SELECT MIN(ArchiveDate) FROM AcmData").Scan(&archiveDate1)
	dbx.CheckErr(err)

	err = db.QueryRow("SELECT MAX(ArchiveDate) FROM AcmData").Scan(&archiveDate2)
	dbx.CheckErr(err)

	return archiveDate1, archiveDate2, nil
}

// GetAcmArticleListByArchiveDates func
func GetAcmArticleListByArchiveDates(dateList []string) ([]hd.AcmArticle, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(dateList)
	query := "SELECT id, archivedate, articlenumber, title, imagesource, journalname, authorname, journaldate, webreference FROM acmdata WHERE archivedate IN " + inPhrase
	rows, err := db.Query(query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetAcmArticleListByArchiveDates(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var acmArticle hd.AcmArticle
	acmArticleList := []hd.AcmArticle{}
	for rows.Next() {
		err := rows.Scan(
			&acmArticle.Id,
			&acmArticle.ArchiveDate,
			&acmArticle.ArticleNumber,
			&acmArticle.Title,
			&acmArticle.ImageSource,
			&acmArticle.JournalName,
			&acmArticle.AuthorName,
			&acmArticle.JournalDate,
			&acmArticle.WebReference,
		)
		if err != nil {
			log.Printf("GetAcmArticleListByArchiveDates(2): %+v\n", err)
			return nil, err
		}
		acmArticleList = append(acmArticleList, acmArticle)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return acmArticleList, err
}

// GetAcmArticleListByDate PostgreSql allows for defining a generic get-all-rows stored proc and appending the WHERE clause to the select instead of defining the WHERE clause inside the stored proc, but it is slower.
func GetAcmArticleListByDate(timeinterval nt.TimeInterval) ([]hd.AcmArticle, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT * FROM GetAcmArticles() WHERE ArchiveDate >= '" + timeinterval.StartDate.StandardDate() + "' AND ArchiveDate <= '" + timeinterval.EndDate.StandardDate() + "'"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	// fields to read
	var id uint32
	var archiveDate, journalDate nt.NullTime
	var articleNumber, title, imageSource, journalName, authorName, webReference, summary string

	var articleList []hd.AcmArticle

	for rows.Next() { // this order follows the \d AcmData description:
		err = rows.Scan(&id, &archiveDate.DT, &articleNumber, &title, &imageSource, &journalName, &authorName, &journalDate.DT, &webReference, &summary)
		dbx.CheckErr(err)
		articleList = append(articleList, hd.AcmArticle{Id: id, ArchiveDate: archiveDate, ArticleNumber: articleNumber, Title: title, ImageSource: imageSource, JournalName: journalName, AuthorName: authorName, JournalDate: journalDate, WebReference: webReference, Summary: summary})
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return articleList, err
}

// GetAcmArticlesByID func
func GetAcmArticlesByID(idMap map[uint32]int, cutoff int) ([]hd.AcmArticle, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	//inPhrase := "(" + strings.Trim(strings.Join(strings.Fields(fmt.Sprint(idMap)), ","), "[]") + ")" // beautiful! but only works with arrays and not maps.
	intlist := make([]string, 0)
	for k, v := range idMap {
		if v >= cutoff {
			intlist = append(intlist, strconv.Itoa(int(k)))
		}
	}
	inPhrase := "(" + strings.Join(intlist, ",") + ")"
	if len(inPhrase) < 4 {
		inPhrase = "(0)"
	}
	SELECT := "SELECT * FROM AcmData WHERE id IN " + inPhrase
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var id uint32
	var archiveDate, journalDate nt.NullTime
	var articleNumber, title, imageSource, journalName, authorName, webReference, summary string
	var articleList []hd.AcmArticle

	for rows.Next() { // this order follows the \d AcmData description:
		err = rows.Scan(&id, &archiveDate.DT, &articleNumber, &title, &imageSource, &journalName, &authorName, &journalDate.DT, &webReference, &summary)
		dbx.CheckErr(err)
		articleList = append(articleList, hd.AcmArticle{Id: id, ArchiveDate: archiveDate, ArticleNumber: articleNumber, Title: title, ImageSource: imageSource, JournalName: journalName, AuthorName: authorName, JournalDate: journalDate, WebReference: webReference, Summary: summary})
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return articleList, err
}

// WordFrequencyList produces Word Frequency table using PostgreSql full-text methods. The exported wordfreq.txt file is currated to remove hex or decimal numbers, web references, etc, then bulk-loaded.
// Postgres full-text searching supports: Stemming, Ranking / Boost, Multiple languages, Fuzzy search for misspelling, Accents.
// 2 PostgreSql data types that support full-text search: tsvector-represents a document in a form optimized for text search; tsquery-represents a text query.
// A tsvector value is a sorted list of distinct lexemes which are words that have been normalized to make different variants of the same word look alike.
// A tsquery value stores lexemes that are to be searched for, and combines them honoring the Boolean operators & (AND), | (OR), and ! (NOT): SELECT to_tsvector('the impossible') @@ to_tsquery('impossible');
// SELECT * FROM ts_stat('SELECT to_tsvector(''simple_english'',summary) from acmdata') ORDER BY nentry DESC, ndoc DESC, word LIMIT 4096;
// The nested SELECT statement can be any select statement that yields a tsvector column, so you could substitute a function that applies the to_tsvector function to any number of text fields, and concatenates them into a single tsvector: SELECT * FROM ts_stat('SELECT to_tsvector(''simple_english'',title) || to_tsvector(''simple_english'',body) from acmdata') ORDER BY nentry DESC;
// http://www.postgresql.org/docs/current/static/textsearch.html
// https://www.w3resource.com/PostgreSQL/postgresql-text-search-function-and-operators.php
// https://www.postgresql.org/docs/9.6/functions-textsearch.html
// https://www.postgresql.org/docs/8.3/textsearch-features.html
// wordFrequencyList assigns Vocabulary.Word,RowCount,Frequency. This is the only reference to [AcmData].
func WordFrequencyList() ([]hd.Vocabulary, error) {
	DB, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer DB.Close()

	//Parameterized form: rows, err := DB.Query("SELECT id, first_name FROM acmdata LIMIT $1", 3)
	//psql: \copy (SELECT * FROM ts_stat('SELECT to_tsvector(''simple_english'',summary) from acmdata ') ORDER BY word, nentry DESC, ndoc DESC) to '/home/david/acm/processed.txt' with csv;
	SELECT := "SELECT * FROM ts_stat('SELECT to_tsvector(''simple_english'',summary) from acmdata') ORDER BY word, nentry DESC, ndoc DESC;"
	rows, err := DB.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	// Use map to record duplicates when found.
	var wordList []hd.Vocabulary
	wordMap := make(map[string]int)
	// fields to read: word | ndoc | nentry
	var word string
	var rowCount, frequency int

	fmt.Print("Duplicates:")
	for rows.Next() {
		err = rows.Scan(&word, &rowCount, &frequency)
		dbx.CheckErr(err)

		newWord, rule := cond.FilteringRules(word)
		if rule < 0 {
			continue
		} else if rule > 0 {
			word = newWord
		}

		if wordMap[word] > 0 { // have duplicate; find previous entry and modify.
			wordList[wordMap[word]-1].RowCount += rowCount
			wordList[wordMap[word]-1].Frequency += frequency
			fmt.Print(" " + word)
		} else { // new entry
			newVocabulary := hd.Vocabulary{Id: 0, Word: word, RowCount: rowCount, Frequency: frequency, WordRank: 0, Probability: 0, SpeechPart: " "} // id, word, rowCount, frequency, wordRank, probability, speechPart)
			wordList = append(wordList, newVocabulary)
			wordMap[word] = len(wordList)
		}

	}
	fmt.Println("")
	err = rows.Err()
	dbx.CheckErr(err)

	return wordList, nil
}

// BulkInsertAcmData includes query to retrieve new Id values and place them into articleList.
// nov-26-2003.html ==> pq: invalid byte sequence for encoding "UTF8": 0xe9 0x67 0xe9
func BulkInsertAcmData(articleList []hd.AcmArticle) (int, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return -1, err
	}
	defer db.Close()

	var maxID uint32 = 0
	sqlStatement := "SELECT MAX(id) FROM acmdata;"
	_ = db.QueryRow(sqlStatement).Scan(&maxID) // row

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	// tableName, field list (except Id)
	stmt, _ := txn.Prepare(pq.CopyIn("acmdata", "archivedate", "articlenumber", "title", "imagesource", "journalname", "authorname", "journaldate", "webreference", "summary"))
	for _, rec := range articleList {
		_, err := stmt.Exec(context.Background(), rec.ArchiveDate.DT, rec.ArticleNumber, rec.Title, rec.ImageSource, rec.JournalName, rec.AuthorName, rec.JournalDate.DT, rec.WebReference, rec.Summary)
		dbx.CheckErr(err)
	}

	_, err = stmt.Exec(context.Background()) // flush needed
	dbx.CheckErr(err)
	err = stmt.Close()
	dbx.CheckErr(err)
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	// update articleList with new Id values
	SELECT := "SELECT id FROM acmdata WHERE id > " + strconv.FormatUint(uint64(maxID), 10) + ";"
	var id, index uint32
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&id)
		dbx.CheckErr(err)
		articleList[index].Id = id
		index++
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)

	return len(articleList), nil
}
