package basicgraph

// graph.go GraphStructure handles graph queries.
// Every hypergraph may be represented by a bipartite graph.
import (
	"context"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	"github.com/jackc/pgx/v4" // The vertices of all graphs in this package are numbered 0..n-1.
)

/*************************************************************************************/

func getTableNames(useTempTable bool) []string {
	tableIndex := 0
	if useTempTable {
		tableIndex = 1
	}
	return []string{[]string{"Simplex", "temp_Simplex"}[tableIndex], []string{"Facet", "temp_Facet"}[tableIndex]}
}

// GetSimplexByNameUserID func : simplexName is case-insensitive.
// Get linked [Simplex] rows using constant {UserID, SimplexName, SimplexType} with varying {Timeinterval}.
func GetSimplexByNameUserID(userID int, simplexName, simplexType string, useTempTable bool) ([]hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	tableNames := getTableNames(useTempTable)
	query := `SELECT s.ID, s.UserID, s.SimplexName, s.SimplexType, s.EulerCharacteristic, s.Dimension, s.FiltrationValue, s.NumSimplices, s.NumVertices, s.BettiNumbers, s.Timeframetype, 
		s.StartDate, s.EndDate, s.Enabled, s.DateCreated, s.DateUpdated, f.ComplexID, f.SourceVertexID, f.TargetVertexID, f.SourceWord, f.TargetWord, f.Weight FROM `
	query += tableNames[0] + " s RIGHT OUTER JOIN " + tableNames[1] + " f ON f.ComplexID=s.ID WHERE s.UserID=" + strconv.Itoa(userID) + " AND LOWER(s.SimplexName)='" +
		strings.ToLower(simplexName) + "' AND s.SimplexType='" + simplexType + "' ORDER BY s.ID, f.SourceVertexID"
	// query returns as many rows per simplex as there are facets.
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexByNameUserID(1): %+v\n", err)
		return nil, err
	}
	defer rows.Close()

	var s hd.SimplexComplex
	var f hd.SimplexFacet
	complexes := make([]hd.SimplexComplex, 0)
	facets := make([]hd.SimplexFacet, 0)
	var timeframetype int
	var startDate, endDate time.Time
	var oldID uint64
	// all simplex values are the same.
	for rows.Next() {
		err := rows.Scan(&s.ID, &s.UserID, &s.SimplexName, &s.SimplexType, &s.EulerCharacteristic, &s.Dimension, &s.FiltrationValue, &s.NumSimplices, &s.NumVertices, &s.BettiNumbers,
			&timeframetype, &startDate, &endDate, &s.Enabled, &s.DateCreated, &s.DateUpdated, &f.ComplexID, &f.SourceVertexID, &f.TargetVertexID, &f.SourceWord, &f.TargetWord, &f.Weight)
		if err != nil {
			log.Printf("GetSimplexByNameUserID(2): %+v\n", err)
			return complexes, err
		}
		s.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		if oldID != s.ID {
			complexes = append(complexes, s)
			oldID = s.ID
		}
		facets = append(facets, f)
	}
	err = rows.Err()
	dbx.CheckErr(err)

	// partition facets into correct simplex
	for ndx, s := range complexes {
		complexes[ndx].FacetVector = make([]hd.SimplexFacet, 0)
		for _, f := range facets {
			if f.ComplexID == s.ID {
				complexes[ndx].FacetVector = append(complexes[ndx].FacetVector, f)
			}
		}
	}

	return complexes, err
}

// GetSimplexListByUserID func fetches all of a user's simplices but not the facets.
func GetSimplexListByUserID(userID int, useTempTable bool) ([]hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return []hd.SimplexComplex{}, err
	}
	defer db.Close()

	tableNames := getTableNames(useTempTable)
	query := `SELECT s.ID, s.UserID, s.SimplexName, s.SimplexType, s.EulerCharacteristic, s.Dimension, s.FiltrationValue, s.NumSimplices, s.NumVertices, s.BettiNumbers, s.Timeframetype, 
		s.StartDate, s.EndDate, s.Enabled, s.DateCreated, s.DateUpdated FROM `
	query += tableNames[0] + " s WHERE s.UserID=" + strconv.Itoa(userID) + " ORDER BY s.ID"

	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexListByUserID(1): %+v\n", err)
		return []hd.SimplexComplex{}, err
	}
	defer rows.Close()

	complexes := make([]hd.SimplexComplex, 0)
	var s hd.SimplexComplex
	var timeframetype int
	var startDate, endDate time.Time

	for rows.Next() {
		err := rows.Scan(&s.ID, &s.UserID, &s.SimplexName, &s.SimplexType, &s.EulerCharacteristic, &s.Dimension, &s.FiltrationValue, &s.NumSimplices, &s.NumVertices, &s.BettiNumbers,
			&timeframetype, &startDate, &endDate, &s.Enabled, &s.DateCreated, &s.DateUpdated)
		if err != nil {
			log.Printf("GetSimplexListByUserID(2): %+v\n", err)
			return complexes, err
		}
		s.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		complexes = append(complexes, s)
	}
	err = rows.Err()
	dbx.CheckErr(err)

	return complexes, err
}

// InsertSimplexComplex func. Insert row into temp_Simplex before inserting []hd.SimplexFacet rows with InsertSimplexFacets().
// Assigns hd.SimplexComplex.ID = each hd.SimplexFacet.ComplexID. temp_Simplex rows are NOT in StartDate order!
func InsertSimplexComplex(sc hd.SimplexComplex) (hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return sc, err
	}
	defer db.Close()

	var id uint64 // DateCreated & DateUpdated use default server time.
	INSERT := `INSERT INTO temp_Simplex (UserID, SimplexName, SimplexType, EulerCharacteristic, Dimension, FiltrationValue, NumSimplices, NumVertices, 
		BettiNumbers, Timeframetype, StartDate, EndDate, Enabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) returning id`

	err = db.QueryRow(context.Background(), INSERT, sc.UserID, sc.SimplexName, sc.SimplexType, sc.EulerCharacteristic, sc.Dimension, sc.FiltrationValue,
		sc.NumSimplices, sc.NumVertices, sc.BettiNumbers, sc.Timeinterval.Timeframetype, sc.Timeinterval.StartDate.DT, sc.Timeinterval.EndDate.DT, sc.Enabled).Scan(&id)
	dbx.CheckErr(err)

	sc.ID = id
	for ndx := 0; ndx < len(sc.FacetVector); ndx++ {
		sc.FacetVector[ndx].ComplexID = sc.ID
	}
	err = BulkInsertSimplexFacets(sc.FacetVector)

	return sc, err
}

// BulkInsertSimplexFacets func inserts into temp_Facet
func BulkInsertSimplexFacets(facets []hd.SimplexFacet) error {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return err
	}
	defer db.Close()

	txn, err := db.Begin(context.Background())
	dbx.CheckErr(err)

	copyCount, err := db.CopyFrom(
		context.Background(),
		pgx.Identifier{"temp_facet"}, // tablename
		[]string{"complexid", "sourcevertexid", "targetvertexid", "sourceword", "targetword", "weight"}, // Must use lowercase column names!
		pgx.CopyFromSlice(len(facets), func(i int) ([]interface{}, error) {
			return []interface{}{facets[i].ComplexID, facets[i].SourceVertexID, facets[i].TargetVertexID, facets[i].SourceWord, facets[i].TargetWord, facets[i].Weight}, nil
		}),
	)

	dbx.CheckErr(err)
	if copyCount == 0 {
		log.Printf("BulkInsertSimplexFacets(1): %+v\n", err)
		fmt.Println("BulkInsertSimplexFacets: no rows inserted")
	}
	err = txn.Commit(context.Background())
	dbx.CheckErr(err)

	return nil
}

// PostSimplexComplex func moves [temp_Simplex] & [temp_Facet] data into [Simplex] & [Facet] tables. Returns new [Simplex].ID values!
func PostSimplexComplex(userID int, simplexName, simplexType string, timeInterval nt.TimeInterval) ([]uint64, error) {
	simplexList, err := GetSimplexByNameUserID(userID, simplexName, simplexType, true) // useTempTable
	dbx.CheckErr(err)

	// Ensure simplexIDs is called in StartDate order.
	sort.Sort(hd.SimplexComplexSorterDate(simplexList))
	simplexIDmap := make(map[uint64]int)
	for ndx := range simplexList {
		simplexIDmap[simplexList[ndx].ID]++
	}
	simplexIDs := make([]uint64, 0, len(simplexIDmap))
	for k := range simplexIDmap {
		simplexIDs = append(simplexIDs, k)
	}

	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var newSimplexID uint64
	idList := make([]uint64, 0)
	for ndx := range simplexIDs { // Use Exec to execute a query that does not return a result set.
		rows, err := db.Query(context.Background(), "SELECT PostSimplexComplex("+strconv.FormatUint(simplexIDs[ndx], 10)+")")
		dbx.CheckErr(err)

		for rows.Next() {
			err = rows.Scan(&newSimplexID)
			dbx.CheckErr(err)
			idList = append(idList, newSimplexID)
		}
		rows.Close()
	}
	//defer rows.Close()

	return idList, nil
}

// GetSimplexWordDifference func returns words that are the same, gained, and lost between two SimplexComplex-Facet sets. Format: word|type={S,G,L}
// CREATE TABLE acmsearch.Word_type (sourceword character varying(32), wordtype char(1) );
func GetSimplexWordDifference(complexid1, complexid2 uint64) ([]hd.KeyValueStringPair, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// PostgreSQL functions invoked with SELECT; stored procs invoked with CALL.
	SELECT := "SELECT WordDifference(" + strconv.FormatUint(complexid1, 10) + "," + strconv.FormatUint(complexid2, 10) + ")"
	rows, err := db.Query(context.Background(), SELECT)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexWordDifference(1): %+v\n", err)
		return []hd.KeyValueStringPair{}, err
	}
	defer rows.Close()

	var str string
	list := make([]hd.KeyValueStringPair, 0)
	for rows.Next() {
		err = rows.Scan(&str) // (stone,G)
		dbx.CheckErr(err)
		index := strings.Index(str, ",")
		list = append(list, hd.KeyValueStringPair{Key: str[1:index], Value: str[index+1 : index+2]})
	}

	return list, nil
}
