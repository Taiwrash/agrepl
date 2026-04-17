package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var pricingCmd = &cobra.Command{
	Use:   "pricing",
	Short: "Show agrepl pricing tiers and features",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("\033[1mAgrepl Pricing Tiers\033[0m")
		fmt.Println("--------------------")
		
		fmt.Println("\n\033[36m1. Local (Free)\033[0m")
		fmt.Println("   - Unlimited Local Record/Replay")
		fmt.Println("   - Basic Visual Diff")
		fmt.Println("   - Local JSON Storage")
		
		fmt.Println("\n\033[32m2. Pro ($19/seat/mo)\033[0m")
		fmt.Println("   - Everything in Local, plus:")
		fmt.Println("   - Team Sync: Share runs via links (push/pull)")
		fmt.Println("   - agrepl test: CI/CD Integration")
		fmt.Println("   - Semantic Diff: LLM-powered drift analysis")
		
		fmt.Println("\n\033[35m3. Enterprise (Custom)\033[0m")
		fmt.Println("   - Everything in Pro, plus:")
		fmt.Println("   - Guardrails: Real-time cost & safety policies")
		fmt.Println("   - Audit Logs: Production trace history")
		fmt.Println("   - SSO & RBAC")
		
		fmt.Println("\nFor more details, visit: \033[4mhttps://agrepl.dev/pricing\033[0m")
	},
}

func init() {
	rootCmd.AddCommand(pricingCmd)
}
