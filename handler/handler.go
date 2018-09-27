package handler

import (
	"github.com/canary-health/kubectl-apply/config"
)

// Handler struct returned by NewHandler
type Handler struct {
	Aws *awsSessHandler
}

// NewHandler is a reusable handler
func NewHandler(c *config.Config) *Handler {
	return &Handler{
		Aws: newAwsSessionHandler(),
	}
}

func setTrue() *bool {
	t := true
	return &t
}

func ptrInt32(i int) *int32 {
	i32 := int32(i)
	return &i32
}

func ptrInt64(i int) *int64 {
	i64 := int64(i)
	return &i64
}
