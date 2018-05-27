package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	from int
	to   string
	tcp  = "tcp"
)

func main() {
	flag.IntVar(&from, "from", 1234, "port from")
	flag.StringVar(&to, "to", "localhost:4321", "port to")

	flag.Parse()

	if _, err := strconv.Atoi(to); err == nil {
		to = fmt.Sprintf("localhost:%s", to)
	}

	fromAddress := fmt.Sprintf("localhost:%d", from)
	toAddres := to

	listen, err := net.Listen(tcp, fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	defer listen.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Print("\nPort forward stopping...")
		os.Exit(0)
	}()

	log.Printf("\nPort forward from %d to %s", from, to)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal(err)
		}

		target, err := net.Dial(tcp, toAddres)
		if err != nil {
			log.Fatal(err)
		}

		go copyIO(conn, target)
		go copyIO(target, conn)

		defer fmt.Println("Clean")
	}
}

func copyIO(src, dst net.Conn) {
	defer src.Close()
	defer dst.Close()

	io.Copy(src, dst)
}
