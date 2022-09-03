package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

const SEARCH_URL = "https://api.boardgameatlas.com/api/search"

type BoardGameAtlas struct {
	// lowercase cause private
	clientId string
}

type Game struct {
	Id 						string 	`json:"id"`
	Name 					string 	`json:"name"`
	Price 				string 	`json:"price"`
	YearPublished uint 		`json:"year_published"`
	Description 	string 	`json:"description"`
	Url						string	`json:"official_url"`
	ImageUrl			string 	`json:"image_url"`
	RulesUrl			string	`json:"rules_url"`
}

type SearchResult struct {
	Games []Game 	`json:"games"`
	Count uint		`json:"count"`
}

func (b *BoardGameAtlas) Search(ctx context.Context, query string, limit uint, skip uint) (*SearchResult, error) {
	// need to return the value using *SearchResult, if not need to manually return SearchResult{} rather than nil

	// create a http client	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, SEARCH_URL, nil)
	if err!= nil {
		// return an error object
		return nil, fmt.Errorf("cannot create http client: %v", err)
	}

	qs := req.URL.Query()
	qs.Add("name", query)
	qs.Add("limit", fmt.Sprintf("%d",limit))
	qs.Add("skip", strconv.Itoa(int(skip)))
	qs.Add("client_id", b.clientId)

	//encode query params, add to back of request
	req.URL.RawQuery = qs.Encode()

	// fmt.Printf("URL = %s\n", req.URL.String())
	res, err := http.DefaultClient.Do(req)

	if err!=nil {
		return nil,fmt.Errorf("cannot create http client for invocation: %v", err)
	}
	
	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("error http status: %s", res.Status)
	}

	var result SearchResult

	//deserialize the JSON payload to struct
	// if statement introduces new scope
	if err:=json.NewDecoder(res.Body).Decode(&result); err!=nil {
		return nil, fmt.Errorf("cannot deserialize JSON payload: %v", err)
	}

	return &result, nil
}

func New(clientId string) BoardGameAtlas {
	// same as js, can structure k:v to k if same var name
	return BoardGameAtlas{clientId}
}
