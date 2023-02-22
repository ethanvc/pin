package pin

type Server struct {
}

func (this *Server) ProcessRequest(protocolReq ProtocolRequest) {
	req := NewRequest(protocolReq)
	protocolReq.InitializeRequest(req)
	req.Status = req.Next()
	protocolReq.FinalizeRequest(req)
}
