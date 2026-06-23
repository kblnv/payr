package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"payr/pkg/plugins"
)

type Config struct {
	Repo     string
	DataPath string
}

type GitHubAPIResponse struct {
	TagName string `json:"tag_name"`
	HtmlUrl string `json:"html_url"`
}

type GitHubReleases struct {
	repo      string
	stateFile string
}

func New(rawConfig json.RawMessage) (plugins.Plugin, error) {
	var config Config

	err := json.Unmarshal(rawConfig, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if config.Repo == "" {
		return nil, fmt.Errorf("repo is required")
	}

	if config.DataPath == "" {
		return nil, fmt.Errorf("dataPath is required")
	}

	return &GitHubReleases{
		repo:      config.Repo,
		stateFile: config.DataPath,
	}, nil
}

func (g *GitHubReleases) Execute(ctx *plugins.Context) (string, error) {
	tag, url, err := g.fetchLatestRelease()
	if err != nil {
		return "", err
	}

	previousTag, err := g.readState()
	if err != nil {
		return "", err
	}

	if previousTag == "" {
		err = g.writeState(tag)
		if err != nil {
			return "", err
		}
		return "", nil
	}

	if previousTag == tag {
		return "", nil
	}

	err = g.writeState(tag)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("New release: %s — %s", tag, url), nil
}

func (g *GitHubReleases) fetchLatestRelease() (string, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/releases/latest", g.repo)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "payr")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("GitHub API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var release GitHubAPIResponse
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse response: %w", err)
	}

	return release.TagName, release.HtmlUrl, nil
}

func (g *GitHubReleases) readState() (string, error) {
	data, err := os.ReadFile(g.stateFile)
	if err != nil {
		if os.IsNotExist(err) {
			return "", nil
		}
		return "", fmt.Errorf("failed to read state file: %w", err)
	}

	return strings.TrimSpace(string(data)), nil
}

func (g *GitHubReleases) writeState(tag string) error {
	return os.WriteFile(g.stateFile, []byte(tag), 0644)
}
