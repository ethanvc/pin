package plog

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"github.com/ethanvc/pin/pin/status"
	"github.com/ethanvc/pin/pin/status/codes"
	"io"
	"net"
	"os"
	"strconv"
	"sync/atomic"
	"time"
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
	if c == nil {
		c = context.Background()
	}
	return context.WithValue(c, contextKeyBasicLoggerContext{}, lc)
}

// GenerateTraceId mac address/ip/timestamp/pid/global counter/random bytes
func GenerateTraceId() string {
	now := time.Now().UnixMicro()
	seq := atomic.AddInt64(&DefaultTraceIdGeneratorInfo.LogIndex, 1)
	var buf bytes.Buffer
	buf.WriteString(strconv.FormatInt(now, 10))
	buf.WriteString(strconv.FormatInt(seq, 10))
	buf.Write(DefaultTraceIdGeneratorInfo.Hash[:])
	return fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
}

type TraceIdGeneratorInfo struct {
	MacAddress string
	Ip         string
	Pid        int
	Hash       [5]byte
	LogIndex   int64
}

var DefaultTraceIdGeneratorInfo TraceIdGeneratorInfo

func init() {
	initTraceIdGenerate()
}

func initTraceIdGenerate() {
	DefaultTraceIdGeneratorInfo.Ip, DefaultTraceIdGeneratorInfo.MacAddress = guessBestIpAndMacAddress()
	DefaultTraceIdGeneratorInfo.Pid = os.Getpid()
	h := md5.New()
	io.WriteString(h, DefaultTraceIdGeneratorInfo.MacAddress)
	io.WriteString(h, DefaultTraceIdGeneratorInfo.Ip)
	io.WriteString(h, strconv.FormatInt(int64(DefaultTraceIdGeneratorInfo.Pid), 10))
	copy(DefaultTraceIdGeneratorInfo.Hash[:], h.Sum(nil))

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
