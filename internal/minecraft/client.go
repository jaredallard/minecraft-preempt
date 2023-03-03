// Copyright (C) 2022 Jared Allard <jared@rgst.io>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package minecraft

import (
	"encoding/json"
	"fmt"

	mcnet "github.com/Tnze/go-mc/net"
	pk "github.com/Tnze/go-mc/net/packet"
	"github.com/pkg/errors"
)

// Client is a minecraft protocol aware connection somewhat
type Client struct {
	*mcnet.Conn

	ProtocolVersion int32
}

// Handshake reads the handshake packet and returns the next state
func (c *Client) Handshake() (nextState int32, original *pk.Packet, err error) {
	var p pk.Packet
	if err := c.ReadPacket(&p); err != nil {
		return -1, nil, err
	}
	if p.ID != 0 {
		return -1, nil, fmt.Errorf("packet ID 0x%X is not handshake", p.ID)
	}

	var (
		sid pk.String
		spt pk.Short
	)
	if err := p.Scan(
		(*pk.VarInt)(&c.ProtocolVersion),
		&sid, &spt,
		(*pk.VarInt)(&nextState)); err != nil {
		return -1, nil, err
	}

	return nextState, &p, nil
}

// SendStatus sends status request and ping packets to the client
func (c *Client) SendStatus(status *Status) error {
	for i := 0; i < 2; i++ {
		var p pk.Packet
		if err := c.ReadPacket(&p); err != nil {
			break
		}

		var err error
		switch p.ID {
		case 0x00: // Status Request
			b, jerr := json.Marshal(status)
			if jerr == nil {
				err = c.WritePacket(pk.Marshal(p.ID, pk.String(b)))
			}
		case 0x01: // Ping
			err = c.WritePacket(p)
		}
		if err != nil {
			return errors.Wrapf(err, "failed to send %d packet", p.ID)
		}
	}
	return nil
}
