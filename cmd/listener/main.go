package main

import (
	"example_listener/internal/config"
	"example_listener/internal/listener"
	"example_listener/internal/server"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "listener",
		Short: "Simple HTTP listener example",
		Run: func(cmd *cobra.Command, args []string) {
			// Connect to Pico server
			srv, err := server.New()
			if err != nil {
				log.Panic(err)
			}
			log.Println("Successfully connected to pico server")

			// Start listener
			err = listener.Start(srv)
			if err != nil {
				log.Panic(err)
			}
		},
	}

	// CLI flags
	rootCmd.Flags().StringVarP(&config.ServerAddr, "server", "s", "", "Server address")
	rootCmd.Flags().StringVarP(&config.ListenerAddr, "listener", "l", "", "Listener address")
	rootCmd.Flags().StringVarP(&config.Token, "token", "t", "", "Authorization token")
	rootCmd.MarkFlagRequired("server")
	rootCmd.MarkFlagRequired("listener")
	rootCmd.MarkFlagRequired("token")

	// Run
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
