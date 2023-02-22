package pin

import "github.com/ethanvc/pin/pin/status"

type Handler struct {
	DirectFunc func() *status.Status
}
