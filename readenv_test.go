package main_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	main "github.com/abhisheksrocks/readme-studio"
	"golang.org/x/exp/slices"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/abhisheksrocks/readme-studio/mocks"
)

type UnitTestReadEnvSuite struct {
	suite.Suite
}

func TestUnitTestReadEnvSuite(t *testing.T) {
	suite.Run(t, new(UnitTestReadEnvSuite))
}

func (uts *UnitTestReadEnvSuite) TestNewReadEnv() {

	const (
		mockFunctionNameFileExist = "FileExist"
		mockFunctionNameLoad      = "Load"
		mockFunctionNameGetenv    = "Getenv"
	)

	const (
		emptyExampleFilePath   = ""
		validExampleFilePath   = "validExampleFileName"
		invalidExampleFilePath = "invalidExampleFileName"
	)

	const (
		emptyFilePath   = ""
		validFilePath   = "validFileName"
		invalidFilePath = "invalidFileName"
	)

	const (
		validKey = "validKey"
		emptyKey = ""
	)

	const (
		validValue = "validValue"
		emptyValue = ""
	)

	const (
		validUsedFor = "validUsedFor"
	)

	var tests = []struct {
		testName        string
		filePath        string
		exampleFilePath string
		keyData         main.EnvKey
		arrangeFunc     func(t *testing.T, m *mocks.ReadEnvEnvironment)
		assertFunc      func(t *testing.T, reVal *main.ReadEnv, err error)
	}{
		{
			testName:        "Only invalidFilePath",
			filePath:        invalidFilePath,
			exampleFilePath: validExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorFileNotFound)
			},
		},
		{
			testName:        "Only emptyFilePath",
			filePath:        emptyFilePath,
			exampleFilePath: validExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.Nil(err)
				asserts.Equal(validValue, reVal.KeyVal.GetCacheValue())
				newValue, err := reVal.KeyVal.GetValue()
				asserts.Nil(err)
				asserts.Equal(validValue, newValue)
			},
		},
		{
			testName:        "emptyFilePath & no default system variable",
			filePath:        emptyFilePath,
			exampleFilePath: validExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			arrangeFunc: func(t *testing.T, m *mocks.ReadEnvEnvironment) {
				m.On(mockFunctionNameGetenv, validKey).Return(emptyValue)
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorValueNotFound)
			},
		},
		{
			testName:        "Only invalidExampleFilePath",
			filePath:        validFilePath,
			exampleFilePath: invalidExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorExampleFileNotFound)
			},
		},
		{
			testName:        "Only emptyExampleFilePath",
			filePath:        validFilePath,
			exampleFilePath: emptyExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.Nil(err)
				asserts.Equal(validValue, reVal.KeyVal.GetCacheValue())
				newValue, err := reVal.KeyVal.GetValue()
				asserts.Nil(err)
				asserts.Equal(validValue, newValue)
			},
		},
		{
			testName:        "emptyFilePath and emptyExampleFilePath",
			filePath:        emptyFilePath,
			exampleFilePath: emptyExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.Nil(err)
				asserts.Equal(validValue, reVal.KeyVal.GetCacheValue())
				newValue, err := reVal.KeyVal.GetValue()
				asserts.Nil(err)
				asserts.Equal(validValue, newValue)
			},
		},
		{
			testName:        "invalidFilePath and invalidExampleFilePath",
			filePath:        invalidFilePath,
			exampleFilePath: invalidExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorExampleFileNotFound)
			},
		},
		{
			testName:        "validKey but emptyValue",
			filePath:        emptyFilePath,
			exampleFilePath: emptyExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			arrangeFunc: func(t *testing.T, m *mocks.ReadEnvEnvironment) {
				m.On(mockFunctionNameGetenv, validKey).Return(emptyValue)
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorValueNotFound)
			},
		},
		{
			testName:        "emptyKey",
			filePath:        emptyFilePath,
			exampleFilePath: emptyExampleFilePath,
			keyData: main.EnvKey{
				Key:     emptyKey,
				UsedFor: validUsedFor,
			},
			arrangeFunc: func(t *testing.T, m *mocks.ReadEnvEnvironment) {
				m.AssertNumberOfCalls(t, mockFunctionNameGetenv, 0)
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorInvalidKey)
			},
		},
		{
			testName:        "success",
			filePath:        validFilePath,
			exampleFilePath: validExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, reVal *main.ReadEnv, err error) {
				asserts := assert.New(t)
				// requires := require.New(t)
				asserts.Nil(err)
				asserts.Equal(validValue, reVal.KeyVal.GetCacheValue())
				newValue, err := reVal.KeyVal.GetValue()
				asserts.Nil(err)
				asserts.Equal(validValue, newValue)
			},
		},
	}

	for _, v := range tests {
		uts.Run(v.testName, func() {
			// t *testing.T
			t := uts.T()
			arrangeFunc := v.arrangeFunc
			assertFunc := v.assertFunc

			path := v.filePath
			examplePath := v.exampleFilePath
			keyData := v.keyData

			t.Parallel()

			// Arrange
			mocker := mocks.ReadEnvEnvironment{}

			if arrangeFunc != nil {
				arrangeFunc(t, &mocker)
			}

			mocker.On(mockFunctionNameFileExist, validExampleFilePath).Return(true)
			mocker.On(mockFunctionNameFileExist, invalidExampleFilePath).Return(false)
			mocker.On(mockFunctionNameLoad, validFilePath).Return(nil)
			mocker.On(mockFunctionNameLoad, invalidFilePath).Return(errors.New("file not found"))
			mocker.On(mockFunctionNameGetenv, validKey).Return(validValue)

			mocker.AssertNotCalled(t, mockFunctionNameGetenv, emptyKey)
			mocker.AssertNotCalled(t, mockFunctionNameFileExist, emptyExampleFilePath)
			mocker.AssertNotCalled(t, mockFunctionNameFileExist, emptyFilePath)
			mocker.AssertNotCalled(t, mockFunctionNameLoad, emptyFilePath)
			mocker.AssertNotCalled(t, mockFunctionNameLoad, validExampleFilePath)
			mocker.AssertNotCalled(t, mockFunctionNameLoad, invalidExampleFilePath)
			mocker.AssertNotCalled(t, mockFunctionNameLoad, emptyExampleFilePath)

			// Act
			reVal, err := main.NewReadEnv(path, examplePath, keyData, &mocker)

			// Assert
			if assertFunc != nil {
				assertFunc(t, reVal, err)
			}
		})

	}

}

type IntTestReadEnvSuite struct {
	suite.Suite
	tempDirectoryPath string
}

func TestIntTestReadEnvSuite(t *testing.T) {
	suite.Run(t, new(IntTestReadEnvSuite))
}

func (its *IntTestReadEnvSuite) SetupSuite() {
	location, err := os.Getwd()
	if err != nil {
		its.T().Log("err:\n", err)
		its.T().Skipf("Cound't get the current working directory!")
		return
	}

	err = os.Mkdir(filepath.Join(location, tmpFileDirectoryName), os.ModePerm)
	if err != nil {
		its.T().Log("err:\n", err)
		its.T().Skipf("Cound't create temporary directory!")
		return
	}
	its.tempDirectoryPath = filepath.Join(location, tmpFileDirectoryName)

	intTests = []intTestData{
		{
			testName:           "Test_OnlyInvalidFilePath",
			envFilePath:        filepath.Join(its.tempDirectoryPath, invalidEnvFilePath),
			exampleEnvFilePath: filepath.Join(its.tempDirectoryPath, validExampleEnvFilePath),
			keyText:            validKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				exampleFile, err := os.Create(intTestData.exampleEnvFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create example file")
					return
				}
				defer exampleFile.Close()

				exampleFile.WriteString(exampleKey + "=" + exampleValue)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.NotNil(err)
				its.EqualError(err, main.ReadEnvErrorFileNotFound)
			},
			// cleanup: func(its *IntTestReadEnvSuite, intTestData intTestData) {
			// 	os.Remove(intTestData.exampleEnvFilePath)
			// },
		},
		{
			testName:           "Test_OnlyEmptyFilePath",
			envFilePath:        "",
			exampleEnvFilePath: filepath.Join(its.tempDirectoryPath, validExampleEnvFilePath),
			keyText:            validKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				exampleFile, err := os.Create(intTestData.exampleEnvFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create example file")
					return
				}
				defer exampleFile.Close()

				exampleFile.WriteString(exampleKey + "=" + exampleValue)
				its.T().Setenv(validKey, validValue)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.Nil(err)
				its.Equal(validValue, result.KeyVal.GetCacheValue())
				newValue, err := result.KeyVal.GetValue()
				its.Nil(err)
				its.Equal(validValue, newValue)
			},
			// cleanup: func(its *IntTestReadEnvSuite, intTestData intTestData) {
			// 	os.Remove(intTestData.exampleEnvFilePath)
			// },
		},
		{
			testName:           "Test_EmptyFilePathButNoDefaultSystemVariables",
			envFilePath:        emptyEnvFilePath,
			exampleEnvFilePath: filepath.Join(its.tempDirectoryPath, validExampleEnvFilePath),
			keyText:            validKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				exampleFile, err := os.Create(intTestData.exampleEnvFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create example file")
					return
				}
				defer exampleFile.Close()
				exampleFile.WriteString(exampleKey + "=" + exampleValue)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.NotNil(err)
				its.EqualError(err, main.ReadEnvErrorValueNotFound)
			},
			// cleanup: func(its *IntTestReadEnvSuite, intTestData intTestData) {
			// 	os.Remove(intTestData.exampleEnvFilePath)
			// },
		},
		{
			testName:           "Test_OnlyInvalidExampleFilePath",
			envFilePath:        filepath.Join(its.tempDirectoryPath, validEnvFilePath),
			exampleEnvFilePath: filepath.Join(its.tempDirectoryPath, invalidExampleEnvFilePath),
			keyText:            validKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				file, err := os.Create(intTestData.envFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create file")
					return
				}
				defer file.Close()
				file.WriteString(validKey + "=" + validValue)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.NotNil(err)
				its.EqualError(err, main.ReadEnvErrorExampleFileNotFound)
			},
		},
		{
			testName:           "Test_OnlyEmptyExampleFilePath",
			envFilePath:        filepath.Join(its.tempDirectoryPath, validEnvFilePath),
			exampleEnvFilePath: emptyExampleEnvFilePath,
			keyText:            validKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				file, err := os.Create(intTestData.envFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create file")
					return
				}
				file.WriteString(intTestData.keyText + "=" + intTestData.valueText)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.Nil(err)
				its.Equal(validValue, result.KeyVal.GetCacheValue())
				newValue, err := result.KeyVal.GetValue()
				its.Nil(err)
				its.Equal(validValue, newValue)
			},
		},
		{
			testName:           "Test_EmptyPathsAndValidKeyAndValidValue",
			envFilePath:        emptyEnvFilePath,
			exampleEnvFilePath: emptyExampleEnvFilePath,
			keyText:            validKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				its.T().Setenv(validKey, validValue)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.Nil(err)
				its.Equal(validValue, result.KeyVal.GetCacheValue())
				newValue, err := result.KeyVal.GetValue()
				its.Nil(err)
				its.Equal(validValue, newValue)
			},
		},
		{
			testName:           "Test_EmptyPathsAndValidKeyButEmptyValue",
			envFilePath:        emptyEnvFilePath,
			exampleEnvFilePath: emptyExampleEnvFilePath,
			keyText:            validKey,
			valueText:          emptyValue,
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.NotNil(err)
				its.EqualError(err, main.ReadEnvErrorValueNotFound)
			},
		},
		{
			testName:           "Test_OnlyEmptyKey",
			envFilePath:        filepath.Join(its.tempDirectoryPath, validEnvFilePath),
			exampleEnvFilePath: filepath.Join(its.tempDirectoryPath, validEnvFilePath),
			keyText:            emptyKey,
			valueText:          validValue,
			beforeTest: func(its *IntTestReadEnvSuite, intTestData intTestData) {
				exampleFile, err := os.Create(intTestData.exampleEnvFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create example file")
					return
				}
				defer exampleFile.Close()
				file, err := os.Create(intTestData.envFilePath)
				if err != nil {
					its.T().Log("err:\n", err)
					its.T().Skipf("Cound't create example file")
					return
				}
				defer file.Close()
				file.WriteString(validKey + "=" + validValue)
				exampleFile.WriteString(exampleKey + "=" + exampleValue)
			},
			assertions: func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error) {
				its.NotNil(err)
				its.EqualError(err, main.ReadEnvErrorInvalidKey)
			},
		},
	}
}

func (its *IntTestReadEnvSuite) SetupTest() {
	os.Mkdir(its.tempDirectoryPath, os.ModePerm)
}

func (its *IntTestReadEnvSuite) TearDownTest() {
	os.RemoveAll(its.tempDirectoryPath)
}

func (its *IntTestReadEnvSuite) TearDownSuite() {
	os.RemoveAll(its.tempDirectoryPath)
}

var intTests []intTestData

type intTestData struct {
	testName           string
	envFilePath        string
	exampleEnvFilePath string
	keyText            string
	valueText          string
	beforeTest         func(its *IntTestReadEnvSuite, intTestData intTestData)
	assertions         func(its *IntTestReadEnvSuite, intTestData intTestData, result *main.ReadEnv, err error)
	cleanup            func(its *IntTestReadEnvSuite, intTestData intTestData)
}

const tmpFileDirectoryName = "tmp"

const (
	validEnvFilePath   = "./validEnvFilePath.env"
	invalidEnvFilePath = "./invalidEnvFilePath.env"
	emptyEnvFilePath   = ""
)

const (
	validExampleEnvFilePath   = "./validExampleEnvFilePath.env"
	invalidExampleEnvFilePath = "./invalidExampleEnvFilePath.env"
	emptyExampleEnvFilePath   = ""
)

const (
	validKey = "VALID_KEY"
	emptyKey = ""
)

const usedForPlaceholder = "some data"

const (
	validValue = "VALID_VALUE"
	emptyValue = ""
)

const (
	exampleKey   = "EXAMPLE_KEY"
	exampleValue = "EXAMPLE_VALUE"
)

func (its *IntTestReadEnvSuite) BeforeTest(suiteName, testName string) {
	for _, v := range intTests {
		if v.testName == testName && v.beforeTest != nil {
			v.beforeTest(its, v)
		}
	}
}

func (its *IntTestReadEnvSuite) AfterTest(suiteName, testName string) {
	for _, v := range intTests {
		if v.testName == testName && v.cleanup != nil {
			v.cleanup(its, v)
		}
	}
}

func (its *IntTestReadEnvSuite) Test_OnlyInvalidFilePath() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_OnlyInvalidFilePath"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_OnlyInvalidFilePath NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_OnlyEmptyFilePath() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_OnlyEmptyFilePath"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_OnlyEmptyFilePath NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_EmptyFilePathButNoDefaultSystemVariables() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_EmptyFilePathButNoDefaultSystemVariables"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_EmptyFilePathButNoDefaultSystemVariables NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_OnlyInvalidExampleFilePath() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_OnlyInvalidExampleFilePath"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_OnlyInvalidExampleFilePath NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_OnlyEmptyExampleFilePath() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_OnlyEmptyExampleFilePath"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_OnlyEmptyExampleFilePath NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_EmptyPathsAndValidKeyAndValidValue() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_EmptyPathsAndValidKeyAndValidValue"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_EmptyPathsAndValidKeyAndValidValue NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_EmptyPathsAndValidKeyButEmptyValue() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_EmptyPathsAndValidKeyButEmptyValue"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_EmptyPathsAndValidKeyButEmptyValue NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}

func (its *IntTestReadEnvSuite) Test_OnlyEmptyKey() {

	testIndex := slices.IndexFunc(intTests, func(c intTestData) bool {
		return c.testName == "Test_OnlyEmptyKey"
	})

	if testIndex == -1 {
		its.T().Skipf("test name -> Test_OnlyEmptyKey NOT FOUND! in [intTests]\n")
		return
	}

	refData := intTests[testIndex]

	reVal, err := main.NewReadEnv(refData.envFilePath, refData.exampleEnvFilePath, main.EnvKey{
		Key:     refData.keyText,
		UsedFor: usedForPlaceholder,
	}, &main.DefReadEnvEnvironment{})

	refData.assertions(its, refData, reVal, err)
}
