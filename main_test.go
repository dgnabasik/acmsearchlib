package main

// go test -v.
import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
)

const (
	SEP = "|"
)

// t.Fatal

/*func TestIntMinTableDriven(t *testing.T) {
    var tests = []struct {
        a, b int
        want int
    }{
        {0, 1, 0},
        {1, 0, 0},
        {2, -2, -2},
        {0, -1, -1},
        {-1, 0, -1},
    }
    for _, tt := range tests {
        testname := fmt.Sprintf("%d,%d", tt.a, tt.b)
        t.Run(testname, func(t *testing.T) {  // t.Run enables running “subtests”, one for each table entry.
            ans := IntMin(tt.a, tt.b)
            if ans != tt.want {
                t.Errorf("got %d, want %d", ans, tt.want)
            }
        })
    }
} */
/* nulltime ************************************************************************************/

// Test_SupportFunctions func
func Test_SupportFunctions(t *testing.T) {
	fmt.Println("Test_SupportFunctions...")
	str := nt.GetShortMonthName(13)
	if str != "" {
		t.Error("a1:Expected '', got ", str)
	}
	str = nt.GetShortMonthName(12)
	if str != "dec" {
		t.Error("a2:Expected 'dec', got ", str)
	}

	data := make([]string, 0)
	ndx, found := nt.StringSliceContains(data, "dec")
	if ndx >= 0 || found {
		t.Error("a3:Expected -1, got ", ndx)
	}
	data = append(data, "dec")
	ndx, found = nt.StringSliceContains(data, "dec")
	if ndx < 0 || !found {
		t.Error("a4:Expected 0, got ", ndx)
	}

}

// Test_NullTime func
func Test_NullTime(t *testing.T) {
	fmt.Println("Test_NullTime...")

	ti := nt.New_NullTime("")
	if ti.StandardDate() != nt.NullDate {
		t.Error("b1: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime("1999-01-01") // Valid=true
	if ti.StandardDate() != nt.NullDate {
		t.Error("b2: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime("2000-13")
	if ti.StandardDate() != nt.NullDate {
		t.Error("b3: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime("2000-01-08")
	if ti.StandardDate() != nt.NullDate {
		t.Error("b4: Expected " + nt.NullDate)
	}

	ti = nt.New_NullTime(nt.VeryFirstDate)
	if ti.StandardDate() != nt.VeryFirstDate {
		t.Error("b5: Expected " + nt.VeryFirstDate)
	}

	ti = nt.New_NullTime1("September 23, 2019")
	if ti.StandardDate() != "2019-09-23" {
		t.Error("b6: Expected 2019-09-23")
	}

	/*ti = nt.NullTimeToday()
	if ti.StandardDate() != "2019-11-25" {
		t.Error("b7: Expected 2019-11-25")
	}*/

	ti = nt.New_NullTimeFromFileName("dec-05-2005.html")
	if ti.StandardDate() != "2005-12-05" {
		t.Error("b8: Expected 2005-12-05")
	}

	if ti.FileSystemDate() != "dec-05-2005" {
		t.Error("b9: Expected dec-05-2005")
	}

	if ti.HtmlArchiveDate() != "2005-12-dec" {
		t.Error("b10: Expected 2005-12-dec")
	}

	if ti.NonStandardDate() != "12/05/05" {
		t.Error("b11: Expected 12/05/05")
	}

	if nt.GetStandardDateForm("dec-05-2005") != "2005-12-05" { // Convert mmm-dd-yyyy to yyyy-mm-dd
		t.Error("b12: Expected 2005-12-05")
	}

	ti.AdvanceNextNullTime()
	if ti.StandardDate() != "2005-12-07" {
		t.Error("b13: Expected 2005-12-07")
	}

	nt1, nt2 := nt.GetStartEndOfWeek(ti)
	if nt1.StandardDate() != "2005-12-04" {
		t.Error("b14: Expected 2005-12-04")
	}
	if nt2.StandardDate() != "2005-12-10" {
		t.Error("b15: Expected 2005-12-10")
	}

	year, month, day, hour, min, sec := nt.NullTimeDiff(nt1, nt2)
	if year != 0 || month != 0 || day != 6 || hour != 0 || min != 0 || sec != 0 {
		t.Error("b16: Expected 6 days diff.")
	}

	// fmt.Printf("%d %d %d %d %d %d\n", year, month, day, hour, min, sec)

	ntx := nt.NullTimeToday()
	dateSet := make([]nt.NullTime, 0)
	dateSet = append(dateSet, ntx)
	dateSet = append(dateSet, ti)
	dateSet = nt.NullTimeSorter(dateSet)
	if len(dateSet) != 2 {
		t.Error("b19: Sorted dateSet != 2 ")
	}

	ntx = nt.New_NullTime("2019-12-01")
	fmt.Println("Testing with " + ntx.StandardDate())
	yes := ntx.IsScheduledDate(nt.TFUnknown)
	if !yes {
		t.Error("b20: Not past 11am.")
	}

	yes = ntx.IsScheduledDate(nt.TFWeek)
	if !yes {
		t.Error("b21: Not start of the week.")
	}

	yes = ntx.IsScheduledDate(nt.TFMonth)
	if !yes {
		t.Error("b22: Not start of the month.")
	}

	ntx = nt.New_NullTime("2020-01-01")
	fmt.Println("Testing with " + ntx.StandardDate())
	yes = ntx.IsScheduledDate(nt.TFQuarter)
	if !yes {
		t.Error("b23: Not start of the Quarter.")
	}

	yes = ntx.IsScheduledDate(nt.TFYear)
	if !yes {
		t.Error("b24: Not start of the year.")
	}

}

/*************************************************************************/

func Test_TimeInterval(t *testing.T) {
	fmt.Println("Test_TimeInterval...")

	// GetTimeFrameFromUnixTimeStamp (uts UnixTimeStamp, timeframetype TimeFrameType) TimeFrame {

	vfd1 := nt.New_NullTime(nt.VeryFirstDate)
	vfd2 := vfd1
	ti := nt.New_TimeInterval(nt.TFUnknown, vfd1, vfd2)
	if ti.StartDate.StandardDate() != ti.EndDate.StandardDate() {
		t.Error("c1: Expected " + ti.EndDate.StandardDate())
	}
	// add one time unitL len(list) always 2 except TFUnknown is 1.
	tiList := nt.GetTimeIntervalDatePartitionList(ti)
	if len(tiList) != 1 && tiList[0].StartDate.StandardDate() != tiList[0].EndDate.StandardDate() {
		t.Error("c2: Expected " + tiList[0].StartDate.StandardDate())
	}

	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(0, 0, 7)) // y,m,d
	ti = nt.New_TimeInterval(nt.TFWeek, vfd1, vfd2)
	tiList = nt.GetTimeIntervalDatePartitionList(ti)
	if len(tiList) != 2 {
		t.Error("c3: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec := nt.NullTimeDiff(vfd1, vfd2)
	if year != 0 || month != 0 || day != 7 || hour != 0 || min != 0 || sec != 0 {
		t.Error("c4: Expected 7 days diff.")
	}

	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(0, 1, 0))
	ti = nt.New_TimeInterval(nt.TFMonth, vfd1, vfd2)
	tiList = nt.GetTimeIntervalDatePartitionList(ti)
	if len(tiList) != 2 {
		t.Error("c5: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = nt.NullTimeDiff(vfd1, vfd2)
	if year != 0 || month != 1 || day != 0 || hour != 0 || min != 0 || sec != 0 {
		t.Error("c6: Expected 1 month diff.")
	}

	// skipped TFQuarter & TFSpan

	vfd2 = nt.New_NullTime2(vfd1.DT.AddDate(1, 0, 0))
	ti = nt.New_TimeInterval(nt.TFYear, vfd1, vfd2)
	tiList = nt.GetTimeIntervalDatePartitionList(ti)
	if len(tiList) != 2 {
		t.Error("c7: Expected 2 tiList items.")
	}
	year, month, day, hour, min, sec = nt.NullTimeDiff(vfd1, vfd2)
	if year != 1 || month != 0 || day != 0 || hour != 0 || min != 0 || sec != 0 {
		t.Error("c8: Expected 1 year diff.")
	}

}

func Test_TimeFrame(t *testing.T) {
	fmt.Println("Test_TimeFrame...")
	vfd1 := nt.New_NullTime(nt.VeryFirstDate) // "2000-01-03"
	vfd2 := nt.New_NullTime2(vfd1.DT.AddDate(1, 0, 0))
	ti := nt.New_TimeInterval(nt.TFYear, vfd1, vfd2)
	tf := nt.New_TimeFrame(ti)

	if tf.GivenDate.StandardDate() != "2000-01-03" {
		t.Error("d1: Expected 2000-01-03.")
	}

	if tf.StartOfWeek.StandardDate() != "2000-01-02" {
		t.Error("d2: Expected 2000-01-02.")
	}

	if tf.EndOfWeek.StandardDate() != "2000-01-08" {
		t.Error("d3: Expected 2000-01-08.")
	}

	if tf.StartOfMonth.StandardDate() != "2000-01-01" {
		t.Error("d4: Expected 2000-01-01.")
	}

	if tf.EndOfMonth.StandardDate() != "2000-01-31" {
		t.Error("d5: Expected 2000-01-31.")
	}

	if tf.StartOfQuarter.StandardDate() != "2000-01-01" {
		t.Error("d6: Expected 2000-01-01.")
	}

	if tf.EndOfQuarter.StandardDate() != "2000-03-31" {
		t.Error("d7: Expected 2000-03-31.")
	}

	if tf.StartOfYear.StandardDate() != "2000-01-01" {
		t.Error("d8: Expected 2000-01-01.")
	}

	if tf.EndOfYear.StandardDate() != "2000-12-31" {
		t.Error("d9: Expected 2000-12-31.")
	}

	if tf.StartOfSpan.StandardDate() != nt.VeryFirstDate {
		t.Error("d10: Expected " + nt.VeryFirstDate)
	}

	today := nt.NullTimeToday()
	if tf.EndOfSpan.StandardDate() != today.StandardDate() {
		t.Error("d11: Expected " + today.StandardDate())
	}

	year, month, day := today.DT.Date()
	tf.Timeframetype = nt.TFSpan
	divisor := float32((year-nt.VeryFirstYear)*148) + float32((int(month)-1)*12) + float32(day)/float32(3)
	if tf.GetDivisor() != divisor {
		t.Error("d12: Expected " + strconv.FormatFloat(float64(divisor), 'E', -1, 32))
	}

	vfd1, vfd2 = tf.GetTimeFrameDates()
	if vfd1.StandardDate() != nt.VeryFirstDate && vfd2.StandardDate() != today.StandardDate() {
		t.Error("d13: Expected " + today.StandardDate())
	}
}

func Test_TimeStamp(t *testing.T) {
	fmt.Println("Test_TimeStamp...")
	timeInUTC := time.Date(nt.VeryFirstYear, 1, 1, 1, 1, 1, 100, time.UTC)

	uts := nt.GetUnixTimeStampFromTime(timeInUTC)
	if uts != 946688461 {
		t.Error("e1: Expected 946688461.")
	}

	utsstr := nt.FormatUnixTimeStampAsString(uts)
	if utsstr != "946688461" {
		t.Error("e2: Expected 946688461.")
	}

	utsstr = nt.FormatUnixTimeStampAsTime(uts)
	if utsstr != "1999-12-31 18:01:01" {
		t.Error("e3: Expected 1999-12-31 18:01:01.")
	}

}

/* Conditional ************************************************************************************/

// BigramPresenceMap struct.
type BigramPresenceMap struct {
	Presence map[string]bool `json:"presence"`
}

// NewBigramPresenceMap func records which of the ORIGINAL word permutations exist or do not exist in that time interval.
func NewBigramPresenceMap(bigrams []string) BigramPresenceMap {
	bpm := new(BigramPresenceMap)
	bpm.Presence = make(map[string]bool)
	for _, bigram := range bigrams {
		bpm.Presence[bigram] = false
	}
	return *bpm
}

/*************************************************************************************************/

// Test_FilteringRules func
func Test_FilteringRules(t *testing.T) { // (word string) (string, int)
	fmt.Println("Test_FilteringRules...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_SelectOccurrenceByDate func
func Test_SelectOccurrenceByDate(t *testing.T) { // (occurrenceList []hd.Occurrence, timeinterval nt.TimeInterval) []hd.Occurrence
	fmt.Println("Test_SelectOccurrenceByDate...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_SelectOccurrenceByID func
func Test_SelectOccurrenceByID(t *testing.T) { // (occurrenceList []hd.Occurrence, acmID uint32) []hd.Occurrence
	fmt.Println("Test_SelectOccurrenceByID...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_SelectOccurrenceByWord func
func Test_SelectOccurrenceByWord(t *testing.T) { // (occurrenceList []hd.Occurrence, word string) []hd.Occurrence
	fmt.Println("Test_SelectOccurrenceByWord...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetOccurrenceListByDate func
func Test_GetOccurrenceListByDate(t *testing.T) { // (timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set, error)
	fmt.Println("Test_GetOccurrenceListByDate...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_CollectWordGrams func
func Test_CollectWordGrams(t *testing.T) { // (wordGrams []string, timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set)
	fmt.Println("Test_CollectWordGrams...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetOccurrencesByAcmid func
func Test_GetOccurrencesByAcmid(t *testing.T) { // (xacmid uint32) ([]hd.Occurrence, error)
	fmt.Println("Test_GetOccurrencesByAcmid...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_WordGramSubset func
func Test_WordGramSubset(t *testing.T) { // (alphaWord string, vocabList []hd.Vocabulary, occurrenceList []hd.Occurrence) []string
	fmt.Println("Test_WordGramSubset...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetDistinctDates func
func Test_GetDistinctDates(t *testing.T) { // (occurrenceList []hd.Occurrence) []nt.NullTime
	fmt.Println("Test_GetDistinctDates...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetDistinctWords func
func Test_GetDistinctWords(t *testing.T) { // (occurrenceList []hd.Occurrence) []string
	fmt.Println("Test_GetDistinctWords...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_BulkInsertConditionalProbability func
func Test_BulkInsertConditionalProbability(t *testing.T) { // (conditionals []hd.ConditionalProbability) error
	fmt.Println("Test_BulkInsertConditionalProbability...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_ExtractKeysFromProbabilityMap func
func Test_ExtractKeysFromProbabilityMap(t *testing.T) { // (wordMap map[string]float32) []string
	fmt.Println("Test_ExtractKeysFromProbabilityMap...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_CalcConditionalProbability func
func Test_CalcConditionalProbability(t *testing.T) { // (startingWordgram string, wordMap map[string]float32, timeinterval nt.TimeInterval) (int, error)
	fmt.Println("Test_CalcConditionalProbability...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// GetConditionalByTimeInterval func
func getConditionalByTimeInterval(bigrams []string, timeInterval nt.TimeInterval, bigramMap map[string]bool, includeTimeframetype bool) ([]hd.ConditionalProbability, error) {
	condProbList := make([]hd.ConditionalProbability, 0)
	db, err := dbx.GetDatabaseReference()
	if err != nil {
		return condProbList, err
	}
	defer db.Close()

	inPhrase := dbx.CompileInClause(bigrams)
	query := "SELECT id, wordlist, probability, timeframetype, startDate, endDate, firstDate, lastDate, pmi, dateUpdated FROM conditional WHERE wordlist IN " + inPhrase +
		" AND " + dbx.CompileDateClause(timeInterval, includeTimeframetype)
	rows, err := db.Query(context.Background(), query)
	dbx.CheckErr(err)
	if err != nil {
		log.Printf("GetConditionalByTimeInterval(1): %+v\n", err)
		return condProbList, err
	}
	defer rows.Close()

	var cProb hd.ConditionalProbability
	var timeframetype int
	var startDate time.Time
	var endDate time.Time

	for rows.Next() {
		err := rows.Scan(&cProb.Id, &cProb.WordList, &cProb.Probability, &timeframetype, &startDate, &endDate, &cProb.FirstDate, &cProb.LastDate, &cProb.Pmi, &cProb.DateUpdated)
		if err != nil {
			log.Printf("GetConditionalByTimeInterval(2): %+v\n", err)
			return condProbList, err
		}
		bigramMap[cProb.WordList] = true
		cProb.Timeinterval = nt.New_TimeInterval(nt.TimeFrameType(timeframetype), nt.New_NullTime2(startDate), nt.New_NullTime2(endDate))
		condProbList = append(condProbList, cProb)
	}

	// get any iteration errors
	err = rows.Err()
	dbx.CheckErr(err)
	return condProbList, err
}

func getWordBigramPermutations(words []string, permute bool) []string {
	result := make([]string, 0)
	if len(words) == 1 {
		return words
	}
	for i := 0; i < len(words); i++ {
		for j := i + 1; j < len(words); j++ {
			result = append(result, words[i]+SEP+words[j])
			if permute {
				result = append(result, words[j]+SEP+words[i])
			}
		}
	}
	return result
}

// Test_GetConditionalByTimeInterval func
func Test_GetConditionalByTimeInterval(t *testing.T) { // (bigrams []string, timeInterval nt.TimeInterval, bigramMap map[string]bool, includeTimeframetype bool) error
	fmt.Println("Test_GetConditionalByTimeInterval...")
	words := []string{"3d", "able", "access"}
	permutations := getWordBigramPermutations(words, true) // permute=true
	fmt.Print("    ")
	fmt.Println(permutations)
	gBigramPresenceMap := NewBigramPresenceMap(permutations)
	timeInterval := nt.New_TimeInterval(nt.TimeFrameType(nt.TFSpan), nt.New_NullTime("2008-01-01"), nt.New_NullTime("2011-12-31"))
	fmt.Print("    ")
	fmt.Println(timeInterval)
	condProbList, err := getConditionalByTimeInterval(permutations, timeInterval, gBigramPresenceMap.Presence, false) // includeTimeframetype=false
	n := len(condProbList)
	if err != nil {
		t.Error("Expected > 0, got ", n)
	}
	fmt.Print("    ")
	for _, c := range condProbList {
		fmt.Print(c.WordList + "  ")
	}
	fmt.Println()
}

// Test_GetConditionalByProbability func
func Test_GetConditionalByProbability(t *testing.T) { // (word string, probabilityCutoff float32, timeInterval nt.TimeInterval, condProbList *[]hd.ConditionalProbability) error
	fmt.Println("Test_GetConditionalByProbability...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetWordBigramPermutations func
func Test_GetWordBigramPermutations(t *testing.T) { // (words []string, permute bool) []string
	fmt.Println("Test_GetWordBigramPermutations...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetConditionalList func
func Test_GetConditionalList(t *testing.T) { // (words []string, timeInterval nt.TimeInterval, permute bool) ([]hd.ConditionalProbability, error)
	fmt.Println("Test_GetConditionalList...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetExistingConditionalBigrams func
func Test_GetExistingConditionalBigrams(t *testing.T) { // (bigrams []string, intervalClause string) ([]string, error)
	fmt.Println("Test_GetExistingConditionalBigrams...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetProbabilityGraph func
func Test_GetProbabilityGraph(t *testing.T) { // (words []string, timeInterval nt.TimeInterval) ([]hd.ConditionalProbability, error)
	fmt.Println("Test_GetProbabilityGraph...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// Test_GetWordgramConditionalsByInterval func
func Test_GetWordgramConditionalsByInterval(t *testing.T) { // (words []string, timeInterval nt.TimeInterval) ([]hd.WordScoreConditionalFlat, error)
	fmt.Println("Test_GetWordgramConditionalsByInterval...")
	incorrect := false
	if incorrect {
		t.Error("Expected 0, got ", 0)
	}
}

// All not-test files init() functions are executed first, then all test files init() functions are executed (hopefully in lexical order).
/* article ************************************************************************************/

/*
article/article.go:func GetArticleCount() (int, error) {
article/article.go:func GetLastDateSavedFromDb() (nt.NullTime, nt.NullTime, error) {
article/article.go:func CallUpdateOccurrence(timeinterval nt.TimeInterval) error {
article/article.go:func CallUpdateTitle(timeinterval nt.TimeInterval) error {
article/article.go:func GetAcmArticleListByArchiveDates(dateList []string) ([]hd.AcmArticle, error) {
article/article.go:func GetAcmArticleListByDate(timeinterval nt.TimeInterval) ([]hd.AcmArticle, error) {
article/article.go:func GetAcmArticlesByID(idMap map[uint32]int, cutoff int) ([]hd.AcmArticle, error) {
article/article.go:func WordFrequencyList() ([]hd.Vocabulary, error) {
article/article.go:func BulkInsertAcmData(articleList []hd.AcmArticle) (int, error) {

func getTableNames(useTempTable bool) []string {
func GetSimplexByNameUserID(userID int, simplexName, simplexType string, useTempTable bool) ([]hd.SimplexComplex, error) {
func GetSimplexListByUserID(userID int, useTempTable bool) ([]hd.SimplexComplex, error) {
func InsertSimplexComplex(sc hd.SimplexComplex) (hd.SimplexComplex, error) {
func BulkInsertSimplexFacets(facets []hd.SimplexFacet) error {
func PostSimplexComplex(userID int, simplexName, simplexType string, timeInterval nt.TimeInterval) ([]uint64, error) {
func GetSimplexWordDifference(complexIdList []uint64) ([]hd.KeyValueStringPair, error) {
func InsertCategoryWords(categoryID uint64, words []string) error {
func InsertWordCategory(description string) (hd.CategoryTable, error) {
func GetSpecialMap(category int) ([]hd.SpecialTable, error) {
func GetCategoryMap() ([]hd.CategoryTable, error) {

conditional/conditional.go:func isHexWord(word string) bool {
conditional/conditional.go:func FilteringRules(word string) (string, int) {
conditional/conditional.go:func SelectOccurrenceByDate(occurrenceList []hd.Occurrence, timeinterval nt.TimeInterval) []hd.Occurrence {
conditional/conditional.go:func SelectOccurrenceByID(occurrenceList []hd.Occurrence, acmID uint32) []hd.Occurrence {
conditional/conditional.go:func SelectOccurrenceByWord(occurrenceList []hd.Occurrence, word string) []hd.Occurrence {
conditional/conditional.go:func GetOccurrenceListByDate(timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set, error) {
conditional/conditional.go:func CollectWordGrams(wordGrams []string, timeinterval nt.TimeInterval) ([]hd.Occurrence, mapset.Set) {
conditional/conditional.go:func GetOccurrencesByAcmid(xacmid uint32) ([]hd.Occurrence, error) {
conditional/conditional.go:func WordGramSubset(alphaWord string, vocabList []hd.Vocabulary, occurrenceList []hd.Occurrence) []string {
conditional/conditional.go:func GetDistinctDates(occurrenceList []hd.Occurrence) []nt.NullTime {
conditional/conditional.go:func GetDistinctWords(occurrenceList []hd.Occurrence) []string {
conditional/conditional.go:func BulkInsertConditionalProbability(conditionals []hd.ConditionalProbability) error {
conditional/conditional.go:func ExtractKeysFromProbabilityMap(wordMap map[string]float32) []string {
conditional/conditional.go:func CalcConditionalProbability(startingWordgram string, wordMap map[string]float32, timeinterval nt.TimeInterval) (int, error) {
conditional/conditional.go:func GetConditionalByTimeInterval(bigrams []string, timeInterval nt.TimeInterval, bigramMap map[string]bool, includeTimeframetype bool) ([]hd.ConditionalProbability, error) {
conditional/conditional.go:func GetConditionalByProbability(word string, probabilityCutoff float32, timeInterval nt.TimeInterval, condProbList *[]hd.ConditionalProbability) error {
conditional/conditional.go:func GetWordBigramPermutations(words []string, permute bool) []string {
conditional/conditional.go:func GetConditionalList(words []string, timeInterval nt.TimeInterval, permute bool) ([]hd.ConditionalProbability, error) {
conditional/conditional.go:func GetExistingConditionalBigrams(bigrams []string, intervalClause string) ([]string, error) {
conditional/conditional.go:func GetProbabilityGraph(words []string, timeInterval nt.TimeInterval) ([]hd.ConditionalProbability, error) {
conditional/conditional.go:func GetWordgramConditionalsByInterval(words []string, timeInterval nt.TimeInterval, dimensions int) ([]hd.WordScoreConditionalFlat, error) {

	func CheckErr(err error) {
	func GetDatabaseConnectionString() string {
	func GetDatabaseReference() (*pgxpool.Pool, error) {
	func TestDbConnection(db *pgxpool.Pool) (*pgxpool.Pool, error) {
	func NoRowsReturned(err error) bool {
	func CompileInClause(words []string) string {
	func GetFormattedDatesForProcedure(timeInterval nt.TimeInterval) string {
	func GetWhereClause(columnName string, wordGrams []string) string {
	func GetSingleDateWhereClause(columnName string, timeInterval nt.TimeInterval) string {
	func CompileDateClause(timeInterval nt.TimeInterval, useTimeframetype bool) string {
	func FormatArrayForStorage(arr []int) []string {
	func (tw *timeWrapper) Scan(in interface{}) error {
	func GetFilePrefixPath() string {
	func ReadDir(dirname string) ([]os.FileInfo, error) {
	func CreateDirectory(dirPath string) error {
	func DeleteDirectory(dirPath string) error {
	func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	func ZipFiles(pathPrefix string, fileExt string, targetFileName string) error {
	func FileExists(filePath string) (bool, error) {
	func ReadFileIntoString(filePath string) (string, error) {
	func ReadTextLines(filePath string, normalizeText bool) ([]string, error) {
	func WriteTextLines(lines []string, filePath string, appendData bool) error {
	func ReadOccurrenceListFromCsvFile(filePath string) ([]hd.Occurrence, error) {
	func GetFileList(filePath string, since nt.NullTime) ([]string, error) {
	func GetMostRecentFileAsNullTime(dirname string) (nt.NullTime, error) {
	func GetFileTime(fileName string) int64 {
	func GetSourceDirectory() string {
	func (fss *FileService) GetTextFile(ctx *gin.Context) {
	func SearchForString(str string, substr string) (int, int) {
	func SearchForStringIndex(str string, substr string) (int, bool) {
	func StringSliceContains(a []string, x string) (int, bool) {
	func StringSetDifference(lines1 []string, lines2 []string) (diff []string) {
	func GetOrderedMap(fieldNames []string) map[int]string {
	func StartNextProgram(pgmName string, args []string) {
	func RemoveDuplicateStrings(stringSlice []string) []string {
	func DeleteStringSliceElement(a []string, str string) []string {
	func RandomHex(n int) string {
	func NewAcmComposite(lenWordList int) AcmComposite {
	func (ac *AcmComposite) UnmarshalJSON(data []byte) error {
	func (aa AcmArticle) Print() {
	func (aa AcmArticle) GetKeyValuePairs() (map[string]string, map[int]string) {
	func (v Vocabulary) GetKeyValuePairs() (map[string]string, map[int]string) {
	func (v Vocabulary) Print() string {
	func (a VocabularySorterFreq) Len() int           { return len(a) }
	func (a VocabularySorterFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a VocabularySorterFreq) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency } // want lowest frequencies first.
	func GetVocabularyItem(word string, vocabList []Vocabulary) int {
	func GetVocabularyItemIndex(word string, vocabList []Vocabulary) int {
	func ReplaceUnicodeCharacters(line string) string {
	func ReplaceSpecialCharacters(line string) string {
	func ReplaceProtected(line string) string {
	func (o Occurrence) Print() string {
	func (o Occurrence) GetKeyValuePairs() (map[string]string, map[int]string) {
	func (a OccurrenceSorterId) Len() int      { return len(a) }
	func (a OccurrenceSorterId) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
	func (a OccurrenceSorterId) Less(i, j int) bool {
	func (a OccurrenceSorterWord) Len() int           { return len(a) }
	func (a OccurrenceSorterWord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a OccurrenceSorterWord) Less(i, j int) bool { return strings.Compare(a[i].Word, a[j].Word) < 0 }
	func (a OccurrenceSorterDate) Len() int      { return len(a) }
	func (a OccurrenceSorterDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
	func (a OccurrenceSorterDate) Less(i, j int) bool {
	func (v WordScore) Print() string {
	func (v WordScore) GetKeyValuePairs() (map[string]string, map[int]string) {
	func New_OrderedArticleMap() OrderedArticleMap {
	func (om OrderedArticleMap) Iterator() func() (string, bool) {
	func (om OrderedArticleMap) FormatTitle(line string) string {
	func (om *OrderedArticleMap) Add(href string, title string) {
	func (om OrderedArticleMap) Get(key string) string {
	func (om OrderedArticleMap) PrintMap() {
	func (a WordScoreConditionalFlatSorter) Len() int           { return len(a) }
	func (a WordScoreConditionalFlatSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
	func (a WordScoreConditionalFlatSorter) Less(i, j int) bool { return a[i].ID < a[j].ID }
	func (a SimplexComplexSorterDate) Len() int      { return len(a) }
	func (a SimplexComplexSorterDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
	func (a SimplexComplexSorterDate) Less(i, j int) bool {

	func GetShortMonthName(m int) string {
	func indexOfMonth(element string, data []string) int {
	func StringSliceContains(a []string, x string) (int, bool) {
	func (tft TimeFrameType) ToString() string {
	func (tft TimeFrameType) ToStrings() []string {
	func TimeframeStrings() []string {
	func New_TimeInterval(timeframetype TimeFrameType, startDate NullTime, endDate NullTime) TimeInterval {
	func GetTimeIntervalDatePartitionList(baseTimeInterval TimeInterval) []TimeInterval {
	func (ti TimeInterval) ToString() string {
	func (tf TimeFrame) ToString() string {
	func (tf TimeFrame) GetDivisor() float32 {
	func (tf TimeFrame) GetTimeFrameDates() (NullTime, NullTime) {
	func (tf TimeFrame) Print() {
	func GetTimeFrameFromUnixTimeStamp(uts UnixTimeStamp, timeframetype TimeFrameType) TimeFrame {
	func (nt *NullTime) Scan(value interface{}) error {
	func (nt *NullTime) AdvanceNextNullTime() {
	func (nt NullTime) FileSystemDate() string {
	func (nt NullTime) HtmlArchiveDate() string {
	func (nt NullTime) StandardDate() string {
	func (nt NullTime) NonStandardDate() string {
	func (nt NullTime) Value() (driver.Value, error) {
	func (nt NullTime) IsScheduledDate(when TimeFrameType) bool {
	func GetStartEndOfWeek(givenDate NullTime) (NullTime, NullTime) {
	func New_TimeFrame(timeInterval TimeInterval) TimeFrame {
	func New_NullTime1(dt string) NullTime {
	func New_NullTime(dt string) NullTime {
	func New_NullTime2(dt time.Time) NullTime {
	func New_NullTimeFromFileName(htmlFile string) NullTime {
	func NullTimeDiff(startDate NullTime, endDate NullTime) (year, month, day, hour, min, sec int) {
	func NullTimeSorter(nullTimes []NullTime) []NullTime {
	func CurrentTimeString() string {
	func NullTimeToday() NullTime {
	func GetStandardDateForm(dt string) string {
	func GetCurrentTimeStamp() *timestamppb.Timestamp { // was *timestamp
	func GetUnixTimeStampFromTime(t time.Time) UnixTimeStamp {
	func GetTimeFromUnixTimeStamp(uts UnixTimeStamp) time.Time {
	func GetCurrentUnixTimeStamp() UnixTimeStamp {
	func FormatUnixTimeStampAsString(uts UnixTimeStamp) string {
	func FormatUnixTimeStampAsTime(uts UnixTimeStamp) string {

	func Encrypt(key, data []byte) ([]byte, error) {
	func Decrypt(key, data []byte) ([]byte, error) {
	func GenerateKey() ([]byte, error) {
	func DeriveKey(password, salt []byte) ([]byte, []byte, error) {
	func EncryptData(password, textdata string) string {
	func DecryptData(password string, ciphertext []byte) (string, error) {
	func GetUserProfile(userName, pwdText string) (hd.UserProfile, error) {
	func InsertUserProfile(userName, userEmail, pwdText string, acmmemberid int) (hd.UserProfile, error) {

	timestampinterval/timestampinterval.go:func GetTimeStampFromUnixTimeStamp(uts nt.UnixTimeStamp) *timestamp.Timestamp {
	timestampinterval/timestampinterval.go:func NewTimeEventRequest(topic string, pbtft MTimeStampInterval_MTimeFrameType) *TimeEventRequest {
	timestampinterval/timestampinterval.go:func NewTimeEventResponse() *TimeEventResponse {
	timestampinterval/timestampinterval.go:func NewTimeStampInterval(timeframetype MTimeStampInterval_MTimeFrameType, startTime nt.UnixTimeStamp, endTime nt.UnixTimeStamp) *MTimeStampInterval {

	func GetVocabularyByWord(wordX string) (hd.Vocabulary, error) {
	func GetVocabularyList(words []string) ([]hd.Vocabulary, error) {
	func getAcmGraphCount() string {
	func GetWordListMap(prefix string) ([]hd.LookupMap, error) {
	func GetVocabularyListByDate(timeinterval nt.TimeInterval) ([]hd.Vocabulary, error) {
	func GetVocabularyMapProbability(wordGrams []string, timeInterval nt.TimeInterval) (map[string]float32, error) {
	func GetTitleWordsBigramInterval(bigrams []string, timeInterval nt.TimeInterval, useOccurrence bool) ([]hd.Occurrence, error) {
	func UpdateVocabulary(recordList []hd.Vocabulary) (int, error) {
	func GetVocabularyMap(fieldName string) (map[string]int, error) {
	func BulkInsertVocabulary(recordList []hd.Vocabulary) (int, error) {
	func CallUpdateVocabulary() error {
	func GetLookupValues(tableName, columnName string) ([]string, error) {

	wordscore/wordscore.go:func GetWordScores(word string) ([]hd.WordScore, error) {
	wordscore/wordscore.go:func GetWordScoreListByTimeInterval(words []string, timeInterval nt.TimeInterval) ([]hd.WordScore, error) {
	wordscore/wordscore.go:func BulkInsertWordScores(wordScoreList []hd.WordScore) error {

*/
/*************************************************************************************/
/*
	conditional/condprob.pb.micro.go:func NewConditionalProbabilityEndpoints() []*api.Endpoint {
	conditional/condprob.pb.micro.go:func NewConditionalProbabilityService(name string, c client.Client) ConditionalProbabilityService {
	conditional/condprob.pb.micro.go:func (c *conditionalProbabilityService) CalcConditionalProbability(ctx context.Context, in *ConditionalProbabilityRequest, opts ...client.CallOption) (*ConditionalProbabilityResponse, error) {
	conditional/condprob.pb.micro.go:func RegisterConditionalProbabilityHandler(s server.Server, hdlr ConditionalProbabilityHandler, opts ...server.HandlerOption) error {
	conditional/condprob.pb.micro.go:func (h *conditionalProbabilityHandler) CalcConditionalProbability(ctx context.Context, in *ConditionalProbabilityRequest, out *ConditionalProbabilityResponse) error {

	timeevent/timeevent.micro.go:func NewTimeEventService(name string, c client.Client) TimeEventService {
	timeevent/timeevent.micro.go:func (c *timeEventService) CreateDay(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) CreateWeek(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) CreateMonth(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) CreateQuarter(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) CreateYear(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) CreateSpan(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) RecordEvent(ctx context.Context, in *TimeEventRequest, opts ...client.CallOption) (*TimeEventResponse, error) {
	timeevent/timeevent.micro.go:func (c *timeEventService) GetTimeEvents(ctx context.Context, in *GetTimeEventRequest, opts ...client.CallOption) (*GetTimeEventResponse, error) {
	timeevent/timeevent.micro.go:func RegisterTimeEventServiceHandler(s server.Server, hdlr TimeEventServiceHandler, opts ...server.HandlerOption) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) CreateDay(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) CreateWeek(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) CreateMonth(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) CreateQuarter(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) CreateYear(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) CreateSpan(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) RecordEvent(ctx context.Context, in *TimeEventRequest, out *TimeEventResponse) error {
	timeevent/timeevent.micro.go:func (h *timeEventServiceHandler) GetTimeEvents(ctx context.Context, in *GetTimeEventRequest, out *GetTimeEventResponse) error {

	webpage/webpage.micro.go:func NewWebpageService(name string, c client.Client) WebpageService {
	webpage/webpage.micro.go:func (c *webpageService) NewWebpage(ctx context.Context, in *WebpageRequest, opts ...client.CallOption) (*WebpageResponse, error) {
	webpage/webpage.micro.go:func RegisterWebpageServiceHandler(s server.Server, hdlr WebpageServiceHandler, opts ...server.HandlerOption) error {
	webpage/webpage.micro.go:func (h *webpageServiceHandler) NewWebpage(ctx context.Context, in *WebpageRequest, out *WebpageResponse) error {

	wordscore/wordscore.micro.go:func NewWordScoreService(name string, c client.Client) WordScoreServiceInterface {
	wordscore/wordscore.micro.go:func (c *WordScoreServiceStruct) GetWordScore(ctx context.Context, in *GetWordScoreRequest, opts ...client.CallOption) (*GetWordScoreResponse, error) {
	wordscore/wordscore.micro.go:func (c *WordScoreServiceStruct) CreateWordScore(ctx context.Context, in *CreateWordScoreRequest, opts ...client.CallOption) (*CreateWordScoreResponse, error) {
	wordscore/wordscore.micro.go:func RegisterWordScoreServiceHandler(s server.Server, hdlr WordScoreServiceHandler, opts ...server.HandlerOption) error {
	wordscore/wordscore.micro.go:func (h *wordScoreServiceHandler) GetWordScore(ctx context.Context, in *GetWordScoreRequest, out *GetWordScoreResponse) error {
	wordscore/wordscore.micro.go:func (h *wordScoreServiceHandler) CreateWordScore(ctx context.Context, in *CreateWordScoreRequest, out *CreateWordScoreResponse) error {
*/
