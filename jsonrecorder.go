
/*
	Package to encapsulate the writing (and reading) (and comparison)
	of JSON records.

	Step 1: get the structure right.
	Step 2: add the JSON logic
	Step 3: profit!

	// TODO(rjkroege): this code needs to be re-organized.
	// I need to have common library code pulled out.
	// And the interceptor needs its own directory / main.
*/
package main

import (
	"fmt"
	"log"
//	"code.google.com/p/goplan9/draw"
//	"image"
//	"syscall"
	"code.google.com/p/goplan9/draw/drawfcall"
	"os"
//	"strings"
//	"sync"
	"io"
)

type JsonRecorder struct {
	c chan string
	complete chan int
}


/*
	I need a type that corresponds to the JSON record.
	Maybe I want to try to separate this change apart
	in some way.
*/
func NewJsonRecorder() *JsonRecorder {
	c := make(chan string, 4)
	complete := make(chan int)
	json := &JsonRecorder{c, complete}
	go json.continuouslyWriteJson()
	return json
}

func (jr *JsonRecorder) WaitToComplete() {
	close(jr.c)
	<-jr.complete
}

/*
	Copies the given devdraw protocol message (on thread) and
	ships the copy into the log channel. Copy logic: I don't need
	to know if devdraw.RPC is mutating the message in some way.

	Also, I want to add additional content prior to encoding.

	TODO(rjkroege): do the JSON stuff.
*/
func (jlog *JsonRecorder) Record(msg *drawfcall.Msg, tag byte) {
	s := msg.String()
	s += ", "
	// How do I add a string thing of an int?
	s +=  fmt.Sprintf("tag: %d", tag)
	jlog.c <- s
}

/*
	Write a message to complete once all messages have been
	encoded and written.
*/
func (json *JsonRecorder) continuouslyWriteJson() {
	filename := os.Getenv("DEVDRAW_LISTENER_OUT")
	if filename == "" {
		filename = "/tmp/devdraw_listener_out";
	}
	fd , err := os.Create(filename)
	if err != nil {
		log.Fatal("openning record ", err)
	}
	// TODO(rjkroege): Need to sort this out once I have the flow working.
	// enc := json.NewEncoder(fd);

	separator := ""

	io.WriteString(fd, "[\n")
	for r := range json.c {
		io.WriteString(fd, separator)
		// need to make this better...
		// err := enc.Encode(r); 
		io.WriteString(fd, r)
		if err != nil {
			log.Fatal("couldn't write the JSON record\n")
		}
		separator = ",\n"
	}
	io.WriteString(fd, "\n]\n")

	// TODO(rjkroege): use defer for this.
	err = fd.Close()
	if err != nil {
		log.Fatal("couldn't close file\n")
	}	

	json.complete <- 1
}