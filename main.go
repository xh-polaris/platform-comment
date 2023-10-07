package main

import (
	"net"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/kitex-contrib/obs-opentelemetry/tracing"
	"github.com/xh-polaris/gopkg/kitex/middleware"
	logx "github.com/xh-polaris/gopkg/util/log"
	"github.com/xh-polaris/service-idl-gen-go/kitex_gen/platform/comment/commentservice"

	"github.com/xh-polaris/platform-comment/biz/infrastructure/util/log"
	"github.com/xh-polaris/platform-comment/provider"
)

func main() {
	klog.SetLogger(logx.NewKlogLogger())
	s, err := provider.NewCommentServerImpl()
	if err != nil {
		panic(err)
	}
	addr, err := net.ResolveTCPAddr("tcp", s.ListenOn)
	if err != nil {
		panic(err)
	}
	svr := commentservice.NewServer(
		s,
		server.WithServiceAddr(addr),
		server.WithSuite(tracing.NewServerSuite()),
		server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{ServiceName: s.Name}),
		server.WithMiddleware(middleware.LogMiddleware(s.Name)),
	)

	err = svr.Run()

	if err != nil {
		log.Error(err.Error())
	}
}
