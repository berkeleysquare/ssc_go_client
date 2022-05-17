package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"net/http"
)

type SscClient struct {
	Client *openapi.APIClient
	Context *context.Context
}

func CreateClient(args *Arguments) (*SscClient, error) {
	//  Create client
	config := openapi.NewConfiguration(args.Url)
	if args.IgnoreCert {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true},}
		httpClient := &http.Client{Transport: tr}
		config.HTTPClient = httpClient
	}

	storCycle := openapi.NewAPIClient(config)

	creds := openapi.ApiCredentials{args.Domain, args.Password, args.Name}
	token, _, err := storCycle.AuthenticationApi.Login(context.TODO(), creds)
	if err != nil {
		fmt.Printf("Failed to get token %v\n", err)
		return nil, fmt.Errorf("could not get token %v\n", err)
	}
	ctx := context.WithValue(context.Background(), openapi.ContextAccessToken, token.Token)
	ret := &SscClient{storCycle, &ctx}
	return ret, nil
}

func NowSchedule() *openapi.ApiProjectSchedule {
	now := "Now"
	return &openapi.ApiProjectSchedule{Period: &now}
}
