package pin

import (
	"github.com/ethanvc/pin/pin/plog"
	"time"
)

type Server struct {
}

func (this *Server) ProcessRequest(protocolReq ProtocolRequest) {
	req := NewRequest(protocolReq)
	protocolReq.InitializeRequest(req)
	req.Status = req.Next()
	protocolReq.FinalizeRequest(req)
	this.logAccessRequest(req)
}

// go test -bench . -benchmem
func (this *Server) logAccessRequest(req *Request) {
	timeMs := time.Now().Sub(req.StartTime).Milliseconds()
	lvl := plog.LevelInfo
	if req.Status.NotOk() {
		lvl = plog.LevelWarn
	}
	req.Logger.Log(lvl, "pin_acc").Str("path", req.PatternPath).Int64("t_ms", timeMs).
		Any("status", req.Status).Any("req", req.Req).Any("resp", req.Resp).Done()
}
