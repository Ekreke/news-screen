package main

import (
	"flag"
	"fmt"
	"os"
	"path"

	"news-screen/internal/conf"
	"news-screen/internal/logger"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/env"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	curEnv string
	id, _  = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../configs", "config path, eg: -conf config.yaml")
	flag.StringVar(&curEnv, "env", "", "environment, eg: -env dev or -env local")
}

func newApp(logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func main() {
	flag.Parse()
	fmt.Println("flagconf", flagconf)
	var sources []config.Source
	baseConfigPath := path.Join(flagconf, "config.yaml")
	sources = append(sources, env.NewSource())
	sources = append(sources, file.NewSource(baseConfigPath))
	log.Info("Current Env: ", curEnv)
	if curEnv != "" {
		envConfigPath := fmt.Sprintf("%s/config-%s.yaml", flagconf, curEnv)
		sources = append(sources, file.NewSource(envConfigPath))
	}
	c := config.New(
		config.WithSource(
			sources...,
		),
	)
	defer c.Close()

	if err := c.Load(); err != nil {
		panic(err)
	}

	var bc conf.Bootstrap
	if err := c.Scan(&bc); err != nil {
		panic(err)
	}

	app, cleanup, err := wireApp(bc.Server, bc.Data, logger.GetZapLogger(bc.Logging))
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}
