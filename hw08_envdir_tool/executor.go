package main

import (
	"log"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, envs Environment) int {
	envStr := envStrings(envs)

	cmdToExec, cmdArgs := executeCmd(cmd)

	execCmd, err := doRun(cmdToExec, cmdArgs, envStr)
	if err != nil {
		log.Println(wrapErr(err, "error while executing command"))
	}

	return returnCode(execCmd)
}

func returnCode(execCmd *exec.Cmd) int {
	return execCmd.ProcessState.ExitCode()
}

func doRun(cmdToExec string, cmdArgs []string, envStr []string) (*exec.Cmd, error) {
	execCmd := exec.Command(cmdToExec, cmdArgs...)
	execCmd.Stdout = os.Stdout
	execCmd.Stderr = os.Stderr
	execCmd.Env = os.Environ()
	execCmd.Env = append(execCmd.Env, envStr...)
	return execCmd, execCmd.Run()
}

func executeCmd(cmds []string) (string, []string) {
	return cmds[0], cmds[1:]
}

func envStrings(envs Environment) []string {
	envStr := make([]string, 0, len(envs))
	for key, val := range envs {
		envStr = append(envStr, key+"="+val.Value)
	}
	return envStr
}
