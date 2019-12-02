package builder

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"bytes"
	
	"github.com/keti-openfx/openfx-cli/cmd/log"
)

func RunImage(image string, handler string, functionName string) (string, error) {
	//Handler Path
	if _, err := os.Stat(handler); err != nil {
		fmt.Printf("Image: %s cannot run.\n", image)

		return "", errors.New(fmt.Sprintf("Unable to run %s, %s is an invalid path\n", image, handler))
	}

	flagSlice := []string{"-d", "-p", "50051:50051", "--name"}
	runCmd := []string{"docker", "run"}
	runCmd = append(runCmd, flagSlice...)
	runCmd = append(runCmd, functionName)
	runCmd = append(runCmd, image)

	targetCmd := exec.Command(runCmd[0], runCmd[1:]...)
	if err := targetCmd.Run(); err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", runCmd)
		return errString, errors.New(errString)
	}

	return "Running container successfully", nil
}

func StopContainer(containerName string) (string, error) {
	stopCmd := []string{"docker", "stop"}
	stopCmd = append(stopCmd, containerName)

	targetCmd := exec.Command(stopCmd[0], stopCmd[1:]...)
	if err := targetCmd.Run(); err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", stopCmd)
		return errString, errors.New(errString)
	}

	return "Stopping container successfully", nil
}

func RemoveContainer(containerName string) (string, error) {
	removeCmd := []string{"docker", "rm"}
	removeCmd = append(removeCmd, containerName)

	targetCmd := exec.Command(removeCmd[0], removeCmd[1:]...)
	if err := targetCmd.Run(); err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", removeCmd)
		return errString, errors.New(errString)
	}

	return "Removing container successfully", nil
}

func CheckImgRunning() string {
	var output bytes.Buffer

	findCmd := []string{"docker", "ps", "-a", "--format", "{{.Names}}"}

	execCmd := exec.Command(findCmd[0], findCmd[1:]...)
	execCmd.Stdout = &output

	if err := execCmd.Run(); err != nil {
		log.Fatal(err)
	}

	return output.String()
}
