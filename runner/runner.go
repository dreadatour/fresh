package runner

import (
	"io"
	"os/exec"
)

func run() bool {
	runnerLog("Running (%s)...", buildPath())

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
		err := cmd.Wait()
		if err != nil && err.Error() != "signal: killed" {
			runnerLog("Running error: %s", err)
			startChannel <- "/"
		}
	}()

	go func() {
		<-stopChannel
		pid := cmd.Process.Pid
		runnerLog("Killing PID %d", pid)
		cmd.Process.Kill()
	}()

	return true
}
