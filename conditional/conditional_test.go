package conditional

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	//cond "github.com/dgnabasik/acmsearchlib/conditional"
	dbx "github.com/dgnabasik/acmsearchlib/database"
	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
)

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
