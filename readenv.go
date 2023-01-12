package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type ReadEnvEnvironment interface {
	Load(filenames ...string) error
	FileExist(filename string) bool
	Getenv(key string) string
}

type DefReadEnvEnvironment struct{}

func (*DefReadEnvEnvironment) Load(filenames ...string) error {
	return godotenv.Load(filenames...)
}

func (*DefReadEnvEnvironment) FileExist(filename string) bool {
	if _, err := os.Stat(filename); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

func (*DefReadEnvEnvironment) Getenv(key string) string {
	return os.Getenv(key)
}

type EnvKey struct {
	Key     string
	UsedFor string
}

type envKeyValue struct {
	KeyData     EnvKey
	val         string
	environment ReadEnvEnvironment
}

const (
	ReadEnvErrorUndefined           = "unknown error"
	ReadEnvErrorInvalidKey          = "invalid key"
	ReadEnvErrorValueNotFound       = "value not found"
	ReadEnvErrorFileNotFound        = "file not found"
	ReadEnvErrorExampleFileNotFound = "example file not found"
)

func (e envKeyValue) GetCacheValue() string {
	return e.val
}

func (e *envKeyValue) GetValue() (string, error) {
	if empty(e.KeyData.Key) {
		return "", errors.New(ReadEnvErrorInvalidKey)
	}
	v := e.environment.Getenv(e.KeyData.Key)
	if empty(v) {
		return "", errors.New(ReadEnvErrorValueNotFound)
	}
	e.val = v
	return v, nil
}

type ReadEnv struct {
	FilePath        string
	ExampleFilePath string
	KeyVal          *envKeyValue
}

func NewReadEnv(envPath string, exampleEnvPath string, key EnvKey, environment ReadEnvEnvironment) (r *ReadEnv, err error) {
	r = &ReadEnv{
		FilePath:        envPath,
		ExampleFilePath: exampleEnvPath,
		KeyVal: &envKeyValue{
			KeyData:     key,
			val:         "",
			environment: environment,
		},
	}

	if notEmpty(r.ExampleFilePath) {
		if !environment.FileExist(r.ExampleFilePath) {
			err = errors.New(ReadEnvErrorExampleFileNotFound)
			return r, err
		}
	}

	if notEmpty(r.FilePath) {
		err = environment.Load(r.FilePath)
		if err != nil {
			err = errors.New(ReadEnvErrorFileNotFound)
			return r, err
		}
	}

	_, err = r.KeyVal.GetValue()

	return r, err
}
