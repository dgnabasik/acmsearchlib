package category

// category.go manages categories. Derived from ~/websites/dropdownlists/golang/listservice.go

import (
	"context"
	"log"
	"strconv"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	pgx "github.com/jackc/pgx/v4"
)

/*************************************************************************************/

// Version func
func Version() string {
	return "1.16.2"
}

// InsertCategoryWords func. 32k statement limit.
func InsertCategoryWords(categoryID uint64, words []string) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	// Must use lowercase column names!
	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"special"},    // tablename
		[]string{"word", "category"}, // dateupdated DEFAULT CURRENT_DATE
		pgx.CopyFromSlice(len(words), func(i int) ([]interface{}, error) {
			return []interface{}{words[i], categoryID}, nil
		}),
	)

	dbx.CheckErr(err)
	if copyCount == 0 {
		log.Printf("InsertCategoryWords: no rows inserted")
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	return nil
}

// InsertWordCategory func
func InsertWordCategory(description string) (hd.CategoryTable, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return hd.CategoryTable{}, err
	}
	defer db.Close()

	var id uint64
	INSERT := "INSERT INTO Wordcategory (description) VALUES ($1) returning id"
	err = db.QueryRow(context.Background(), INSERT, description).Scan(&id)
	dbx.CheckErr(err)

	categoryTable := hd.CategoryTable{Id: id, Description: description, DateUpdated: time.Now().UTC()}
	return categoryTable, nil
}

// GetSpecialMap func filters by category
func GetSpecialMap(category int) ([]hd.SpecialTable, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT id, word, category, dateupdated FROM Special WHERE category=" + strconv.Itoa(category) + " ORDER BY word"
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var special hd.SpecialTable
	specialMap := []hd.SpecialTable{}
	for rows.Next() {
		err = rows.Scan(&special.Id, &special.Word, &special.Category, &special.DateUpdated)
		dbx.CheckErr(err)
		specialMap = append(specialMap, special)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return specialMap, nil
}

// GetCategoryMap func
func GetCategoryMap() ([]hd.CategoryTable, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	SELECT := "SELECT id, description, dateupdated FROM Wordcategory"
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var catnap hd.CategoryTable
	categoryMap := []hd.CategoryTable{}
	for rows.Next() {
		err = rows.Scan(&catnap.Id, &catnap.Description, &catnap.DateUpdated)
		dbx.CheckErr(err)
		categoryMap = append(categoryMap, catnap)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return categoryMap, nil
}
