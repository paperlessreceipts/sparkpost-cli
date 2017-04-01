package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	sp "github.com/SparkPost/gosparkpost"
	"github.com/codegangsta/cli"
)

// Column mapping for Mandrill Blacklist
const (
	MandrillEmailCol      = 0
	MandrillReasonCol     = 1
	MandrillDetailCol     = 2
	MandrillCreatedCol    = 3
	MandrillExpiresAtCol  = 4
	MandrillLastEventCol  = 5
	MandrillExpiresAt2Col = 6
	MandrillSubAccountCol = 7
)

// Column mapping for SendGrid Blacklist
const (
	SendgridEmailCol = 0
	SendgridCreated  = 1
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {

	ValidParameters := []string{
		"to", "from", "domain", "cursor", "limit", "per_page", "page", "sources", "types", "description",
	}

	app := cli.NewApp()

	app.Version = "0.0.2"
	app.Name = "suppression-sparkpost-cli"
	app.Usage = "SparkPost suppression list cli\n\n\tSee https://developers.sparkpost.com/api/suppression-list.html"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "baseurl, u",
			Value:  "https://api.sparkpost.com",
			Usage:  "Optional baseUrl for SparkPost.",
			EnvVar: "SPARKPOST_BASEURL",
		},
		cli.StringFlag{
			Name:   "apikey, k",
			Value:  "",
			Usage:  "Required SparkPost API key",
			EnvVar: "SPARKPOST_API_KEY",
		},
		cli.StringFlag{
			Name:  "verbose",
			Value: "false",
			Usage: "Dumps additional information to console",
		},
		cli.StringFlag{
			Name:  "file, f",
			Value: "",
			Usage: "Compatible blacklist CSV file. See README.md for more info.",
		},
		cli.StringFlag{
			Name:  "command",
			Value: "list",
			Usage: "Optional one of list, retrieve, search, delete, mandrill, sendgrid",
		},
		cli.StringFlag{
			Name:  "recipient",
			Value: "",
			Usage: "Recipient email address. Example rcpt_1@example.com",
		},

		// Search Parameters
		cli.StringFlag{
			Name:  "from",
			Value: "",
			Usage: "Optional datetime the entries were last updated, in the format of YYYY-MM-DDTHH:mm:ssZ (2015-04-10T00:00:00)",
		},
		cli.StringFlag{
			Name:  "to",
			Value: "",
			Usage: "Optional datetime the entries were last updated, in the format YYYY-MM-DDTHH:mm:ssZ (2015-04-10T00:00:00)",
		},
		cli.StringFlag{
			Name:  "types",
			Value: "",
			Usage: "Optional types of entries to include in the search, i.e. entries with \"transactional\" and/or \"non_transactional\" keys set to true",
		},
		cli.StringFlag{
			Name:  "limit",
			Value: "",
			Usage: "Optional maximum number of results to return. Must be between 1 and 100000. Default value is 100000",
		},
	}
	app.Action = func(c *cli.Context) {

		if c.String("apikey") == "" {
			log.Fatalf("Error: SparkPost API key must be set\n")
			return
		}

		isVerbose := false
		if c.String("verbose") == "true" {
			isVerbose = true
		}

		cfg := &sp.Config{
			BaseUrl:    c.String("baseurl"),
			ApiKey:     c.String("apikey"),
			ApiVersion: 1,
			Verbose:    isVerbose,
		}

		var client sp.Client
		err := client.Init(cfg)
		if err != nil {
			log.Fatalf("SparkPost client init failed: %s\n", err)
			return
		}

		parameters := make(map[string]string)

		for i, val := range ValidParameters {

			if c.String(ValidParameters[i]) != "" {
				parameters[val] = c.String(val)
			}
		}

		switch c.String("command") {
		case "list":
			listWrapper, _, err := client.SuppressionList()

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}
			csvEntryPrinter(listWrapper, true)

		case "retrieve":
			recpipient := c.String("recipient")
			if recpipient == "" {
				log.Fatalf("ERROR: The `retrieve` command requires a recipient.")
				return
			}

			listWrapper, _, err := client.SuppressionRetrieve(recpipient)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}
			csvEntryPrinter(listWrapper, false)

		case "search":
			parameters := make(map[string]string)

			for i, val := range ValidParameters {

				if c.String(ValidParameters[i]) != "" {
					parameters[val] = c.String(val)
				}
			}

			listWrapper, _, err := client.SuppressionSearch(parameters)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}
			csvEntryPrinter(listWrapper, true)

		case "delete":
			recpipient := c.String("recipient")
			if recpipient == "" {
				log.Fatalf("ERROR: The `delete` command requires a recipient.")
				return
			}

			_, err := client.SuppressionDelete(recpipient)

			if err != nil {
				log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
				return
			}
			fmt.Println("OK")

		case "mandrill":
			fmt.Printf("Processing: %s\n", c.String("file"))
			file := c.String("file")
			if file == "" {
				log.Fatalf("ERROR: The `mandrill` command requires a CSV file.")
				return
			}

			f, err := os.Open(file)
			check(err)

			var entries = []sp.SuppressionEntry{}

			batchCount := 1

			blackListRow := csv.NewReader(bufio.NewReader(f))
			blackListRow.FieldsPerRecord = 8

			for {
				record, err := blackListRow.Read()
				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatalf("ERROR: Failed to process '%s':\n\t%s", file, err)

					return
				}

				if record[MandrillEmailCol] == "email" {
					// Skip over header row
					continue
				}

				if record[MandrillReasonCol] != "hard-bounce" {
					// Ignore soft-bounce
					continue
				}

				if strings.Count(record[MandrillEmailCol], "@") != 1 {
					fmt.Printf("WARN: Ignoring '%s'. It is not a valid email address.\n", record[MandrillEmailCol])
					continue
				}

				entry := sp.SuppressionEntry{}

				if record[MandrillEmailCol] == "" {
					// Must have email as it is suppression list primary key
					continue
				}

				entry.Email = record[MandrillEmailCol]
				entry.Transactional = false
				entry.NonTransactional = true
				entry.Description = fmt.Sprintf("MBL: %s", record[MandrillDetailCol])

				entries = append(entries, entry)

				if len(entries) > (1024 * 100) {
					fmt.Printf("Uploading batch %d\n", batchCount)
					_, err := client.SuppressionUpsert(entries)

					if err != nil {
						log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
						return
					}
					entries = []sp.SuppressionEntry{}
					batchCount++
				}
			}

			if len(entries) > 0 {
				fmt.Printf("Uploading batch %d\n", batchCount)
				_, err := client.SuppressionUpsert(entries)

				if err != nil {
					log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
					return
				}
			}
			fmt.Println("DONE")

		case "sendgrid":
			file := c.String("file")
			if file == "" {
				log.Fatalf("ERROR: The `sendgrid` command requires a CSV file.")
				return
			}

			f, err := os.Open(file)
			check(err)

			var entries = []sp.SuppressionEntry{}

			batchCount := 1

			blackListRow := csv.NewReader(bufio.NewReader(f))
			blackListRow.FieldsPerRecord = 2

			for {
				record, err := blackListRow.Read()
				if err == io.EOF {
					break
				}

				if err != nil {
					log.Fatalf("ERROR: Failed to process '%s':\n\t%s", file, err)

					return
				}

				if record[SendgridEmailCol] == "email" {
					// Skip over header row
					continue
				}

				entry := sp.SuppressionEntry{}

				if record[SendgridEmailCol] == "" {
					// Must have email as it is suppression list primary key
					continue
				}

				// SendGrid suppression lists are very dirty and tend to have invalid data. Some examples of invalid addresses are:
				// 	#02232014, gmail.com, To, 8/27/2015, name@yahoo.comett@domain.com"
				if strings.Count(record[SendgridEmailCol], "@") != 1 {
					fmt.Printf("WARN: Ignoring '%s'. It is not a valid email address.\n", record[SendgridEmailCol])
					continue
				}

				entry.Email = record[SendgridEmailCol]
				entry.Transactional = false
				entry.NonTransactional = true
				entry.Description = fmt.Sprintf("SBL: imported from SendGrid")

				entries = append(entries, entry)

				if len(entries) > (1024 * 100) {
					fmt.Printf("Uploading batch %d\n", batchCount)
					_, err := client.SuppressionUpsert(entries)

					if err != nil {
						log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
						return
					}
					entries = []sp.SuppressionEntry{}
					batchCount++
				}

			}

			if len(entries) > 0 {
				fmt.Printf("Uploading batch %d\n", batchCount)
				_, err := client.SuppressionUpsert(entries)

				if err != nil {
					log.Fatalf("ERROR: %s\n\nFor additional information try using `--verbose true`\n\n\n", err)
					return
				}
			}
			fmt.Println("DONE")

		default:
			fmt.Printf("\n\nERROR: Unknown Commnad[%s]\n\n", c.String("command"))

			return
		}

	}
	app.Run(os.Args)

}

func csvEntryPrinter(suppressionList *sp.SuppressionListWrapper, summary bool) {
	entries := suppressionList.Results

	if summary {
		fmt.Printf("Recipient, Transactional, NonTransactional, Source, Updated, Created\n")
	} else {
		fmt.Printf("Recipient, Transactional, NonTransactional, Source, Updated, Created, Description\n")
	}

	for i := range entries {
		entry := entries[i]
		if summary {
			fmt.Printf("%s, %t, %t, %s, %s, %s\n", entry.Recipient, entry.Transactional, entry.NonTransactional, entry.Source, entry.Updated, entry.Created)
		} else {
			fmt.Printf("%s, %t, %t, %s,%s, %s, %s\n", entry.Recipient, entry.Transactional, entry.NonTransactional, entry.Source, entry.Updated, entry.Created, sanatize(entry.Description))
		}
	}
}

func sanatize(str string) string {

	return stripchars(str, ",\n\r")
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}
