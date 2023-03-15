package errorx

import "google.golang.org/grpc/status"

var (
	ErrNoSuchPost            = status.Error(10301, "no such post")
	ErrInvalidObjectId       = status.Error(10302, "invalid objectId")
	ErrPaginatorTokenExpired = status.Error(10303, "paginator token has been expired")
)
