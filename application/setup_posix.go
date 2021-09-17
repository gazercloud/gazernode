package application

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func setupPosix() {
	log.Println(" *** Gazer setup procedure ***")
	log.Println("stopping service")
	StopService()
	log.Println("stopping service - complete")

	log.Println("copying to /usr/local/bin")
	sourceFile := os.Args[0]
	destinationFile := "/usr/local/bin/gazernode"
	input, err := ioutil.ReadFile(sourceFile)
	if err != nil {
		log.Println("Error (read):", err)
		return
	}

	err = ioutil.WriteFile(destinationFile, input, 0755)
	if err != nil {
		log.Println("Error (write):", err)
		return
	}
	log.Println("copying to /usr/local/bin - complete")
	log.Println("installing service ...")
	cmdInstall := exec.Command(destinationFile, "-install")
	err = cmdInstall.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("installing service ... complete")
	log.Println("starting service ...")
	cmdStart := exec.Command(destinationFile, "-start")
	err = cmdStart.Run()
	if err != nil {
		log.Println(err)
	}
	log.Println("starting service ... complete")
	log.Println(" *** Gazer setup procedure: complete ***")
}
