package alertmyshow

import (
	"fmt"
	"net/http"
)

type searchErrorResponse struct {
	Msg string `json:"msg"`
}

type Venue struct {
	ID      int    `json:"id"`
	Pid     int    `json:"pid"`
	Name    string `json:"name"`
	Address string `json:"address"`
}

type MovieSearchResult struct {
	Meta struct {
		Movies []struct {
			ID     string `json:"id"`
			FrmtID string `json:"frmtId"`
			Name   string `json:"name"`
			Lang   string `json:"lang"`
		} `json:"movies"`
		Venues []Venue `json:"cinemas"`
	} `json:"meta"`
}

// Get all available venue for the movie in city
func (m Movie) GetVenues() ([]Venue, error) {
	var searchErrorResponse searchErrorResponse
	var movieSearchResult MovieSearchResult

	u := fmt.Sprintf("https://apiproxy.paytm.com/v3/movies/search/movie?meta=1&reqData=1&city=%s&movieCode=%s&date=%s&version=3&site_id=1&channel=web&child_site_id=1", m.City, m.ID.GroupID, m.Date)
	client := &http.Client{}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("error: unable to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil, fmt.Errorf("error: empty response")
	}

	if resp.StatusCode == http.StatusBadRequest {
		if err := readJSON(resp.Body, &searchErrorResponse); err != nil {
			return nil, fmt.Errorf("error decoding search: %v", err)
		}
		return nil, fmt.Errorf("search error: %s", searchErrorResponse.Msg)
	}

	if resp.StatusCode == http.StatusOK {
		if err := readJSON(resp.Body, &movieSearchResult); err != nil {
			return nil, fmt.Errorf("error parsing search result: %v", err)
		}

		for _, movie := range movieSearchResult.Meta.Movies {
			if movie.ID == m.ID.Code {
				return movieSearchResult.Meta.Venues, nil
			}

		}
		return nil, fmt.Errorf("error: no venues found!")
	}

	return nil, fmt.Errorf("unknown search error: status Code: %d", resp.StatusCode)

}
