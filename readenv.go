package main

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type EnvKeyValue struct {
	Key     string
	UsedFor string

	Val string
}

const (
	ReadEnvErrorUndefined           = "unknown error"
	ReadEnvErrorValueNotFound       = "value not found"
	ReadEnvErrorFileNotFound        = "file not found"
	ReadEnvErrorExampleFileNotFound = "example file not found"
)

func (e ReadEnv) GetCacheValue() string {
	return e.KeyVal.Val
}

func (e *ReadEnv) GetValue() (string, error) {
	v := os.Getenv(e.KeyVal.Key)
	if v == "" {
		return "", errors.New(ReadEnvErrorValueNotFound)
	}
	e.KeyVal.Val = v
	return v, nil
}

type ReadEnv struct {
	FilePath        string
	ExampleFilePath string
	KeyVal          *EnvKeyValue
}

func NewReadEnv(envPath string, exampleEnvPath string, key string, usedFor string) (r *ReadEnv, err error) {
	r = &ReadEnv{
		FilePath:        envPath,
		ExampleFilePath: exampleEnvPath,
		KeyVal: &EnvKeyValue{
			Key:     key,
			UsedFor: usedFor,
			Val:     "",
		},
	}

	if notEmpty(r.ExampleFilePath) {
		if _, err := os.Stat(r.ExampleFilePath); errors.Is(err, os.ErrNotExist) {
			err = errors.New(ReadEnvErrorExampleFileNotFound)
			return r, err
		}
	}

	if notEmpty(r.FilePath) {
		err = godotenv.Load(r.FilePath)
		if err != nil {
			err = errors.New(ReadEnvErrorFileNotFound)
			return r, err
		}
	}

	_, err = r.GetValue()

	return r, err
}
