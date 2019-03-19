package builder

import (
	"errors"
	"fmt"
	"os/exec"
)

func RenameImage(image string, registry string, functionName string) (string, error) {
	renameCmd := []string{"docker", "tag"}
	renameCmd = append(renameCmd, image)
	renameCmd = append(renameCmd, registry+"/"+functionName)

	targetCmd := exec.Command(renameCmd[0], renameCmd[1:]...)
	if err := targetCmd.Run(); err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", renameCmd)
		return errString, errors.New(errString)
	}

	return "Renaming Image successfully", nil
}

func RemoveOldImage(image string) (string, error) {
	removeCmd := []string{"docker", "rmi"}
	removeCmd = append(removeCmd, image)

	targetCmd := exec.Command(removeCmd[0], removeCmd[1:]...)
	if err := targetCmd.Run(); err != nil {
		errString := fmt.Sprintf("ERROR - Could not execute command: %s", removeCmd)
		return errString, errors.New(errString)
	}

	return "Removing previous image successfully", nil
}
