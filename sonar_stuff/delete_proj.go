package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
)

const projectKey = "kubra-hq_clinton-sonar-cci2" // hard-coded target

func main() {

	err := DeleteProject(projectKey)
	if err != nil {
		fmt.Printf("Error deleting project: %v\n", err)
	} else {
		fmt.Println("Successfully deleted project:", projectKey)
	}
}

// DeleteProject deletes a SonarQube project by key.
// Auth: KUBRA_SONARCLOUD_TOKEN env var (passed as the basic-auth username).
// https://sonarqube.us/web_api/api/projects/delete
func DeleteProject(projectKey string) error {

	req, err := http.NewRequest(http.MethodPost,
		"https://sonarqube.us/api/projects/delete",
		nil)
	if err != nil {
		return err
	}
	req.URL.RawQuery = url.Values{"project": {projectKey}}.Encode()
	req.SetBasicAuth(os.Getenv("KUBRA_SONARCLOUD_TOKEN"), "")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("delete %q: %s", projectKey, resp.Status)
	}
	return nil
}
