package version

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/badie/bdev/pkg/ui"
)

// Version is set at build time via ldflags
var Version = "3.0.0"

// BuildTime is set at build time via ldflags
var BuildTime = "development"

// NewCommand returns the version cobra command
func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version number of B.DEV",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(ui.Primary("B.DEV CLI"))
			fmt.Println(ui.Muted(ui.SeparatorLight))
			fmt.Printf("  %s %s\n", ui.Bold("Version:"), Version)
			fmt.Printf("  %s %s\n", ui.Bold("Go:"), runtime.Version())
			fmt.Printf("  %s %s/%s\n", ui.Bold("Platform:"), runtime.GOOS, runtime.GOARCH)
			fmt.Printf("  %s %s\n", ui.Bold("Built:"), BuildTime)
			fmt.Println(ui.Muted(ui.SeparatorLight))
		},
	}
}
