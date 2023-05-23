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
	"fmt"
	cdkgo "github.com/vanus-labs/cdk-go"
	"github.com/vanus-labs/connector/source/chatai/chat"
)

var _ cdkgo.SourceConfigAccessor = &whatsAppConfig{}

type whatsAppConfig struct {
	cdkgo.SourceConfig `json:",inline" yaml:",inline"`
	VanusCloud         bool `json:"vanus_cloud" yaml:"vanus_cloud"`

	*chat.ChatConfig `json:",inline" yaml:",inline"`
	EnableChatAi     bool `json:"enable_chatai" yaml:"enable_chatai"`
}

func NewConfig() cdkgo.SourceConfigAccessor {
	return &whatsAppConfig{}
}

func (c *whatsAppConfig) Validate() error {
	if c.EnableChatAi {
		if c.ChatConfig == nil {
			return fmt.Errorf("enable chat but chat config is empty")
		}
		err := c.ChatConfig.Validate()
		if err != nil {
			return err
		}
	}
	return c.SourceConfig.Validate()
}

