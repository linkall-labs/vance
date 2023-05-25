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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"

	ce "github.com/cloudevents/sdk-go/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/skip2/go-qrcode"
	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/cdk-go/log"
	"github.com/vanus-labs/connector/source/chatai/chat"
)

var _ cdkgo.Source = &whatsAppSource{}

func NewWhatsAppSource() cdkgo.Source {
	return &whatsAppSource{
		events: make(chan *cdkgo.Tuple, 100),
	}
}

type whatsAppSource struct {
	config      *whatsAppConfig
	events      chan *cdkgo.Tuple
	client      *whatsmeow.Client
	chatService *chat.ChatService
}

func (s *whatsAppSource) Initialize(ctx context.Context, cfg cdkgo.ConfigAccessor) error {
	s.config = cfg.(*whatsAppConfig)
	dbFileName := "store.db"
	if s.config.Data != "" {
		dbBytes, err := base64.StdEncoding.DecodeString(s.config.Data)
		if err != nil {
			return err
		}
		err = os.WriteFile(dbFileName, dbBytes, 0644)
		if err != nil {
			return err
		}
		log.Info("Database restored successfully.", nil)
	}

	if s.config.EnableChatAi {
		s.chatService = chat.NewChatService(*s.config.ChatConfig)
	}

	dbLog := waLog.Stdout("Database", "DEBUG", true)
	container, err := sqlstore.New("sqlite3", fmt.Sprintf("file:%s?_foreign_keys=on", dbFileName), dbLog)
	if err != nil {
		return err
	}

	deviceStore, err := container.GetFirstDevice()
	if err != nil {
		return err
	}

	clientLog := waLog.Stdout("Client", "DEBUG", true)
	s.client = whatsmeow.NewClient(deviceStore, clientLog)

	s.events = make(chan *cdkgo.Tuple, 100)
	s.client.AddEventHandler(func(evt interface{}) {

		switch v := evt.(type) {
		case *events.Message:

			info := v.Info
			message := v.Message.GetConversation()
			if message != "" {
				if v.Info.IsFromMe {
					if v.Info.Sender.User == v.Info.Chat.User {
						event := s.makeEvent(info, message)
						s.events <- &cdkgo.Tuple{
							Event: event,
							Success: func() {
								// TODO
								b, _ := json.Marshal(event)
								log.Info("send event success: "+string(b), nil)
							},
							Failed: func(err error) {
								// TODO
								b, _ := json.Marshal(event)
								log.Info("send event failed: "+string(b), map[string]interface{}{
									log.KeyError: err,
								})
							},
						}
					}

				} else {
					event := s.makeEvent(info, message)
					s.events <- &cdkgo.Tuple{
						Event: event,
						Success: func() {
							// TODO
							b, _ := json.Marshal(event)
							log.Info("send event success: "+string(b), nil)
						},
						Failed: func(err error) {
							// TODO
							b, _ := json.Marshal(event)
							log.Info("send event failed: "+string(b), map[string]interface{}{
								log.KeyError: err,
							})
						},
					}
				}
			}
		}

	})

	if s.client.Store.ID == nil {
		if s.config.Data != "" {
			return fmt.Errorf("data exist but store id is nil")
		}
		qrChan, _ := s.client.GetQRChannel(context.Background())
		err = s.client.Connect()
		if err != nil {
			return err
		}
		for evt := range qrChan {
			if evt.Event == "code" {
				err := qrcode.WriteFile(evt.Code, qrcode.Medium, 256, "qr.png")

				if err != nil {
					fmt.Println("Failed to generate QR code:", err)
				} else {
					fmt.Println("QR code generated successfully")
				}

			} else {
				fmt.Println("Login event:", evt.Event)
			}
		}
	} else {
		err = s.client.Connect()
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *whatsAppSource) Name() string {
	return "WhatsAppSource"
}

func (s *whatsAppSource) Destroy() error {
	if s.client != nil {
		s.client.Disconnect()
	}
	if s.chatService != nil {
		s.chatService.Close()
	}
	return nil
}

func (s *whatsAppSource) Chan() <-chan *cdkgo.Tuple {
	return s.events
}

func (s *whatsAppSource) makeEvent(info types.MessageInfo, message string) *ce.Event {
	event := ce.NewEvent()
	event.SetID(info.ID)
	event.SetSource("whatsapp")
	event.SetType(info.Type)
	if s.chatService != nil {
		resp, err := s.chatService.ChatCompletion(context.Background(), s.config.ChatConfig.DefaultChatMode, info.Sender.User, message)
		if err != nil {
			log.Warning("failed to get content from Chat", map[string]interface{}{
				log.KeyError: err,
			})
		}
		event.SetData(ce.ApplicationJSON, map[string]interface{}{
			"info":    info,
			"message": resp,
		})
	} else {
		event.SetData(ce.ApplicationJSON, map[string]interface{}{
			"info":    info,
			"message": message,
		})
	}
	return &event
}
