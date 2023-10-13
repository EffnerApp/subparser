package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/mkideal/cli"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"subparser/dsb"
	"subparser/parsers"
)

const (
	ExitInvalidArgs       = 1
	ExitParserNotFound    = 2
	ExitFileReadFailed    = 3
	ExitDSBLoginFailed    = 4
	ExitLoadingFailed     = 5
	ExitParsingFailed     = 5
	ExitFileWritingFailed = 6
)

type arguments struct {
	cli.Helper
	Parser string `cli:"parser,P" usage:"name of the parser to use (default: effner) [effner, effner-de]"`
	Input  string `cli:"input,i" usage:"input file (required if source is file)"`
	Source string `cli:"source,s" usage:"source for the data (default: file) [file,dsb,effner]"`
	User   string `cli:"user,u" usage:"username"`
	Pass   string `cli:"pass,p" usage:"password"`
	Output string `cli:"output,o" usage:"file to output parsed data to (if unset SYSOUT is used)"`
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

		dsbUser := argv.User
		dsbPass := argv.Pass

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

		if argv.Source == "file" {
			// load data from file
			content, err := os.ReadFile(argv.Input)

			if err != nil {
				handle(err, ExitFileReadFailed)
			}

			data = string(content)
		} else if argv.Source == "dsb" {
			if argv.User == "" || argv.Pass == "" {
				fmt.Println("Error: DSB-Source requires credentials!")
				os.Exit(ExitInvalidArgs)
			}

			// load the data from DSB
			dsbInstance := dsb.NewDSB(dsbUser, dsbPass)

			err := dsbInstance.Login()

			if err != nil {
				handle(err, ExitDSBLoginFailed)
			}

			err = dsbInstance.LoadTimetables()

			if err != nil {
				handle(err, ExitLoadingFailed)
			}

			// get the document from the dsb instance
			document := dsbInstance.Documents[0].Children[0]
			content, err := document.Download()

			if err != nil {
				handle(err, ExitLoadingFailed)
			}

			data = string(content)
		} else if argv.Source == "effner" {
			if argv.Pass == "" {
				fmt.Println("Error: Effner-Source requires password!")
				os.Exit(ExitInvalidArgs)
			}

			if argv.Parser != "effner-de" {
				fmt.Println("Error: Effner-Source only works with effner-de parser!")
				os.Exit(ExitInvalidArgs)
			}

			// TODO Move this somewhere else.
			form := url.Values{}
			form.Add("post_password", argv.Pass)

			// load the html from effner.de
			req, err := http.NewRequest("POST", "https://effner.de/wp-login.php?action=postpass", strings.NewReader(form.Encode()))

			if err != nil {
				handle(err, ExitLoadingFailed)
			}

			req.Header.Set("Referer", "https://effner.de/service/vertretungsplan/")
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			jar, err := cookiejar.New(nil)

			httpClient := &http.Client{
				Jar: jar,
			}

			res, err := httpClient.Do(req)

			if err != nil {
				handle(err, ExitLoadingFailed)
			}

			link := res.Header.Get("Link")
			linkParts := strings.Split(link, ",")

			final := strings.Split(linkParts[1], ";")[0]
			final = final[2 : len(final)-1]

			req, err = http.NewRequest("GET", final, nil)
			res, err = httpClient.Do(req)

			body, err := io.ReadAll(res.Body)

			if err != nil {
				handle(err, ExitLoadingFailed)
			}

			data = string(body)
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
