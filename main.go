package main

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

var OutputFile string
var ScriptFile string

var Quiet bool
var Version bool

var rootCmd = &cobra.Command{
	Use:   "bat",
	Short: "Wordlist generator tool",
	Long:  "Powerfull tool for wordlist generating based on custom scripts",
	Run: func(cmd *cobra.Command, args []string) {
		if Version {
			ShowVersion()
			return
		}

		// Show Information
		if !Quiet {
			fmt.Printf("Wordlist: %s\n", OutputFile)
			fmt.Printf("BatFile:  %s\n", ScriptFile)
		}

		lexer := NewLexer(ScriptFile)
		lexer.Lex()
		if len(lexer.tokens) <= 0 {
			return
		}
		parser := NewParser(lexer.tokens)
		parser.Parse(0, parser.TokensLength)

		if !Quiet {
			fmt.Printf("Generated Lines: %d", Writer.lines)
		}
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&OutputFile, "out", "o", "wordlist.txt", "Wordlist output name")
	rootCmd.Flags().StringVarP(&ScriptFile, "input", "i", "main.sbt", "Input Bat Scirpt name")
	rootCmd.Flags().BoolVarP(&Quiet, "quiet", "q", false, "Quiet Mode")
	rootCmd.Flags().BoolVarP(&Version, "version", "v", false, "Show bat version")
}

func main() {

	log.SetFlags(0)
	Execute()

}
