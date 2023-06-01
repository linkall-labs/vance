// Copyright 2023 Linkall Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package internal

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/slack-go/slack"
	"github.com/slack-go/slack/slackevents"

	ce "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"

	cdkgo "github.com/vanus-labs/cdk-go"
)

var (
	errMissingHubSignatureHeader = errors.New("missing signature Header")
	errInvalidSignatureHeader    = errors.New("invalid signature Header")
	errInvalidContentTypeHeader  = errors.New("only support application/json Content-Type Header")
	errReadPayload               = errors.New("error read payload")
	errVerificationFailed        = errors.New("signature verification failed")
	errVerificationTokenFailed   = errors.New("token verification failed")
)

func (s *slackSource) verifyRequestSignature(req *http.Request, body []byte) error {
	sv, err := slack.NewSecretsVerifier(req.Header, s.config.SigningSecret)
	if err != nil {
		if err == slack.ErrMissingHeaders {
			return errMissingHubSignatureHeader
		}
		s.logger.Info().Err(err).Msg("new secret verifier failed")
		return errInvalidSignatureHeader
	}
	_, _ = sv.Write(body)
	if err := sv.Ensure(); err != nil {
		return errVerificationFailed
	}
	return nil
}

func (s *slackSource) event(w http.ResponseWriter, req *http.Request) error {
	contentType := req.Header.Get("Content-Type")
	if contentType != "application/json" {
		return errInvalidContentTypeHeader
	}
	body, err := ioutil.ReadAll(req.Body)
	if err != nil || len(body) == 0 {
		return errReadPayload
	}
	err = s.verifyRequestSignature(req, body)
	if err != nil {
		return err
	}
	var event map[string]interface{}
	err = json.Unmarshal(body, &event)
	if err != nil {
		return err
	}
	token, err := getStringValue(event, "token")
	if err != nil {
		return err
	}
	if !s.verifyToken.Verify(token) {
		return errVerificationTokenFailed
	}
	eventType, err := getStringValue(event, "type")
	if err != nil {
		return err
	}
	if eventType == slackevents.CallbackEvent {
		// https://api.slack.com/apis/connections/events-api#events-JSON
		s.makeEvent(event)
		w.Write([]byte("accepted"))
		return nil
	}
	if eventType == slackevents.URLVerification {
		// https://api.slack.com/apis/connections/events-api#verification
		challenge, err := getStringValue(event, "challenge")
		if err != nil {
			return err
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte(challenge))
		return nil
	}
	return nil
}

func getEventType(body map[string]interface{}) string {
	event, ok := body["event"].(map[string]interface{})
	if !ok {
		return ""
	}
	eventType, _ := getStringValue(event, "type")
	return eventType
}

func getEventUser(body map[string]interface{}) string {
	event, ok := body["event"].(map[string]interface{})
	if !ok {
		return ""
	}
	text, _ := getStringValue(event, "user")
	return text
}

func getEventText(body map[string]interface{}) string {
	event, ok := body["event"].(map[string]interface{})
	if !ok {
		return ""
	}
	text, _ := getStringValue(event, "text")
	return text
}

func (s *slackSource) makeEvent(body map[string]interface{}) error {
	eventType := getEventType(body)
	event := ce.NewEvent()
	event.SetID(uuid.NewString())
	event.SetSource(fmt.Sprintf("https://github.com/vanus-labs/vanus-connect/connectors/source-slack"))
	event.SetType("event_callback")
	event.SetExtension("eventtype", eventType)
	delete(body, "token")
	if s.chatService != nil && eventType == "app_mention" {
		go func(event *ce.Event, body map[string]interface{}) {
			text := getEventText(body)
			// <@U04M1L7L64U> msg
			arr := strings.SplitN(text, " ", 2)
			var content string
			if len(arr) == 2 {
				content = arr[1]
			}
			user := getEventUser(body)
			resp, err := s.chatService.ChatCompletion(context.Background(), s.config.ChatConfig.DefaultChatMode, user, content)
			if err != nil {
				s.logger.Warn().Err(err).Msg("failed to get content from Chat")
			}
			body["result"] = resp
			s.pushEvent(event, body)
		}(&event, body)
	} else {
		s.pushEvent(&event, body)
	}
	return nil
}

func (s *slackSource) pushEvent(event *ce.Event, body map[string]interface{}) {
	event.SetData(ce.ApplicationJSON, body)
	s.ch <- &cdkgo.Tuple{
		Event: event,
		Success: func() {
			s.logger.Info().Msg("send event success")
		},
		Failed: func(err error) {
			s.logger.Info().Err(err).Msg("send event failed")
		},
	}
}
