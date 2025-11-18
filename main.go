/*
   Copyright [2025] [0xf55]

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

package main

import (
	"log"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// cobra variables
var OutputFile string
var BatFiles []string
var Version bool
var Quiet bool

const banner = `⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢀⣀⣀⣀⣀⠀
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣠⣴⣶⣿⣿⣿⣿⣿⠟⠉
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⢿⣿⣿⣿⣿⣿⣿⣇
⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⣿⣶⡀⠈⣿⣿⣿⣿⣿⠟⠋
⠀⠀⠀⠀⠀⣀⣄⠀⠀⢿⣶⣾⣿⣿⣇⣼⣿⣿⣿⣿⡃   
⠀⠀⢀⣰⣾⣿⣿⣷⣄⣈⣉⣿⣿⣿⣿⣿⣿⣿⣿⣿⠆  By 0xf55
⠀⢠⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡅
⢀⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⣿⡿⣦
⢸⡟⠋⠉⠉⠙⠋⠀⠀⠉⠛⠋⠉⠛⠿⣿⡟⠻⠅⠈⠷⠆
⠘⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠀⠈⣧⠄`

var rootCmd = &cobra.Command{
	Use:   "bat",
	Short: "Wordlist generator tool",
	Long:  "Powerfull tool for wordlist generating based on custom files",
	Run: func(cmd *cobra.Command, args []string) {

		Writer = NewWriter()

		// show version
		if Version {
			return
		}

		SearchBT()

		ncolor := color.New()
		ncolor.AddRGB(38, 139, 210)

		// Show Information
		if !Quiet {
			ncolor.Println(`> The art of generating "wordlists"`)
			ncolor.Printf("> Bat\t\t%d.%d.%d\n", Major, Minor, Patch)
			ncolor.Printf("> Wordlist:\t%s\n", OutputFile)
			ncolor.Printf("> BatFiles:\t%s\n", BatFiles[:])
		}

		RunAll()

		if !Quiet {
			ncolor.Printf("\r> Lines:  \t%d\n", Writer.lines)
		}

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {

	new := color.New()
	new.AddRGB(7, 54, 66)
	if !Quiet {
		new.Println(banner)
	}

	ReversedCache = make(map[string]string)
	BatFiles = make([]string, 0, 2)
	Charset = Lower + Upper // defualt values
	Special = "#!_*@="      // ....
	Numbers = "0123456789"

	rootCmd.Flags().StringSliceVarP(&BatFiles, "input", "i", nil, "Input Bat Scirpt name")
	rootCmd.Flags().StringVarP(&OutputFile, "out", "o", "wordlist.txt", "Wordlist output name")
	rootCmd.Flags().BoolVarP(&Quiet, "quiet", "q", false, "Quiet Mode")
	rootCmd.Flags().BoolVarP(&Version, "version", "v", false, "Show bat version")

}

func main() {

	log.SetFlags(0)

	Execute()

}
