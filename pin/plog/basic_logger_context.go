package plog

import (
	"context"
	"github.com/ethanvc/pin/pin/status"
	"github.com/ethanvc/pin/pin/status/codes"
	"net"
	"os"
)

type BasicLoggerContext struct {
	TraceId string
}

type contextKeyBasicLoggerContext struct{}

func BasicLoggerContextFromCtx(c context.Context) BasicLoggerContext {
	lc, _ := c.Value(contextKeyBasicLoggerContext{}).(BasicLoggerContext)
	return lc
}

func WithBasicLoggerContext(c context.Context, lc BasicLoggerContext) context.Context {
	return context.WithValue(c, contextKeyBasicLoggerContext{}, lc)
}

// GenerateTraceId mac address/ip/timestamp/pid/global counter/random bytes
func GenerateTraceId() string {
	return ""
}

type TraceIdGeneratorInfo struct {
	MacAddress string
	Ip         string
	Pid        int
}

var DefaultTraceIdGeneratorInfo TraceIdGeneratorInfo

func init() {
	initTraceIdGenerate()
}

func initTraceIdGenerate() {
	DefaultTraceIdGeneratorInfo.Ip, DefaultTraceIdGeneratorInfo.MacAddress = guessBestIpAndMacAddress()
	DefaultTraceIdGeneratorInfo.Pid = os.Getpid()
}

func guessBestIpAndMacAddress() (string, string) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		panic(status.NewStatus(codes.Internal, "GetIpAddrFailed"))
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	ifas, err := net.Interfaces()
	if err != nil {
		panic(status.NewStatus(codes.Internal, "GetIpAddrFailed"))
	}
	for _, ifa := range ifas {
		ifaddrs, err := ifa.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range ifaddrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			default:
				continue
			}
			if !ip.IsGlobalUnicast() {
				continue
			}
			if ip.String() != localAddr.IP.String() {
				continue
			}
			return localAddr.IP.String(), ifa.HardwareAddr.String()
		}
	}
	return localAddr.String(), ""
}
