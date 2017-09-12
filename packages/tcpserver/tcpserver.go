// Copyright 2016 The go-daylight Authors
// This file is part of the go-daylight library.
//
// The go-daylight library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-daylight library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-daylight library. If not, see <http://www.gnu.org/licenses/>.

package tcpserver

import (
	"flag"
	"fmt"

	"sync/atomic"

	"io"

	"github.com/EGaaS/go-egaas-mvp/packages/consts"
	logger "github.com/EGaaS/go-egaas-mvp/packages/log"
)

var (
	counter int64
)

func init() {
	flag.Parse()
}

// HandleTCPRequest proceed TCP requests
func HandleTCPRequest(rw io.ReadWriter) {
	logger.LogDebug(consts.FuncStarted, "")
	defer func() {
		atomic.AddInt64(&counter, -1)
	}()

	count := atomic.AddInt64(&counter, +1)
	if count > 20 {
		return
	}

	dType := &TransactionType{}
	err := ReadRequest(dType, rw)
	if err != nil {
		logger.LogError(consts.IOError, err)
		return
	}

	logger.LogDebug(consts.DebugMessage, fmt.Sprintf("tcpservers: got request type: %d", dType.Type))
	var response interface{}

	switch dType.Type {
	case 1:
		req := &DisRequest{}
		err = ReadRequest(req, rw)
		if err == nil {
			err = Type1(req, rw)
		}

	case 2:
		req := &DisRequest{}
		err = ReadRequest(req, rw)
		if err == nil {
			response, err = Type2(req)
		}

	case 4:
		req := &ConfirmRequest{}
		err = ReadRequest(req, rw)
		if err == nil {
			response, err = Type4(req)
		}

	case 7:
		req := &GetBodyRequest{}
		err = ReadRequest(req, rw)
		if err == nil {
			response, err = Type7(req)
		}

	case 10:
		response, err = Type10()
	}

	if err != nil {
		logger.LogError(consts.DBError, fmt.Sprintf("tcpserver: parse request error: %s", err))
		return
	}
	if response == nil {
		return
	}

	logger.LogDebug(consts.DebugMessage, fmt.Sprintf("tcpserver response: %+v", response))
	err = SendRequest(response, rw)
	if err != nil {
		logger.LogError(consts.IOError, fmt.Sprintf("tcpserver handle error: %s", err))
	}
}
