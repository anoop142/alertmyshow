package alertmyshow


import(
	"net/http"
	"fmt"
	"strings"
)

type searchErrorResponse struct {
	Msg    string `json:"msg"`
}

type MovieSearchResult struct {
	Meta       struct {
		Movies []struct {
			ID         string   `json:"id"`
			FrmtID     string   `json:"frmtId"`
			Name       string   `json:"name"`
			Lang      string `json:"lang"`
		} `json:"movies"`
	Cinemas []struct {
			ID                    int     `json:"id"`
			Pid                   int     `json:"pid"`
			Name                  string  `json:"name"`
			Address               string  `json:"address"`
	} `json:"cinemas"`
	} `json:"meta"`
}


func (m Movie)IsBookingStarted(venues []string)(bool, error){
	var searchErrorResponse searchErrorResponse
	var movieSearchResult MovieSearchResult

	u := fmt.Sprintf("https://apiproxy.paytm.com/v3/movies/search/movie?meta=1&reqData=1&city=%s&movieCode=%s&date=%s&version=3&site_id=1&channel=web&child_site_id=1", m.City, m.ID.GroupID, m.Date)
	client := &http.Client{}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil{
		return false, fmt.Errorf("error: unable to create request: %v", err)
	}

	resp, err := client.Do(req)
	if err != nil{
		return false, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent{
		return false, nil
	}

	if resp.StatusCode == http.StatusBadRequest{
		if err:= readJSON(resp.Body, &searchErrorResponse); err != nil{
			return false, fmt.Errorf("error decoding search error: %v", err)
		}
		return false, fmt.Errorf("search error: %s", searchErrorResponse.Msg)
	}

	if resp.StatusCode == http.StatusOK{
		if err:= readJSON(resp.Body, &movieSearchResult); err!= nil{
			return false, fmt.Errorf("error parsing search result: %v", err)
		}

		for _,movie:= range movieSearchResult.Meta.Movies{
			if movie.ID == m.ID.Code{
				// check for venue, string contains match
				for _, cinema := range movieSearchResult.Meta.Cinemas{
					for _,venue := range venues{
						if strings.Contains(strings.ToLower(cinema.Name), venue){
							return true, nil
						}
					}
				}
				break
			}

		}
		return false, nil
	}

	return false, fmt.Errorf("unknown search error: status Code: %d", resp.StatusCode)

}







