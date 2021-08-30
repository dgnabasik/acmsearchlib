package headers

// Provide structs and their methods. These structs may or may not relect table schemas!

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

// constants
const (
	floatFormatter = "%.12f"
	dateFormatter  = "01-02-2006"
	Unknown        = "Unknown"
	HREF           = "<a href="
)

// SearchForString Searching/Filtering: Do not use strings.Contains unless you need exact matching rather than language-correct string searches!
// Example: start, end := SearchForString('foobar', 'bar')
func SearchForString(str string, substr string) (int, int) {
	m := search.New(language.English, search.IgnoreCase)
	return m.IndexString(str, substr)
}

// SearchForStringIndex Example: index, found := SearchForStringIndex('foobar', 'bar')
func SearchForStringIndex(str string, substr string) (int, bool) {
	m := search.New(language.English, search.IgnoreCase)
	start, _ := m.IndexString(str, substr)
	if start == -1 {
		return -1, false
	}
	return start, true
}

// StringSliceContains return index else -1
func StringSliceContains(a []string, x string) (int, bool) {
	for ndx, n := range a {
		if x == n {
			return ndx, true
		}
	}
	return -1, false
}

// StringSetDifference func returns the elements in lines1 but not in lines2, i.e., set difference of two arrays. Case-sensitive comparison.
func StringSetDifference(lines1 []string, lines2 []string) (diff []string) {
	m := make(map[string]bool)
	for _, item := range lines2 {
		m[item] = true
	}
	for _, item := range lines1 {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return // diff
}

// GetOrderedMap func
func GetOrderedMap(fieldNames []string) map[int]string {
	orderedMap := make(map[int]string, len(fieldNames))
	for ndx := 0; ndx < len(fieldNames); ndx++ {
		orderedMap[ndx] = fieldNames[ndx]
	}
	return orderedMap
}

// LookupMap struct handles unique key(string)-value(string) pairs, but must be returned to client as an array and not a map.
type LookupMap struct {
	Value int    `json:"value" binding:"required"`
	Label string `json:"label" binding:"required"`
}

// TimeEventService | WebpageService | ArticleService >>> WordScoreService | ConditionalService | GraphService
func StartNextProgram(pgmName string, args []string) {
	program := "../" + pgmName + "/" + pgmName + " "
	argument := strings.Join(args, " ")
	fmt.Print("Press Enter to execute: " + program + argument)
	os.Stdin.Read([]byte{0})
	cmd := exec.Command(program, argument)
	cmd.Start() // asynchronous
}

// RemoveDuplicateStrings func also removes empty strings.
func RemoveDuplicateStrings(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			if len(strings.TrimSpace(entry)) > 0 {
				keys[entry] = true
				list = append(list, entry)
			}
		}
	}
	return list
}

// DeleteStringSliceElement func maintains order. Removes only the first matching element. See RemoveDuplicateStrings() above.
func DeleteStringSliceElement(a []string, str string) []string {
	ndx := -1
	for i := 0; i < len(a); i++ {
		if str == a[i] {
			ndx = i
			break
		}
	}
	if ndx == 0 {
		return a[1:]
	}
	if ndx > 0 && ndx < len(a)-1 {
		return append(a[:ndx-1], a[ndx+1:]...)
	}
	if ndx > 0 && ndx == len(a)-1 {
		return append(a[:ndx-1], a[ndx-1])
	}
	return a
}

// RandomHex func returns max 128 bits. Returns lowercase unless error. Use n=12 for [aidata.keycode]
func RandomHex(n int) string {
	if n > 128 {
		n = 128
	} else if n <= 0 {
		n = 16
	}
	bytes := make([]byte, n)
	if _, err := rand.Read(bytes); err != nil {
		return "0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF0123456789ABCDEF"[0:n]
	}
	return hex.EncodeToString(bytes)
}

/*************************************************************************************************/

// AcmComposite struct. NOT USED. Retrieved separately by axios concurrent requests.
type AcmComposite struct { // filtered by startDate+endDate
	Vocabulary             []Vocabulary             `json:"vocabulary"`
	WordScore              []WordScore              `json:"wordscore"`
	ConditionalProbability []ConditionalProbability `json:"conditionalprobability"`
}

// NewAcmComposite funcAcmComposite
func NewAcmComposite(lenWordList int) AcmComposite {
	composite := new(AcmComposite)
	composite.Vocabulary = make([]Vocabulary, lenWordList)
	composite.WordScore = make([]WordScore, lenWordList)
	composite.ConditionalProbability = make([]ConditionalProbability, 0)
	return *composite
}

// UnmarshalJSON custom method for AcmComposite. Beautiful!
func (ac *AcmComposite) UnmarshalJSON(data []byte) error {
	array := [...]interface{}{&ac.Vocabulary, &ac.WordScore, &ac.ConditionalProbability}
	return json.Unmarshal(data, &array)
}

// Authorization struct
type Authorization struct {
	Auth0Domain       string `json:"auth0domain" binding:"required"`
	Auth0ClientID     string `json:"auth0clientid" binding:"required"`
	Auth0Audience     string `json:"auth0audience" binding:"required"`
	Auth0Callback     string `json:"auth0callback" binding:"required"`
	Auth0ClientSecret string `json:"auth0clientsecret" binding:"required"`
}

// AcmArticle struct Wrap the nullable cols in sql statements with a COALESCE(fieldName, '')
type AcmArticle struct {
	Id            uint32      `json:"id"`
	ArchiveDate   nt.NullTime `json:"archivedate"` // type date in db; nullable
	ArticleNumber string      `json:"articlenumber"`
	Title         string      `json:"title"`
	ImageSource   string      `json:"imagesource"`
	JournalName   string      `json:"journalname"`
	AuthorName    string      `json:"authorname"`
	JournalDate   nt.NullTime `json:"journaldate"` // type date in db; nullable
	WebReference  string      `json:"webreference"`
	Summary       string      `json:"summary"`
}

// Print func
func (aa AcmArticle) Print() {
	fmt.Println(aa.ArchiveDate.StandardDate() + " " + aa.ArticleNumber + ": " + aa.Title)
	fmt.Println(aa.JournalDate.StandardDate() + " " + aa.JournalName + ": " + aa.AuthorName)
	fmt.Println(aa.ImageSource)
	fmt.Println(aa.WebReference)
	fmt.Println(aa.Summary)
	fmt.Println("")
}

// GetKeyValuePairs Does not include for summary. 2nd map is used for ordering.
func (aa AcmArticle) GetKeyValuePairs() (map[string]string, map[int]string) {
	fieldNames := []string{"Id", "ArchiveDate", "ArticleNumber", "Title", "ImageSource", "JournalName", "AuthorName", "JournalDate", "WebReference"}
	orderedMap := GetOrderedMap(fieldNames)

	predicateMap := make(map[string]string, len(fieldNames))
	predicateMap[fieldNames[0]] = strconv.FormatUint(uint64(aa.Id), 10)
	predicateMap[fieldNames[1]] = aa.ArchiveDate.StandardDate()
	predicateMap[fieldNames[2]] = aa.ArticleNumber
	predicateMap[fieldNames[3]] = aa.Title
	predicateMap[fieldNames[4]] = aa.ImageSource
	predicateMap[fieldNames[5]] = aa.JournalName
	predicateMap[fieldNames[6]] = aa.AuthorName
	predicateMap[fieldNames[7]] = aa.JournalDate.StandardDate()
	predicateMap[fieldNames[8]] = aa.WebReference

	return predicateMap, orderedMap
}

/*************************************************************************************************/

// Vocabulary struct
type Vocabulary struct {
	Id              uint32    `json:"id"`
	Word            string    `json:"word"`
	RowCount        int       `json:"rowcount"`
	Frequency       int       `json:"frequency"`
	WordRank        int       `json:"wordrank"`
	Probability     float32   `json:"probability"` // Probability of word at rank.
	SpeechPart      string    `json:"speechpart"`  // Assign using BulkInsert_Vocabulary_Speechpart().
	OccurrenceCount int       `json:"occurrencecount"`
	Stem            string    `json:"stem"` // Assign using libstemmer program
	DateUpdated     time.Time `json:"dateupdated"`
}

// GetKeyValuePairs method
func (v Vocabulary) GetKeyValuePairs() (map[string]string, map[int]string) {
	fieldNames := []string{"Id", "Word", "RowCount", "Frequency", "WordRank", "Probability", "SpeechPart", "OccurrenceCount", "Stem", "DateUpdated"}
	orderedMap := GetOrderedMap(fieldNames)

	predicateMap := make(map[string]string, len(fieldNames))
	predicateMap[fieldNames[0]] = strconv.FormatUint(uint64(v.Id), 10)
	predicateMap[fieldNames[1]] = v.Word
	predicateMap[fieldNames[2]] = strconv.Itoa(v.RowCount)
	predicateMap[fieldNames[3]] = strconv.Itoa(v.Frequency)
	predicateMap[fieldNames[4]] = strconv.Itoa(v.WordRank)
	predicateMap[fieldNames[5]] = fmt.Sprintf(floatFormatter, v.Probability)
	predicateMap[fieldNames[6]] = v.SpeechPart
	predicateMap[fieldNames[7]] = strconv.Itoa(v.OccurrenceCount)
	predicateMap[fieldNames[8]] = v.Stem
	predicateMap[fieldNames[9]] = v.DateUpdated.Format(dateFormatter)

	return predicateMap, orderedMap
}

// Print method
func (v Vocabulary) Print() string {
	return fmt.Sprintf("%d : %s : %d : %d : %d : %f : %s : %d : %s : %s", v.Id, v.Word, v.RowCount, v.Frequency, v.WordRank, v.Probability, v.SpeechPart, v.OccurrenceCount, v.Stem, v.DateUpdated.Format(dateFormatter))
}

// VocabularySorterFreq Sort interface by Frequency. Len() is the number of elements in the collection.
type VocabularySorterFreq []Vocabulary

func (a VocabularySorterFreq) Len() int           { return len(a) }
func (a VocabularySorterFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VocabularySorterFreq) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency } // want lowest frequencies first.

type VocabularySorterWord []Vocabulary

func (a VocabularySorterWord) Len() int           { return len(a) }
func (a VocabularySorterWord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VocabularySorterWord) Less(i, j int) bool { return strings.Compare(a[i].Word, a[j].Word) < 0 }

// GetVocabularyItem Find existing item by Word, Return index else -1.
// for i := range myconfig{.} using a range loop on the index avoids copying the entire item.
func GetVocabularyItem(word string, vocabList []Vocabulary) int {
	for ndx := range vocabList {
		if word == vocabList[ndx].Word {
			return ndx
		}
	}
	return -1
}

// GetVocabularyItemIndex is concurrent version of GetVocabularyItem().
func GetVocabularyItemIndex(word string, vocabList []Vocabulary) int {
	numCPU := runtime.GOMAXPROCS(0)
	c := make(chan int, numCPU) // Buffering optional but sensible.
	lv := len(vocabList)
	var wg sync.WaitGroup
	ndx := -1

	for i := 0; i < numCPU; i++ {
		wg.Add(1)
		items := vocabList[i*lv/numCPU : (i+1)*lv/numCPU]
		go func(word string, items []Vocabulary, c chan int, i int) {
			defer wg.Done() // Decrement the counter when the goroutine completes.
			for j := range items {
				if word == items[j].Word {
					c <- j + i*lv/numCPU // + 1
					break
				}
			}
		}(word, items, c, i)
	}

	go func() {
		wg.Wait()
		close(c)
	}()

	for i := range c {
		if i >= 0 {
			ndx = i
		}
	}

	return ndx
}

/*************************************************************************************************/

// ReplaceSpecialCharacters for database storage. See https://www.starr.net/is/type/htmlcodes.html
func ReplaceSpecialCharacters(line string) string {
	r := strings.NewReplacer(
		"'", "",
		"\n", " ",
		"<a", "",
		"/", " ",
		"&ndash;", " ",
		"&mdash;", " ",
		"&shy;", " ",
		"&nbsp;", " ",
		"&rsquo;", "`",
		"&lsquo;", "`",
		"&ldquo;", "`",
		"&rdquo;", "`",
		"&sup1;", "`",
		"&#151;", " -- ",
		"&rdquo;", "",
		"&ldquo;", "",
		"&ecirc;", "ê",
		"&egrave;", "è",
		"&Eacute;", "É",
		"&eacute;", "é",
		"&aacute;", "á",
		"&oacute;", "ó",
		"&aring;", "å",
		"&szlig;", "ß",
		"&uuml;", "ü",
		"&auml;", "ä",
		"&euml;", "ë",
		"&ouml;", "ö",
		"&oslash;", "ø",
		"&hellip;", " ",
		"&amp;", "🙵",
		"&pound;", "£",
		"&euro;", "€",
		"&ntilde;", "ñ",
		",", "",
		";", "",
		"?", "",
		":", "",
		"!", "",
		"$", "",
		")", "",
		"(", "",
		"]", "",
		"[", "",
		"}", "",
		"{", "",
		">", "",
		"<", "",
		"`", "",
		"#", "",
		"*", "",
		"@", "",
		"--", "-",
	)
	result := r.Replace(line)
	return result
}

// ReplaceProtected func
func ReplaceProtected(line string) string {
	const PROTECTED = "protected]"
	result := line
	index1, found1 := SearchForStringIndex(strings.ToLower(result), HREF)
	index2, found2 := SearchForStringIndex(result, PROTECTED)
	if found1 && found2 {
		if index1 > 0 {
			result = line[0:index1-1] + " " + Unknown + line[index2+len(PROTECTED):]
		} else {
			result = Unknown + line[index2+len(PROTECTED):]
		}
	}
	return result
}

/*************************************************************************************************/

// Occurrence struct is summary based, not sentence based.
type Occurrence struct {
	AcmId       uint32      `json:"acmid"`
	ArchiveDate nt.NullTime `json:"archivedate"`
	Word        string      `json:"word"`
	Nentry      int         `json:"nentry"`
}

// Print method
func (o Occurrence) Print() string {
	return fmt.Sprintf("%d : %s : %s : %d", o.AcmId, o.Word, o.ArchiveDate.StandardDate(), o.Nentry)
}

// GetKeyValuePairs method
func (o Occurrence) GetKeyValuePairs() (map[string]string, map[int]string) {
	fieldNames := []string{"AcmId", "ArchiveDate", "Word", "Nentry"}
	orderedMap := GetOrderedMap(fieldNames)

	predicateMap := make(map[string]string, len(fieldNames))
	predicateMap[fieldNames[0]] = strconv.FormatUint(uint64(o.AcmId), 10)
	predicateMap[fieldNames[1]] = o.ArchiveDate.StandardDate()
	predicateMap[fieldNames[2]] = o.Word
	predicateMap[fieldNames[3]] = strconv.Itoa(o.Nentry)

	return predicateMap, orderedMap
}

// OccurrenceSorterId Sort interface by AcmId+ArchiveDate.
type OccurrenceSorterId []Occurrence

func (a OccurrenceSorterId) Len() int      { return len(a) }
func (a OccurrenceSorterId) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a OccurrenceSorterId) Less(i, j int) bool {
	return a[i].AcmId < a[j].AcmId && a[i].ArchiveDate.DT.Before(a[j].ArchiveDate.DT)
}

// OccurrenceSorterWord Sort interface by Word
type OccurrenceSorterWord []Occurrence

func (a OccurrenceSorterWord) Len() int           { return len(a) }
func (a OccurrenceSorterWord) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a OccurrenceSorterWord) Less(i, j int) bool { return strings.Compare(a[i].Word, a[j].Word) < 0 }

// OccurrenceSorterDate Sort interface by ArchiveDate+AcmId
type OccurrenceSorterDate []Occurrence

func (a OccurrenceSorterDate) Len() int      { return len(a) }
func (a OccurrenceSorterDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a OccurrenceSorterDate) Less(i, j int) bool {
	return a[i].ArchiveDate.DT.Before(a[j].ArchiveDate.DT) && a[i].AcmId < a[j].AcmId
}

/*************************************************************************************************/

// WordScore struct
type WordScore struct {
	Id           uint64          `json:"id"`
	Word         string          `json:"word"`
	Timeinterval nt.TimeInterval `json:"timeinterval"`
	Density      float32         `json:"density"`
	Linkage      float32         `json:"linkage"`
	Growth       float32         `json:"growth"`
	Score        float32         `json:"score"`
	Wstfidf      float32         `json:"wstfidf"`
}

// Print method
func (v WordScore) Print() string {
	return fmt.Sprintf("%s : %s : %f : %f : %f : %f : %f", v.Word, v.Timeinterval.ToString(), v.Density, v.Linkage, v.Growth, v.Score, v.Wstfidf)
}

// GetKeyValuePairs method
func (v WordScore) GetKeyValuePairs() (map[string]string, map[int]string) {
	fieldNames := []string{"Id", "Word", "Timeinterval", "Density", "Linkage", "Growth", "Score", "Wstfidf"}
	orderedMap := GetOrderedMap(fieldNames)

	predicateMap := make(map[string]string, len(fieldNames))
	predicateMap[fieldNames[0]] = strconv.FormatUint(v.Id, 10)
	predicateMap[fieldNames[1]] = v.Word
	predicateMap[fieldNames[2]] = v.Timeinterval.ToString()
	predicateMap[fieldNames[3]] = fmt.Sprintf(floatFormatter, v.Density)
	predicateMap[fieldNames[4]] = fmt.Sprintf(floatFormatter, v.Linkage)
	predicateMap[fieldNames[5]] = fmt.Sprintf(floatFormatter, v.Growth)
	predicateMap[fieldNames[6]] = fmt.Sprintf(floatFormatter, v.Score)
	predicateMap[fieldNames[7]] = fmt.Sprintf(floatFormatter, v.Wstfidf)

	return predicateMap, orderedMap
}

/*************************************************************************************************/

// OrderedArticleMap struct for ordering titles.
type OrderedArticleMap struct {
	articleMap   map[string]string
	articleNames []string
}

// New_OrderedArticleMap func
func New_OrderedArticleMap() OrderedArticleMap {
	p := new(OrderedArticleMap)
	p.articleMap = make(map[string]string)
	p.articleNames = make([]string, 0)
	return *p
}

// Iterator method returns the next articleName using closure iterator.
// Usage: iter := s.Iterator(); for i, ok := iter(); ok; i, ok = iter() {  }
func (om OrderedArticleMap) Iterator() func() (string, bool) {
	i := -1
	return func() (string, bool) {
		i++
		if i == len(om.articleNames) {
			return "", false
		}
		return om.articleNames[i], true
	}
}

// FormatTitle method
func (om OrderedArticleMap) FormatTitle(line string) string {
	result := ReplaceProtected(line)
	result = strings.ReplaceAll(result, "\"", "")
	result = strings.ReplaceAll(result, "%", " Percent")
	result = ReplaceSpecialCharacters(result)
	return result
}

// Add method: No need to order articleNames, but could. Modifies self.
func (om *OrderedArticleMap) Add(href string, title string) {
	om.articleMap[href] = om.FormatTitle(title)
	om.articleNames = append(om.articleNames, href)
}

// Get method
func (om OrderedArticleMap) Get(key string) string {
	return om.articleMap[key]
}

// PrintMap method
func (om OrderedArticleMap) PrintMap() {
	for _, key := range om.articleNames {
		fmt.Println(key + ": " + om.articleMap[key])
	}
	fmt.Println("")
}

/*************************************************************************************************/

// ConditionalProbability struct does NOT include the wordarray text[] column in [Conditional].
type ConditionalProbability struct {
	Id           uint64          `json:"id"`
	WordList     string          `json:"wordlist"`    // concatenated("|")
	Probability  float32         `json:"probability"` // Conditional Probability
	ReverseProb  float32         `json:"reverseprob"`
	Tfidf        float32         `json:"tfidf"`
	Timeinterval nt.TimeInterval `json:"timeinterval"` // declared in nulltime.go;
	//FirstDate    nt.NullTime     `json:"firstdate"`
	//LastDate     nt.NullTime     `json:"lastdate"`
	Pmi         float32   `json:"pmi"` // point mutual information.
	DateUpdated time.Time `json:"dateupdated"`
}

// SpecialTable struct
type SpecialTable struct {
	Id          uint64    `json:"id"`
	Word        string    `json:"word"`
	Category    int       `json:"category"`
	DateUpdated time.Time `json:"dateupdated"`
}

// CategoryTable struct
type CategoryTable struct {
	Id          uint64    `json:"id"`
	Description string    `json:"description"`
	DateUpdated time.Time `json:"dateupdated"`
}

// WordScoreConditionalFlat struct has extracted Timeinterval.
type WordScoreConditionalFlat struct {
	ID            int       `json:"id"` // negative values
	WordArray     []string  `json:"wordarray"`
	Wordlist      string    `json:"wordlist"`
	Score         float32   `json:"score"`
	Probability   float32   `json:"probability"`
	ReverseProb   float32   `json:"reverseprob"`
	Tfidf         float32   `json:"tfidf"`
	Pmi           float32   `json:"pmi"`
	Timeframetype int       `json:"timeframetype"`
	StartDate     time.Time `json:"startdate"`
	EndDate       time.Time `json:"enddate"`
	//FirstDate     time.Time `json:"firstdate"`
	//LastDate      time.Time `json:"lastdate"`
	Common bool `json:"common"` // intersection; not in database.
}

// WordScoreConditionalFlatSorter sort interface by ID.
type WordScoreConditionalFlatSorter []WordScoreConditionalFlat

func (a WordScoreConditionalFlatSorter) Len() int           { return len(a) }
func (a WordScoreConditionalFlatSorter) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a WordScoreConditionalFlatSorter) Less(i, j int) bool { return a[i].ID < a[j].ID }

// Print method
func (wscf WordScoreConditionalFlat) Print() string {
	ti := nt.New_TimeInterval(nt.TimeFrameType(wscf.Timeframetype), nt.New_NullTime2(wscf.StartDate), nt.New_NullTime2(wscf.EndDate))
	//fd := nt.New_NullTime2(wscf.FirstDate)
	//ld := nt.New_NullTime2(wscf.LastDate)	// , fd.StandardDate(), ld.StandardDate()	 : %s : %s
	str := fmt.Sprintf("%s : %f : %f : %f : %f : %f : %s", wscf.Wordlist, wscf.Score, wscf.Probability, wscf.ReverseProb, wscf.ReverseProb, wscf.Pmi, ti.ToString())
	return str
}

type UserProfile struct {
	ID          int       `json:"id"`
	UserName    string    `json:"username"`
	UserEmail   string    `json:"useremail"`
	Password    string    `json:"password"`
	AcmMemberId int       `json:"acmmemberid"`
	DateUpdated time.Time `json:"dateupdated"`
}

// GraphNode struct reflects IGraphNode interface in react-app-env.d.ts. Does not include OccurrenceCount or DateUpdated fields.
type GraphNode struct {
	NodeID       int             `json:"nodeid"`       // Vertices of all graphs are uniquely numbered 0..n-1.
	ID           int             `json:"id"`           // Vocabulary.Id	but this requires slight change in D3 to use nodeid instead of default id!
	Word         string          `json:"word"`         // Vocabulary
	RowCount     int             `json:"rowcount"`     // Vocabulary
	Frequency    int             `json:"frequency"`    // Vocabulary
	WordRank     int             `json:"wordrank"`     // Vocabulary
	Probability  float32         `json:"probability"`  // Vocabulary
	SpeechPart   string          `json:"speechpart"`   // Vocabulary
	Timeinterval nt.TimeInterval `json:"timeinterval"` // WordScore
	Density      float32         `json:"density"`      // WordScore
	Linkage      float32         `json:"linkage"`      // WordScore
	Growth       float32         `json:"growth"`       // WordScore
	Score        float32         `json:"score"`        // WordScore
	Wstfidf      float32         `json:"wstfidf"`      // WordScore
}

// GraphLink struct consolidates 2 ConditionalProbability objects.
type GraphLink struct {
	SourceNodeID int     `json:"source"` // json must be named 'source' to be D3-compatible.
	TargetNodeID int     `json:"target"` // json must be named 'target' to be D3-compatible.
	Level        int     `json:"level"`
	WordList1    string  `json:"wordlist1"` // concatenated("|")
	CondProb1    float32 `json:"condprob1"` // P(wordA|wordB)
	WordList2    string  `json:"wordlist2"`
	CondProb2    float32 `json:"condprob2"` // P(wordB|wordA)
	Pmi          float32 `json:"pmi"`       // point mutual information.
	Tfidf        float32 `json:"tfidf"`
	//FirstDate  time.Time `json:"firstdate"`
	//LastDate   time.Time `json:"lastdate"`
	SameDateList []time.Time `json:"samedatelist"`
}

// TitleSummary struct for display.
type TitleSummary struct {
	ID          int       `json:"id"`
	ArchiveDate time.Time `json:"archivedate"`
	Word        string    `json:"word"`
	Title       string    `json:"title"`
	Summary     string    `json:"summary"`
}

// KeyValuePairInterface interface for AcmArticle, Vocabulary, Occurrence, WordScore structs.
type KeyValuePairInterface interface {
	Print() string
	GetKeyValuePairs() (map[string]string, map[int]string)
}

// SimplexFacet struct.	Vertices of all graphs are uniquely numbered 0..n-1.  Undirected.
type SimplexFacet struct {
	ComplexID      uint64  `json:"complexid"` // FK to [SimplexComplex]
	SourceVertexID int     `json:"source"`    // json must be named 'source' to be D3-compatible.
	TargetVertexID int     `json:"target"`    // json must be named 'target' to be D3-compatible.
	SourceWord     string  `json:"sourceword"`
	TargetWord     string  `json:"targetword"`
	Weight         float32 `json:"weight"` // usually Pmi
}

// SimplexComplex struct implements ISimplexInterface. Reflects [Simplex] & [Facet] tables.
type SimplexComplex struct {
	ID                  uint64          `json:"id"`
	UserID              int             `json:"userid"`      // FK to [User] table; default 0.
	SimplexName         string          `json:"simplexname"` // assigned by user
	SimplexType         string          `json:"simplextype"` // could be chosen by user: {Rips, Čech, Alpha, Cubical, Hasse}
	EulerCharacteristic int             `json:"eulercharacteristic"`
	Dimension           int             `json:"dimension"`
	FiltrationValue     float32         `json:"filtrationvalue"`
	NumSimplices        int             `json:"numsimplices"`
	NumVertices         int             `json:"numvertices"`
	BettiNumbers        []int           `json:"bettinumbers"` // max(3)
	Timeinterval        nt.TimeInterval `json:"timeinterval"` // split out in db table.
	Enabled             int             `json:"enabled"`      // 0 is disabled, >0 is enabled.
	DateCreated         time.Time       `json:"datecreated"`  // server time
	DateUpdated         time.Time       `json:"dateupdated"`  // server time
	FacetVector         []SimplexFacet  `json:"facetvector"`
}

// SimplexComplexSorterDate Sort interface by StartDate.
type SimplexComplexSorterDate []SimplexComplex

func (a SimplexComplexSorterDate) Len() int      { return len(a) }
func (a SimplexComplexSorterDate) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a SimplexComplexSorterDate) Less(i, j int) bool {
	return a[i].Timeinterval.StartDate.DT.Before(a[j].Timeinterval.StartDate.DT)
}

// Print method
func (sc SimplexComplex) Print() string {
	str := fmt.Sprintf("%s : %s : %d : %d : %d : %d : %f : %s", sc.SimplexName, sc.SimplexType, sc.EulerCharacteristic, sc.Dimension, sc.NumSimplices, sc.NumVertices, sc.FiltrationValue, sc.Timeinterval.ToString())
	return str
}

// CreateSimplexComplex func does NOT assign ID; assigned when saved to database.
func CreateSimplexComplex(scName, scType string, facets []SimplexFacet, timeinterval nt.TimeInterval,
	userID, eulerCharacteristic, dimension, numSimplices, numVertices int, filtrationValue float32) SimplexComplex {

	complex := SimplexComplex{
		UserID:              userID,
		SimplexName:         scName,
		SimplexType:         scType,
		EulerCharacteristic: eulerCharacteristic,
		Dimension:           dimension,
		FiltrationValue:     filtrationValue,
		NumSimplices:        numSimplices,
		NumVertices:         numVertices,
		BettiNumbers:        make([]int, 3),
		FacetVector:         make([]SimplexFacet, len(facets)),
		Timeinterval:        timeinterval,
		Enabled:             1,
		DateCreated:         time.Now().UTC(),
		DateUpdated:         time.Now().UTC(),
	}
	copy(complex.FacetVector, facets) // (dst,src)
	return complex
}

// SimplexBarcode struct
type SimplexBarcode struct {
	ComplexID           uint64          `json:"complexid"` // FK to [SimplexComplex]
	ConnectedComponents int             `json:"connectedcomponents"`
	NumberHoles         int             `json:"numberholes"`
	ScaleParameter      float32         `json:"scaleparameter"` // neighborhood radius
	Timeinterval        nt.TimeInterval `json:"timeinterval"`   // this orders []SimplexBarcode
}

// KeyValueStringPair struct does not enforce unique Key.
type KeyValueStringPair struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}

// WordChanges struct
type WordChanges struct {
	QueryWords    []string  `json:"querywords"`
	SameWords     []string  `json:"samewords"`
	LossWords     []string  `json:"losswords"`
	GainWords     []string  `json:"gainwords"`
	Timeframetype int       `json:"timeframetype"`
	StartDate     time.Time `json:"startdate"`  // earlier period
	EndDate       time.Time `json:"enddate"`    // earlier period
	BeginDate     time.Time `json:"begindate"`  // later period
	FinishDate    time.Time `json:"finishdate"` // later period
	ChangeRate    float32   `json:"changerate"`
}

// CreateWordChangesStruct func where Rate of change = (Loss + Gain)/(Loss + Gain + Same)
func CreateWordChangesStruct(queryWords []string, kvsp []KeyValueStringPair, timeinterval nt.TimeInterval, begindate, finishdate time.Time) WordChanges {
	wc := WordChanges{Timeframetype: int(timeinterval.Timeframetype), StartDate: timeinterval.StartDate.DT, EndDate: timeinterval.EndDate.DT, BeginDate: begindate, FinishDate: finishdate}
	wc.QueryWords = make([]string, len(queryWords))
	copy(wc.QueryWords, queryWords)
	wc.SameWords = make([]string, 0)
	wc.GainWords = make([]string, 0)
	wc.LossWords = make([]string, 0)
	for _, kvp := range kvsp {
		switch kvp.Value {
		case "S":
			wc.SameWords = append(wc.SameWords, kvp.Key)
		case "G":
			wc.GainWords = append(wc.GainWords, kvp.Key)
		case "L":
			wc.LossWords = append(wc.LossWords, kvp.Key)
		}
	}
	if len(kvsp) > 0 {
		wc.ChangeRate = float32(len(wc.LossWords)+len(wc.GainWords)) / float32(len(wc.SameWords)+len(wc.LossWords)+len(wc.GainWords))
	}
	return wc
}

func isHexWord(word string) bool {
	_, err := hex.DecodeString(word)
	return len(word) >= 10 && err == nil
}

// FilteringRules filters output from Postgres ts_stat select. Include 3d prefixes.
// Return 0 for ok, -1 to completely ignore, 1 for modified word.
func FilteringRules(word string) (string, int) {
	if len(strings.TrimSpace(word)) <= 1 {
		return word, -1
	}

	ignore := strings.HasPrefix(word, "0") || strings.HasPrefix(word, "1") || strings.HasPrefix(word, "2") || (strings.HasPrefix(word, "3") && !strings.HasPrefix(word, "3d")) || strings.HasPrefix(word, "4") || strings.HasPrefix(word, "5") || strings.HasPrefix(word, "6") || strings.HasPrefix(word, "7") || strings.HasPrefix(word, "8") || strings.HasPrefix(word, "9") || strings.HasPrefix(word, "-") || strings.HasPrefix(word, "+") || strings.Count(word, "/") > 1 || strings.Count(word, "_") > 1 || strings.HasPrefix(word, "www.") || strings.HasSuffix(word, ".com") || strings.HasSuffix(word, ".org") || isHexWord(word)
	if ignore {
		return word, -1
	}

	newWord := word // Remove leading/trailing . /
	if strings.HasPrefix(newWord, ".") || strings.HasPrefix(newWord, "/") || strings.HasPrefix(newWord, "`") || strings.HasPrefix(newWord, "€") || strings.HasPrefix(newWord, "£") {
		newWord = newWord[1:]
	}
	if strings.HasSuffix(newWord, ".") || strings.HasSuffix(newWord, "/") || strings.HasSuffix(newWord, "`") || strings.HasSuffix(newWord, ";") {
		newWord = newWord[:len(newWord)-1]
	}

	if newWord != word {
		return newWord, 1
	}

	return word, 0
}
