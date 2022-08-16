package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	flag.Parse()

	envDirPath := mustEnvDirPath()

	env, err := ReadDir(envDirPath)
	if err != nil {
		doFatal(err.Error())
	}

	cmd := buildCmd()

	code := RunCmd(cmd, env)

	os.Exit(code)
}

func buildCmd() []string {
	return flag.Args()[1:]
}

func mustEnvDirPath() string {
	envDirPath := flag.Arg(0)

	if envDirPath == "" {
		doFatal("env dir not passed")
	}

	return envDirPath
}

func doFatal(str string) {
	log.Fatalf(str)
}

func wrapErr(err error, str string) error {
	return fmt.Errorf("%s: %w", str, err)
}
