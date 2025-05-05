package mid

import (
	"context"
	"net/http"

	"github.com/Farber98/garage-sale/app/api/errs"
	"github.com/Farber98/garage-sale/app/api/mid"
	"github.com/Farber98/garage-sale/foundation/logger"
	"github.com/Farber98/garage-sale/foundation/web"
)

var codeStatus [19]int

func init() {
	codeStatus[errs.OK.Value()] = http.StatusOK
	codeStatus[errs.Canceled.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.Unknown.Value()] = http.StatusInternalServerError
	codeStatus[errs.InvalidArgument.Value()] = http.StatusBadRequest
	codeStatus[errs.DeadlineExceeded.Value()] = http.StatusGatewayTimeout
	codeStatus[errs.NotFound.Value()] = http.StatusNotFound
	codeStatus[errs.AlreadyExists.Value()] = http.StatusConflict
	codeStatus[errs.PermissionDenied.Value()] = http.StatusForbidden
	codeStatus[errs.ResourceExhausted.Value()] = http.StatusTooManyRequests
	codeStatus[errs.FailedPrecondition.Value()] = http.StatusBadRequest
	codeStatus[errs.Aborted.Value()] = http.StatusConflict
	codeStatus[errs.OutOfRange.Value()] = http.StatusBadRequest
	codeStatus[errs.Unimplemented.Value()] = http.StatusNotImplemented
	codeStatus[errs.Internal.Value()] = http.StatusInternalServerError
	codeStatus[errs.DataLoss.Value()] = http.StatusInternalServerError
	codeStatus[errs.Unauthenticated.Value()] = http.StatusUnauthorized
}

// Error executes the error middleware functionality.
func Error(log *logger.Logger) web.MidHandler {
	m := func(handler web.HandlerFunc) web.HandlerFunc {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			hdl := func(ctx context.Context) error {
				return handler(ctx, w, r)
			}

			if err := mid.Error(ctx, log, hdl); err != nil {
				// Safely check if the returned error is an errs.Error
				errs := err.(errs.Error)
				if err := web.Respond(ctx, w, errs, codeStatus[errs.Code.Value()]); err != nil {
					// If we receive the shutdown err we need to return it back to base handler and shut down service.
					if web.IsShutdown(err) {
						return err
					}
				}
			}

			return nil
		}
		return h
	}
	return m
}
