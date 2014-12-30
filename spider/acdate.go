package main

import (
	"time"
)

type Season struct {
	Year    int
	Quarter int
}

func getRecentSeason() Season {
	season := Season{}
	season.Year = time.Now().Year()
	month := (int)(time.Now().Month())
	if month <= 3 {
		season.Quarter = 1
	} else if month <= 6 {
		season.Quarter = 2
	} else if month <= 9 {
		season.Quarter = 3
	} else {
		season.Quarter = 1
	}
	return season
}

func getAllSeason() []Season {
	seasons := []Season{}

	var begin int

	begin = Config.DownloadFlag.YearFrom
	nowSeason := getRecentSeason()

	for i := begin; i < nowSeason.Year; i++ {
		seasons = append(seasons, Season{Year: i, Quarter: 1})
		seasons = append(seasons, Season{Year: i, Quarter: 2})
		seasons = append(seasons, Season{Year: i, Quarter: 3})
		seasons = append(seasons, Season{Year: i, Quarter: 4})
	}

	for i := 1; i <= nowSeason.Quarter; i++ {
		seasons = append(seasons, Season{Year: nowSeason.Year, Quarter: i})
	}

	return seasons
}
