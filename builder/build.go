package builder

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func BuildImage(image string, handler string, functionName string, registry string, runtime string, nocache bool, buildArgMap map[string]string, buildOptions []string, verbose bool) (string, error) {

	//Handler Path
	if _, err := os.Stat(handler); err != nil {
		fmt.Printf("Image: %s not built.\n", image)

		return "", errors.New(fmt.Sprintf("Unable to build %s, %s is an invalid path\n", image, handler))

	}

	if strings.ToLower(runtime) == "dockerfile" {
		fmt.Printf("Building: %s with Dockerfile. Please wait..\n", image)
	}

	var flagSlice []string
	flagSlice = append(flagSlice, "--build-arg", fmt.Sprintf("%s=%s", "REGISTRY", registry))
	for k, v := range buildArgMap {
		flagSlice = append(flagSlice, "--build-arg", fmt.Sprintf("%s=%s", k, v))
	}

	if len(buildOptions) > 0 {
		flagSlice = append(flagSlice, "--build-arg", fmt.Sprintf("%s=%s", "ADDITIONAL_PACKAGE", strings.Join(buildOptions, " ")))
	}

	if nocache {
		flagSlice = append(flagSlice, "--no-cache")
	}

	buildCmd := []string{"docker", "build"}
	buildCmd = append(buildCmd, flagSlice...)
	buildCmd = append(buildCmd, "-t", image, "--network=host", "../")

	var err error
	var result string
	if verbose {
		if err := ExecCommandPipe(handler, buildCmd, os.Stdout, os.Stderr); err != nil {
			return result, err
		}
	} else {
		if result, err = ExecCommand(handler, buildCmd); err != nil {
			return result, err
		}
	}
	fmt.Printf("Image: %s built in local environment.\n", image)
	return result, nil

}

// ExecCommand run a system command
func ExecCommand(tempPath string, cmdSlice []string) (string, error) {
	targetCmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	targetCmd.Dir = tempPath
	stdoutStderr, err := targetCmd.CombinedOutput()
	if err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", cmdSlice)
		return string(stdoutStderr), errors.New(errString)
	}
	return string(stdoutStderr), nil
}

func ExecCommandPipe(tempPath string, cmdSlice []string, stdout, stderr io.Writer) error {
	targetCmd := exec.Command(cmdSlice[0], cmdSlice[1:]...)
	targetCmd.Dir = tempPath
	targetCmd.Stdout = stdout
	targetCmd.Stderr = stderr
	targetCmd.Start()
	err := targetCmd.Wait()
	if err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", cmdSlice)
		return errors.New(errString)
	}
	return nil
}
