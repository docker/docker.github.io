// Copyright 2013 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package ipv6

import "net"

// A payloadHandler represents the IPv6 datagram payload handler.
type payloadHandler struct {
	net.PacketConn
	rawOpt
}

func (c *payloadHandler) ok() bool { return c != nil && c.PacketConn != nil }
