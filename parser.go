package main

import (
	"bufio"
	"crypto/md5"
	"encoding/xml"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

func parseIOS(fileInfo FileInfo) (MatchedResults, error) {
	f, err := os.Open(fileInfo.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var matchedData MatchedResults
	reader := bufio.NewReaderSize(f, 4096)
	i := 1
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		}

		s := string(line)
		r := regexp.MustCompile(`"(.+?)"\s*?=\s*?"(.+?)";`)
		results := r.FindAllStringSubmatch(s, -1)
		if len(results) == 0 {
			i += 1
			continue
		}
		h := md5.New()
		io.WriteString(h, fileInfo.FilePath)
		lineData := &LocalizableStringRaw{
			FileInfo: fileInfo,
			Line:     i,
			Key:      results[0][1],
			Value:    results[0][2],
			Checked:  false,
		}
		matchedData = append(matchedData, lineData)
		i += 1
	}

	return matchedData, nil
}

type AndroidStringRoot struct {
	XMLName      xml.Name      `xml:"resources"`
	Strings      []String      `xml:"string,omitempty"`
	StringsArray []StringArray `xml:"string-array,omitempty"`
}

type String struct {
	XMLName xml.Name `xml:"string"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:",chardata"`
}

type StringArray struct {
	XMLName xml.Name `xml:"string-array"`
	Name    string   `xml:"name,attr"`
	Items   []string `xml:"item"`
}

func parseAndroid(fileInfo FileInfo) (MatchedResults, error) {
	f, err := os.Open(fileInfo.FilePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var matchedData MatchedResults
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	androidStringRoot := AndroidStringRoot{}
	err = xml.Unmarshal(data, &androidStringRoot)
	if err != nil {
		return nil, err
	}
	for _, str := range androidStringRoot.Strings {
		lineData := &LocalizableStringRaw{
			FileInfo: fileInfo,
			Key:      str.Name,
			Value:    str.Value,
			Checked:  false,
		}
		matchedData = append(matchedData, lineData)
	}
	for _, strArray := range androidStringRoot.StringsArray {
		for itemIndex, item := range strArray.Items {
			lineData := &LocalizableStringRaw{
				FileInfo: fileInfo,
				Key:      strArray.Name + "." + strconv.Itoa(itemIndex),
				Value:    item,
				Checked:  false,
			}
			matchedData = append(matchedData, lineData)
		}
	}

	return matchedData, nil
}

func parse(fileInfo FileInfo) (MatchedResults, error) {
	if fileInfo.PlatformType == IOS {
		return parseIOS(fileInfo)
	} else {
		return parseAndroid(fileInfo)
	}
}

func parseFiles(fileInfos []FileInfo) ([]MatchedResults, error) {
	var results []MatchedResults
	for _, fileInfo := range fileInfos {
		result, err := parse(fileInfo)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
