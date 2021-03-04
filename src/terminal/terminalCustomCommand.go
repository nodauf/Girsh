package terminal

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"

	utils "nc-shell/src/utils"
)

const fsMaxbufsize = 4096

// Upload a file on linux reverse shell with base64
func (terminal *Terminal) Upload(cmd string) {
	// Clean the command upload ....
	terminal.cleanCmd(cmd)
	filePath := strings.Split(cmd, " ")[1]
	printTerm(filePath)
	f, err := os.Open(filePath)

	if err != nil {
		log.Println("File " + filePath + " does not exist")
		return
	}

	defer f.Close()
	statinfo, err := f.Stat()
	if err != nil {
		log.Println("stat() failure for the file: " + filePath)
		return
	}

	filename := statinfo.Name()

	buf := make([]byte, utils.Min(fsMaxbufsize, statinfo.Size()))
	n := 0
	for err == nil {
		n, err = f.Read(buf)
		encodedStr := base64.StdEncoding.EncodeToString(buf[0:n])
		cmd := "echo " + encodedStr + "|base64 -d>> " + filename
		printTerm(string(buf))
		terminal.execute(cmd)
	}

}

// Download future function to download a file
func (terminal *Terminal) Download(cmd string) {

	fmt.Print("\033[999D")
	fmt.Println("\nDownload the file " + cmd)
}

//print the text a the beginning of the line
func printTerm(str string) {
	fmt.Print("\033[999D")
	fmt.Println("\n" + str)
}
