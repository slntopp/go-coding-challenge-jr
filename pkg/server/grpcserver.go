package server

import (
	"challenge/pkg/api"
	"challenge/util"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type GRPCServer struct{}

func (s *GRPCServer) MakeShortLink(ctx context.Context, link *api.Link) (*api.Link, error) {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	client := &http.Client{}

	var data = fmt.Sprintf("{ \"long_url\": \"%s\", \"domain\": \"bit.ly\" }", link.GetData())
	var body = strings.NewReader(data)
	req, err := http.NewRequest("POST", "https://api-ssl.bitly.com/v4/shorten", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", config.BitlyOauthToken))

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var jsonBody map[string]interface{}
	err = json.Unmarshal(bodyText, &jsonBody)
	if err != nil {
		log.Fatal(err)
	}

	if shortLink, linkExist := jsonBody["link"].(string); linkExist {
		return &api.Link{Data: shortLink}, nil
	} else {
		err = fmt.Errorf("%v", jsonBody["message"])
		return nil, err
	}
}

func (s *GRPCServer) ReadMetadata(ctx context.Context, placeholder *api.Placeholder) (*api.Placeholder, error) {
	return &api.Placeholder{Data: fmt.Sprintf("%v", ctx.Value("i-am-random-key"))}, nil
}

func (s *GRPCServer) StartTimer(timer *api.Timer, server api.ChallengeService_StartTimerServer) error {
	client := &http.Client{}

	CheckTime := func(timerName string) int {
		req, err := http.NewRequest("GET", "https://timercheck.io"+"/"+timerName, nil)
		if err != nil {
			log.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()

		bodyText, err := ioutil.ReadAll(resp.Body)
		var jsonBody map[string]interface{}
		err = json.Unmarshal(bodyText, &jsonBody)
		if err != nil {
			log.Fatal(err)
		}

		if secondsRemaining, timerExist := jsonBody["seconds_remaining"].(float64); timerExist {
			return int(secondsRemaining)
		} else {
			return 0
		}
	}

	if CheckTime(timer.GetName()) == 0 {
		req, err := http.NewRequest("GET", "https://timercheck.io"+"/"+timer.GetName()+"/"+strconv.Itoa(int(timer.GetSeconds())), nil)
		if err != nil {
			log.Fatal(err)
		}

		_, err = client.Do(req)
		if err != nil {
			log.Fatal(err)
		}
	}

	for {
		secondsRemaining := CheckTime(timer.GetName())
		if secondsRemaining == 0 {
			break
		}
		timer.Seconds = int64(secondsRemaining)
		server.Send(timer)
		frequencyTimer := time.NewTimer(time.Second * time.Duration(timer.GetFrequency()))
		<-frequencyTimer.C
	}

	return nil
}
