package runtime

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/keti-openfx/openfx-cli/config"
	"github.com/keti-openfx/openfx-cli/versioncontrol"
	"github.com/spf13/cobra"
)

var (
	runtimeDir string
)

var RuntimeCmd = &cobra.Command{
	Use:     "runtime SUBCOMMAND",
	Aliases: []string{"r"},
	Short:   "runtime specific operations",
	Long: `
	runtime command create, list the runtime used in creation openfx function
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func init() {

	RuntimeCmd.PersistentFlags().StringVarP(&runtimeDir, "dir", "d", "./runtime", "")
	RuntimeCmd.AddCommand(listCmd)
}

type Runtimes struct {
	Runtimes map[string]Runtime `yaml:runtimes`
}

type Runtime struct {
	Name    string         `yaml:"-"`
	Dir     string         `yaml:"dir"`
	Handler config.Handler `yaml:"handler,omitempty"`
	//ARG in Dockerfile
	BuildPackages []string `yaml:"build_packages,omitempty"`
	BuildArgs     []string `yaml:"build_args,omitempty"`
}

func ValidateEnv(val string) (string, error) {
	arr := strings.Split(val, "=")
	if arr[0] == "" {
		return "", fmt.Errorf("invalid environment variable: %s", val)
	}
	if len(arr) > 1 {
		return val, nil
	}
	return fmt.Sprintf("%s=%s", val, os.Getenv(val)), nil
}

func ExistFileOrDir(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}

func DownloadRuntimes(path, runtimeURL string) error {
	args := map[string]string{"dir": path, "repo": runtimeURL}
        fmt.Println(runtimeURL)
	if err := versioncontrol.GitClone.Invoke(".", args); err != nil {
		return err
	}

	return nil
}

func GetRuntime(functionName, runtimeName, runtimeDir string) (*Runtime, error) {
	runtimes, err := ReadRuntimeList(runtimeDir)
	if err != nil {
		return nil, err
	}

	for name, v := range runtimes.Runtimes {
		v.Name = name
		if v.Name == runtimeName {
			runtimePath := filepath.Join(runtimeDir, v.Dir)
			err := copyDir(runtimePath, "./"+functionName)
			if err != nil {
				return nil, err
			}
			return &v, nil
		}

	}

	return nil, errors.New("Not found runtime:" + runtimeName)
}

func ReadRuntimeList(dir string) (*Runtimes, error) {
	file := filepath.Join(dir, "list.yml")
	fileData, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var runtimes Runtimes
	err = yaml.Unmarshal(fileData, &runtimes)
	if err != nil {
		fmt.Printf("Error with YAML Config file\n")
		return nil, err
	}

	return &runtimes, nil
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}
