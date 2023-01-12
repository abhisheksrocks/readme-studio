package main_test

import (
	"errors"
	"testing"

	main "github.com/abhisheksrocks/readme-studio"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/abhisheksrocks/readme-studio/mocks"
)

func TestNewReadEnv(t *testing.T) {

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
		assertFunc      func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error)
	}{
		{
			testName:        "Only invalidFilePath",
			filePath:        invalidFilePath,
			exampleFilePath: validExampleFilePath,
			keyData: main.EnvKey{
				Key:     validKey,
				UsedFor: validUsedFor,
			},
			assertFunc: func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error) {
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorFileNotFound)
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
			assertFunc: func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error) {
				asserts.NotNil(err)
				asserts.EqualError(err, main.ReadEnvErrorExampleFileNotFound)
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
			assertFunc: func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error) {
				asserts.Nil(err)
				asserts.Equal(validValue, reVal.KeyVal.GetCacheValue())
				newValue, err := reVal.KeyVal.GetValue()
				asserts.Nil(err)
				asserts.Equal(validValue, newValue)
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
			assertFunc: func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error) {
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
			assertFunc: func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error) {
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
			assertFunc: func(t *testing.T, asserts *assert.Assertions, requires *require.Assertions, reVal *main.ReadEnv, err error) {
				asserts.Nil(err)
				asserts.Equal(validValue, reVal.KeyVal.GetCacheValue())
				newValue, err := reVal.KeyVal.GetValue()
				asserts.Nil(err)
				asserts.Equal(validValue, newValue)
			},
		},
	}

	for _, v := range tests {
		t.Run(v.testName, func(t *testing.T) {
			arrangeFunc := v.arrangeFunc
			assertFunc := v.assertFunc

			path := v.filePath
			examplePath := v.exampleFilePath
			keyData := v.keyData

			t.Parallel()
			requires := require.New(t)
			asserts := assert.New(t)

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
				assertFunc(t, asserts, requires, reVal, err)
			}
		})

	}

}
