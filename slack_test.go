// Copyright 2015 Bowery, Inc.

package slack

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	testClient     *Client
	testChannel    = "#testing"
	testBadChannel = "#foobar"
	testText       = "trying this out"
	testUsername   = "drizzy drake"
)

func init() {
	testClient = NewClient("some-token")
}

func TestSendMessageSuccessful(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(sendMessageHandlerOK))
	defer server.Close()
	slackAddr = server.URL

	err := testClient.SendMessage(testChannel, testText, testUsername)
	if err != nil {
		t.Error(err)
	}
}

func sendMessageHandlerOK(rw http.ResponseWriter, req *http.Request) {
	res := &slackPostMessageRes{Ok: true}
	body, _ := json.Marshal(res)
	rw.Write(body)
}

func TestSendMessageMissingArgument(t *testing.T) {
	t.Parallel()
	err := testClient.SendMessage("", testText, testUsername)
	if err == nil {
		t.Error("should have failed, channel missing")
	}
}

func TestSendMessageBadResponse(t *testing.T) {
	t.Parallel()
	server := httptest.NewServer(http.HandlerFunc(sendMessageHandlerBad))
	defer server.Close()
	slackAddr = server.URL

	err := testClient.SendMessage(testBadChannel, testText, testUsername)
	if err == nil {
		t.Error("should have failed, invalid channel")
	}
	if err.Error() != "channel_not_found" {
		t.Error("received unexpected error")
	}
}

func sendMessageHandlerBad(rw http.ResponseWriter, req *http.Request) {
	res := &slackPostMessageRes{Ok: false, Error: "channel_not_found"}
	body, _ := json.Marshal(res)
	rw.Write(body)
}
