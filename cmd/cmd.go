package cmd

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/gazercloud/gazernode/gazer_client"
	"github.com/gazercloud/gazernode/utilities/paths"
	"io/ioutil"
	"os"
	"strings"
)

type Session struct {
	currentUnitId string
	//currentUnitName string
	currentItem string
	client      *gazer_client.GazerNodeClient
}

type SessionSettings struct {
	Address    string
	SessionKey string
}

func (c *Session) load() {
	homeFolder := paths.HomeFolder()
	bs, err := ioutil.ReadFile(homeFolder + "/gazer_termo.json")
	if err != nil {
		return
	}
	var settings SessionSettings
	err = json.Unmarshal(bs, &settings)
	if err != nil {
		return
	}
	if settings.Address != "" {
		c.client = gazer_client.NewWithSession(settings.Address, settings.SessionKey)
		fmt.Println("Session loaded:", settings.Address)
	}
}

func (c *Session) save() {
	homeFolder := paths.HomeFolder()
	var settings SessionSettings
	settings.Address = c.client.Address()
	settings.SessionKey = c.client.SessionToken()
	bs, err := json.Marshal(settings)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(homeFolder+"/gazer_termo.json", bs, 0660)
	if err != nil {
		return
	}
}

func (c *Session) currentPath() string {
	result := "/"
	if c.currentUnitId != "" {
		result = "/" + c.currentUnitId
		if c.currentItem != "" {
			result = "/" + c.currentItem
		}
	}
	return result
}

func (c *Session) currentPathIsUnit() bool {
	if c.currentUnitId != "" && c.currentItem == "" {
		return true
	}
	return false
}

func (c *Session) currentPathIsItem() bool {
	if c.currentUnitId != "" && c.currentItem != "" {
		return true
	}
	return false
}

func (c *Session) execLine(session *Session, line string) bool {
	var err error

	line = strings.Trim(line, "\r\n\t")
	parts := strings.FieldsFunc(line, func(r rune) bool {
		return r == ' ' || r == '\r' || r == '\n'
	})

	if len(parts) < 1 {
		return false
	}

	if parts[0] == "quit" || parts[0] == "exit" {
		return true
	}

	params := parts[1:]

	switch parts[0] {
	case "quit":
		return true
	case "exit":
		return true
	case "help":
		err = c.cmdHelp(params)
	case "connect":
		err = c.cmdConnect(params)
	case "disconnect":
		err = c.cmdDisconnect(params)
	case "l":
		fallthrough
	case "ls":
		err = c.cmdLs(params)
	case "cd":
		err = c.cmdCd(params)
	case "item":
		err = c.cmdItem(params)
	case "unit":
		err = c.cmdUnit(params)
	case "cloud":
		err = c.cmdCloud(params)
	default:
		err = errors.New("wong command")
	}

	if err != nil {
		color.Set(color.FgRed)
		fmt.Println("ERROR:", err)
		color.Unset()
	}

	return false
}

func Console() {
	var err error

	session := &Session{}
	session.load()

	commandLine := ""
	in := bufio.NewReader(os.Stdin)
	for true {
		prompt := ">"
		if session.client != nil {
			prompt = session.client.Address() + session.currentPath() + ">"
		}
		fmt.Print(prompt)
		commandLine, err = in.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}
		if session.execLine(session, commandLine) {
			break
		}
	}
	fmt.Println("exit")
}
