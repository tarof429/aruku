/*
Copyright © 2020 Taro Fukunaga <tarof429@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/manifoldco/promptui"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/tarof429/aruku/aruku"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: aruku-bin <cmdset file>")
			os.Exit(1)
		}

		cmdSetFile := strings.Trim(args[0], "")

		var a aruku.App

		err := a.Load(cmdSetFile)

		if err != nil {
			fmt.Println(err)
		}

		fmt.Printf("\n\n%v\n\n", a.Description)

		var promptItems []string

		for _, list := range a.CmdList {
			promptItems = append(promptItems, list.Description)
			//fmt.Printf("Adding: %v\n", list.Description)
		}

		promptItems = append(promptItems, "Exit")

		prompt := promptui.Select{
			Label: "Select an item from the menu",
			Items: promptItems,
		}

		_, result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		if result == "Exit" {
			os.Exit(0)
		}

		a.SetCmdList(result)

		total := a.TotalCmds()

		for i := 0; a.HasNextCmd(); {

			fmt.Printf("\nCommand: %v\n\n", a.GetCurrentCmd().Description)

			time.Sleep((time.Millisecond * 500))

			var prompt promptui.Select

			if a.HasNextCmd() && a.HasPreviousCmd() {
				prompt = promptui.Select{
					Label: "Select Command",
					Items: []string{"Run", "Back", "Next", "Exit"},
				}
			} else if a.HasNextCmd() {
				prompt = promptui.Select{
					Label: "Select Command",
					Items: []string{"Run", "Next", "Exit"},
				}
			} else {
				prompt = promptui.Select{
					Label: "Select Command",
					Items: []string{"Run", "Back", "Exit"},
				}
			}

			_, result, err := prompt.Run()

			switch result {
			case "Run":
				time.Sleep((time.Millisecond * 500))
				a.Run()
				a.ShowCurrentCommandOutput()
				a.PointToNextCmd()
				i = i + 1
			case "Back":
				time.Sleep((time.Millisecond * 500))
				i = i - 1
				a.PointToPreviousCmd()
				continue
			case "Next":
				i = i + 1
				a.PointToNextCmd()
				continue
			case "Exit":
				fmt.Println("Exit")
				os.Exit(0)
			}

			if err != nil {
				fmt.Printf("Prompt failed %v\n", err)
				return
			}

			footer := "\nShowing " + strconv.Itoa(i) + "/" + strconv.Itoa(total) + " commands\n"
			pterm.DefaultCenter.Println(footer)

			time.Sleep((time.Millisecond * 500))

			if a.HasNextCmd() == false {
				prompt = promptui.Select{
					Label: "Select Command",
					Items: []string{"Back", "Exit"},
				}
				_, result, err := prompt.Run()

				switch result {
				case "Back":
					time.Sleep((time.Millisecond * 500))
					i = i - 1
					a.PointToPreviousCmd()
					continue
				case "Exit":
					fmt.Println("Exit")
					os.Exit(0)
				}
				if err != nil {
					fmt.Printf("Prompt failed %v\n", err)
					return
				}
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(loadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// loadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// loadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
