package main

import (
	"log"
	"tpi-mon/api"
	"tpi-mon/tpi"
)

func main() {
	c, err := tpi.NewLocalClient("127.0.0.1", 9751, "aBcDe1")
	if err != nil {
		log.Panicln(err)
	}

	api.Run(c, 9750)
}

// func main() {

// 	c, err := tpi.NewLocalClient("127.0.0.1", 9751, "aBcDe1")
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	// g, err := gocui.NewGui(gocui.OutputNormal)
// 	// if err != nil {
// 	// 	log.Panicln(err)
// 	// }
// 	// defer g.Close()
// 	// g.SetManagerFunc(layoutTerminal)

// 	// if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
// 	// 	log.Panicln(err)
// 	// }

// 	// if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
// 	// 	log.Panicln(err)
// 	// }

// 	f, err := os.OpenFile("tpi-logs.txt", os.O_RDWR|os.O_APPEND, 0666)
// 	if err != nil {
// 		log.Panicln(err)
// 	}
// 	fmt.Fprintln(f, "-------")

// 	readCommands(c)

// 	for {
// 		select {
// 		case e := <-c.GetEventCh():
// 			fmt.Fprintln(f, e)
// 		case err := <-c.GetErrorCh():
// 			log.Panicln(err)
// 		}
// 	}
// }

// // func layoutTerminal(g *gocui.Gui) error {
// // 	maxX, maxY := g.Size()
// // 	// if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
// // 	// 	if err != gocui.ErrUnknownView {
// // 	// 		return err
// // 	// 	}
// // 	// 	fmt.Fprintln(v, "Hello world!!!!")
// // 	// }

// // 	var err error
// // 	var logsView *gocui.View
// // 	var cmdView *gocui.View

// // 	cmdViewStartY := maxY - 4

// // 	if logsView, err = g.SetView("logs", 0, 0, maxX-1, cmdViewStartY); err != nil {
// // 		if err != gocui.ErrUnknownView {
// // 			return err
// // 		}
// // 	}
// // 	fmt.Fprintln(logsView, "Hello Logsss!!!!", time.Now())

// // 	if cmdView, err = g.SetView("cmd", 0, cmdViewStartY, maxX-1, maxY-1); err != nil {
// // 		if err != gocui.ErrUnknownView {
// // 			return err
// // 		}
// // 	}
// // 	fmt.Fprintln(cmdView, "Hello Cmds!!!!", time.Now())

// // 	return nil
// // }

// // func quit(g *gocui.Gui, v *gocui.View) error {
// // 	return gocui.ErrQuit
// // }

// func readCommands(c tpi.Client) {
// 	fmt.Println("Press CTRL-D to exit.")

// 	rl, err := readline.New("> ")
// 	if err != nil {
// 		log.Panicln(err)
// 	}

// 	go func() {
// 		for {
// 			line, err := rl.Readline()
// 			if err != nil { // io.EOF
// 				os.Exit(0)
// 			}
// 			execCommand(c, rl, line)
// 		}
// 	}()
// }

// func execCommand(c tpi.Client, rl *readline.Instance, cmdLine string) {

// 	split := strings.Split(cmdLine, " ")
// 	verb := split[0]
// 	args := split[1:]

// 	if len(verb) == 0 {
// 		return
// 	}

// 	switch verb {
// 	case "AwayArm":
// 		c.AwayArm(args[0])
// 	case "StayArm":
// 		c.StayArm(args[0])
// 	case "ZeroEntryDelayArm":
// 		c.ZeroEntryDelayArm(args[0])
// 	case "ArmWithCode":
// 		c.ArmWithCode(args[0], args[1])
// 	case "Disarm":
// 		pin, err := rl.ReadPassword("PIN:")
// 		if err != nil {
// 			panic(err)
// 		}
// 		c.Disarm(args[0], string(pin))
// 	case "Panic":
// 		c.PanicAlarm(args[0])
// 	default:
// 		fmt.Printf("Error: invalid command %s\n", verb)
// 	}
// }

// func logEvents(c *tpi.Client) {
// }
