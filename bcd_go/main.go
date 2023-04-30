package main

import (
	"bcd_go/config"
	"bcd_go/user"
	"bcd_go/video"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{}

func main() {
	config.InitRedis()

	rootCmd.AddCommand(user.Cmd())
	rootCmd.AddCommand(video.Cmd())
	rootCmd.AddCommand(ServerCmd())

	rootCmd.Execute()
}
