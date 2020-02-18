package main

import (
	"bufio"
	"fmt"
	"github.com/bedoron/ygui/treeBuilder"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
	"io/ioutil"
	"log"
	"os"
)

func readPipe() ([]byte, error){
	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	if info.Mode()&os.ModeCharDevice != 0 && info.Mode()&os.ModeNamedPipe != 0 {
		fmt.Println("The command is intended to work with pipes.")
		fmt.Println("Usage: echo \"some yaml...\" | ygui")
		return nil, fmt.Errorf("no pipe")
	}

	reader := bufio.NewReader(os.Stdin)

	return ioutil.ReadAll(reader)
}

func main() {
	var fileContent []byte
	var err error
	if len(os.Args) == 2 {
		fileName := os.Args[1]
		fileContent, err = ioutil.ReadFile(fileName)
	} else {
		fileContent, err = readPipe()
	}

	if err != nil {
		fmt.Println("error ", err)
		return
	}

	fmt.Println(string(fileContent))

	builder, err := treeBuilder.NewBuilder(fileContent)
	if err != nil {
		fmt.Println("Failed building ", err)
		return
	}

	fmt.Println(builder)
	nodes := builder.Nodes()

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}
	defer ui.Close()

	l := widgets.NewTree()
	l.TextStyle = ui.NewStyle(ui.ColorYellow)
	l.WrapText = false
	l.SetNodes(nodes)

	x, y := ui.TerminalDimensions()

	l.SetRect(0, 0, x, y)

	ui.Render(l)

	previousKey := ""
	uiEvents := ui.PollEvents()
	for {
		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		case "<Home>":
			l.ScrollTop()
		case "<Enter>":
			l.ToggleExpand()
		case "G", "<End>":
			l.ScrollBottom()
		case "E":
			l.ExpandAll()
		case "C":
			l.CollapseAll()
		case "<Resize>":
			x, y := ui.TerminalDimensions()
			l.SetRect(0, 0, x, y)
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
		}

		ui.Render(l)
	}
}
