package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/go-echarts/go-echarts/v2/render"
)

//go test -bench=BenchmarkUkvs -run=^\$ | tee ukvs_bench.log

func main() {

	switch {
	case len(os.Args) != 2 || os.Args[1] == "":
		log.Fatal("Filename must be specified")
	}

	var f, err = os.Open(os.Args[1])
	if err != nil {
		log.Fatal("Failed to open file", err)
	}

	var r Report
	if err = r.Fill(f); err != nil {
		log.Fatal("Failed to parse benchmark result", err)
	}
	if err = r.Aggregate(); err != nil {
		log.Fatal("Failed to aggregate benchmark result", err)
	}

	var vw render.Renderer
	if vw, err = Visualize(&r); err != nil {
		log.Fatal("Failed to generate visualization", err)
	}

	var l net.Listener
	if l, err = net.Listen("tcp", ":8111"); err != nil {
		log.Fatal("Failed to start listening TCP 8111 port", err)
	}

	log.Println("TCP port 8111 is listened")

	var m = http.NewServeMux()
	m.Handle("/", serverHandler(vw))

	if err = http.Serve(l, m); err != nil {
		log.Fatal("Failed to serve WEB server", err)
	}
}

func serverHandler(vw render.Renderer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		if err := vw.Render(w); err != nil {
			log.Fatal("Failed to render", err)
		}
	})
}
