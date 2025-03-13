package main

import (
	"encoding/xml"
	"os"
	"log/slog"
	"reflect"
	"strings"
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestGetTestSuiteByName_WhenTestSuitesIsEmpty(t *testing.T) {
	testSuites := TestSuites{
		XMLName:   xml.Name{},
		Time:      "0",
		Tests:     0,
		Errors:    0,
		Failures:  0,
		TestSuite: []TestSuite{},
	}
	testSuiteTimestamp := "2023-01-01T12:00:00"
	fileName := "test_file"
	recordName := "test_record"
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	expectedTestSuite := TestSuite{
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
	expectedIndex := -1

	actualTestSuite, actualIndex := getTestSuiteByName(testSuites, testSuiteTimestamp, fileName, recordName, *logger)

	if !reflect.DeepEqual(actualTestSuite, expectedTestSuite) {
		t.Errorf("expected test suite: %v, got: %v", expectedTestSuite, actualTestSuite)
	}
	if actualIndex != expectedIndex {
		t.Errorf("expected index: %d, got: %d", expectedIndex, actualIndex)
	}
}

func TestGetTestSuiteByName_ExistingTestSuite(t *testing.T) {
	testSuites := TestSuites{
		Time:      "0",
		Tests:     0,
		Errors:    0,
		Failures:  0,
		TestSuite: []TestSuite{
			{
				Name:      "filename_recordname",
				Timestamp: "2022-01-01T12:00:00",
				Time:      "0",
				Tests:     10,
				Errors:    0,
				Failures:  0,
				Skipped:   0,
				Properties: []Property{},
				TestCases:  []TestCase{},
			},
		},
	}
	testSuiteTimestamp := "2022-01-01T12:00:00"
	fileName := "filename"
	recordName := "recordname"
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	resultTestSuite, _ := getTestSuiteByName(testSuites, testSuiteTimestamp, fileName, recordName, *logger)

	expectedTestSuite := TestSuite{
		Name:      "filename_recordname",
		Timestamp: "2022-01-01T12:00:00",
		Time:      "0",
		Tests:     10,
		Errors:    0,
		Failures:  0,
		Skipped:   0,
		Properties: []Property{},
		TestCases:  []TestCase{},
	}
	assert.Equal(t, expectedTestSuite, resultTestSuite, "gettestsuitebyname should return an existing test suite")
}

func TestGetTestSuiteByName_NotExistingTestSuite(t *testing.T) {
	testSuites := TestSuites{
		Time:      "0",
		Tests:     0,
		Errors:    0,
		Failures:  0,
		TestSuite: []TestSuite{
			{
				Name:      "filename_recordname_1",
				Timestamp: "2022-01-01T12:00:00",
				Time:      "0",
				Tests:     10,
				Errors:    0,
				Failures:  0,
				Skipped:   0,
				Properties: []Property{},
				TestCases:  []TestCase{},
			},
		},
	}
	testSuiteTimestamp := "2022-01-01T12:00:00"
	fileName := "filename"
	recordName := "recordname"
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	resultTestSuite, _ := getTestSuiteByName(testSuites, testSuiteTimestamp, fileName, recordName, *logger)

	expectedTestSuite := TestSuite{
		Name:      "filename_recordname",
		Timestamp: "2022-01-01T12:00:00",
		Time:      "0",
		Tests:     0,
		Errors:    0,
		Failures:  0,
		Skipped:   0,
		Properties: []Property{},
		TestCases:  []TestCase{},
	}
	assert.Equal(t, expectedTestSuite, resultTestSuite, "gettestsuitebyname should return an existing test suite")
}

func TestGetTestCaseByName_ExistingTestCase(t *testing.T) {
	testSuite := TestSuite{
		TestCases: []TestCase{
			{
				Name: "TestCase1",
			},
			{
				Name: "TestCase2",
			},
		},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	testCaseName := "TestCase2"

	testCase, index := getTestCaseByName(testSuite, testCaseName, *logger)

	if index != 1 {
		t.Errorf("Expected index 1, but got %d", index)
	}
	if testCase.Name != testCaseName {
		t.Errorf("Expected testCaseName %s, but got %s", testCaseName, testCase.Name)
	}
}

func TestGetTestCaseByName_CreatesNewTestCase_WhenNameDoesNotMatch(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	testSuite := TestSuite{
		TestCases: []TestCase{
			{Name: "TestCase1"},
			{Name: "TestCase2"},
		},
	}
	testCaseName := "TestCase3"

	testCase, _ := getTestCaseByName(testSuite, testCaseName, *logger)

	if testCase.Name != testCaseName {
		t.Errorf("Expected testCase.Name to be %s, but got %s", testCaseName, testCase.Name)
	}
}

func TestGetTestCaseByName_EmptyTestCaseName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	testSuite := TestSuite{
		TestCases: []TestCase{},
	}
	testCaseName := ""

	testCase, index := getTestCaseByName(testSuite, testCaseName, *logger)

	if index != -1 {
		t.Errorf("Expected index to be -1, but got %d", index)
	}
	if testCase.Name != "" {
		t.Errorf("Expected empty testCase.Name, but got %s", testCase.Name)
	}
}

func TestRecordInSkipCategory(t *testing.T) {
	skipCategories := []string{"Предупреждение"}

    tests := []struct {
        input  []string
        output bool
    }{
		{input: []string{"Предупреждение", "Предупреждение", "Предупреждение"}, output: true},
		{input: []string{"A1", "B2", "B2"}, output: false},
    }

    for _, tt := range tests {
        result := recordInSkipCategory(tt.input, skipCategories)
        if result != tt.output {
            t.Errorf("recordInSkipCategory(%q) = %v, want %v", tt.input, result, tt.output)
        }
    }
}

func TestRecordInSkipObject(t *testing.T) {
	skipObjects := []string{"Справочник.Номенклатура.МодульОбъекта", ".Удалить"}

	tests := []struct {
		input  []string
		output bool
	}{
		{input: []string{"A1", "A1", "A1", "A1","A1", "A1"}, output: false},
		{input: []string{"A1", "A1", "A1", "A1","A1", "Справочник.Номенклатура.МодульОбъекта"}, output: true},
		{input: []string{"A1", "A1", "A1", "A1","A1", "Справочник.Удалить_Номенклатура.МодульОбъекта"}, output: true},
		{input: []string{"A1", "A1", "A1", "A1","A1", "Справочник.УдалитьНоменклатура.МодульОбъекта"}, output: true},
	}

	for _, tt := range tests {
		result := recordInSkipObject(tt.input, skipObjects)
		if result != tt.output {
			t.Errorf("recordInSkipObject(%q) = %v, want %v", tt.input, result, tt.output)
		}
	}
}

func TestRecordInSkipList(t *testing.T) {
	skipCategories := []string{"Предупреждение"}
	skipObjects := []string{"Справочник.Номенклатура.МодульОбъекта", ".Удалить"}
	skipSignificanteCategories := []string{"Значительная_Переносимость"}
	skipErrorText := []string{"Неподдерживаемый оператор [Web-клиент]"}

	tests := []struct {
		input  []string
		output bool
	}{
		{input: []string{"A1", "A1", "A1", "A1","A1", "A1", "", ""}, output: false},
		{input: []string{"A1", "A1", "A1", "A1","A1", "A1", "", "Неподдерживаемый оператор [Web-клиент]"}, output: true},
		{input: []string{"A1", "A1", "Предупреждение", "A1","A1", "A1", "", ""}, output: true},
		{input: []string{"A1", "Значительная", "Переносимость", "A1","A1", "A1", "", ""}, output: true},
		{input: []string{"A1", "A1", "A1", "A1","A1", "Справочник.Номенклатура.МодульОбъекта", "", ""}, output: true},
		{input: []string{"A1", "A1", "A1", "A1","A1", "Справочник.Удалить_Номенклатура.МодульОбъекта", "", ""}, output: true},
		{input: []string{"A1", "A1", "A1", "A1","A1", "Справочник.УдалитьНоменклатура.МодульОбъекта", "", ""}, output: true},
	}

	for _, tt := range tests {
		result := recordInSkipList(tt.input, skipObjects, skipCategories, skipSignificanteCategories, skipErrorText)
		if result != tt.output {
			t.Errorf("recordInSkipList(%q) = %v, want %v", tt.input, result, tt.output)
		}
	}
}

func TestRecordInSkipErrorsList(t *testing.T) {
	record := []string{"A1", "A1", "A1", "A1","A1", "A1", "A1", "A1"}
	key := strings.Join(append(record[1:2], record[3:]...), "\t")
	parentErrorsKeys := make(map[string]struct{})
	parentErrorsKeys[key] = struct{}{}

	tests := []struct {
		input  []string
		output bool
	}{
		{input: []string{"A1", "A1", "A1", "A1","A1", "A1", "A1", "A1"}, output: true},
		{input: []string{"A1", "A2", "A1", "A1","A1", "A1", "A1", "A1"}, output: false},
	}

	for _, tt := range tests {
		result := recordInSkipErrorsList(tt.input, parentErrorsKeys)
		if result != tt.output {
			t.Errorf("recordInSkipErrorsList(%q) = %v, want %v", tt.input, result, tt.output)
		}
	}
}

func TestReadTSVFile_FileNotFound(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	filePath := "non_existent_file.tsv"

	_, _, err := readTSVFile(filePath, logger)

	if err == nil {
		t.Errorf("expected error, but got nil")
	}
	if !os.IsNotExist(err) {
		t.Errorf("expected error to be os.IsNotExist, but got: %v", err)
	}
}

func TestReadTSVFile_EmptyInputFile(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	filePath := "empty_file.tsv"
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed creating empty file: %v", err)
	}
	file.Close()

	records, _, err := readTSVFile(filePath, logger)

	if err != nil {
		t.Fatalf("unexpected error reading empty file: %v", err)
	}
	if len(records) != 0 {
		t.Fatalf("expected empty records, got %d records", len(records))
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Fatalf("Can't remove file: %v", err)
		}
	}()
}

func TestReadTSVFileWithLines(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	filePath := "test_data_with_empty_lines.tsv"
	file, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("failed creating empty file: %v", err)
	}

	data := "2024-07-26T15:12:49+0300\tТривиальная\tСтандарты кодирования\tcf\tcom.e1c.v8codestyle.bsl:doc-comment-field-in-description-suggestion\tОбщийМодуль.WebAPI_Локализация.Модуль\tстрока 13\tВозможно Поле указано в описании\n2024-07-26T15:12:49+0300\tТривиальная\tСтандарты кодирования\tcf\tcom.e1c.v8codestyle.bsl:doc-comment-field-in-description-suggestion\tОбщийМодуль.WebAPI_Локализация.Модуль\tстрока 15\tВозможно Поле указано в описании\n2024-07-26T15:12:49+0300\tТривиальная\tСтандарты кодирования\tcf\tcom.e1c.v8codestyle.bsl:doc-comment-field-in-description-suggestion\tОбщийМодуль.WebAPI_Локализация.Модуль\tстрока 294\tВозможно Поле указано в описании\n"
	_, err = file.WriteString(data)
	if err != nil {
		t.Fatalf("failed writing to file: %v", err)
	}
	file.Close()

	records, _, err := readTSVFile(filePath, logger)

	if err != nil {
		t.Fatalf("readTSVFile should not return an error, but got: %v", err)
	}

	expectedRecords := [][]string{
		{"2024-07-26T15:12:49+0300","Тривиальная","Стандарты кодирования","cf","com.e1c.v8codestyle.bsl:doc-comment-field-in-description-suggestion","ОбщийМодуль.WebAPI_Локализация.Модуль","строка 13","Возможно Поле указано в описании"},
		{"2024-07-26T15:12:49+0300","Тривиальная","Стандарты кодирования","cf","com.e1c.v8codestyle.bsl:doc-comment-field-in-description-suggestion","ОбщийМодуль.WebAPI_Локализация.Модуль","строка 15","Возможно Поле указано в описании"},
		{"2024-07-26T15:12:49+0300","Тривиальная","Стандарты кодирования","cf","com.e1c.v8codestyle.bsl:doc-comment-field-in-description-suggestion","ОбщийМодуль.WebAPI_Локализация.Модуль","строка 294","Возможно Поле указано в описании"},
	}

	if !reflect.DeepEqual(records, expectedRecords) {
		t.Errorf("readTSVFile should return expected records, but got: %v", records)
	}

	defer func() {
		if err := os.Remove(filePath); err != nil {
			t.Fatalf("Can't remove file: %v", err)
		}
	}()
}
