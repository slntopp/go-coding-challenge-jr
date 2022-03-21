package server

import (
	"challenge/pkg/api"
	"challenge/util"
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
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
	return &api.Link{Data: string(bodyText)}, nil
}

func (s *GRPCServer) ReadMetadata(ctx context.Context, placeholder *api.Placeholder) (*api.Placeholder, error) {
	return &api.Placeholder{Data: fmt.Sprintf("%v", ctx.Value("i-am-random-key"))}, nil
}

func (s *GRPCServer) StartTimer(imer *api.Timer, server api.ChallengeService_StartTimerServer) error {
	return nil
}
