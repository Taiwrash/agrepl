package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:   "uninstall",
	Short: "Uninstall agrepl and optionally clear its data",
	Long:  `This command downloads and runs the uninstallation script to remove the agrepl binary and its configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\033[36mRunning uninstallation script...\033[0m")

		// We execute the uninstall script via curl | bash to ensure it's the latest version
		// similar to how install works, or we can look for it locally if cloned.
		// For consistency with install.sh, we'll suggest the curl command.
		
		scriptURL := "https://raw.githubusercontent.com/taiwrash/agrepl/main/scripts/uninstall.sh"
		
		fmt.Printf("To complete uninstallation, please run:\n\n")
		fmt.Printf("\033[1mcurl -sSL %s | bash\033[0m\n\n", scriptURL)
		
		// Optionally, if we are in the repo, we can try to run it directly
		if _, err := os.Stat("scripts/uninstall.sh"); err == nil {
			fmt.Printf("Or run the local script:\n")
			fmt.Printf("\033[1m./scripts/uninstall.sh\033[0m\n")
			return
		}

		// Try to run it automatically if user wants? 
		// Actually, standard CLI uninstall often just gives the command or handles it if it has permissions.
		// Since it needs sudo potentially, it's safer to let user run it.
	},
}

func init() {
	rootCmd.AddCommand(uninstallCmd)
}

func runCommand(name string, args ...string) error {
	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
