package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
)

const githubTeamsAPI = "https://api.github.com"

// githubToken is a personal access token with read:org scope (SSO-authorized
// for the org if it enforces SAML SSO), read from the GITHUB_TOKEN env var.
var githubToken = os.Getenv("GITHUB_TOKEN")

func main() {

	org := "kubra-hq"

	teams, err := listTeams(org)
	if err != nil {
		slog.Error("Error listing teams", "org", org, "error", err)
		os.Exit(1)
	}

	printTeamTree(buildTeamTree(teams), 0)
}

// Team is the subset of GitHub's team object needed to reconstruct the
// parent/child hierarchy from a flat list.
type Team struct {
	ID     int    `json:"id"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
	Parent *struct {
		ID int `json:"id"`
	} `json:"parent"`

	Children []*Team `json:"-"`
}

// listTeams fetches every team in the org. GitHub returns child teams inline
// in this same paginated list, distinguished only by their "parent" field.
func listTeams(org string) ([]*Team, error) {
	if githubToken == "" || githubToken == "PASTE_TOKEN_HERE" {
		return nil, fmt.Errorf("set githubToken in get_teams.go to a personal access token first")
	}

	var all []*Team
	for page := 1; ; page++ {
		url := fmt.Sprintf("%s/orgs/%s/teams?per_page=100&page=%d", githubTeamsAPI, org, page)
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+githubToken)
		req.Header.Set("Accept", "application/vnd.github+json")
		req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("list teams for %s failed: %s: %s", org, resp.Status, strings.TrimSpace(string(body)))
		}

		var batch []*Team
		if err := json.NewDecoder(resp.Body).Decode(&batch); err != nil {
			return nil, err
		}
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
	}

	return all, nil
}

// buildTeamTree groups the flat team list by parent, returning the top-level
// (parentless) teams with their descendants attached via Children.
func buildTeamTree(teams []*Team) []*Team {
	byID := make(map[int]*Team, len(teams))
	for _, t := range teams {
		byID[t.ID] = t
	}

	var roots []*Team
	for _, t := range teams {
		if t.Parent == nil {
			roots = append(roots, t)
			continue
		}
		if parent, ok := byID[t.Parent.ID]; ok {
			parent.Children = append(parent.Children, t)
		} else {
			roots = append(roots, t) // parent outside the fetched set; treat as root
		}
	}
	return roots
}

func printTeamTree(teams []*Team, depth int) {
	for _, t := range teams {
		fmt.Printf("%s- %s (%s)\n", strings.Repeat("  ", depth), t.Name, t.Slug)
		printTeamTree(t.Children, depth+1)
	}
}
