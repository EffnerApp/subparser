package main

import (
	"encoding/json"
	"fmt"
	"github.com/mkideal/cli"
	"os"
	"subparser/dsb"
	"subparser/parsers"
)

const (
	ExitInvalidArgs       = 1
	ExitParserNotFound    = 2
	ExitFileReadFailed    = 3
	ExitDSBLoginFailed    = 4
	ExitDSBLoadingFailed  = 5
	ExitParsingFailed     = 5
	ExitFileWritingFailed = 6
)

type arguments struct {
	cli.Helper
	Parser  string `cli:"parser,P" usage:"name of the parser to use (default: effner)"`
	Input   string `cli:"input,i" usage:"input file. If unset, DSB user and pass are required"`
	DSBUser string `cli:"user,u" usage:"dsb username"`
	DSBPass string `cli:"pass,p" usage:"dsb password"`
	Output  string `cli:"output,o" usage:"file to output parsed data to (if unset SYSOUT is used)"`
}

func getParser(parser string) parsers.Parser {
	switch parser {
	case "effner":
		return &parsers.EffnerParser{}
	default:
		return nil
	}
}

func handle(error error, exit int) {
	fmt.Println("Error: ", error.Error())
	os.Exit(exit)
}

func main() {
	os.Exit(cli.Run(new(arguments), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*arguments)

		parserName := argv.Parser

		inputFile := argv.Input
		dsbUser := argv.DSBUser
		dsbPass := argv.DSBPass

		// validate input file and dsb credentials
		if inputFile == "" && dsbUser == "" && dsbPass == "" {
			fmt.Println("Error: Either input file or dsb credentials must be set!")
			os.Exit(ExitInvalidArgs)
		}

		if inputFile != "" && (dsbUser != "" || dsbPass != "") {
			fmt.Println("Error: Input file and DSB credentials cannot be defined at the same time!")
			os.Exit(ExitInvalidArgs)
		}

		// validate and load parser
		if parserName == "" {
			parserName = "effner"
		}
		parser := getParser(parserName)
		if parser == nil {
			fmt.Println("Error: Parser not found! Allowed: effner")
			os.Exit(ExitParserNotFound)
		}

		// we can only log to SYSOUT if we don't use it to transport the result
		canLog := argv.Output != ""

		// prepare the data to parse
		var data string

		if argv.Input != "" {
			// load data from file
			content, err := os.ReadFile(argv.Input)

			if err != nil {
				handle(err, ExitFileReadFailed)
			}

			data = string(content)
		} else {
			// load the data from DSB
			dsbInstance := dsb.NewDSB(dsbUser, dsbPass)

			err := dsbInstance.Login()

			if err != nil {
				handle(err, ExitDSBLoginFailed)
			}

			err = dsbInstance.LoadTimetables()

			if err != nil {
				handle(err, ExitDSBLoadingFailed)
			}

			// get the document from the dsb instance
			document := dsbInstance.Documents[0].Children[0]
			content, err := document.Download()

			if err != nil {
				handle(err, ExitDSBLoadingFailed)
			}

			data = string(content)
		}

		plans, err := parser.Parse(data)

		if err != nil {
			handle(err, ExitParsingFailed)
		}

		if canLog {
			fmt.Println("Parsing completed, Yay!")
		}

		plansJson, err := json.Marshal(plans)

		if err != nil {
			handle(err, ExitParsingFailed)
		}

		if argv.Output == "" {
			fmt.Println(string(plansJson))
			os.Exit(0)
		}

		// write output to file
		file, err := os.Create(argv.Output)

		if err != nil {
			handle(err, ExitFileWritingFailed)
		}
		defer file.Close()

		_, err = file.Write(plansJson)
		if err != nil {
			handle(err, ExitFileWritingFailed)
		}

		return nil
	}))
}
