package middleware

import (
	"net/http"

	"github.com/elazarl/goproxy"
	log "github.com/sirupsen/logrus"

	"github.com/9seconds/crawlera-headless-proxy/config"
)

type proxyRequestMiddleware struct {
	UniqBase
}

func (p *proxyRequestMiddleware) OnRequest() ReqType {
	return p.BaseOnRequest(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		rstate := GetRequestState(ctx)

		log.WithFields(log.Fields{
			"request-id":     rstate.ID,
			"client-id":      rstate.ClientID,
			"method":         req.Method,
			"url":            req.URL,
			"proto":          req.Proto,
			"content-length": req.ContentLength,
			"remote-addr":    req.RemoteAddr,
			"headers":        req.Header,
		}).Debug("HTTP request to sent to Crawlera")

		rstate.StartCrawleraRequest()

		return req, nil
	})
}

func (p *proxyRequestMiddleware) OnResponse() RespType {
	return p.BaseOnResponse(func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		GetRequestState(ctx).FinishCrawleraRequest()

		return resp
	})
}

func NewProxyRequestMiddleware(conf *config.Config, proxy *goproxy.ProxyHttpServer) *proxyRequestMiddleware {
	ware := &proxyRequestMiddleware{}
	ware.conf = conf
	ware.proxy = proxy
	ware.mtype = middlewareTypeProxyRequest

	return ware
}