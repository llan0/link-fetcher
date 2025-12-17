package config

import (
	"encoding/json"
	"os"
	"time"
)

type Accounts struct {
	YoutubeChannelID  string `json:"youtube_channel_id"`
	TikTokUsername    string `json:"tiktok_username"`
	InstagramUsername string `json:"instagram_username"`
	InstagramUserID   string `json:"instagram_user_id"`
}

type APIKeys struct {
	Youtube       string `json:"youtube"`
	RapidAPIKey   string `json:"rapid_api_key"`
	TikTokHost    string `json:"tiktok_host"`
	InstagramHost string `json:"instagram_host"`
}

type Config struct {
	StartDateStr string    `json:"start_date"`
	EndDateStr   string    `json:"end_date"`
	Accounts     Accounts  `json:"accounts"`
	APIKeys      APIKeys   `json:"api_keys"`
	StartDate    time.Time `json:"-"`
	EndDate      time.Time `json:"-"`
}

func Load(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cfg Config
	if err := json.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, err
	}

	layout := "2006-01-02"
	start, err := time.Parse(layout, cfg.StartDateStr)
	if err != nil {
		return nil, err
	}
	end, err := time.Parse(layout, cfg.EndDateStr)
	if err != nil {
		return nil, err
	}

	cfg.StartDate = start
	cfg.EndDate = end.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	return &cfg, nil
}
