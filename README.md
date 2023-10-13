# subparser
Internal tool for the EffnerApp that can load substitution plans from different sources in HTML format, parse them and export them into JSON format.

## Usage
```
subparser.exe -s/--source <source> [-d/--destination <destination>] [-P/--parser <parser>] [-u/--user <dsb username>] [-p/--password <dsb/effner.de password] [-i/--input <input file path>] [-o/--o <output file path>]
```

## Sources
- local .html file
- DSB from heineking media
- effner.de/service/vertretungsplan

## Parsers
- Effner DSB
- Effner.DE

## Destinations
- file
- sysout

## Work efford
- 13.10. about 5h
