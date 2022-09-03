package main

import (
	"context"
	"flag"
	"fmt"
	"github/wongshennan/boardgame/api"
	"log"
	"strings"
	"time"

	"github.com/fatih/color"
)

// bga --query "ticket to ride" --clientId abc123 --skip 10 --limit 5
func main(){

	// define the cli commands
	query := flag.String("query", "", "Boardgame name to search")
	clientId := flag.String("clientId", "", "The id of the client")
	skip := flag.Uint("skip", 0, "Skips the number of results provided. It's generally used for paging results.")
	limit := flag.Uint("limit", 30, "Limits the number of results returned. The max limit is 100. The default is 30.")
	timeout := flag.Uint("timeout", 10, "Timeout")

	// ensure the inputs are parsed
	flag.Parse()

	
	if isNull(*query) {
		log.Fatalln("Please use --query to set the boardgame name to search")
	}
	
	if isNull(*clientId) {
		log.Fatalln("Please use --clientId to set the boardgame name to search")
	}

	fmt.Printf("query=%s, clientId=%s, limit=%d, skip=%d\n", *query, *clientId, *limit, *skip)

	// instantiate new api object
	bga := api.New(*clientId)

	// only possible because we attach a receiver
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(*timeout * uint(time.Second)))
	defer cancel()

	result, err := bga.Search(ctx, *query, *limit, *skip)

	if err != nil {
		log.Fatalf("cannot search for boardgame: %v", err)
	}

	// construct colros
	boldGreen := color.New(color.Bold).Add(color.FgHiGreen).SprintFunc()

	boldBlue := color.New(color.FgBlue).SprintFunc()

	for _,g := range result.Games {
		fmt.Printf("%s: %s\n", boldGreen("Name"),g.Name)
		fmt.Printf("%s: %s\n", boldGreen("Description"), g.Description)
		fmt.Printf("%s: %s\n", boldGreen("Price"), g.Price)
		fmt.Printf("%s: %d\n", boldGreen("Year Published"), g.YearPublished)
		fmt.Printf("URL: %s\n\n", boldBlue(g.Url))
	}
}

func isNull(s string) bool {
	return len(strings.TrimSpace(s)) <=0
}