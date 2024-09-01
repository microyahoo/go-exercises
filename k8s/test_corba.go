package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	var cmd = &cobra.Command{
		Use:   "viper-test",
		Short: "testing viper",
		Run: func(command *cobra.Command, args []string) {
			fmt.Printf("thing1: %q\n", viper.GetString("thing1"))
			fmt.Printf("thing2: %q\n", viper.GetString("thing2"))
		},
	}

	viper.AutomaticEnv()
	flags := cmd.Flags()
	flags.String("thing1", "", "The first thing")
	viper.BindPFlag("thing1", flags.Lookup("thing1"))
	flags.String("thing2", "", "The second thing")
	viper.BindPFlag("thing2", flags.Lookup("thing2"))

	cmd.Execute()
}
