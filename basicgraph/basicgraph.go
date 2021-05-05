package basicgraph

// graph.go GraphStructure handles graph queries.
// Every hypergraph may be represented by a bipartite graph.
import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	"github.com/jackc/pgx/v4"
	"github.com/yourbasic/graph" // The vertices of all graphs in this package are numbered 0..n-1.
)

// GraphInterface interface functions are not placed into acmsearchlib.
type GraphInterface interface {
	BuildUndirectedGraph(condProbList []hd.ConditionalProbability) *graph.Mutable
}

// GraphStructure struct implements GraphInterface;
type GraphStructure struct {
	Nodes []hd.GraphNode
	Links []hd.GraphLink
}

/*************************************************************************************/

// BuildUndirectedGraph func
func (bgs *GraphStructure) BuildUndirectedGraph(gs GraphStructure) *graph.Mutable {
	g := graph.New(4)
	g.AddBoth(0, 1) //  0 -- 1
	g.AddBoth(0, 2) //  |    |
	g.AddBoth(2, 3) //  2 -- 3
	g.AddBoth(1, 3)

	// Visit all edges of a graph.
	for v := 0; v < g.Order(); v++ {
		g.Visit(v, func(w int, c int64) (skip bool) {
			// Visiting edge (v, w) of cost c.
			return
		})
	}

	// The immutable data structure created by Sort has an Iterator that returns neighbors in increasing order.
	// Visit the edges in order.
	for v := 0; v < g.Order(); v++ {
		graph.Sort(g).Visit(v, func(w int, c int64) (skip bool) {
			// Visiting edge (v, w) of cost c.
			return
		})
	}

	// The return value of an iterator function is used to break out of the iteration. Visit, in turn, returns a boolean indicating if it was aborted.
	// Skip the iteration at the first edge (v, w) with w equal to 3.
	for v := 0; v < g.Order(); v++ {
		aborted := graph.Sort(g).Visit(v, func(w int, c int64) (skip bool) {
			fmt.Println(v, w)
			if w == 3 {
				skip = true // Aborts the call to Visit.
			}
			return
		})
		if aborted {
			break
		}
	}

	return g
}

/*************************************************************************************/

// GetSimplexByNameUserID func : simplexName is case-insensitive.
func GetSimplexByNameUserID(simplexName string, userID int) (hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return hd.SimplexComplex{}, err
	}
	defer db.Close()

	query := `SELECT s.ID, s.UserID, s.SimplexName, s.SimplexType, s.EulerCharacteristic, s.Dimension, s.FiltrationValue, s.NumSimplices, s.BettiNumbers, 
		s.Timeframetype, s.StartDate, s.EndDate, s.Enabled, s.DateCreated, s.DateUpdated, f.ComplexID, f.SourceVertexID, f.TargetVertexID, f.SourceWord, 
		f.TargetWord, f.Weight FROM Simplex s RIGHT OUTER JOIN Facet f ON f.ComplexID=s.ID`
	query += " WHERE s.UserID=" + strconv.Itoa(userID) + " AND LOWER(s.SimplexName)='" + strings.ToLower(simplexName) + "'"

	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexByNameUserID(1): %+v\n", err)
		return hd.SimplexComplex{}, err
	}
	defer rows.Close()

	var s hd.SimplexComplex
	var f hd.SimplexFacet
	facets := make([]hd.SimplexFacet, 0)
	var timeframetype int
	var startDate, endDate time.Time
	// all simplex values are the same.
	for rows.Next() {
		err := rows.Scan(&s.ID, &s.UserID, &s.SimplexName, &s.SimplexType, &s.EulerCharacteristic, &s.Dimension, &s.FiltrationValue, &s.NumSimplices, &s.BettiNumbers,
			&timeframetype, &startDate, &endDate, &s.Enabled, &s.DateCreated, &s.DateUpdated, &f.ComplexID, &f.SourceVertexID, &f.TargetVertexID, &f.SourceWord, &f.TargetWord, &f.Weight)
		if err != nil {
			log.Printf("GetSimplexByNameUserID(2): %+v\n", err)
			return s, err
		}
		s.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		facets = append(facets, f)
	}
	err = rows.Err()
	dbx.CheckErr(err)

	s.FacetVector = make([]hd.SimplexFacet, len(facets))
	copy(s.FacetVector, facets) // (dst,src)
	return s, err
}

// GetSimplexListByUserID func gets all of a user's simplices.
func GetSimplexListByUserID(userID int) ([]hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return []hd.SimplexComplex{}, err
	}
	defer db.Close()

	query := `SELECT s.ID, s.UserID, s.SimplexName, s.SimplexType, s.EulerCharacteristic, s.Dimension, s.FiltrationValue, s.NumSimplices, s.BettiNumbers, 
		s.Timeframetype, s.StartDate, s.EndDate, s.Enabled, s.DateCreated, s.DateUpdated, f.ComplexID, f.SourceVertexID, f.TargetVertexID, f.SourceWord, 
		f.TargetWord, f.Weight FROM Simplex s RIGHT OUTER JOIN Facet f ON f.ComplexID=s.ID`
	query += " WHERE s.UserID=" + strconv.Itoa(userID) + " ORDER BY s.SimplexName, f.SourceVertexID"
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetSimplexListByUserID(1): %+v\n", err)
		return []hd.SimplexComplex{}, err
	}
	defer rows.Close()

	var s hd.SimplexComplex
	var f hd.SimplexFacet
	var oldID uint64 // 0
	complexes := make([]hd.SimplexComplex, 0)
	facets := make([]hd.SimplexFacet, 0)
	var timeframetype int
	var startDate, endDate time.Time

	for rows.Next() {
		err := rows.Scan(&s.ID, &s.UserID, &s.SimplexName, &s.SimplexType, &s.EulerCharacteristic, &s.Dimension, &s.FiltrationValue, &s.NumSimplices, &s.BettiNumbers,
			&timeframetype, &startDate, &endDate, &s.Enabled, &s.DateCreated, &s.DateUpdated, &f.ComplexID, &f.SourceVertexID, &f.TargetVertexID, &f.SourceWord, &f.TargetWord, &f.Weight)
		if err != nil {
			log.Printf("GetSimplexListByUserID(2): %+v\n", err)
			return complexes, err
		}
		s.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		if s.ID != oldID {
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

// InsertSimplexComplex func. Insert row before inserting []hd.SimplexFacet rows with InsertSimplexFacets().
// Assigns hd.SimplexComplex.ID = each hd.SimplexFacet.ComplexID.
func InsertSimplexComplex(sc hd.SimplexComplex) (hd.SimplexComplex, error) {
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return sc, err
	}
	defer db.Close()

	var id uint64 // BettiNumbers: '{0,1,0}'		DateCreated & DateUpdated use default server time.
	INSERT := `INSERT INTO Simplex (UserID, SimplexName, SimplexType, EulerCharacteristic, Dimension, FiltrationValue, NumSimplices, BettiNumbers, 
		Timeframetype, StartDate, EndDate, Enabled) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) returning id`
	bettiNumbers := dbx.FormatArrayForStorage(sc.BettiNumbers)

	err = db.QueryRow(context.Background(), INSERT, sc.UserID, sc.SimplexName, sc.SimplexType, sc.EulerCharacteristic, sc.Dimension, sc.FiltrationValue,
		sc.NumSimplices, bettiNumbers, sc.Timeinterval.Timeframetype, sc.Timeinterval.StartDate.DT, sc.Timeinterval.EndDate.DT, sc.Enabled).Scan(&id)
	dbx.CheckErr(err)

	sc.ID = id
	for ndx := 0; ndx < len(sc.FacetVector); ndx++ {
		sc.FacetVector[ndx].ComplexID = sc.ID
	}
	err = BulkInsertSimplexFacets(sc.FacetVector)

	return sc, err
}

// BulkInsertSimplexFacets func.
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
		pgx.Identifier{"facet"}, // tablename
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
