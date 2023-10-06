package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"imuslab.com/utm/pkg/utils"
)

func UptimeMonitorInit() error {
	log.Println("-- Uptime Monitor Started --")
	if !utils.FileExists(configFilepath) {
		log.Println("config.json not found. Template created.")
		template := Config{
			Targets:  []*Target{&exampleTarget},
			Interval: 60,
		}
		js, _ := json.MarshalIndent(template, "", " ")
		os.WriteFile(configFilepath, js, 0775)
		os.Exit(0)
	}

	c, err := os.ReadFile(configFilepath)
	if err != nil {
		return (err)
	}

	parsedConfig := Config{}
	err = json.Unmarshal(c, &parsedConfig)
	if err != nil {
		return (err)
	}

	usingConfig = &parsedConfig

	//Start the endpoint listener
	ticker := time.NewTicker(time.Duration(usingConfig.Interval) * time.Second)
	done := make(chan bool)

	go func() {
		//Start the uptime check once first before entering loop
		log.Println("Started initial uptime check. Might take a while before any results is shown on the web UI")
		ExecuteUptimeCheck()
		log.Println("Initial uptime check completed")
		for {
			select {
			case <-done:
				return
			case t := <-ticker.C:
				log.Println("Uptime updated - ", t.Unix())
				ExecuteUptimeCheck()
			}
		}
	}()

	return nil
}

func ExecuteUptimeCheck() {
	for _, target := range usingConfig.Targets {
		//For each target to check online, do the following
		var thisRecord Record
		if target.Protocol == "http" || target.Protocol == "https" {
			log.Println("Updating uptime status for " + target.Name)
			online, laterncy, statusCode := getWebsiteStatusWithLatency(target.URL)
			thisRecord = Record{
				Timestamp:  time.Now().Unix(),
				ID:         target.ID,
				Name:       target.Name,
				URL:        target.URL,
				Protocol:   target.Protocol,
				Online:     online,
				StatusCode: statusCode,
				Latency:    laterncy,
			}

			//fmt.Println(thisRecord)

		} else {
			log.Println("Unknown protocol: " + target.Protocol + ". Skipping")
			continue
		}

		thisRecords, ok := onlineStatusLog[target.ID]
		if !ok {
			//First record. Create the array
			onlineStatusLog[target.ID] = []*Record{&thisRecord}
		} else {
			//Append to the previous record
			thisRecords = append(thisRecords, &thisRecord)

			//Check if the record is longer than the logged record. If yes, clear out the old records
			if len(thisRecords) > usingConfig.RecordsInJson {
				thisRecords = thisRecords[1:]
			}

			onlineStatusLog[target.ID] = thisRecords

		}
	}

	//Write the results to a json file
	if usingConfig.LogToFile {
		//Log to file
		js, _ := json.MarshalIndent(onlineStatusLog, "", " ")
		os.WriteFile(logFilepath, js, 0775)
	}

}

/*
	Web Interface Handler
*/

func HandleUptimeLogRead(w http.ResponseWriter, r *http.Request) {
	id, _ := utils.GetPara(r, "id")
	if id == "" {
		js, _ := json.Marshal(onlineStatusLog)
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	} else {
		//Check if that id exists
		log, ok := onlineStatusLog[id]
		if !ok {
			http.NotFound(w, r)
			return
		}

		js, _ := json.MarshalIndent(log, "", " ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}

}

/*
	Utilities
*/

// Get website stauts with latency given URL, return is conn succ and its latency and status code
func getWebsiteStatusWithLatency(url string) (bool, int64, int) {
	start := time.Now().UnixNano() / int64(time.Millisecond)
	statusCode, err := getWebsiteStatus(url)
	end := time.Now().UnixNano() / int64(time.Millisecond)
	if err != nil {
		log.Println(err.Error())
		return false, 0, 0
	} else {
		diff := end - start
		succ := false
		if statusCode >= 200 && statusCode < 300 {
			//OK
			succ = true
		} else if statusCode >= 300 && statusCode < 400 {
			//Redirection code
			succ = true
		} else {
			succ = false
		}

		return succ, diff, statusCode
	}

}
func getWebsiteStatus(url string) (int, error) {
	resp, err := http.Get(url)
	if err != nil {
		return 0, err
	}
	status_code := resp.StatusCode
	resp.Body.Close()
	return status_code, nil
}
