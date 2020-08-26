package gateway

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/gostarer/gateway/dao"
	"github.com/gostarer/gateway/router"
	proxy_router "github.com/gostarer/gateway/router/proxy_router"

	"github.com/gostarer/gateway/infra/lib"
)

var (
	endpoint = flag.String("endpoint", "", "input endpoint dashboard or server")
	config   = flag.String("config", "", "input config file")
)

func main() {
	flag.Parse()
	if *endpoint == "" || *config == "" {
		flag.Usage()
		os.Exit(1)
	}
	if *endpoint == "dsshboard" {
		lib.InitModule(*config)
		defer lib.Destroy()
		router.HttpServerRun()
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		router.HttpServerStop()
	} else {
		lib.InitModule(*config)
		defer lib.Destroy()

		dao.ServiceManagerHandler.LoadOnce()
		dao.AppManagerHandler.LoadOnce()

		go func() {
			proxy_router.HttpServerRun()
		}()
		go func() {
			proxy_router.HttpsServerRun()
		}()
		go func() {
			proxy_router.TcpServerRun()
		}()

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		proxy_router.TcpServerStop()
		proxy_router.HttpServerStop()
		proxy_router.HttpsServerStop()
	}
}
