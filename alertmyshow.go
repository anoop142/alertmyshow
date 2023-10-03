package alertmyshow

import(
	"fmt"
	"net/http"
	"strings"
)



type MovieID struct{
	GroupID string
	Code string
}

type groupedMovie struct{
			Label             string   `json:"label"`
			ScrnFmt           []string `json:"scrnFmt"`
			Rnr               struct {
				HasUReview int `json:"hasUReview"`
				HasContent int `json:"hasContent"`
			} `json:"rnr"`
			BookingStatus        string `json:"bookingStatus,omitempty"`
			PopularityRank       int    `json:"popularityRank,omitempty"`
			LanguageFormatGroups []struct {
				Lang              string `json:"lang"`
				FmtGrpID          string `json:"fmtGrpId"`
				TotalSessionCount int    `json:"totalSessionCount"`
				ScreenFormats     []struct {
					MovieCode         string `json:"movieCode"`
					ScrnFmt           string `json:"scrnFmt"`
					NextAvailableDate string `json:"nextAvailableDate"`
				} `json:"screenFormats"`
			} `json:"languageFormatGroups"`
			ReleaseDate string `json:"releaseDate"`
}

type MoviesNow struct {
	Data struct {
		GroupedMovies []groupedMovie `json:"groupedMovies"`
	} `json:"data"`
}




func  extractMovieID(groupedMovies []groupedMovie, title, language, screen string)(MovieID, error){
	var movieID MovieID
	for _, movie := range groupedMovies{
		if strings.ToLower(movie.Label) == title{
			for _, l := range movie.LanguageFormatGroups{
				if strings.ToLower(l.Lang) == language{
					movieID.GroupID = l.FmtGrpID
					for _, s := range l.ScreenFormats{
						if strings.ToLower(s.ScrnFmt) == screen{
							movieID.Code = s.MovieCode
							break
						}
					}
					if movieID.Code == ""{
						return movieID, fmt.Errorf("error extracting movieID: screen format %s not found", screen)
					}
				}
			}
		}
	}
	if movieID.GroupID == ""{
		return movieID, fmt.Errorf("error extracting movieID: title %s not found", title)
	}

	return movieID, nil

}
			


func GetMovieID(title, language, screen, city string) (MovieID, error){
	var moviesNow MoviesNow
	var movieID MovieID

	client := &http.Client{}
	u := fmt.Sprintf("https://apiproxy.paytm.com/v3/movies/search/movies?version=3&site_id=1&channel=web&child_site_id=1&city=%s&mdp=1", city)

	req, err := http.NewRequest("GET", u, nil)
	if err != nil{
		return movieID, fmt.Errorf("error creating request %v",err)
	}

	resp, err := client.Do(req)

	if err != nil{
		return movieID, err
	}

	if resp.StatusCode != http.StatusOK{
		return movieID, fmt.Errorf("status code : %v while fetching movie id", resp.StatusCode)
	}

	defer resp.Body.Close()


	err = readJSON(resp.Body, &moviesNow)

	if err != nil{
		return movieID, fmt.Errorf("error unmarshaling moviesNow: %v", err)

	}

	return extractMovieID(moviesNow.Data.GroupedMovies, title, language, screen)


}



