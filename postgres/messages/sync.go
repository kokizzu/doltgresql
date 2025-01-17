// Copyright 2023 Dolthub, Inc.
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

package messages

import "github.com/dolthub/doltgresql/postgres/connection"

func init() {
	connection.InitializeDefaultMessage(Sync{})
	connection.AddMessageHeader(Sync{})
}

// Sync represents a PostgreSQL message.
type Sync struct{}

var syncDefault = connection.MessageFormat{
	Name: "Sync",
	Fields: connection.FieldGroup{
		{
			Name:  "Header",
			Type:  connection.Byte1,
			Flags: connection.Header,
			Data:  int32('S'),
		},
		{
			Name:  "MessageLength",
			Type:  connection.Int32,
			Flags: connection.MessageLengthInclusive,
			Data:  int32(4),
		},
	},
}

var _ connection.Message = Sync{}

// Encode implements the interface connection.Message.
func (m Sync) Encode() (connection.MessageFormat, error) {
	return m.DefaultMessage().Copy(), nil
}

// Decode implements the interface connection.Message.
func (m Sync) Decode(s connection.MessageFormat) (connection.Message, error) {
	if err := s.MatchesStructure(*m.DefaultMessage()); err != nil {
		return nil, err
	}
	return Sync{}, nil
}

// DefaultMessage implements the interface connection.Message.
func (m Sync) DefaultMessage() *connection.MessageFormat {
	return &syncDefault
}
