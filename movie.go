package alertmyshow

import(
	"time"
	"strings"
	"fmt"
)

type Movie struct{
	ID MovieID
	Title string
	Language string
	City	string
	Screen string
	Date	string
}



func NewMovie(title, language, city, screen, date string)(Movie, error){
	var movie Movie
	availableScreens := []string{"2d", "3d", "imax 2d", "imax 3d"}
	movie.Title = strings.ToLower(title)
	movie.Language = strings.ToLower(language)
	movie.City = strings.ToLower(city)
	movie.Screen = strings.ToLower(screen)
	screenCheck := false

	for _, s:= range availableScreens{
		if s == movie.Screen{
			screenCheck = true
			break
		}

	}


	if !screenCheck{
		return movie, fmt.Errorf("error unknown screen : %s", screen)
	}

	_, err := time.Parse("2006-01-02", date)
	if err != nil{
		return movie, err
	}
	movie.Date = date

	movie.ID, err = GetMovieID(title, language, screen, city)
	if err != nil{
		return movie, fmt.Errorf("error getting movieID: %w", err)
	}

	return movie, nil
}




	


