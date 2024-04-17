package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"math"
	"os"
)

type FileStats struct {
	Path  string
	Bytes int
	Lines int
	Words int
	Chars int
	Err   *error
}

func (fileStats *FileStats) String(l, w, c, m bool, maxLen int) string {
	var buffer bytes.Buffer
	if fileStats.Err != nil {
		buffer.WriteString(fmt.Sprintf("wc: %s: %s", fileStats.Path, *fileStats.Err))
		return buffer.String()
	}
	if l {
		buffer.WriteString(fmt.Sprintf("%*d ", maxLen, fileStats.Lines))
	}
	if w {
		buffer.WriteString(fmt.Sprintf("%*d ", maxLen, fileStats.Words))
	}
	if m {
		buffer.WriteString(fmt.Sprintf("%*d ", maxLen, fileStats.Chars))
	}
	if c {
		buffer.WriteString(fmt.Sprintf("%*d ", maxLen, fileStats.Bytes))
	}
	if fileStats.Path != "/dev/stdin" {
		buffer.WriteString(fileStats.Path)
	}
	return buffer.String()
}

func GetFileStats(countLines, countWords, countBytes, countChars bool, data []byte, filePath string) (*FileStats, error) {

	fileStats := &FileStats{}
	if countBytes {
		fileStats.Bytes = len(data)
	}

	if countLines {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Split(bufio.ScanLines)
		for scanner.Scan() {
			fileStats.Lines++
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	if countWords {
		scanner := bufio.NewScanner(bytes.NewReader(data))
		scanner.Split(bufio.ScanWords)
		for scanner.Scan() {
			fileStats.Words++
		}
		if err := scanner.Err(); err != nil {
			return nil, err
		}
	}

	if countChars {
		fileStats.Chars = len(string(data))
	}
	fileStats.Path = filePath
	return fileStats, nil
}

func main() {

	var countBytes bool
	var countLines bool
	var countWords bool
	var countChars bool
	flag.BoolVar(&countBytes, "c", false, "count bytes in the input")
	flag.BoolVar(&countLines, "l", false, "count lines in the input")
	flag.BoolVar(&countWords, "w", false, "count words in the input")
	flag.BoolVar(&countChars, "m", false, "count characters in the input")
	flag.Parse()

	//If no flags are provided, count Bytes, Lines and Words
	if !countBytes && !countLines && !countWords {
		countBytes = true
		countLines = true
		countWords = true
	}

	//If no file is provided, read from stdin
	filepaths := flag.Args()
	if len(filepaths) == 0 {
		filepaths = append(filepaths, "/dev/stdin")
	}

	totalCounts := &FileStats{
		Path: "total",
	}
	var fileCounts []*FileStats

	for _, filepath := range filepaths {
		fileStat := &FileStats{}
		data, err := os.ReadFile(filepath)
		if err != nil {
			fileStat = &FileStats{
				Path: filepath,
				Err:  &err,
			}
			fileCounts = append(fileCounts, fileStat)
			continue
		}

		fileStat, err = GetFileStats(countLines, countWords, countBytes, countChars, data, filepath)
		if err != nil {
			fileStat = &FileStats{
				Path: filepath,
				Err:  &err,
			}
			fileCounts = append(fileCounts, fileStat)
			continue
		}

		fileCounts = append(fileCounts, fileStat)
		totalCounts.Bytes += fileStat.Bytes
		totalCounts.Lines += fileStat.Lines
		totalCounts.Words += fileStat.Words
		totalCounts.Chars += fileStat.Chars
	}

	//Find the maximum length of the counts
	maximumCount := math.Max(
		math.Max(float64(totalCounts.Bytes), float64(totalCounts.Lines)),
		math.Max(float64(totalCounts.Words), float64(totalCounts.Chars)),
	)
	maxDigits := int(math.Floor(math.Log10(maximumCount))) + 1

	for _, fileCount := range fileCounts {
		fmt.Println(fileCount.String(countLines, countWords, countBytes, countChars, maxDigits))
	}

	if len(fileCounts) > 1 {
		fmt.Println(totalCounts.String(countLines, countWords, countBytes, countChars, maxDigits))
	}

}
