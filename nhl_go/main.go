package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"nhl_go/models"
	"strconv"
)

func main() {
	nhlurl := "https://api-web.nhle.com/v1/schedule/2026-05-13"

	res, err := http.Get(nhlurl)

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var nhl_data models.NHLSchedule
	json_err := json.NewDecoder(res.Body).Decode(&nhl_data)

	if json_err != nil {
		panic(json_err)
	}

	games := process_teams(nhl_data)

	send_to_calendar(games)

}

func send_to_calendar(games map[string]map[string]any) {

	cal_data := make(map[string]models.CalEvent)

	for id, game := range games {
		cal_data[id] = models.CalEvent{
			ID:    id,
			Title: game["home_team"].(string) + " vs " + game["away_team"].(string),
		}
	}

	println(fmt.Sprintf("Mock send Calendar Events: %+v", cal_data))

}

func process_teams(json_data models.NHLSchedule) map[string]map[string]any {

	game_data := make(map[string]any)
	games := make(map[string]map[string]any)

	for _, gameweek := range json_data.GameWeek {
		game_date := gameweek.Date
		for _, game := range gameweek.Games {

			game_data["home_team"] = game.HomeTeam.HomeTeamPlaceName.Default
			game_data["away_team"] = game.AwayTeam.PlaceName.Default
			game_data["gameDate"] = game_date
			game_data["homeScore"] = game.HomeTeam.Score
			game_data["awayScore"] = game.AwayTeam.Score
			games[strconv.Itoa(game.ID)] = game_data

		}
	}

	return games

}
