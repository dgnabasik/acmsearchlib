package article

//  manages articles.

import (
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	// comment
	_ "github.com/lib/pq"
)

func Version() string {
	return "1.0.10"
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
