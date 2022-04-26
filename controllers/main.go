package controllers

import (
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/restapi"
)

var listener *restapi.HttpListener
var logger = log.Get()

func Init() {
	listener = restapi.GetHttpListener()
	listener.AddJsonContent().AddLogger().AddHealthCheck()

	listener.AddController(ValidateHeadersController, "/security/headers", "POST")
	listener.Start()
}
