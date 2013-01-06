package goaspen

import (
	"fmt"
	"net/http"
	"strings"
)

const (
	internalAcceptHeader = "X-GoAspen-Accept"
	pathTransHeader      = "X-HTTP-Path-Translated"
)

type handlerFuncRegistration struct {
	RequestPath    string
	HandlerFunc    http.HandlerFunc
	Negotiated     bool
	Virtual        bool
	Regexp         bool
	RegWithNetHTTP bool

	w *Website
}

func UpdateContextFromVirtualPaths(ctx *map[string]interface{},
	requestPath, vPath string) {

	realCtx := *ctx

	rpParts := strings.Split(requestPath, "/")
	vpParts := strings.Split(vPath, "/")

	if len(rpParts) != len(vpParts) {
		debugf("Request and virtual paths have different "+
			"part counts, so not updating request context: %q, %q",
			requestPath, vPath)
		return
	}

	for i, vPart := range vpParts {
		if len(vPart) < 1 {
			continue
		}

		if vPart[0] == '%' {
			realCtx[vPart[1:]] = rpParts[i]
		}
	}
}

func serve404(w http.ResponseWriter, req *http.Request) {
	charset := req.Header.Get("X-GoAspen-CharsetDynamic")
	if len(charset) == 0 {
		charset = "utf-8"
	}

	w.Header().Set("Content-Type", fmt.Sprintf("text/html; charset=%v", charset))
	w.WriteHeader(http.StatusNotFound)
	w.Write(http404Response)
}
