package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/jjh2kiss/rawfish/config"
	"github.com/jjh2kiss/rawfish/math"
	"github.com/jjh2kiss/rawfish/rawfishnet"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "rawfish-go"
	app.Usage = "HTTP Server for RAW HTTP DATA"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "addr, a",
			Value: "0.0.0.0",
			Usage: "Bind address",
		},
		cli.StringFlag{
			Name:  "port, p",
			Value: "80",
			Usage: "port number",
		},
		cli.StringFlag{
			Name:  "root, r",
			Usage: "Directory for Service Root",
		},
		cli.IntFlag{
			Name:  "read-timeout",
			Value: 10,
			Usage: "timeout for read",
		},
		cli.IntFlag{
			Name:  "write-timeout",
			Value: 10,
			Usage: "timeout for write",
		},
		cli.BoolFlag{
			Name:  "force-200-ok, f",
			Usage: "reply 200 OK when does not have page",
		},
		cli.IntFlag{
			Name:  "force-200-ok-content-size",
			Value: 0,
			Usage: "content size for force 200 ok",
		},
		cli.BoolFlag{
			Name:  "https",
			Usage: "enable HTTPS",
		},
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "verbose mode",
		},
		cli.StringFlag{
			Name:  "pemfile",
			Usage: "pemfile for HTTPS",
		},
		cli.IntFlag{
			Name:  "process",
			Value: 1,
			Usage: "Count for Processes",
		},
		cli.IntFlag{
			Name:  "rate",
			Value: 0,
			Usage: "send at in bytes/sec(Bps), default is unlimited(0)",
		},
	}

	app.Action = func(c *cli.Context) {
		config := config.Config{}
		config.Addr = c.String("addr")
		config.Port = c.String("port")
		config.ReadTimeout = c.Int("read-timeout")
		config.WriteTimeout = c.Int("write-timeout")
		config.Https = c.Bool("https")
		config.Force200Ok = c.Bool("force-200-ok")
		config.Force200OkSize = c.Int("force-200-ok-content-size")
		config.Pemfile = c.String("pemfile")
		config.Rate = c.Int("rate")

		root := c.String("root")
		root, err := filepath.Abs(root)
		if err != nil {
			log.Fatal("Fail to get abs path for root:", c.String("root"))
		}
		config.Root = root

		process := c.Int("process")
		config.Process = math.IntMin(process, runtime.NumCPU())
		runtime.GOMAXPROCS(process)

		if c.Bool("verbose") == false {
			log.SetOutput(ioutil.Discard)
		}

		mux := http.NewServeMux()
		server := &http.Server{
			Addr:           config.FullAddress(),
			Handler:        mux,
			ReadTimeout:    time.Duration(config.ReadTimeout) * time.Second,
			WriteTimeout:   time.Duration(config.WriteTimeout) * time.Second,
			MaxHeaderBytes: 1 << 20,
		}

		handler := rawfishnet.NewRawfishHandler(&config, "")
		mux.Handle("/", handler)

		log.Printf("About to listen on %s Go to http://%s/",
			config.Port,
			config.FullAddress(),
		)

		if config.Https {
			err := server.ListenAndServeTLS(
				config.Pemfile,
				config.Pemfile,
			)

			if err != nil {
				log.Fatal("ListenAndServeTLS: ", err.Error())
				return
			}
		} else {
			server.ListenAndServe()
		}
	}

	app.Run(os.Args)
}
