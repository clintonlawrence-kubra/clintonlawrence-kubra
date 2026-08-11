// package main

// import (
// 	"crypto"
// 	"crypto/rand"
// 	"crypto/rsa"
// 	"crypto/sha256"
// 	"crypto/x509"
// 	"encoding/base64"
// 	"encoding/json"
// 	"encoding/pem"
// 	"fmt"
// 	"io"
// 	"log/slog"
// 	"net/http"
// 	"os"
// 	"strings"
// 	"time"

// 	"github.com/joho/godotenv"
// )

// const githubAPI = "https://api.github.com"

// func main() {

// 	repo_name := "clinton-cci-test"

// 	//Load local .env if exists
// 	if _, err := os.Stat(".env"); err == nil {
// 		err := godotenv.Load(".env")
// 		if err != nil {
// 			slog.Error("Error loading .env file", "error", err)
// 		}
// 	}

// 	repos := []string{
// 		repo_name,
// 		repo_name + "-config",
// 		repo_name + "-infra",
// 	}

// 	for _, repo := range repos {
// 		err := deleteRepo(repo)
// 		if err != nil {
// 			if strings.Contains(err.Error(), "404 Not Found") {
// 				slog.Warn("Repo not found, skipping deletion", "repo", repo)
// 				continue
// 			}
// 			slog.Error("Error deleting repo", "repo", repo, "error", err)
// 		} else {
// 			slog.Info("Successfully deleted repo", "repo", repo)
// 		}
// 	}

// }

// // cached installation token + owner, reused across deleteRepo calls
// var cachedToken string

// func deleteRepo(repo string) error {

// 	token, err := installationAuth()
// 	if err != nil {
// 		return err
// 	}

// 	repoOwner := "kubra-sandbox"

// 	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/repos/%s/%s", githubAPI, repoOwner, repo), nil)
// 	if err != nil {
// 		return err
// 	}
// 	req.Header.Set("Authorization", "Bearer "+token)
// 	req.Header.Set("Accept", "application/vnd.github+json")
// 	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusNoContent {
// 		body, _ := io.ReadAll(resp.Body)
// 		return fmt.Errorf("delete repo %s/%s failed: %s: %s", repoOwner, repo, resp.Status, strings.TrimSpace(string(body)))
// 	}
// 	return nil
// }

// // installationAuth mints (and caches) an installation access token from the
// // GitHub App credentials in the environment, returning the token and the owner
// // (login of the account the app is installed on).
// func installationAuth() (installToken string, err error) {
// 	if cachedToken != "" {
// 		return cachedToken, nil
// 	}

// 	clientID := os.Getenv("GITHUB_APP_CLIENT_ID")
// 	installationID := os.Getenv("GITHUB_INSTALLATION_ID")
// 	keyB64 := os.Getenv("GITHUB_APP_PRIVATE_KEY_B64")
// 	if clientID == "" || installationID == "" || keyB64 == "" {
// 		return "", fmt.Errorf("missing GITHUB_APP_CLIENT_ID, GITHUB_INSTALLATION_ID, or GITHUB_APP_PRIVATE_KEY_B64")
// 	}

// 	key, err := parsePrivateKey(keyB64)
// 	if err != nil {
// 		return "", fmt.Errorf("parse private key: %w", err)
// 	}

// 	jwtToken, err := generateJWT(clientID, key)
// 	if err != nil {
// 		return "", fmt.Errorf("generate jwt: %w", err)
// 	}

// 	installToken, err = installationToken(jwtToken, installationID)
// 	if err != nil {
// 		return "", fmt.Errorf("create installation token: %w", err)
// 	}

// 	cachedToken = installToken
// 	return installToken, nil
// }

// func parsePrivateKey(keyB64 string) (*rsa.PrivateKey, error) {
// 	pemBytes, err := base64.StdEncoding.DecodeString(strings.TrimSpace(keyB64))
// 	if err != nil {
// 		return nil, fmt.Errorf("base64 decode: %w", err)
// 	}

// 	block, _ := pem.Decode(pemBytes)
// 	if block == nil {
// 		return nil, fmt.Errorf("no PEM block found in decoded key")
// 	}

// 	// GitHub App keys are PKCS#1 ("RSA PRIVATE KEY"); fall back to PKCS#8.
// 	if key, err := x509.ParsePKCS1PrivateKey(block.Bytes); err == nil {
// 		return key, nil
// 	}
// 	parsed, err := x509.ParsePKCS8PrivateKey(block.Bytes)
// 	if err != nil {
// 		return nil, err
// 	}
// 	key, ok := parsed.(*rsa.PrivateKey)
// 	if !ok {
// 		return nil, fmt.Errorf("private key is not RSA")
// 	}
// 	return key, nil
// }

// func generateJWT(issuer string, key *rsa.PrivateKey) (string, error) {
// 	now := time.Now()
// 	header := map[string]string{"alg": "RS256", "typ": "JWT"}
// 	claims := map[string]any{
// 		"iat": now.Add(-60 * time.Second).Unix(), // allow for clock drift
// 		"exp": now.Add(9 * time.Minute).Unix(),   // GitHub max is 10 minutes
// 		"iss": issuer,
// 	}

// 	headerJSON, err := json.Marshal(header)
// 	if err != nil {
// 		return "", err
// 	}
// 	claimsJSON, err := json.Marshal(claims)
// 	if err != nil {
// 		return "", err
// 	}

// 	enc := base64.RawURLEncoding
// 	signingInput := enc.EncodeToString(headerJSON) + "." + enc.EncodeToString(claimsJSON)

// 	hashed := sha256.Sum256([]byte(signingInput))
// 	sig, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA256, hashed[:])
// 	if err != nil {
// 		return "", err
// 	}

// 	return signingInput + "." + enc.EncodeToString(sig), nil
// }

// func installationToken(jwtToken, installationID string) (string, error) {
// 	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/app/installations/%s/access_tokens", githubAPI, installationID), nil)
// 	if err != nil {
// 		return "", err
// 	}
// 	req.Header.Set("Authorization", "Bearer "+jwtToken)
// 	req.Header.Set("Accept", "application/vnd.github+json")
// 	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

// 	resp, err := http.DefaultClient.Do(req)
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusCreated {
// 		body, _ := io.ReadAll(resp.Body)
// 		return "", fmt.Errorf("%s: %s", resp.Status, strings.TrimSpace(string(body)))
// 	}

// 	var out struct {
// 		Token string `json:"token"`
// 	}
// 	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
// 		return "", err
// 	}
// 	if out.Token == "" {
// 		return "", fmt.Errorf("no token in access_tokens response")
// 	}
// 	return out.Token, nil
// }
