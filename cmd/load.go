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
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"github.com/tarof429/aruku/aruku"
)

// loadCmd represents the load command
var loadCmd = &cobra.Command{
	Use:   "load",
	Short: "Load directory",
	Long:  `Load a directory that contains commands and any scripts`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Println("Usage: aruku-bin load <cmdset directory>")
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

		for a.HasCmd() {

			if a.GetCurrentCmd().CommandType == aruku.ExecuteCommandType {
				fmt.Printf("\nTask: %v\n\n", a.GetCurrentCmd().Description)

				time.Sleep((time.Millisecond * 100))

				time.Sleep((time.Millisecond * 100))
				s := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
				s.Start()

				// Run command
				a.Run()

				s.Stop()
				a.ShowCurrentCommandOutput()
				//a.PointToNextCmd()

				var prompt promptui.Select

				if a.HasNextCmd() {
					prompt = promptui.Select{
						Label: "Select Command",
						Items: []string{"Next", "Exit"},
					}
				} else {
					prompt = promptui.Select{
						Label: "Select Command",
						Items: []string{"Exit"},
					}
				}

				_, result, err := prompt.Run()

				switch result {
				case "Next":
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
			} else {
				done := false

				for done == false {
					fmt.Printf("%v: ", a.GetCurrentCmd().Description)
					done = a.Run()
				}
				a.PointToNextCmd()
			}

			time.Sleep((time.Millisecond * 100))

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
