package filesystem

// deploy using 'go install'
import (
	"archive/zip"
	"bufio"
	"database/sql"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	hd "github.com/dgnabasik/acmsearchlib/headers"
	nt "github.com/dgnabasik/acmsearchlib/nulltime"
	"github.com/gin-gonic/gin"
)

/*************************************************************************************/

// FileServiceInterface interface functions are not placed into acmsearchlib.
type FileServiceInterface interface {
	GetTextFile(ctx *gin.Context)
}

// FileService struct implements FileServiceInterface.
type FileService struct {
	//tableController TableController
}

// TableController struct
type TableController struct {
	DB *sql.DB
}

// constants
const (
	PostfixHTML = ".html"
)

// GetFilePrefixPath func returns html file location.
func GetFilePrefixPath() string {
	envVar := os.Getenv("ACM_FILE_PREFIX")
	if envVar == "" {
		log.Printf("ACM_FILE_PREFIX env var not found...using default\n")
		envVar = "/home/david/"
	}
	return envVar

}

// ReadDir reads the directory named by dirname and returns a list of FileInfo entries [sorted by filename.]
func ReadDir(dirname string) ([]os.FileInfo, error) {
	f, err := os.Open(dirname)
	if err != nil {
		return nil, err
	}
	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return nil, err
	}
	//sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })
	return list, nil
}

// CreateDirectory func
func CreateDirectory(dirPath string) error {
	_, err := os.Stat(dirPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(dirPath, 0777)
	}
	if err != nil {
		str := "filesystem.CreateDirectory(" + dirPath + "): "
		log.Printf(str+"%s%+v\n", err)
	}
	return err
}

// DeleteDirectory func
func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		str := "filesystem.DeleteDirectory(" + dirPath + "): "
		log.Printf(str+"%s%+v\n", err)
	}
	return err
}

// AddFileToZip func
func AddFileToZip(zipWriter *zip.Writer, filename string) error {
	fileToZip, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer fileToZip.Close()

	// Get the file information
	info, err := fileToZip.Stat()
	if err != nil {
		return err
	}

	header, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	// Using FileInfoHeader() above only uses the basename of the file. If we want
	// to preserve the folder structure we can overwrite this with the full path.
	header.Name = filename

	// Change to deflate to gain better compression: http://golang.org/pkg/archive/zip/#pkg-constants
	header.Method = zip.Deflate

	writer, err := zipWriter.CreateHeader(header)
	if err != nil {
		return err
	}
	_, err = io.Copy(writer, fileToZip)
	return err
}

// ZipFiles func pathPrefix+fileTemplate must be well-formed. pathPrefix must include trailing slash.
// Produces zip file contains full pathPrefix. Returns zip filename.
func ZipFiles(pathPrefix string, fileExt string, targetFileName string) (string, error) {
	fileInfo, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		str := "filesystem.ZipFiles(" + pathPrefix + "): "
		log.Printf(str+"%s%+v\n", err)
		return "", err
	}

	var fileList []string

	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), fileExt) {
			fileList = append(fileList, pathPrefix+file.Name())
		}
	}

	if len(fileList) == 0 {
		err = errors.New("There are no matching files in " + pathPrefix + "*" + fileExt)
		log.Printf("filesystem.ZipFiles: %+v\n", err)
		return "", err
	}

	TargetFileName := pathPrefix + targetFileName
	newZipFile, err := os.Create(TargetFileName)
	if err != nil {
		str := "filesystem.ZipFiles(" + TargetFileName + "): "
		log.Printf(str+"%s%+v\n", err)
		return "", err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range fileList {
		err = AddFileToZip(zipWriter, file)
		if err != nil {
			str := "filesystem.AddFileToZip(" + file + "): "
			log.Printf(str+"%s%+v\n", err)
			return TargetFileName, err
		}
	}

	return TargetFileName, nil
}

// FileExists Returns false if directory.
func FileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false, nil
	}
	return !info.IsDir(), err
}

// ReadFileIntoString func
func ReadFileIntoString(filePath string) (string, error) {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(content), err
}

// ReadTextLines reads a whole file into memory and returns a slice of its lines. Applys .ToLower(). Skips empty lines if normalizeText is true.
func ReadTextLines(filePath string, normalizeText bool) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("filesystem.ReadTextLines: %+v\n", err)
		return nil, err
	}
	fi, err := file.Stat()
	if err != nil {
		log.Printf("filesystem.ReadTextLines: %+v\n", err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	maxCapacity := fi.Size() + 1
	buf := make([]byte, maxCapacity)
	scanner.Buffer(buf, int(maxCapacity))

	for scanner.Scan() {
		if normalizeText {
			str := strings.ToLower(strings.TrimSpace(scanner.Text()))
			if len(str) > 0 {
				lines = append(lines, str)
			}
		} else {
			lines = append(lines, scanner.Text())
		}
	}
	return lines, scanner.Err()
}

// WriteTextLines writes/appends the lines to the given file.
func WriteTextLines(lines []string, filePath string, appendData bool) error {
	if !appendData {
		itExists, _ := FileExists(filePath)
		if itExists {
			_ = os.Remove(filePath)
		}
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("filesystem.WriteTextLines: %+v\n", err)
		return err
	}
	defer file.Close()

	contents := strings.Join(lines, "\n")
	if _, err := file.WriteString(contents); err != nil {
		log.Printf("filesystem.WriteTextLines: %+v\n", err)
	}

	return err
}

// for file sorting
type fileSort struct {
	FileName string
	FileTime int64
}

// GetFileList returns []string of file pathnames ordered by logical datetime. Built for ~/acmFiles.
func GetFileList(filePath string, since nt.NullTime) ([]string, error) {
	var files []fileSort
	var fileList []string
	cutoff := (since.DT.UnixNano() / 1000000) - 100

	fileInfo, err := ioutil.ReadDir(filePath)
	if err != nil {
		return fileList, err
	}
	for _, file := range fileInfo {
		if strings.HasSuffix(strings.ToLower(file.Name()), PostfixHTML) {
			fileTime := GetFileTime(file.Name())
			if fileTime >= cutoff {
				files = append(files, fileSort{FileName: GetFilePrefixPath() + file.Name(), FileTime: fileTime})
			}
		}
	}

	// order by logical datetime
	sort.SliceStable(files, func(i, j int) bool {
		return files[i].FileTime < files[j].FileTime
	})

	// extract FileNames
	for _, file := range files {
		fileList = append(fileList, file.FileName)
	}
	return fileList, nil
}

// GetMostRecentFileAsNullTime reads a directory and return the most recent (valid) file as a NullTime. Built for ~/acmFiles.
func GetMostRecentFileAsNullTime(dirname string) (nt.NullTime, error) {
	f, err := os.Open(dirname)
	if err != nil {
		log.Printf("Unable to open folder: %s %+v\n", dirname, err)
		return nt.NullTimeToday(), err
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Printf("Unable to open file: %+v\n", err)
		return nt.NullTimeToday(), err
	}

	// process template files (case-sensitive) and translate file format jul-12-2019.html to yyyy-mm-dd
	regPattern := "\\D\\D\\D-\\d\\d-\\d\\d\\d\\d" + PostfixHTML
	layoutISO := "2006-01-02"
	var mostRecentFile = ""
	var latest int64

	for _, file := range files {
		match, _ := regexp.MatchString(regPattern, file.Name())
		if match && file.Size() > 0 {
			date := nt.GetStandardDateForm(file.Name())
			t, _ := time.Parse(layoutISO, date)
			if t.Unix() > latest {
				mostRecentFile = date
				latest = t.Unix()
			}
		}
	}
	return nt.New_NullTime(mostRecentFile), err
}

// GetFileTime converts fileName to milliseconds.
func GetFileTime(fileName string) int64 {
	file := strings.Replace(filepath.Base(fileName), PostfixHTML, "", 1) // s, old, new string, n int
	sdt := nt.New_NullTime(nt.GetStandardDateForm(file))
	return sdt.DT.UnixNano() / 1000000
}

// GetSourceDirectory func sources local acm.env file.
func GetSourceDirectory() string {
	envVar := os.Getenv("REACT_ACM_SOURCE_DIR")
	if envVar == "" {
		log.Printf("REACT_ACM_SOURCE_DIR env var not found...using default\n")
		envVar = "/home/david/websites/acmsearch/docs/"
	}
	return envVar
}

/*************************************************************************************/

// GetTextFile method assigns key values starting at 1.
func (fss *FileService) GetTextFile(ctx *gin.Context) {
	filename := ctx.Param("name")
	words, err := ReadTextLines(GetSourceDirectory()+filename, true) // applys toLower()

	if err != nil {
		log.Printf("FileService.GetTextFile: %+v\n", err)
		ctx.JSON(404, gin.H{
			"message": "FileService.GetTextFile: " + err.Error(),
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
