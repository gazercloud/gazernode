package application

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

func setupPosix() {
	log.Println("stopping service")
	StopService()
	log.Println("stopping service - complete")

	log.Println("copying to /usr/local/bin")
	sourceFile := filepath.Dir(os.Args[0])
	destinationFile := "/usr/local/bin/gazer_node"
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Println("Error (read):", err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0644)
	if err != nil {
		log.Println("Error (write):", err)
		return
	}
	log.Println("copying to /usr/local/bin - complete")
	log.Println("installing service ...")
	exec.Command(destinationFile, "-install")
	log.Println("installing service ... complete")
	log.Println("starting service ...")
	exec.Command(destinationFile, "-start")
	log.Println("starting service ... complete")
}
