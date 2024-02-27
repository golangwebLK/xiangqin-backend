package main

import (
	"fmt"
	"log"
	"xiangqin-backend/cmd/migrate"
	"xiangqin-backend/cmd/server"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:           "xiangqin-backend",
	Short:         "xiangqin-backend",
	Long:          "xiangqin-backend",
	SilenceErrors: true,
	SilenceUsage:  true,
	Args: func(cmd *cobra.Command, args []string) error {
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		tip()
	},
}

func tip() {
	fmt.Println("欢迎使用条件匹配系统")
}

func init() {
	rootCmd.AddCommand(server.StartCmd)
	rootCmd.AddCommand(migrate.StartCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
