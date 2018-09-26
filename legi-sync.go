package main

import (
	"fmt"
	"github.com/jlaffaye/ftp"
	"io"
	"os"
	"strings"
)

var (
	host = "echanges.dila.gouv.fr"
	port = 21
	login = "anonymous"
	password = "anonymous"
	directory = "LEGI"
	prefix = "LEGI/LEGI_"
	suffix = ".tar.gz"
)

/**
  * Sync FTP directory in current local directory
  * returns 0 if new files are downloaded
  */
func main() {
	found := 1
	defer func() {
		os.Exit(found)
	}()
	client, err := ftp.Dial(fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}
	
	if err := client.Login(login, password); err != nil {
		panic(err)		
	}

	entries, err := client.NameList(directory)
	if err != nil {
		panic(err)
	}

	for _, entry := range entries {
		if strings.HasPrefix(entry, prefix) && 
			strings.HasSuffix(entry, suffix) {
			if _, err := os.Stat(entry); os.IsNotExist(err) {
				found = 0
				fmt.Println(entry)
				Copy(*client, entry)
			}
		}
	}
}

func Copy(client ftp.ServerConn, entry string) {
	response, err := client.Retr(entry)
	if err != nil {
		panic(err)
	}
	defer response.Close()			

	fo, err := os.Create(entry)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	buf := make ([]byte, 1024)
	for {
		n, err := response.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}
		if _, err := fo.Write(buf[:n]); err != nil {
			panic(err)
		}
	}
	
}
