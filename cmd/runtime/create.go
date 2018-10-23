package runtime

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	from string
)

func init() {
	createCmd.Flags().StringVarP(&from, "from-runtime", "r", "dockerfile", "")
}

var createCmd = &cobra.Command{
	Use:   `create <RUNTIME_NAME>`,
	Short: "create a OpenFx custom runtime",
	Long:  ``,
	RunE:  runCreate,
}

func runCreate(cmd *cobra.Command, args []string) error {
	fmt.Println(from)
	return nil
}
