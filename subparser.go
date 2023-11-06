package main

import (
	"crypto/tls"
	"fmt"
	"github.com/mkideal/cli"
	"net/http"
	"os"
	"subparser/destination"
	"subparser/parsers"
	"subparser/source"
)

const (
	ExitInvalidArgs    = 1
	ExitParserNotFound = 2
	ExitFileReadFailed = 3
	ExitDSBLoginFailed = 4
	ExitLoadingFailed  = 5
	ExitParsingFailed  = 5
	ExitWritingFailed  = 6
)

type arguments struct {
	cli.Helper
	Parser      string `cli:"parser,P" usage:"name of the parser to use (default: effner) [effner, effner-de]"`
	Input       string `cli:"input,i" usage:"input file (required if source is file)"`
	Source      string `cli:"source,s" usage:"source for the data (default: file) [file,dsb,effner]"`
	Destination string `cli:"destination,d" usage:"destination where to put the parsed data, (default: sysout) [file]"`
	User        string `cli:"user,u" usage:"username"`
	Pass        string `cli:"pass,p" usage:"password"`
	Output      string `cli:"output,o" usage:"file to output parsed data to (if unset SYSOUT is used)"`
}

func getParser(parser string) parsers.Parser {
	switch parser {
	case "effner":
		return &parsers.EffnerDSBParser{}
	case "effner-de":
		return &parsers.EffnerDEParser{}
	default:
		return nil
	}
}

func getSource(argv *arguments) source.Source {
	switch argv.Source {
	case "dsb":
		return &source.DSBSource{
			User: argv.User,
			Pass: argv.Pass,
		}
	case "file":
		return &source.FileSource{
			Path: argv.Input,
		}
	case "effner":
		return &source.EffnerDESource{
			Password: argv.Pass,
		}
	default:
		return nil
	}
}

func getDestination(argv *arguments) destination.Destination {
	switch argv.Destination {
	case "file":
		return &destination.FileDestination{
			Path: argv.Output,
		}
	default:
		return &destination.SysoutDestination{}
	}
}

func handle(error error, exit int) {
	fmt.Println("Error: ", error.Error())
	os.Exit(exit)
}

func main() {
	// TODO REMOVE THIS REMOVE THIS THIS IS NOT WHAT WE DO USUALLY I SWEAR
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	os.Exit(cli.Run(new(arguments), func(ctx *cli.Context) error {
		argv := ctx.Argv().(*arguments)

		parserName := argv.Parser

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
		canLog := argv.Destination != "sysout" && argv.Destination != ""

		src := getSource(argv)

		if src == nil {
			fmt.Println("Error: Source not found! Allowed: effner, dsb, file")
			os.Exit(ExitInvalidArgs)
		}

		dest := getDestination(argv)

		if dest == nil {
			fmt.Println("Error: Destination not found! Allowed: file")
			os.Exit(ExitInvalidArgs)
		}

		// prepare the data to parse
		data, err := src.Load()

		if err != nil {
			handle(err, ExitLoadingFailed)
		}

		plans, err := parser.Parse(data)

		if err != nil {
			handle(err, ExitParsingFailed)
		}

		if canLog {
			fmt.Println("Parsing completed, Yay!")
		}

		err = dest.Write(plans)

		if err != nil {
			handle(err, ExitWritingFailed)
		}

		return nil
	}))
}
