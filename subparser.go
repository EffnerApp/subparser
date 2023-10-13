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
	EXITDSBLoadingFailed  = 5
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
				handle(err, EXITDSBLoadingFailed)
			}

			// get the document from the dsb instance
			document := dsbInstance.Documents[0].Children[0]
			content, err := document.Download()

			if err != nil {
				handle(err, EXITDSBLoadingFailed)
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

	//parser := parsers.EffnerParser{}
	//
	//plans, err := parser.Parse("<!DOCTYPE html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\">\n<html>\n<head><link rel=\"stylesheet\" type=\"text/css\" href=\"willi.css\"></link><script src=\"willi.js\" type=\"text/javascript\"></script><title>WILLI</title></head>\n<body>\n<a name=\"oben\"/><h1>Vertretungspl&auml;ne f&uuml;r </h1><br />\n<a href=\"#13.10.2023\">13.10.2023</a><br />\n<a href=\"#16.10.2023\">16.10.2023</a><br />\n<a name=\"13.10.2023\"><hr /></a>\n<p class=\"seite\" style=\"text-align:left\">\n<h2>Vertretungsplan f&uuml;r Freitag, 13.10.2023 <h4>(erstellt: 13.10.2023 um 8:23 Uhr) <font size=\"+3\"> </p>\n<p class=\"seite\" style=\"text-align:left\">\n<h4>Abwesende Klassen:</h4> <table class=\"K\"><colgroup><col width=\"150\"/><col width=\"300\"/></colgroup> <tbody class=\"K\"><tr class=\"K\"><th rowspan=\"1\" class=\"K\">\n8E</th>\n<td>\n1.-11. </td>\n</tr></tbody>\n</table>\n</p>\n<p class=\"seite\" style=\"text-align:left\">\n<h4>Vertretungen:</h4> <table class=\"k\" border=\"3\"><tr><th width=\"100\">\nKlasse </th>\n<th width=\"150\">\nLehrkraft</th>\n<th width=\"100\">\nStd.</th>\n<th width=\"100\">\n&nbsp;vertreten durch</th>\n<th width=\"100\">\n&nbsp;Raum</th>\n<th width=\"400\">\n</th>\n</tr>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n5A</th>\n<td>\nmu</td>\n<td>\n3</td>\n<td>\n&nbsp;wi</td>\n<td>\n&nbsp;232</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n5B</th>\n<td>\nfp</td>\n<td>\n1</td>\n<td>\n&nbsp;gn</td>\n<td>\n&nbsp;213</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nfp</td>\n<td>\n2</td>\n<td>\n&nbsp;h&uuml;</td>\n<td>\n&nbsp;213</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n5C</th>\n<td>\nls</td>\n<td>\n1</td>\n<td>\n&nbsp;fk</td>\n<td>\n&nbsp;242</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nls</td>\n<td>\n2</td>\n<td>\n&nbsp;li</td>\n<td>\n&nbsp;242</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nkv</td>\n<td>\n6</td>\n<td>\n&nbsp;ln</td>\n<td>\n&nbsp;242</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n5E</th>\n<td>\nch</td>\n<td>\n4</td>\n<td>\n&nbsp;nd</td>\n<td>\n&nbsp;233</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nch</td>\n<td>\n5</td>\n<td>\n&nbsp;sp</td>\n<td>\n&nbsp;216</td>\n<td>\n&nbsp;zusammen mit 5E(DInt)</td>\n</tr><tr class=\"k\">\n<td>\nch</td>\n<td>\n6</td>\n<td>\n&nbsp;sp</td>\n<td>\n&nbsp;216</td>\n<td>\n&nbsp;zusammen mit 5E(DInt)</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n5F</th>\n<td>\nga</td>\n<td>\n3</td>\n<td>\n&nbsp;zs</td>\n<td>\n&nbsp;234</td>\n<td>\n&nbsp;AA </td>\n</tr><tr class=\"k\">\n<td>\nga</td>\n<td>\n4</td>\n<td>\n&nbsp;hg</td>\n<td>\n&nbsp;234</td>\n<td>\n&nbsp;AA</td>\n</tr><tr class=\"k\">\n<td>\nkv</td>\n<td>\n5</td>\n<td>\n&nbsp;dm</td>\n<td>\n&nbsp;234</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n5I</th>\n<td>\nhs</td>\n<td>\n4</td>\n<td>\n&nbsp;lt</td>\n<td>\n&nbsp;S105</td>\n<td>\n&nbsp;Pausenaufs. anderw. vergeb.</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n6A</th>\n<td>\nhs</td>\n<td>\n1</td>\n<td>\n&nbsp;pp</td>\n<td>\n&nbsp;S009</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nba</td>\n<td>\n3</td>\n<td>\n&nbsp;hd</td>\n<td>\n&nbsp;211</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nba</td>\n<td>\n4</td>\n<td>\n&nbsp;sp</td>\n<td>\n&nbsp;211</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n6D</th>\n<td>\nba</td>\n<td>\n3</td>\n<td>\n&nbsp;hd</td>\n<td>\n&nbsp;211</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nba</td>\n<td>\n4</td>\n<td>\n&nbsp;sp</td>\n<td>\n&nbsp;211</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n6E</th>\n<td>\nsh</td>\n<td>\n2</td>\n<td>\n&nbsp;nu</td>\n<td>\n&nbsp;S209</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nhs</td>\n<td>\n5</td>\n<td>\n&nbsp;an</td>\n<td>\n&nbsp;S213</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\nws</td>\n<td>\n6</td>\n<td>\n&nbsp;sg</td>\n<td>\n&nbsp;S213</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n6F</th>\n<td>\nsh</td>\n<td>\n2</td>\n<td>\n&nbsp;nu</td>\n<td>\n&nbsp;S209</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nlp</td>\n<td>\n3</td>\n<td>\n&nbsp;ms</td>\n<td>\n&nbsp;124</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nlp</td>\n<td>\n4</td>\n<td>\n&nbsp;pn</td>\n<td>\n&nbsp;124</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n6H</th>\n<td>\nsh</td>\n<td>\n2</td>\n<td>\n&nbsp;nu</td>\n<td>\n&nbsp;S209</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n6I</th>\n<td>\nsh</td>\n<td>\n2</td>\n<td>\n&nbsp;nu</td>\n<td>\n&nbsp;S209</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\ngb</td>\n<td>\n4</td>\n<td>\n&nbsp;bc</td>\n<td>\n&nbsp;S205</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\ngb</td>\n<td>\n5</td>\n<td>\n&nbsp;tb</td>\n<td>\n&nbsp;S205</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n6K</th>\n<td>\nsh</td>\n<td>\n2</td>\n<td>\n&nbsp;nu</td>\n<td>\n&nbsp;S209</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n7A</th>\n<td>\nws</td>\n<td>\n3</td>\n<td>\n&nbsp;bi</td>\n<td>\n&nbsp;212</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nws</td>\n<td>\n4</td>\n<td>\n&nbsp;bo</td>\n<td>\n&nbsp;224</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n7B</th>\n<td>\ndb</td>\n<td>\n5</td>\n<td>\n&nbsp;bl</td>\n<td>\n&nbsp;241</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\ndb</td>\n<td>\n6</td>\n<td>\n&nbsp;cc</td>\n<td>\n&nbsp;241</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n7F</th>\n<td>\ndb</td>\n<td>\n5</td>\n<td>\n&nbsp;bl</td>\n<td>\n&nbsp;241</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\ndb</td>\n<td>\n6</td>\n<td>\n&nbsp;cc</td>\n<td>\n&nbsp;241</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n7H</th>\n<td>\nkv</td>\n<td>\n1</td>\n<td>\n&nbsp;mv</td>\n<td>\n&nbsp;S211</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n8A</th>\n<td>\nba</td>\n<td>\n1</td>\n<td>\n&nbsp;nd</td>\n<td>\n&nbsp;225</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nba</td>\n<td>\n2</td>\n<td>\n&nbsp;tr</td>\n<td>\n&nbsp;225</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n8B</th>\n<td>\nhh</td>\n<td>\n2</td>\n<td>\n&nbsp;he</td>\n<td>\n&nbsp;137</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"5\" class=\"k\">\n8C</th>\n<td>\nga</td>\n<td>\n1</td>\n<td>\n&nbsp;rk</td>\n<td>\n&nbsp;226</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\nga</td>\n<td>\n2</td>\n<td>\n&nbsp;vg</td>\n<td>\n&nbsp;226</td>\n<td>\n&nbsp;AA</td>\n</tr><tr class=\"k\">\n<td>\nsh</td>\n<td>\n3</td>\n<td>\n&nbsp;so</td>\n<td>\n&nbsp;244</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nkc</td>\n<td>\n4</td>\n<td>\n&nbsp;hu</td>\n<td>\n&nbsp;226</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nsu</td>\n<td>\n5</td>\n<td>\n&nbsp;kd</td>\n<td>\n&nbsp;226</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n8D</th>\n<td>\nsh</td>\n<td>\n3</td>\n<td>\n&nbsp;so</td>\n<td>\n&nbsp;244</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nba</td>\n<td>\n5</td>\n<td>\n&nbsp;uh</td>\n<td>\n&nbsp;244</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n8F</th>\n<td>\nsh</td>\n<td>\n3</td>\n<td>\n&nbsp;so</td>\n<td>\n&nbsp;244</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nsg</td>\n<td>\n4</td>\n<td>\n&nbsp;bi</td>\n<td>\n&nbsp;223</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nsu</td>\n<td>\n6</td>\n<td>\n&nbsp;wd</td>\n<td>\n&nbsp;223</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n9A</th>\n<td>\nsg</td>\n<td>\n3</td>\n<td>\n&nbsp;sg</td>\n<td>\n&nbsp;217</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nbi</td>\n<td>\n4</td>\n<td>\n&nbsp;sg</td>\n<td>\n&nbsp;217</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n9B</th>\n<td>\nnp</td>\n<td>\n3</td>\n<td>\n&nbsp;np</td>\n<td>\n&nbsp;021</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nwi</td>\n<td>\n4</td>\n<td>\n&nbsp;np</td>\n<td>\n&nbsp;045</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n9C</th>\n<td>\nsl</td>\n<td>\n3</td>\n<td>\n&nbsp;sl</td>\n<td>\n&nbsp;215</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\npn</td>\n<td>\n4</td>\n<td>\n&nbsp;sl</td>\n<td>\n&nbsp;215</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nsh</td>\n<td>\n5</td>\n<td>\n&nbsp;hr</td>\n<td>\n&nbsp;135</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n9D</th>\n<td>\nsm</td>\n<td>\n3</td>\n<td>\n&nbsp;sm</td>\n<td>\n&nbsp;117</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nsm</td>\n<td>\n4</td>\n<td>\n&nbsp;sm</td>\n<td>\n&nbsp;117</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n9E</th>\n<td>\nkj</td>\n<td>\n3</td>\n<td>\n&nbsp;kj</td>\n<td>\n&nbsp;132</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nsl</td>\n<td>\n4</td>\n<td>\n&nbsp;kj</td>\n<td>\n&nbsp;132</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nsh</td>\n<td>\n5</td>\n<td>\n&nbsp;hr</td>\n<td>\n&nbsp;135</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n9F</th>\n<td>\nbs</td>\n<td>\n3</td>\n<td>\n&nbsp;bs</td>\n<td>\n&nbsp;214</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nhg</td>\n<td>\n4</td>\n<td>\n&nbsp;bs</td>\n<td>\n&nbsp;214</td>\n<td>\n&nbsp;Projekt Ice-Breaker</td>\n</tr><tr class=\"k\">\n<td>\nsh</td>\n<td>\n5</td>\n<td>\n&nbsp;hr</td>\n<td>\n&nbsp;135</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n10B</th>\n<td>\nws</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nws</td>\n<td>\n2</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10C</th>\n<td>\nls</td>\n<td>\n5</td>\n<td>\n&nbsp;sc</td>\n<td>\n&nbsp;146</td>\n<td>\n&nbsp;LdK </td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10D</th>\n<td>\ngb</td>\n<td>\n2</td>\n<td>\n&nbsp;ka</td>\n<td>\n&nbsp;231</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10E</th>\n<td>\neb</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;verlegt auf Mo(9.10.) 3.St. </td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10F</th>\n<td>\nkv</td>\n<td>\n4</td>\n<td>\n&nbsp;ri</td>\n<td>\n&nbsp;235</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n10GL</th>\n<td>\nws</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nws</td>\n<td>\n2</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nlp</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n10GS</th>\n<td>\nws</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nws</td>\n<td>\n2</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nlp</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n11B</th>\n<td>\nri</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;verlegt auf Di(10.10.) 3.St. </td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n11C</th>\n<td>\ngb</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n11FS</th>\n<td>\nba</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n12Q</th>\n<td>\nhh</td>\n<td>\n5</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nhh</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n12Q2</th>\n<td>\nnp</td>\n<td>\n4</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;Bibliotheksarbeit</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n12Q3</th>\n<td>\nnp</td>\n<td>\n4</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;Bibliotheksarbeit</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n12Q5</th>\n<td>\ndb</td>\n<td>\n2</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n12Q6</th>\n<td>\ndb</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\ndb</td>\n<td>\n2</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\nWku</th>\n<td>\nsu</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\nILV9</th>\n<td>\ndb</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\nILV10</th>\n<td>\nhh</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n</table>\n</p>\n<p class=\"seite\" style=\"text-align:left\">\n <table class=\"F\" border=\"3\"><colgroup><col width=\"899\"/></colgroup> </table>\n</marquee> </p>\n<hr />\n<a name=\"16.10.2023\"><hr /></a>\n<p class=\"seite\" style=\"text-align:left\">\n<h2>Vertretungsplan f&uuml;r Montag, 16.10.2023 <h4>(erstellt: 13.10.2023 um 8:23 Uhr) <font size=\"+3\"> </p>\n<p class=\"seite\" style=\"text-align:left\">\n<h4>Abwesende Klassen:</h4> <table class=\"K\"><colgroup><col width=\"150\"/><col width=\"300\"/></colgroup> <tbody class=\"K\"><tr class=\"K\"><th rowspan=\"1\" class=\"K\">\n8C</th>\n<td>\n1.-11. </td>\n</tr></tbody>\n</table>\n</p>\n<p class=\"seite\" style=\"text-align:left\">\n<h4>Vertretungen:</h4> <table class=\"k\" border=\"3\"><tr><th width=\"100\">\nKlasse </th>\n<th width=\"150\">\nLehrkraft</th>\n<th width=\"100\">\nStd.</th>\n<th width=\"100\">\n&nbsp;vertreten durch</th>\n<th width=\"100\">\n&nbsp;Raum</th>\n<th width=\"400\">\n</th>\n</tr>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n5A</th>\n<td>\nhi</td>\n<td>\n3</td>\n<td>\n&nbsp;h&uuml;</td>\n<td>\n&nbsp;232</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nhi</td>\n<td>\n4</td>\n<td>\n&nbsp;ws</td>\n<td>\n&nbsp;232</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n5C</th>\n<td>\nls</td>\n<td>\n4</td>\n<td>\n&nbsp;la</td>\n<td>\n&nbsp;242</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nls</td>\n<td>\n5</td>\n<td>\n&nbsp;ef</td>\n<td>\n&nbsp;242</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n5E</th>\n<td>\nzs</td>\n<td>\n6</td>\n<td>\n&nbsp;rk</td>\n<td>\n&nbsp;233</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"5\" class=\"k\">\n6H</th>\n<td>\ngu</td>\n<td>\n1</td>\n<td>\n&nbsp;bg</td>\n<td>\n&nbsp;S014</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\ngu</td>\n<td>\n2</td>\n<td>\n&nbsp;bg</td>\n<td>\n&nbsp;S014</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\nbh</td>\n<td>\n8</td>\n<td>\n&nbsp;pp</td>\n<td>\n&nbsp;S211</td>\n<td>\n&nbsp;zusammen mit 6H(Inf)</td>\n</tr><tr class=\"k\">\n<td>\nbh</td>\n<td>\n9</td>\n<td>\n&nbsp;mv</td>\n<td>\n&nbsp;S204</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nbh</td>\n<td>\n10</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n6I</th>\n<td>\nhi</td>\n<td>\n1</td>\n<td>\n&nbsp;hd</td>\n<td>\n&nbsp;S205</td>\n<td>\n&nbsp;AA</td>\n</tr><tr class=\"k\">\n<td>\nhi</td>\n<td>\n2</td>\n<td>\n&nbsp;fo</td>\n<td>\n&nbsp;S205</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n7E</th>\n<td>\ngu</td>\n<td>\n3</td>\n<td>\n&nbsp;cc</td>\n<td>\n&nbsp;121</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n8A</th>\n<td>\nhh</td>\n<td>\n1</td>\n<td>\n&nbsp;no</td>\n<td>\n&nbsp;225</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\ntr</td>\n<td>\n6</td>\n<td>\n&nbsp;ka</td>\n<td>\n&nbsp;131</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"4\" class=\"k\">\n8B</th>\n<td>\ntr</td>\n<td>\n6</td>\n<td>\n&nbsp;ka</td>\n<td>\n&nbsp;131</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nrf</td>\n<td>\n2</td>\n<td>\n&nbsp;ga</td>\n<td>\n&nbsp;137</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nrf</td>\n<td>\n3</td>\n<td>\n&nbsp;so</td>\n<td>\n&nbsp;137</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\ngu</td>\n<td>\n5</td>\n<td>\n&nbsp;hm</td>\n<td>\n&nbsp;022</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n8E</th>\n<td>\nsu</td>\n<td>\n1</td>\n<td>\n&nbsp;kr</td>\n<td>\n&nbsp;141</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n8F</th>\n<td>\nsu</td>\n<td>\n4</td>\n<td>\n&nbsp;fk</td>\n<td>\n&nbsp;223</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n9C</th>\n<td>\nsu</td>\n<td>\n2</td>\n<td>\n&nbsp;sl</td>\n<td>\n&nbsp;215</td>\n<td>\n&nbsp;statt 6.St. </td>\n</tr><tr class=\"k\">\n<td>\nsl</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;verlegt auf 2.St. </td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n9E</th>\n<td>\nzs</td>\n<td>\n2</td>\n<td>\n&nbsp;mc</td>\n<td>\n&nbsp;132</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n9F</th>\n<td>\nzs</td>\n<td>\n5</td>\n<td>\n&nbsp;bs</td>\n<td>\n&nbsp;214</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n10A</th>\n<td>\nzs</td>\n<td>\n4</td>\n<td>\n&nbsp;kv</td>\n<td>\n&nbsp;142</td>\n<td>\n&nbsp;</td>\n</tr><tr class=\"k\">\n<td>\nma</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10B</th>\n<td>\nnd</td>\n<td>\n2</td>\n<td>\n&nbsp;gr</td>\n<td>\n&nbsp;136</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10C</th>\n<td>\nls</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n10E</th>\n<td>\nzs</td>\n<td>\n8</td>\n<td>\n&nbsp;rw</td>\n<td>\n&nbsp;018</td>\n<td>\n&nbsp;zusammen mit 10E(Ph&Uuml;)</td>\n</tr><tr class=\"k\">\n<td>\nzs</td>\n<td>\n9</td>\n<td>\n&nbsp;rw</td>\n<td>\n&nbsp;018</td>\n<td>\n&nbsp;zusammen mit 10E(Ph&Uuml;)</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10F</th>\n<td>\ntr</td>\n<td>\n3</td>\n<td>\n&nbsp;wb</td>\n<td>\n&nbsp;111</td>\n<td>\n&nbsp;</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n10GL</th>\n<td>\ntr</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n10GS</th>\n<td>\ntr</td>\n<td>\n1</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nrf</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nrf</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"10\" class=\"k\">\n11P</th>\n<td>\noe</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nno</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nef</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nws</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\npn</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nno</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nhh</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\npn</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nef</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\noe</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n11D</th>\n<td>\ngu</td>\n<td>\n4</td>\n<td>\n&nbsp;he</td>\n<td>\n&nbsp;019</td>\n<td>\n&nbsp;LdK</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"6\" class=\"k\">\n12Q</th>\n<td>\ndb</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nuh</td>\n<td>\n8</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\ndb</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nuh</td>\n<td>\n9</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nef</td>\n<td>\n10</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\nef</td>\n<td>\n11</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"3\" class=\"k\">\n12Q2</th>\n<td>\nks</td>\n<td>\n1</td>\n<td>\n&nbsp;bc</td>\n<td>\n&nbsp;147</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\nks</td>\n<td>\n2</td>\n<td>\n&nbsp;bc</td>\n<td>\n&nbsp;147</td>\n<td>\n&nbsp;LdK</td>\n</tr><tr class=\"k\">\n<td>\nma</td>\n<td>\n5</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;147</td>\n<td>\n&nbsp;Bibliotheksarbeit</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n12Q3</th>\n<td>\nma</td>\n<td>\n5</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;147</td>\n<td>\n&nbsp;Bibliotheksarbeit</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n12Q4</th>\n<td>\nsu</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"2\" class=\"k\">\n12Q5</th>\n<td>\nsu</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr><tr class=\"k\">\n<td>\ntr</td>\n<td>\n5</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;113</td>\n<td>\n&nbsp;Bibliotheksarbeit</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\n12Q6</th>\n<td>\nsu</td>\n<td>\n6</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\nFoer</th>\n<td>\nws</td>\n<td>\n7</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n<tbody class=\"k\"><tr class=\"k\"><th rowspan=\"1\" class=\"k\">\nPlus</th>\n<td>\noe</td>\n<td>\n7</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;</td>\n<td>\n&nbsp;entf&auml;llt</td>\n</tr></tbody>\n</table>\n</p>\n<p class=\"seite\" style=\"text-align:left\">\n <table class=\"F\" border=\"3\"><colgroup><col width=\"899\"/></colgroup> </table>\n</marquee> </p>\n<hr />\n</body></html>")
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//jsonBytes, err := json.Marshal(plans)
	//_ = os.WriteFile("subs.json", jsonBytes, 0644)
}
