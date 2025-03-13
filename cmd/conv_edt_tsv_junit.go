package main

import (
	"encoding/xml"
	"encoding/csv"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/azheval/conv_edt_tsv_junit/pkg/logging"
	"github.com/azheval/conv_edt_tsv_junit/pkg/config"
)

var (
	version = ""
	build   = ""
)

type ErrorRecord struct {
	Date        string `xml:"date,attr"`
	Priority    string `xml:"priority,attr"`
	CheckType   string `xml:"checkType,attr"`
	Project     string `xml:"project,attr"`
	Standard    string `xml:"standard,attr"`
	ErrorModule string `xml:"errorModule,attr"`
	ErrorLine   string `xml:"errorLine,attr"`
	ErrorText   string `xml:"errorText,attr"`
}

type TestSuites struct {
	XMLName   xml.Name  `xml:"testsuites"`
	Time      string    `xml:"time,attr"`
	Tests     int       `xml:"tests,attr"`
	Errors    int       `xml:"errors,attr"`
	Failures  int       `xml:"failures,attr"`
	TestSuite []TestSuite `xml:"testsuite"`
}

type TestSuite struct {
	XMLName    xml.Name   `xml:"testsuite"`
	Name       string     `xml:"name,attr"`
	Timestamp  string     `xml:"timestamp,attr"`
	Time       string     `xml:"time,attr"`
	Tests      int        `xml:"tests,attr"`
	Errors     int        `xml:"errors,attr"`
	Failures   int        `xml:"failures,attr"`
	Skipped    int        `xml:"skipped,attr"`
	Properties []Property `xml:"properties>property"`
	TestCases  []TestCase `xml:"testcase"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type TestCase struct {
	XMLName   xml.Name  `xml:"testcase"`
	ClassName string    `xml:"classname,attr"`
	Name      string    `xml:"name,attr"`
	Time      string    `xml:"time,attr"`
	Failures  []Failure `xml:"failure"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
	Text    string `xml:",chardata"`
}



func getTestSuiteByName(testSuites TestSuites, testSuiteTimestamp string, fileName string, recordName string, logger slog.Logger) (TestSuite, int) {
	if len(testSuites.TestSuite) == 0 {
		ts := TestSuite{
			Name:       fileName + "_" + recordName,
			Timestamp:  testSuiteTimestamp,
			Time:       "0",
			Tests:      0,
			Errors:     0,
			Failures:   0,
			Skipped:    0,
			Properties: []Property{},
			TestCases:  []TestCase{},
		}
		logger.Debug("new test suite", "name", fileName+"_"+recordName)
		return ts, -1
	} else {
		for index, ts := range testSuites.TestSuite {
			if ts.Name == fileName+"_"+recordName {
				logger.Debug("found test suite", "name", fileName+"_"+recordName)
				return ts, index
			}
		}
		ts := TestSuite{
			Name:       fileName + "_" + recordName,
			Timestamp:  testSuiteTimestamp,
			Time:       "0",
			Tests:      0,
			Errors:     0,
			Failures:   0,
			Skipped:    0,
			Properties: []Property{},
			TestCases:  []TestCase{},
		}
		logger.Debug("created test suite", "name", fileName+"_"+recordName)
		return ts, -1
	}
}

func getTestCaseByName(testSuite TestSuite, testCaseName string, logger slog.Logger) (TestCase, int) {
	for index, tc := range testSuite.TestCases {
		if tc.Name == testCaseName {
			logger.Debug("found test case", "name", testCaseName)
			return tc, index
		}
	}
	tc := TestCase{
		ClassName: "",
		Name:      testCaseName,
		Time:      fmt.Sprintf("%f", 0.01),
		Failures:  []Failure{},
	}
	logger.Debug("created test case", "name", testCaseName)
	return tc, -1
}

func main() {
	workspace, _ := os.Getwd()
	currentTime := time.Now()
	testSuiteTimestamp := currentTime.Format("2006-01-02T15:04:05")

	var settingsFilePath string
	flag.StringVar(&settingsFilePath, "settings_file", "config.json", "Путь к файлу настроек проекта")
	versionFlag := flag.Bool("version", false, "version number and exit")
	debugFlag := flag.Bool("debug", false, "show debug messages")
	flag.Parse()

	if *versionFlag {
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Build: %s\n", build)
		os.Exit(0)
	}

	configApp := loadConfigFromFile(settingsFilePath)
	
	logger := logging.CreateLogger(filepath.Join(workspace, configApp.OutputFileFolder, "edt_validator.log"), debugFlag)

	logger.Info("start application", "version", version, "build", build)

	files, err := os.ReadDir(filepath.Join(workspace, configApp.InputFileFolder))
	if err != nil {
		logger.Error("failed reading input file folder", "error", err.Error())
		return
	}

	isParentErrors := bool(configApp.SkipErrorsFile != "")
	var parentErrors [][]string
	var parentErrorsKeys map[string]struct{}
	if isParentErrors {
		parentErrors, parentErrorsKeys, err = readTSVFile(filepath.Join(workspace, configApp.InputFileFolder, configApp.SkipErrorsFile), logger)
		if err != nil {
			logger.Error("failed reading parent errors file", "error", err.Error())
			return
		}

		parentFileExtension := filepath.Ext(configApp.SkipErrorsFile)
		parentFileName := strings.TrimSuffix(configApp.SkipErrorsFile, parentFileExtension)
		createNewTestSuites(parentErrors, logger, false, parentErrorsKeys, testSuiteTimestamp, parentFileName, configApp.OutputFileFolder, configApp.SkipObjects, configApp.SkipCategories, configApp.SkipSignificanceCcategories, configApp.SkipErrorText)
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".tsv") {
			fileName := strings.TrimSuffix(file.Name(), ".tsv")

			logger.Debug("start processing file", "file", file.Name())

			tsvFile, err := os.Open(filepath.Join(configApp.InputFileFolder, file.Name()))
			if err != nil {
				logger.Error("failed opening tsv file", "file", file.Name(), "error", err.Error())
				panic(err)
			}
			defer tsvFile.Close()

			reader := csv.NewReader(tsvFile)
			reader.Comma = '\t'
			reader.LazyQuotes = true
			reader.FieldsPerRecord = -1
			
			records, err := reader.ReadAll()
			if err != nil {
				logger.Error("failed reading tsv file", "file", file.Name(), "error", err.Error())
				panic(err)
			}

			createNewTestSuites(records, logger, isParentErrors, parentErrorsKeys, testSuiteTimestamp, fileName, configApp.OutputFileFolder, configApp.SkipObjects, configApp.SkipCategories, configApp.SkipSignificanceCcategories, configApp.SkipErrorText)
		}
	}
	logger.Info("end application")
}

func loadConfigFromFile(settingsFilePath string) *config.AppConfig {
	configApp := config.NewAppConfig()
	config.LoadConfig(configApp, settingsFilePath)
	return configApp
}

func createNewTestSuites(records [][]string, logger *slog.Logger, isParentErrors bool, parentErrorsKeys map[string]struct{}, testSuiteTimestamp string, fileName string, outputFileFolder string, skipObjects []string, skipCategories []string, skipSignificanceCategories []string, skipErrorText []string) {

	testSuites := TestSuites{
		Time:      "0",
		Tests:     0,
		Errors:    0,
		Failures:  0,
		TestSuite: []TestSuite{},
	}

	for _, record := range records {
		if recordInSkipList(record, skipObjects, skipCategories, skipSignificanceCategories, skipErrorText) {
			logger.Debug("record in skip list", "record", record)
			continue
		}

		if isParentErrors {
			if recordInSkipErrorsList(record, parentErrorsKeys) {
				logger.Debug("record in skip errors list", "record", record)
				continue
			}
		}
		testSuite, indexTestSuite := getTestSuiteByName(testSuites, testSuiteTimestamp, fileName, record[1]+"_"+record[2], *logger)
		testCase, indexTestCase := getTestCaseByName(testSuite, record[5], *logger)

		failure := Failure{}
		failure.Type = record[2]
		failure.Message = record[1] + "; " + record[2] + "; " + record[4]
		failure.Text = record[5] + "; " + record[6] + "; " + record[7]
		logger.Debug("added failure", "type", failure.Type, "message", failure.Message, "text", failure.Text)
		testCase.Failures = append(testCase.Failures, failure)

		if indexTestCase == -1 {
			testSuite.TestCases = append(testSuite.TestCases, testCase)
		} else {
			testSuite.TestCases[indexTestCase] = testCase
		}
		testSuite.Tests++
		testSuite.Failures++

		if indexTestSuite == -1 {
			testSuites.TestSuite = append(testSuites.TestSuite, testSuite)
		} else {
			testSuites.TestSuite[indexTestSuite] = testSuite
		}
		testSuites.Tests++
		testSuites.Failures++
	}

	for index_ts, ts := range testSuites.TestSuite {
		var newTestCases []TestCase
		for _, tc := range ts.TestCases {
			if len(tc.Failures) > 1 {
				for i, f := range tc.Failures {
					newTestCase := TestCase{
						ClassName: tc.ClassName + "_unique_" + strconv.Itoa(i),
						Name:      tc.Name,
						Time:      tc.Time,
						Failures:  []Failure{f},
					}
					newTestCases = append(newTestCases, newTestCase)
					logger.Debug("added new test case", "name", newTestCase.Name, "class", newTestCase.ClassName)
				}
			} else {
				newTestCases = append(newTestCases, tc)
			}

		}
		testSuites.TestSuite[index_ts].TestCases = newTestCases
	}

	writeXMLData(logger, testSuites, fileName, outputFileFolder)
}

func writeXMLData(logger *slog.Logger, testSuites TestSuites, fileName string, outputFileFolder string) {
	xmlData, err := xml.MarshalIndent(testSuites, "", "    ")
	if err != nil {
		logger.Error("failed marshaling xml", "error", err.Error())
		panic(err)
	}

	xmlFile, err := os.Create(filepath.Join(outputFileFolder, fileName+".xml"))
	if err != nil {
		logger.Error("failed creating xml", "error", err.Error())
		panic(err)
	}
	defer xmlFile.Close()

	_, err = xmlFile.Write([]byte(xml.Header))
	if err != nil {
		logger.Error("failed writing header", "error", err.Error())
		panic(err)
	}

	_, err = xmlFile.Write(xmlData)
	if err != nil {
		logger.Error("failed writing xml data", "error", err.Error())
		panic(err)
	}
}

func recordInSkipErrorsList(record []string, parentErrorsKeys map[string]struct{}) bool {
	newSlice := make([]string, len(record))
	copy(newSlice, record)
	
	key := strings.Join(append(newSlice[1:2], newSlice[3:]...), "\t")
	if _, found := parentErrorsKeys[key]; !found {
		return false
	}
	return true
}

func recordInSkipList(record []string, skipObjects []string, skipCategories []string, skipSignificanceCategories []string, skipErrorText []string) bool {
	return recordInSkipObject(record, skipObjects) || recordInSkipCategory(record, skipCategories) || recordInSkipSignificanteCategories(record, skipSignificanceCategories) || recordInSkipErrorText(record, skipErrorText)
}

func recordInSkipCategory(record []string, skipCategories []string) bool {
	for _, skip := range skipCategories {
		if skip == record[2] {
			return true
		}
	}
	return false
}

func recordInSkipObject(record []string, skipObjects []string) bool {
	for _, skip := range skipObjects {
		if skip == record[5] {
			return true
		}

		if strings.Contains(record[5], skip) {
            return true
        }
	}
	return false
}

func recordInSkipSignificanteCategories(record []string, skipSignificanceCategories []string) bool {
	for _, skip := range skipSignificanceCategories {
		if skip == fmt.Sprintf("%s_%s", record[1], record[2]) {
			return true
		}
	}
	return false
}

func recordInSkipErrorText(record []string, skipErrorText []string) bool {
	for _, skip := range skipErrorText {
		if skip == record[7] {
			return true
		}
	}
	return false
}

func readTSVFile(filePath string, logger *slog.Logger) ([][]string, map[string]struct{}, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	
	reader := csv.NewReader(file)
	reader.Comma = '\t'
	reader.LazyQuotes = true
	reader.FieldsPerRecord = -1

	records, err := reader.ReadAll()
	if err != nil {
		logger.Error("failed reading parent tsv", "error", err.Error())
		panic(err)
	}

	lines := make(map[string]struct{})
	for _, record := range records {
		newSlice := make([]string, len(record))
		copy(newSlice, record)
		key := strings.Join(append(newSlice[1:2], newSlice[3:]...), "\t")
		lines[key] = struct{}{}
	}

	return records, lines, nil
}
