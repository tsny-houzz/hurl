package main

import (
	"bytes"
	"encoding/json"
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

var verbose bool

func main() {
	app := &cli.App{
		Name:        "hurl",
		Description: "Basic auth is handled by env vars STG_HOUZZ_USER and STG_HOUZZ_PASS",
		Usage:       "Curl substitute for stghouzz routing and testing",
		UsageText:   "EXAMPLE: hurl -b -c codespace=tsny http://prismic-cms-main.stghouzz.stg-main-eks.stghouzz.com/prismic-cms",
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
			&cli.BoolFlag{
				Name:  "v",
				Usage: "Verbose",
			},
			&cli.BoolFlag{
				Name:  "no-auth",
				Usage: "Don't use basic auth with env vars: STG_HOUZZ_USER and STG_HOUZZ_PASS",
			},
			&cli.BoolFlag{
				Name:  "mc",
				Usage: "Mimic a browser user-agent",
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
				if !strings.Contains(cookie, "=") {
					cookie = fmt.Sprintf("jkdebug=%s", cookie)
					fmt.Printf("Using default cookie: %s\n", cookie)
				}
			}

			displayOnly := c.Bool("d")
			printBody := c.Bool("b")
			verbose = c.Bool("v")

			headers := strings.Split(getEnv("DISPLAY_HEADERS", defaultHeaders), " ")

			client := &http.Client{
				Timeout: 15 * time.Second,
				CheckRedirect: func(req *http.Request, via []*http.Request) error {
					if req.Response != nil {
						printResponse(req.Response, headers, displayOnly)
						fmt.Printf("\n <%v>\n", strings.Repeat("-", 122))
					} else {
						println("warn: redirect had no response")
					}
					return nil
				},
			}

			if !strings.HasPrefix(url, "http") {
				url = "http://" + url
			}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				return err
			}

			if !c.Bool("no-auth") {
				user, pass := os.Getenv("STG_HOUZZ_USER"), os.Getenv("STG_HOUZZ_PASS")
				if user != "" && pass != "" {
					req.SetBasicAuth(user, pass)
				}
			}

			if c.Bool("mc") {
				agent := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36"
				req.Header.Set("User-Agent", agent)
			}

			setCookie(cookie, req)

			resp, err := client.Do(req)
			if err != nil {
				return err
			}
			defer resp.Body.Close()

			printResponse(resp, headers, displayOnly)

			if printBody {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err.Error())
				}
				if strings.HasPrefix(resp.Header.Get("Content-Type"), "application/json") {
					var prettyJSON bytes.Buffer
					if err := json.Indent(&prettyJSON, body, "", "  "); err != nil {
						log.Fatal("Failed to parse JSON: ", err)
					}
					fmt.Println("\n", prettyJSON.String())
				} else {
					fmt.Println("\n", strings.TrimSpace(string(body)))
				}
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
		first, second, found := strings.Cut(c, "=")
		if !found {
			log.Fatalf("error: %v: invalid cookie format", c)
		}
		cookieName = first
		cookieValue = second
	} else {
		cookieName = "jkdebug"
		cookieValue = c
	}

	// Check if the cookie is already set
	for _, cookie := range req.Cookies() {
		if cookie.Name == cookieName {
			return
		}
	}

	// Add the encoded cookie
	req.AddCookie(&http.Cookie{Name: cookieName, Value: cookieValue})
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

// printResponse displays headers + more, highlighting specified ones if needed.
func printResponse(resp *http.Response, displayHeaders []string, displayOnlyDesiredHeaders bool) {
	var highlightedOutput bytes.Buffer
	var normalOutput bytes.Buffer

	fmt.Println()
	fmt.Println("+", resp.Request.URL)
	emoji := "✅"
	if resp.StatusCode == 404 {
		emoji = "⚠️"
	}
	if resp.StatusCode > 404 {
		emoji = "🟥"
	}
	fmt.Println("-", resp.Request.URL.Scheme, resp.Status, emoji)
	headers := sortHeaders(resp.Header)

	if verbose {
		for k, v := range resp.Request.Header {
			fmt.Println("> H:", k, v)
		}
	}

	for key, values := range headers {
		line := fmt.Sprintf("< H: %s: %s\n", key, strings.Join(values, ", "))
		if displayOnlyDesiredHeaders {
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

	out := strings.TrimSpace(highlightedOutput.String())
	if out != "" {
		fmt.Println(out)
	}
	out = strings.TrimSpace(normalOutput.String())
	if out != "" {
		fmt.Println(out)
	}
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
