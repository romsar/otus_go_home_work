package main

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

const (
	cutSet         = " 	"
	forbiddenChars = "="
)

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envFiles, err := os.ReadDir(dir)
	if err != nil {
		return nil, wrapErr(err, "cannot read dir: "+dir)
	}

	envs := make(Environment, len(envFiles))

	for _, envFile := range envFiles {
		path := filePath(dir, envFile)
		fileContent, err := os.ReadFile(path)
		if err != nil {
			return nil, wrapErr(err, "cannot read env file: "+path)
		}

		line := firstLine(fileContent)

		envKey, envValue := keyValue(envFile, line)

		if strings.ContainsAny(envValue, forbiddenChars) {
			return nil, wrapErr(
				err,
				fmt.Sprintf("env file %s contains forbidden characters: %s", envFile.Name(), forbiddenChars),
			)
		}

		envValue = prepareEnvValue(envValue)
		needRm := needRemove(envValue)

		envs[envKey] = EnvValue{
			Value:      envValue,
			NeedRemove: needRm,
		}
	}

	return envs, nil
}

func needRemove(envValue string) bool {
	return envValue == ""
}

func prepareEnvValue(envValue string) string {
	return replaceNilsWithLineBreak(envValue)
}

func replaceNilsWithLineBreak(envValue string) string {
	return string(bytes.ReplaceAll([]byte(envValue), []byte{0x00}, []byte("\n")))
}

func keyValue(envFile os.DirEntry, line string) (string, string) {
	envKey, envValue := envFile.Name(), strings.TrimRight(line, cutSet)
	return envKey, envValue
}

func firstLine(fileContent []byte) string {
	val := string(fileContent)
	lines := strings.Split(val, "\n")
	firstLine := lines[0]
	return firstLine
}

func filePath(dir string, envFile os.DirEntry) string {
	return filepath.Join(dir, envFile.Name())
}
