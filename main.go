package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/urfave/cli"
)

const (
	defaultHeaders = "x-istio-vs hz-serverid location"
	highlightColor = "\033[33m"
	resetColor     = "\033[0m"
)

func main() {
	app := &cli.App{
		Name:      "hurl",
		Usage:     "Curl substitute for stghouzz routing and testing",
		UsageText: "EXAMPLE: hurl -b -c codespace=tsny http://prismic-cms-main.stghouzz.stg-main-eks.stghouzz.com/prismic-cms",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "d",
				Usage: "Display only specific headers",
			},
			&cli.StringFlag{
				Name:  "c",
				Usage: "Set a cookie in the format 'name=value'. Defaults to 'jkdebug=value' if '=' is missing.",
			},
			&cli.BoolFlag{
				Name:  "b",
				Usage: "Whether to print the final response body to stdout",
			},
		},
		Action: func(c *cli.Context) error {

			url := c.Args().First()
			if url == "" {
				_ = cli.ShowAppHelp(c)
				fmt.Println()
				return fmt.Errorf("arg 1 required: no URL provided")
			}

			cookie := c.String("c")
			if cookie != "" {
				if strings.Contains(cookie, "=") {
					fmt.Printf("Using cookie: %s\n", cookie)
				} else {
					cookie = fmt.Sprintf("jkdebug=%s", cookie)
					fmt.Printf("Using default cookie: %s\n", cookie)
				}
			}

			displayOnly := c.Bool("d")
			printBody := c.Bool("b")

			displayHeaders := strings.Split(getEnv("DISPLAY_HEADERS", defaultHeaders), " ")

			client := &http.Client{
				Timeout: 15 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					// Add jkdebug cookie for each redirect if provided
					// if *cookie != "" {
					// 	req.AddCookie(&http.Cookie{Name: "jkdebug", Value: *cookie})
					// }
					if req.Response != nil {
						printHeaders(req.Response, displayHeaders, displayOnly)
					}
					return nil
				},
			}

			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}

			user, pass := os.Getenv("STG_HOUZZ_USER"), os.Getenv("STG_HOUZZ_PASS")
			if user != "" && pass != "" {
				req.SetBasicAuth(user, pass)
			}

			setCookie(cookie, req)

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			printHeaders(resp, displayHeaders, displayOnly)

			if printBody {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err.Error())
				}
				fmt.Println("\n", strings.TrimSpace(string(body)))
			}

			return nil
		},
	}

	// Run the application
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func setCookie(c string, req *http.Request) {
	if c == "" {
		return
	}
	var cookieName, cookieValue string
	if strings.Contains(c, "=") {
		parts := strings.SplitN(c, "=", 2)
		cookieName = parts[0]
		cookieValue = parts[1]
	} else {
		cookieName = "jkdebug"
		cookieValue = c
	}

	cookie := &http.Cookie{Name: cookieName, Value: cookieValue}
	fmt.Printf("> %v=%v\n", cookieName, cookieValue)
	req.AddCookie(cookie)
}

func sortHeaders(headers http.Header) http.Header {
	keys := make([]string, 0, len(headers))
	for key := range headers {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	sortedHeaders := http.Header{}

	for _, key := range keys {
		sortedHeaders[key] = headers[key]
	}
	return sortedHeaders
	// Optionally, clear the original headers and copy sorted headers
	// headers = sortedHeaders // Uncomment if you want to reassign
}

// printHeaders displays headers, highlighting specified ones if needed.
func printHeaders(resp *http.Response, displayHeaders []string, displayOnly bool) {
	var highlightedOutput bytes.Buffer
	var normalOutput bytes.Buffer

	fmt.Println()
	fmt.Println(">", resp.Request.URL)
	fmt.Println("-", resp.Request.URL.Scheme, resp.Status)
	headers := sortHeaders(resp.Header)

	for key, values := range headers {
		line := fmt.Sprintf("%s: %s\n", key, strings.Join(values, ", "))
		if displayOnly {
			// Display only specific headers if -d is set
			if containsIgnoreCase(displayHeaders, key) {
				highlightedOutput.WriteString(highlight(line))
			}
		} else {
			// Highlight specific headers, but display all
			if containsIgnoreCase(displayHeaders, key) {
				highlightedOutput.WriteString(highlight(line))
			} else {
				normalOutput.WriteString(line)
			}
		}
	}

	// Print highlighted headers first, then normal headers
	fmt.Print(highlightedOutput.String())
	fmt.Print(normalOutput.String())
}

// highlight adds color to a header line.
func highlight(line string) string {
	return highlightColor + line + resetColor
}

// containsIgnoreCase checks if a slice contains a string (case-insensitive).
func containsIgnoreCase(slice []string, item string) bool {
	itemLower := strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == itemLower {
			return true
		}
	}
	return false
}

// getEnv retrieves an environment variable or a default value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
