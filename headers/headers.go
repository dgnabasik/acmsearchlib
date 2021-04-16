package headers

// Provide structs and their methods. These structs may or may not relect table schemas!

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	nt "github.com/dgnabasik/acmsearchlib/nulltime"

	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

func Version() string {
	return "1.16.2"
}

// constants
const (
	floatFormatter = "%.12f"
	Unknown        = "Unknown"
	HREF           = "<a href="
)

// SearchForString Searching/Filtering: Do not use strings.Contains unless you need exact matching rather than language-correct string searches!
// Example: start, end := SearchForString('foobar', 'bar')
//	if start != -1 && end != -1 {    fmt.Println("found at", start, end) }
func SearchForString(str string, substr string) (int, int) {
	m := search.New(language.English, search.IgnoreCase)
	return m.IndexString(str, substr)
}

// SearchForStringIndex Example: index, found := SearchForStringIndex('foobar', 'bar')
//	if found { fmt.Println("match starts at", index) }
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

/*************************************************************************************************/

// KeyValuePairInterface interface for AcmArticle, Vocabulary, Occurrence, WordScore structs.
type KeyValuePairInterface interface {
	Print() string
	GetKeyValuePairs() (map[string]string, map[int]string)
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
	Id          uint32  `json:"id"`
	Word        string  `json:"word"`
	RowCount    int     `json:"rowcount"`
	Frequency   int     `json:"frequency"`
	WordRank    int     `json:"wordrank"`
	Probability float32 `json:"probability"` // Probability of word at rank.
	SpeechPart  string  `json:"speechpart"`  // Assign using BulkInsert_Vocabulary_Speechpart().
	Stem        string  `json:"stem"`        // Assign using libstemmer program
}

// GetKeyValuePairs method
func (v Vocabulary) GetKeyValuePairs() (map[string]string, map[int]string) {
	fieldNames := []string{"Id", "Word", "RowCount", "Frequency", "WordRank", "Probability", "SpeechPart", "Stem"}
	orderedMap := GetOrderedMap(fieldNames)

	predicateMap := make(map[string]string, len(fieldNames))
	predicateMap[fieldNames[0]] = strconv.FormatUint(uint64(v.Id), 10)
	predicateMap[fieldNames[1]] = v.Word
	predicateMap[fieldNames[2]] = strconv.Itoa(v.RowCount)
	predicateMap[fieldNames[3]] = strconv.Itoa(v.Frequency)
	predicateMap[fieldNames[4]] = strconv.Itoa(v.WordRank)
	predicateMap[fieldNames[5]] = fmt.Sprintf(floatFormatter, v.Probability)
	predicateMap[fieldNames[6]] = v.SpeechPart
	predicateMap[fieldNames[7]] = v.Stem

	return predicateMap, orderedMap
}

// Print method
func (v Vocabulary) Print() string {
	return fmt.Sprintf("%d : %s : %d : %d : %d : %f : %s : %s", v.Id, v.Word, v.RowCount, v.Frequency, v.WordRank, v.Probability, v.SpeechPart, v.Stem)
}

// VocabularySorterFreq Sort interface by Frequency. Len() is the number of elements in the collection.
type VocabularySorterFreq []Vocabulary

func (a VocabularySorterFreq) Len() int           { return len(a) }
func (a VocabularySorterFreq) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a VocabularySorterFreq) Less(i, j int) bool { return a[i].Frequency > a[j].Frequency } // want lowest frequencies first.

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
					c <- j + i*lv/numCPU + 1
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
			break
		}
	}

	return ndx
}

/*************************************************************************************************/

// ReplaceUnicodeCharacters func
func ReplaceUnicodeCharacters(line string) string {
	//const accent1 = "\xe9\x67\xe9"	// é
	//const accent2 = "\xe8\x6d\x65"	// è
	result := line
	if !utf8.ValidString(result) {
		bstr := []byte(result)
		for index, b := range bstr {
			if b == '\xe9' || b == '\xe8' {
				result = result[:index] + "e" + result[index+1:]
			}
		}
	}
	return result
}

// ReplaceSpecialCharacters for database storage. See https://www.starr.net/is/type/htmlcodes.html
func ReplaceSpecialCharacters(line string) string {
	r := strings.NewReplacer(
		"'", "",
		"\n", " ",
		"<a", "",
		"(", "",
		")", "",
		"-", " ",
		"/", " ",
		"&ndash;", " ",
		"&mdash;", " ",
		"&shy;", " ",
		"&nbsp;", " ",
		"&rsquo;", "`",
		"&lsquo;", "`",
		"&ldquo;", "`",
		"&rdquo;", "`",
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
		"&sup1;", "`",
		"&hellip;", " ",
		"&amp;", "🙵",
		"&pound;", "£",
		"&euro;", "€",
		"&ntilde;", "ñ",
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
	return fmt.Sprintf("%d:%s:%s:%d", o.AcmId, o.Word, o.ArchiveDate.StandardDate(), o.Nentry)
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
}

// New_WordScore func
func New_WordScore(id uint64, word string, timeframetype nt.TimeFrameType, startDate nt.NullTime, endDate nt.NullTime, density float32, linkage float32, growth float32, score float32) WordScore {
	timeInterval := nt.TimeInterval{Timeframetype: timeframetype, StartDate: startDate, EndDate: endDate}
	return WordScore{id, word, timeInterval, density, linkage, growth, score}
}

// Print method
func (v WordScore) Print() string {
	return fmt.Sprintf("%s : %s : %f : %f : %f : %f", v.Word, v.Timeinterval.ToString(), v.Density, v.Linkage, v.Growth, v.Score)
}

// GetKeyValuePairs method
func (v WordScore) GetKeyValuePairs() (map[string]string, map[int]string) {
	fieldNames := []string{"Id", "Word", "Timeinterval", "Density", "Linkage", "Growth", "Score"}
	orderedMap := GetOrderedMap(fieldNames)

	predicateMap := make(map[string]string, len(fieldNames))
	predicateMap[fieldNames[0]] = strconv.FormatUint(v.Id, 10)
	predicateMap[fieldNames[1]] = v.Word
	predicateMap[fieldNames[2]] = v.Timeinterval.ToString()
	predicateMap[fieldNames[3]] = fmt.Sprintf(floatFormatter, v.Density)
	predicateMap[fieldNames[4]] = fmt.Sprintf(floatFormatter, v.Linkage)
	predicateMap[fieldNames[5]] = fmt.Sprintf(floatFormatter, v.Growth)
	predicateMap[fieldNames[6]] = fmt.Sprintf(floatFormatter, v.Score)

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
// Usage: iter := s.Iterator(); for i, ok := iter(); ok; i, ok = iter() { fmt.Println(i) }
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
	result = ReplaceUnicodeCharacters(result)
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

// ConditionalProbability struct
type ConditionalProbability struct {
	Id           uint64          `json:"id"`
	WordList     string          `json:"wordlist"`     // concatenated("|")
	Probability  float32         `json:"probability"`  // Conditional Probability
	Timeinterval nt.TimeInterval `json:"timeinterval"` // declared in nulltime.go;
	FirstDate    nt.NullTime     `json:"firstdate"`
	LastDate     nt.NullTime     `json:"lastdate"`
	Pmi          float32         `json:"pmi"` // point mutual information.
	DateUpdated  nt.NullTime     `json:"dateupdated"`
}

/*************************************************************************************************/

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

/*************************************************************************************************/

// WordScoreConditionalFlat struct has extracted Timeinterval.
type WordScoreConditionalFlat struct {
	ID            int       `json:"id"` // negative values
	WordArray     []string  `json:"wordarray"`
	Wordlist      string    `json:"wordlist"`
	Score         float32   `json:"score"`
	Pmi           float32   `json:"pmi"`
	Timeframetype int       `json:"timeframetype"`
	StartDate     time.Time `json:"startdate"`
	EndDate       time.Time `json:"enddate"`
	FirstDate     time.Time `json:"firstdate"`
	LastDate      time.Time `json:"lastdate"`
	Common        bool      `json:"common"` // intersection; not in database.
}
