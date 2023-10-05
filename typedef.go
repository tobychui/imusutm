package main

/*
	Uptime Monitor Configuration and Types
*/
type Record struct {
	Timestamp  int64
	ID         string
	Name       string
	URL        string
	Protocol   string
	Online     bool
	StatusCode int
	Latency    int64
}

type Target struct {
	ID       string
	Name     string
	URL      string
	Protocol string
}
type Config struct {
	Targets       []*Target
	Interval      int
	RecordsInJson int
	LogToFile     bool
}

/*
	TOTP Configs
*/
type TotpEntry struct {
	Name   string
	Secret string
	Link   string
}

type TotpCode struct {
	Name     string
	Code     string
	Link     string
	Succ     bool
	ValidFor int //How many seconds this code will be valid for
}

type TotpConfig struct {
	Entries []*TotpEntry
}
