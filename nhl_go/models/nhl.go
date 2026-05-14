package models

import "time"

type NHLSchedule struct {
	PreviousStartDate      string         `json:"previousStartDate"`
	GameWeek               []GameWeek     `json:"gameWeek"`
	OddsPartners           []OddsPartners `json:"oddsPartners"`
	PreSeasonStartDate     string         `json:"preSeasonStartDate"`
	RegularSeasonStartDate string         `json:"regularSeasonStartDate"`
	RegularSeasonEndDate   string         `json:"regularSeasonEndDate"`
	PlayoffEndDate         string         `json:"playoffEndDate"`
	NumberOfGames          int            `json:"numberOfGames"`
}
type Venue struct {
	Default string `json:"default"`
}
type TvBroadcasts struct {
	ID             int    `json:"id"`
	Market         string `json:"market"`
	CountryCode    string `json:"countryCode"`
	Network        string `json:"network"`
	SequenceNumber int    `json:"sequenceNumber"`
}
type CommonName struct {
	Default string `json:"default"`
}
type PlaceName struct {
	Default string `json:"default"`
}
type PlaceNameWithPreposition struct {
	Default string `json:"default"`
	Fr      string `json:"fr"`
}
type AwayTeam struct {
	ID                       int                      `json:"id"`
	CommonName               CommonName               `json:"commonName"`
	PlaceName                PlaceName                `json:"placeName"`
	PlaceNameWithPreposition PlaceNameWithPreposition `json:"placeNameWithPreposition"`
	Abbrev                   string                   `json:"abbrev"`
	Logo                     string                   `json:"logo"`
	DarkLogo                 string                   `json:"darkLogo"`
	AwaySplitSquad           bool                     `json:"awaySplitSquad"`
	Score                    int                      `json:"score"`
}
type HomeTeamCommonName struct {
	Default string `json:"default"`
}
type HomeTeamPlaceName struct {
	Default string `json:"default"`
}
type HomeTeamPlaceNameWithPreposition struct {
	Default string `json:"default"`
	Fr      string `json:"fr"`
}
type HomeTeam struct {
	ID                               int                              `json:"id"`
	HomeTeamCommonName               HomeTeamCommonName               `json:"commonName"`
	HomeTeamPlaceName                HomeTeamPlaceName                `json:"placeName"`
	HomeTeamPlaceNameWithPreposition HomeTeamPlaceNameWithPreposition `json:"placeNameWithPreposition"`
	Abbrev                           string                           `json:"abbrev"`
	Logo                             string                           `json:"logo"`
	DarkLogo                         string                           `json:"darkLogo"`
	HomeSplitSquad                   bool                             `json:"homeSplitSquad"`
	Score                            int                              `json:"score"`
}
type PeriodDescriptor struct {
	Number               int    `json:"number"`
	PeriodType           string `json:"periodType"`
	MaxRegulationPeriods int    `json:"maxRegulationPeriods"`
}
type GameOutcome struct {
	LastPeriodType string `json:"lastPeriodType"`
}
type FirstInitial struct {
	Default string `json:"default"`
}
type LastName struct {
	Default string `json:"default"`
}
type WinningGoalie struct {
	PlayerID     int          `json:"playerId"`
	FirstInitial FirstInitial `json:"firstInitial"`
	LastName     LastName     `json:"lastName"`
}
type WinningGoalScorerFirstInitial struct {
	Default string `json:"default"`
}
type WinningGoalScorerLastName struct {
	Default string `json:"default"`
}
type WinningGoalScorer struct {
	PlayerID                      int                           `json:"playerId"`
	WinningGoalScorerFirstInitial WinningGoalScorerFirstInitial `json:"firstInitial"`
	WinningGoalScorerLastName     WinningGoalScorerLastName     `json:"lastName"`
}
type SeriesStatus struct {
	Round                int    `json:"round"`
	SeriesAbbrev         string `json:"seriesAbbrev"`
	SeriesTitle          string `json:"seriesTitle"`
	SeriesLetter         string `json:"seriesLetter"`
	NeededToWin          int    `json:"neededToWin"`
	TopSeedTeamAbbrev    string `json:"topSeedTeamAbbrev"`
	TopSeedWins          int    `json:"topSeedWins"`
	BottomSeedTeamAbbrev string `json:"bottomSeedTeamAbbrev"`
	BottomSeedWins       int    `json:"bottomSeedWins"`
	GameNumberOfSeries   int    `json:"gameNumberOfSeries"`
}
type Games struct {
	ID                int               `json:"id"`
	Season            int               `json:"season"`
	GameType          int               `json:"gameType"`
	Venue             Venue             `json:"venue"`
	NeutralSite       bool              `json:"neutralSite"`
	StartTimeUTC      time.Time         `json:"startTimeUTC"`
	EasternUTCOffset  string            `json:"easternUTCOffset"`
	VenueUTCOffset    string            `json:"venueUTCOffset"`
	VenueTimezone     string            `json:"venueTimezone"`
	GameState         string            `json:"gameState"`
	GameScheduleState string            `json:"gameScheduleState"`
	TvBroadcasts      []TvBroadcasts    `json:"tvBroadcasts"`
	AwayTeam          AwayTeam          `json:"awayTeam"`
	HomeTeam          HomeTeam          `json:"homeTeam"`
	PeriodDescriptor  PeriodDescriptor  `json:"periodDescriptor"`
	GameOutcome       GameOutcome       `json:"gameOutcome"`
	WinningGoalie     WinningGoalie     `json:"winningGoalie"`
	WinningGoalScorer WinningGoalScorer `json:"winningGoalScorer"`
	SeriesStatus      SeriesStatus      `json:"seriesStatus"`
	SeriesURL         string            `json:"seriesUrl"`
	ThreeMinRecap     string            `json:"threeMinRecap"`
	ThreeMinRecapFr   string            `json:"threeMinRecapFr"`
	CondensedGame     string            `json:"condensedGame"`
	CondensedGameFr   string            `json:"condensedGameFr"`
	GameCenterLink    string            `json:"gameCenterLink"`
}
type GameWeek struct {
	Date          string  `json:"date"`
	DayAbbrev     string  `json:"dayAbbrev"`
	NumberOfGames int     `json:"numberOfGames"`
	DatePromo     []any   `json:"datePromo"`
	Games         []Games `json:"games"`
}
type OddsPartners struct {
	PartnerID   int    `json:"partnerId"`
	Country     string `json:"country"`
	Name        string `json:"name"`
	ImageURL    string `json:"imageUrl"`
	SiteURL     string `json:"siteUrl"`
	BgColor     string `json:"bgColor"`
	TextColor   string `json:"textColor"`
	AccentColor string `json:"accentColor"`
}
