package istio

import "net/http"

// GetHttpForwardHeader 获取Http跟踪Header
func GetHttpForwardHeader(req *http.Request) http.Header {
	header := http.Header{}
	incomingHeaders := []string{
		"x-request-id",
		"x-b3-traceid",
		"x-b3-spanid",
		"x-b3-parentspanid",
		"x-b3-sampled",
		"x-b3-flags",
		"x-ot-span-context",
	}
	for _, key := range incomingHeaders {
		header.Set(key, req.Header.Get(key))
	}
	return header
}
