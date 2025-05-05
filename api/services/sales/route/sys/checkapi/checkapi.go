package checkapi

import (
	"context"
	"math/rand/v2"
	"net/http"

	"github.com/Farber98/garage-sale/app/api/errs"
	"github.com/Farber98/garage-sale/foundation/web"
)

func liveness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}

func readiness(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}

func testerror(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.IntN(100); n%2 == 0 {
		return errs.Newf(errs.FailedPrecondition, "this message is trusted")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}

func testpanic(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.IntN(100); n%2 == 0 {
		panic("WE ARE PANICKING!!!!")
	}

	status := struct {
		Status string
	}{
		Status: "OK",
	}

	return web.Respond(ctx, w, status, http.StatusOK)
}
