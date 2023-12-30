/*
Anoop S
2023

*/

package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/anoop142/alertmyshow"
)

var (
	version = "dev"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s -t title -l language -screen screen_type -d yyyy-mm-dd -list | -v theatres(comma separated) -c city [-poll poll_in_minutes]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	title := flag.String("t", "", "movie title")
	language := flag.String("l", "", "movie language")
	screen := flag.String("s", "", "screen type")
	city := flag.String("c", "", "city")
	date := flag.String("d", "", "date")
	showVersion := flag.Bool("version", false, "version")
	list := flag.Bool("list", false, "list venues")

	var venues []string
	flag.Func("v", "venues", func(val string) error {
		for _, s := range strings.Split(val, ",") {
			venues = append(venues, strings.TrimSpace(s))
		}

		return nil
	})
	poll := flag.Int("poll", 0, "poll time in minutes")

	flag.Usage = printUsage
	flag.Parse()

	if *showVersion {
		fmt.Println(version)
		os.Exit(0)
	}

	required := []string{"t", "l", "s", "c", "d"}
	seen := make(map[string]bool)
	flag.Visit(func(f *flag.Flag) { seen[f.Name] = true })
	for _, req := range required {
		if !seen[req] {
			fmt.Fprintf(os.Stderr, "missing argument -%s\n\n", req)
			flag.Usage()
			os.Exit(1)
		}
	}

	movie, err := alertmyshow.NewMovie(*title, *language, *city, *screen, *date)

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	venuesFound := false

	// poll mode
	if *poll != 0 {
		// spinner
		go func() {
			for {
				for _, r := range `-\|/` {
					fmt.Fprintf(os.Stdout, "\rLooking for tickets %c", r)
					time.Sleep(100 * time.Millisecond)
				}
			}
		}()

		for {

			availVenues, err := movie.GetVenues()
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				continue
			}
			for _, v := range availVenues {
				// go over the supplied venues
				for _, venue := range venues {
					if strings.Contains(strings.ToLower(v.Name), venue) {
						venuesFound = true
						fmt.Println(v.Name)
					}
				}
			}
			if venuesFound {
				os.Exit(0)
			}

			time.Sleep(time.Duration(*poll) * time.Minute)

		}

	} else if *list {
		availVenues, err := movie.GetVenues()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, v := range availVenues {
			fmt.Println(v.Name)
		}

	} else {
		availVenues, err := movie.GetVenues()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		for _, v := range availVenues {
			// go over the supplied venues
			for _, venue := range venues {
				if strings.Contains(strings.ToLower(v.Name), venue) {
					venuesFound = true
					fmt.Println(v.Name)
				}
			}
		}
		if !venuesFound {
			fmt.Println("No match found!")
			os.Exit(1)
		}

	}

}
