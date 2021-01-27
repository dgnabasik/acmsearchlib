package article

//  manages articles.

import (
	"log"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

func Version() string {
	return "1.0.10"
}

// GetArticleCount func
func GetArticleCount() int {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM AcmData").Scan(&count)
	dbx.CheckErr(err)
	return count
}

// GetLastDateSavedFromDb returns the earliest and latest AcmData.ArchiveDate values else default time.
func GetLastDateSavedFromDb() (nt.NullTime, nt.NullTime, error) {
	articleCount := GetArticleCount()
	if articleCount == 0 {
		return nt.New_NullTime(""), nt.New_NullTime(""), nil // default time.
	}

	db, err := dbx.GetDatabaseReference()
	defer db.Close()

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
