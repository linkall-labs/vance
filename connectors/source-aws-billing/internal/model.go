// Copyright 2022 Linkall Inc.
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

type BillingData struct {
	Date             string `json:"date"`
	Service          string `json:"service"`
	AmortizedCost    Cost   `json:"amortizedCost"`
	BlendedCost      Cost   `json:"blendedCost"`
	NetAmortizedCost Cost   `json:"netAmortizedCost"`
	NetUnblendedCost Cost   `json:"netUnblendedCost"`
	UnblendedCost    Cost   `json:"unblendedCost"`
}

type Cost struct {
	Amount *string `json:"amount"`
	Unit   *string `json:"unit"`
}
