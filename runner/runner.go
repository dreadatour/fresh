package runner

import (
	"fmt"
	"io"
	"os/exec"
)

func run() bool {
	runnerLog(fmt.Sprintf("Running (%s)...", buildPath()))

	cmd := exec.Command(buildPath())

	stderr, err := cmd.StderrPipe()
	if err != nil {
		fatal(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fatal(err)
	}

	err = cmd.Start()
	if err != nil {
		fatal(err)
	}

	go io.Copy(appLogWriter{}, stderr)
	go io.Copy(appLogWriter{}, stdout)

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)

		err := cmd.Process.Kill()
		if err != nil {
			runnerLog("Killing process error: %s", err)
		}

		err = cmd.Wait()
		if err != nil {
			runnerLog("Waiting for process killed error: %s", err)
		}
	}()

	return true
}
