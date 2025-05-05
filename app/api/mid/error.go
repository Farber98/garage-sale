package mid

import (
	"context"

	"github.com/Farber98/garage-sale/app/api/errs"
	"github.com/Farber98/garage-sale/foundation/logger"
)

func Error(ctx context.Context, log *logger.Logger, handler Handler) error {
	err := handler(ctx)
	if err == nil {
		return nil
	}

	log.Error(ctx, "message", "ERROR", err.Error())

	// Check if the error is already one of our custom Error types and return it directly.
	if errs.IsError(err) {
		return errs.GetError(err)
	}

	// If it's not a known Error type, wrap it as an Unknown error.
	// Pass the result of errs.Unknown.String() as the message
	return errs.Newf(errs.Unknown, errs.Unknown.String())
}
