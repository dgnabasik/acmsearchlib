package category

// category.go manages categories. Derived from ~/websites/dropdownlists/golang/listservice.go

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	dbx "github.com/dgnabasik/acmsearchlib/database"
	fs "github.com/dgnabasik/acmsearchlib/filesystem"
	hd "github.com/dgnabasik/acmsearchlib/headers"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

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

// CategoryInput struct comes in as FormData, so need form tags. See SetWordCategory().
type CategoryInput struct {
	Description string `json:"description" form:"description" binding:"required"`
	FileName    string `json:"filename" form:"filename" binding:"required"`
}

// GraphServiceInterface interface functions are not placed into acmsearchlib.
type GraphServiceInterface interface {
	GetSpecialMap(ctx *gin.Context)
	GetCategoryMap(ctx *gin.Context)
}

// ListService struct implements GraphServiceInterface.
type ListService struct {
	tableController TableController
}

// FileServiceInterface interface functions are not placed into acmsearchlib.
type FileServiceInterface interface {
	GetTextFile(ctx *gin.Context)
	SetWordCategory(ctx *gin.Context)
}

// FileService struct implements FileServiceInterface.
type FileService struct {
	tableController TableController
}

// TableController struct
type TableController struct {
	DB *sql.DB
}

/*************************************************************************************/

// InsertCategoryWords func. 32k statement limit.
func InsertCategoryWords(categoryID uint64, words []string) error {
	dateupdated := time.Now()

	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	txn, err := db.Begin()
	dbx.CheckErr(err)

	// Must use lowercase column names! First param is table name.
	stmt, err := txn.Prepare(pq.CopyIn("special", "word", "category", "dateupdated"))
	dbx.CheckErr(err)

	for _, word := range words {
		_, err = stmt.Exec(word, categoryID, dateupdated)
		dbx.CheckErr(err)
	}

	_, err = stmt.Exec()
	dbx.CheckErr(err)

	err = stmt.Close()
	dbx.CheckErr(err)

	err = txn.Commit()
	dbx.CheckErr(err)

	return nil
}

// InsertWordCategory func
func InsertWordCategory(description string) (CategoryTable, error) {
	dateupdated := time.Now()
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	var id uint64
	INSERT := "INSERT INTO wordcategory (description, dateupdated) VALUES ($1, $2) returning id"
	err = db.QueryRow(INSERT, description, dateupdated).Scan(&id)
	dbx.CheckErr(err)

	categoryTable := CategoryTable{Id: id, Description: description, DateUpdated: dateupdated}
	return categoryTable, nil
}

// GetSpecialMap func filters by category
func GetSpecialMap(category int) ([]hd.LookupMap, error) {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT id, word FROM special WHERE category=" + strconv.Itoa(category) + " ORDER BY word"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var lookup hd.LookupMap
	lookupMap := []hd.LookupMap{}
	for rows.Next() {
		err = rows.Scan(&lookup.Value, &lookup.Label)
		dbx.CheckErr(err)
		lookupMap = append(lookupMap, lookup)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return lookupMap, nil
}

// GetCategoryMap func
func GetCategoryMap() ([]hd.LookupMap, error) {
	db, err := dbx.GetDatabaseReference()
	defer db.Close()

	SELECT := "SELECT id, description FROM wordcategory"
	rows, err := db.Query(SELECT)
	dbx.CheckErr(err)
	defer rows.Close()

	var lookup hd.LookupMap
	lookupMap := []hd.LookupMap{}
	for rows.Next() {
		err = rows.Scan(&lookup.Value, &lookup.Label)
		dbx.CheckErr(err)
		lookupMap = append(lookupMap, lookup)
	}

	err = rows.Err()
	dbx.CheckErr(err)

	return lookupMap, nil
}

/*************************************************************************************/

// GetSpecialMap method
func (qs *ListService) GetSpecialMap(ctx *gin.Context) {
	category, err := strconv.Atoi(ctx.Param("category"))
	lookupMap, err := GetSpecialMap(category)

	if err != nil {
		log.Printf("ListService.GetSpecialMap: %+v\n", err)
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("ListService.GetSpecialMap: " + err.Error()),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"LookupMap": lookupMap,
	})
}

// GetCategoryMap method
func (qs *ListService) GetCategoryMap(ctx *gin.Context) {
	lookupMap, err := GetCategoryMap()
	if err != nil {
		log.Printf("ListService.GetCategoryMap: %+v\n", err)
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("ListService.GetCategoryMap: " + err.Error()),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"LookupMap": lookupMap,
	})
}

/*************************************************************************************/

// GetSourceDirectory func sources local acm.env file.
func GetSourceDirectory() string {
	return os.Getenv("REACT_ACM_SOURCE_DIR")
}

// GetTextFile method assigns key values starting at 1.
func (fss *FileService) GetTextFile(ctx *gin.Context) {
	filename := ctx.Param("name")
	words, err := fs.ReadTextLines(GetSourceDirectory()+filename, true) // applys toLower()

	if err != nil {
		log.Printf("ListService.GetTextFile: %+v\n", err)
		ctx.JSON(404, gin.H{
			"message": fmt.Sprintf("ListService.GetTextFile: " + err.Error()),
		})
		return
	}

	// first remove duplicates.
	amap := make(map[int]string, len(words))
	for ndx, word := range words {
		amap[ndx] = word
	}

	// ndx starts at 0 but lookupMap.Value starts at 1.
	lookupMap := make([]hd.LookupMap, len(amap))
	for ndx := range amap {
		if len(amap[ndx]) > 0 {
			lookupMap[ndx] = hd.LookupMap{Value: ndx + 1, Label: amap[ndx]}
		}
	}

	ctx.JSON(200, gin.H{
		"LookupMap": lookupMap,
	})
}

// SetWordCategory method returns CategoryTable struct if successful.
func (fss *FileService) SetWordCategory(ctx *gin.Context) {
	var categoryInput CategoryInput
	err := ctx.Bind(&categoryInput)
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	// these are arbitrary HTML response codes.
	if len(categoryInput.Description) == 0 || len(categoryInput.FileName) == 0 {
		ctx.JSON(401, gin.H{
			"message": fmt.Sprintf("ListService.SetWordCategory: missing description or filename"),
		})
		return
	}

	words, err := fs.ReadTextLines(GetSourceDirectory()+categoryInput.FileName, true) //  normalizeText
	if err != nil {
		ctx.JSON(403, gin.H{
			"error": err.Error(),
		})
		return
	}

	wordCategory, err := InsertWordCategory(categoryInput.Description)
	if err == nil {
		err = InsertCategoryWords(wordCategory.Id, words)
		if err != nil {
			ctx.JSON(404, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	ctx.JSON(200, gin.H{
		"CategoryTable": wordCategory,
	})
}

/*************************************************************************************/

/* InitializeRoutes func: cannot rely upon the order of execution of init() functions!
// Query string parameters are parsed using the existing underlying request object.
func InitializeRoutes(qs *ListService, fss *FileService) *gin.Engine {
	defer qs.tableController.DB.Close()

	gin.SetMode(gin.ReleaseMode) // Switch to "release" mode in production; or export GIN_MODE=release
	router := gin.Default()

	// Credential is not supported if the CORS header ‘Access-Control-Allow-Origin’ is ‘*’
	// The wildcard asterisk only works for AllowedOrigins. Using the asterisk in AllowedMethods and AllowedHeaders will have no affect.
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"User-Agent", "Referrer", "Host", "Token", "Accept", "Content-Type", "Origin", "Content-Length", "X-Requested-With", "Accept-Encoding"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		ExposeHeaders:    []string{"Content-Length", "Content-Type"},
		AllowOriginFunc: func(origin string) bool {
			return true // origin == hostName
		},
		MaxAge: 86400,
	}))

	router.Use(ErrorHandler)
	router.Static("/static", "./build/static") // use the loaded source
	router.Use(static.Serve("/", static.LocalFile("./build", true)))

	// Direct all routes to index.html:
	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "index.html", nil)
	})

	list := router.Group("/list")
	{
		list.GET("", qs.GetCategoryMap)
		list.GET("/:category", qs.GetSpecialMap) // number
	}

	file := router.Group("/file")
	{
		file.GET("/:name", fss.GetTextFile) // string
		file.OPTIONS("", ContextOptions)
		file.POST("", fss.SetWordCategory) // append formdata{description, filename}
	}

	apiPort := GetPort()
	api := "Handling REST-API calls on " + GetHost() + ":" + apiPort
	fmt.Println(api)
	fmt.Println("  GET /list")
	fmt.Println("  GET /list/:category")
	fmt.Println("  GET /file/:name")

	router.Run(":" + apiPort)
	return router
} */
