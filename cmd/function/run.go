package function

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"syscall"
	"strings"
	
	"github.com/keti-openfx/openfx-cli/builder"
	"github.com/keti-openfx/openfx-cli/cmd/log"
	"github.com/keti-openfx/openfx-cli/config"
	watcher "github.com/keti-openfx/openfx-cli/pb"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

func init() {
	runCmd.Flags().StringVarP(&configFile, "config", "f", "", "Path to YAML config file describing function(s)")
}

var runCmd = &cobra.Command{
	Use:   `run <FUNCTION_NAME>`,
	Short: "Run OpenFx Image in local",
	Long: `
	Run OpenFx Image that created when execute "build" command for debugging in local
	`,
	Example: `
	openfx-cli function run echo-service -f config.yaml
	echo "hi" | openfx-cli function run echo-service
	`,
	PreRunE: preRunRun,
	RunE: func(cmd *cobra.Command, args []string) error {
		functionName = args[0]

		if err := runRun(); err != nil {
			fmt.Println(err.Error())
		}
		return nil
	},
}

func preRunRun(cmd *cobra.Command, args []string) error {
	if configFile == "" {
		files, err := ioutil.ReadDir("./")
		if err != nil {
			log.Fatal(err)
		}

		for _, f := range files{
			if strings.Contains(f.Name(), "yaml") {
				configFile = f.Name()
			}
		}

		if err := parseConfigFile(); err != nil {
			return err
		}
	} else {
		if err := parseConfigFile(); err != nil {
			return err
		}
	}

	if len(args) < 1 {
		log.Fatal("please provide a name of the function\n")
	}

	functionName = args[0]

	return nil
}

func run(function config.Function) error {
	result, err := builder.RunImage(function.Image, function.Handler.Dir, function.Name)
	if err != nil {
		log.Print(result)
		return err
	}

	return nil
}

func stop(function config.Function) error {
	result, err := builder.StopContainer(function.Name)
	if err != nil {
		log.Print(result)
		return err
	}

	return nil
}

func remove(function config.Function) error {
	result, err := builder.RemoveContainer(function.Name)
	if err != nil {
		log.Print(result)
		return err
	}

	return nil
}

func Call(address string, input []byte, function config.Function) string {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		if grpcerr := stop(function); grpcerr != nil {
			return grpcerr.Error()
		}

		if grpcerr := remove(function); grpcerr != nil {
			return grpcerr.Error()
		}

		log.Fatal("%v\n", err)
	}
	defer conn.Close()

	client := watcher.NewFxWatcherClient(conn)
	ctx := context.Background()

	r, err := client.Call(ctx, &watcher.Request{Input: input}, grpc.WaitForReady(true))
	if err != nil {

		if grpcerr := stop(function); grpcerr != nil {
			return grpcerr.Error()
		}

		if grpcerr := remove(function); grpcerr != nil {
			return grpcerr.Error()
		}

		log.Fatal("%v\n", err)
	}

	return r.Output
}

func runRun() error {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	
	if len(fxServices.Functions) <= 0 {
		return errors.New("")
	}

	for name, function := range fxServices.Functions {
		go func() {
			<-c
			os.Exit(1)
		}()
		function.Name = name

		if len(functionName) < 1 || function.Name != functionName {
			log.Fatal("Invalid function name. please describe name of function correctly\n")
		}

		runningContainer := builder.CheckImgRunning()

		if strings.Contains(runningContainer, function.Name) {
			if err := stop(function); err != nil {
				return err
			}

			if err := remove(function); err != nil {
				return err
			}
		}

		log.Info("Running image (%s) in local\n", function.Image)
		log.Info("Starting FxWatcher Server ...\n")
		log.Info("Call %s in user's local\n", function.Name)
		if err := run(function); err != nil {
			return err
		}

		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) != 0 {
			fmt.Fprintf(os.Stderr, "Reading from STDIN - hit (Control + D) to stop.\n")
		}

		functionInput, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("unable to read standard input: %s\n", err.Error())
		}
		fmt.Printf("Handler request: %s\n", string(functionInput))
		res := Call("localhost:50051", functionInput, function)

		if res != "" {
			fmt.Printf("Handler reply: ")
			os.Stdout.WriteString(res)
		}

		if err := stop(function); err != nil {
			return err
		}

		if err := remove(function); err != nil {
			return err
		}

	}

	return nil
}
