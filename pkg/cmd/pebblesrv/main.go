package main

import (
	"runtime"

	"github.com/tachunwu/distcomm/pkg/pebblesrv"
)

func main() {
	srv := pebblesrv.NewPebbleServer()
	srv.Init(*addr, *dataDir)
	srv.Serve()
	defer srv.Close()
	for {
		runtime.Gosched()
	}
}
