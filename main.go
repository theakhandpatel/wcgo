package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"runtime"
	"sync"
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

func outputStats(out io.Writer, fileCounts []FileStats, l, w, c, m bool) {
	totalCounts := FileStats{}
	for _, fileCount := range fileCounts {
		totalCounts.Bytes += fileCount.Bytes
		totalCounts.Lines += fileCount.Lines
		totalCounts.Words += fileCount.Words
		totalCounts.Chars += fileCount.Chars
	}

	//Find the maximum length of the counts
	maximumCount := math.Max(
		math.Max(float64(totalCounts.Bytes), float64(totalCounts.Lines)),
		math.Max(float64(totalCounts.Words), float64(totalCounts.Chars)),
	)
	maxDigits := int(math.Floor(math.Log10(maximumCount))) + 1

	for _, fileCount := range fileCounts {
		fmt.Fprintln(out, fileCount.String(l, w, c, m, maxDigits))
	}

	if len(fileCounts) > 1 {
		fmt.Fprintln(out, totalCounts.String(l, w, c, m, maxDigits))
	}
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
	if !countBytes && !countLines && !countWords && !countChars {
		countBytes = true
		countLines = true
		countWords = true
	}

	//If no file is provided, read from stdin
	filepaths := flag.Args()
	if len(filepaths) == 0 {
		filepaths = append(filepaths, "/dev/stdin")
	}

	resCh := make(chan []FileStats)
	doneCh := make(chan struct{})
	filesCh := make(chan string)

	go func() {
		defer close(filesCh)
		for _, filepath := range filepaths {
			filesCh <- filepath
		}
	}()

	wg := sync.WaitGroup{}

	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for filepath := range filesCh {
				fileStat := &FileStats{}
				data, err := os.ReadFile(filepath)

				if err != nil {
					fileStat = &FileStats{
						Path: filepath,
						Err:  &err,
					}
					resCh <- []FileStats{*fileStat}
					continue
				}

				fileStat, err = GetFileStats(countLines, countWords, countBytes, countChars, data, filepath)
				if err != nil {
					fileStat = &FileStats{
						Path: filepath,
						Err:  &err,
					}
					// fileCounts = append(fileCounts, fileStat)
					resCh <- []FileStats{*fileStat}
					continue
				}

				resCh <- []FileStats{*fileStat}
			}

		}()
	}

	go func() {
		wg.Wait()
		close(doneCh)
	}()

	fileCounts := []FileStats{}

	for {
		select {
		case data := <-resCh:
			fileCounts = append(fileCounts, data...)
		case <-doneCh:
			outputStats(os.Stdout, fileCounts, countLines, countWords, countBytes, countChars)
			return
		}
	}

}
