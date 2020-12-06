package filesystem

// deploy using 'go install'
import (
	hd "acmsearchlib/headers"
	nt "acmsearchlib/nulltime"
	"archive/zip"
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

// constants
const (
	FileSystemPrefix = "/home/david/"
	PrefixFilePath   = FileSystemPrefix + "acmFiles/"
	PrefixProcessed  = FileSystemPrefix + "acm/"
	PostfixHTML      = ".html"
)

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
		fmt.Print("CreateDirectory(" + dirPath + "): ")
		fmt.Println(err)
	}
	return err
}

// DeleteDirectory func
func DeleteDirectory(dirPath string) error {
	err := os.RemoveAll(dirPath)
	if err != nil {
		fmt.Print("DeleteDirectory(" + dirPath + "): ")
		fmt.Println(err)
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

// ZipFiles func pathPrefix+fileTemplate must be well-formed.
func ZipFiles(pathPrefix string, fileExt string, targetFileName string) error {
	fileInfo, err := ioutil.ReadDir(pathPrefix)
	if err != nil {
		fmt.Print("ZipFiles(" + pathPrefix + "): ")
		fmt.Println(err)
		return err
	}

	var fileList []string

	for _, file := range fileInfo {
		if strings.HasSuffix(file.Name(), fileExt) {
			fileList = append(fileList, pathPrefix+file.Name())
		}
	}

	if len(fileList) == 0 {
		fmt.Println("There are no matching files in " + pathPrefix + "*" + fileExt)
		return nil
	}

	newZipFile, err := os.Create(targetFileName)
	if err != nil {
		fmt.Print("ZipFiles(" + targetFileName + "): ")
		fmt.Println(err)
		return err
	}
	defer newZipFile.Close()

	zipWriter := zip.NewWriter(newZipFile)
	defer zipWriter.Close()

	// Add files to zip
	for _, file := range fileList {
		err = AddFileToZip(zipWriter, file)
		if err != nil {
			fmt.Print("AddFileToZip(" + file + "): ")
			fmt.Println(err)
			return err
		}
	}

	return nil
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

// ReadTextLines reads a whole file into memory and returns a slice of its lines. Applys .ToLower()
func ReadTextLines(filePath string, normalizeText bool) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print("ReadTextLines(): ")
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if normalizeText {
			lines = append(lines, strings.ToLower(strings.TrimSpace(scanner.Text())))
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
		fmt.Print("WriteTextLines(): ")
		fmt.Println(err)
		return err
	}
	defer file.Close()

	contents := strings.Join(lines, "\n")
	if _, err := file.WriteString(contents); err != nil {
		log.Println(err)
	}

	return err
}

// ReadOccurrenceListFromCsvFile caller must assign Id, ArchiveDate.	24889 | 2009-09-02
func ReadOccurrenceListFromCsvFile(filePath string) ([]hd.Occurrence, error) {
	var occurrenceList []hd.Occurrence
	source, err := ReadTextLines(filePath, true)
	if err != nil {
		fmt.Println("Could not open " + filePath)
		return occurrenceList, err
	}

	archiveDate := nt.NullTimeToday()
	for _, line := range source {
		tokens := strings.Split(line, ",")
		nentry, _ := strconv.Atoi(tokens[2])
		item := hd.Occurrence{AcmId: 0, ArchiveDate: archiveDate, Word: tokens[0], Nentry: nentry}
		occurrenceList = append(occurrenceList, item)
	}

	return occurrenceList, nil
}

// for file sorting
type fileSort struct {
	FileName string
	FileTime int64
}

// GetFileList returns []string of file pathnames ordered by logical datetime.
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
				files = append(files, fileSort{FileName: PrefixFilePath + file.Name(), FileTime: fileTime})
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

// GetMostRecentFileAsNullTime reads a directory and return the most recent (valid) file as a NullTime.
func GetMostRecentFileAsNullTime(dirname string) (nt.NullTime, error) {
	f, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}
	files, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		log.Fatal(err)
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
