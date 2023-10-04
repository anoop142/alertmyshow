package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
	"strings"

	"github.com/anoop142/alertmyshow"
)

func printUsage() {
	fmt.Println("Usage:")
	fmt.Printf("%s -t title -l language -screen screen_type -d yyyy-mm-dd -v theatres(comma separated) -c city [-poll poll_in_minutes]\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	title := flag.String("t", "", "movie title")
	language := flag.String("l", "", "movie language")
	screen := flag.String("s", "", "screen type")
	city := flag.String("c", "", "city")
	date := flag.String("d", "", "date")
	var venues []string
	flag.Func("v", "venues", func(val string)error{
		for _, s := range strings.Split(val, ","){
			venues = append(venues, strings.TrimSpace(s))
		}

		return nil
	})

	poll := flag.Int("poll", 0, "poll time in minutes")

	flag.Usage = printUsage
	flag.Parse()

	required := []string{"t", "l", "s", "c", "d", "v"}
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
			bookingStarted, err := movie.IsBookingStarted(venues)
			if err != nil {
				log.Println(err)
			}
			if bookingStarted {
				fmt.Println("\rTickets Available!         ")
				break
			}
			time.Sleep(time.Duration(*poll) * time.Minute)

		}

	} else {
		bookingStarted, err := movie.IsBookingStarted(venues)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		if bookingStarted {
			fmt.Println("Tickets Available!")
		} else {
			fmt.Fprintln(os.Stderr, "Tickets Not Available.")
			os.Exit(2)
		}
	}

}
