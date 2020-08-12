package version

import (
	"github.com/spf13/cobra"
	"fmt"
)

const (
	flagLong = "long"
)

var (

	// VersionCmd prints out the current sdk version
	VersionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the app version",
		RunE: func(_ *cobra.Command, _ []string) error {
			verInfo := newVersionInfo()
			fmt.Println(verInfo)
			return nil
		},
	}
)

func init() {
	VersionCmd.Flags().Bool(flagLong, false, "Print long version information")
}
