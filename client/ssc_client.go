package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/SpectraLogic/ssc_go_client/openapi"
	"net/http"
)

type SscClient struct {
	Client      *openapi.APIClient
	Credentials *openapi.ApiCredentials
	Context     *context.Context
}

func CreateClient(args *Arguments) (*SscClient, error) {
	//  Create client
	config := openapi.NewConfiguration(args.Url)
	if args.IgnoreCert {
		tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
		httpClient := &http.Client{Transport: tr}
		config.HTTPClient = httpClient
	}

	storCycle := openapi.NewAPIClient(config)
	plainPassword := args.Password
	var err error
	if args.Encrypted {
		plainPassword, err = Decrypt(args.Password)
		if err != nil {
			return nil, fmt.Errorf("could not decrypt password %v\n", err)
		}
	}

	creds := openapi.ApiCredentials{args.Domain, plainPassword, args.Name}
	token, _, err := storCycle.AuthenticationApi.Login(context.TODO(), creds)
	if err != nil {
		fmt.Printf("Failed to get token %v\n", err)
		return nil, fmt.Errorf("could not get token %v\n", err)
	}
	ctx := context.WithValue(context.Background(), openapi.ContextAccessToken, token.Token)
	ret := &SscClient{storCycle, &creds, &ctx}
	return ret, nil
}

/*
Context should be request-scoped. But many simple, single request calls use the Context member
This method is preferred (and can update the token), but the legacy will work
*/
func (sscClient SscClient) getContext(updateToken bool) (*context.Context, error) {
	if !updateToken && sscClient.Context != nil {
		return sscClient.Context, nil
	}
	token, _, err := sscClient.Client.AuthenticationApi.Login(context.TODO(), *sscClient.Credentials)
	if err != nil {
		return nil, fmt.Errorf("could not get token for user %s %v\n", sscClient.Credentials.Username, err)
	}
	ret := context.WithValue(context.Background(), openapi.ContextAccessToken, token.Token)
	return &ret, nil
}

func (sscClient SscClient) updateToken() (*SscClient, error) {
	newCtx, err := sscClient.getContext(true)
	if err != nil {
		return nil, fmt.Errorf("could not get new token %v\n", err)
	}
	return &SscClient{sscClient.Client, sscClient.Credentials, newCtx}, nil
}

func (sscClient SscClient) getToken() interface{} {
	ctx := *sscClient.Context
	return ctx.Value(openapi.ContextAccessToken)
}

func NowSchedule() *openapi.ApiProjectSchedule {
	now := "Now"
	return &openapi.ApiProjectSchedule{Period: &now}
}

// get validation errors available in GenericOpenAPIError
func ExpandOpenApiErr(e error) error {
	detailError, ok := e.(openapi.GenericOpenAPIError)
	if ok {
		return fmt.Errorf("%v\n%s", detailError, detailError.Body())
	}
	return e
}
