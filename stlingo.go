package main

import (
	"bufio"
	"crypto/md5"
	"fmt"
	"github.com/ararog/verbo"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
)

type PlatformType int

const (
	Android PlatformType = iota
	IOS
)

type FileInfo struct {
	FilePath     string
	FilePathHash string
	Index        int
	PlatformType PlatformType
}

type StringLine struct {
	FileInfo FileInfo
	Line     int
	Key      string
	Value    string
	Checked  bool
}

type MatchedResults []*StringLine

func (m MatchedResults) Len() int {
	return len(m)
}

func (m MatchedResults) Swap(i, j int) {
	m[i], m[j] = m[j], m[i]
}

func (m MatchedResults) Less(i, j int) bool {
	if strings.Compare(m[i].Value, m[j].Value) == -1 {
		return true
	} else {
		return false
	}
}

type ScoredStringLine struct {
	*StringLine
	Score int
	Index int
}

type ScoredStringLineList []ScoredStringLine

func (l ScoredStringLineList) Len() int {
	return len(l)
}

func (l ScoredStringLineList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l ScoredStringLineList) Less(i, j int) bool {
	return l[i].Score < l[j].Score
}

func quote(str string) string {
	return "\"" + strings.TrimSpace(str) + "\""
}

func analyze(matchedData MatchedResults, diffScore int) error {
	fmt.Printf(",PlatformType, score(or X), val, key, line, file_index\n")

	// rank similarity
	for i := 0; i < len(matchedData); i++ {
		lineDataI := matchedData[i]
		if lineDataI.Checked == true {
			continue
		}

		fmt.Printf(",%v,X,%s,%s,%d,%s\n", lineDataI.FileInfo.PlatformType, quote(lineDataI.Value), lineDataI.Key, lineDataI.Line, lineDataI.FileInfo.FilePath)

		var scores ScoredStringLineList
		for j := i + 1; j < len(matchedData); j++ {
			lineDataJ := matchedData[j]

			score := verbo.Levenshtein(lineDataI.Value, lineDataJ.Value)
			scores = append(scores, ScoredStringLine{
				StringLine: lineDataJ,
				Score:      score,
				Index:      j,
			})
		}
		//sort.Sort(sort.Reverse(scores))
		sort.Sort(scores)

		for _, item := range scores {
			if item.Score > diffScore {
				break
			}
			if item.StringLine.Checked == true {
				continue
			}

			if item.Score == 0 {
				matchedData[item.Index].Checked = true
			}

			fmt.Printf(",%v,%d,%s,%s,%d,%d\n", item.StringLine.FileInfo.PlatformType, item.Score, quote(item.StringLine.Value), item.StringLine.Key, item.StringLine.Line, item.StringLine.FileInfo.Index)
		}
	}

	return nil
}

func readInputTxt() ([]string, error) {
	var f *os.File
	fmt.Printf(">> read file: %s\n", os.Args[1])
	f, err := os.Open(os.Args[1])
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := bufio.NewReaderSize(f, 4096)
	rc := regexp.MustCompile(`//.*?`)
	rs := regexp.MustCompile(`\s*?\t*?.*?\n`)
	var results []string
	for {
		line, _, err := reader.ReadLine()
		//fmt.Printf("%s\n", line)
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		s := string(line)
		if len(s) == 0 {
			continue
		}
		if rc.MatchString(s) {
			continue
		}
		if rs.MatchString(s) {
			continue
		}
		results = append(results, s)
	}
	return results, nil
}

func mapFileInfos(files []string) []FileInfo {
	var fileInfos []FileInfo
	fmt.Printf(",index, platformType, hash, filePath\n")
	rxml := regexp.MustCompile(`.*?\.xml`)
	for index, filePath := range files {
		h := md5.New()
		io.WriteString(h, filePath)
		hash := fmt.Sprintf("%x", h.Sum(nil))
		platformType := IOS
		if rxml.MatchString(filePath) {
			platformType = Android
		}
		fi := FileInfo{
			FilePath:     filePath,
			FilePathHash: hash,
			Index:        index,
			PlatformType: platformType,
		}
		fileInfos = append(fileInfos, fi)
		fmt.Printf(", %d, %v, %s, %s\n", index, platformType, hash, filePath)
	}
	fmt.Printf("\n\n")
	return fileInfos
}

func main() {
	files, err := readInputTxt()
	if err != nil {
		panic(err)
	}

	fileInfos := mapFileInfos(files)

	results, err := parseFiles(fileInfos)
	if err != nil {
		panic(err)
	}

	var flattenMatchedResults MatchedResults
	for _, matchedResults := range results {
		for _, stringLine := range matchedResults {
			flattenMatchedResults = append(flattenMatchedResults, stringLine)
		}
	}
	sort.Sort(flattenMatchedResults)
	analyze(flattenMatchedResults, 3)
}
