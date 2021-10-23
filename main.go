package main

import (
	"fmt"
	"github.com/go-vgo/robotgo"
	"github.com/lxn/win"
	"math/rand"
	"syscall"
)

var (
	windowsNames = []string{"Hearthstone", "하스스톤", "《爐石戰記》", "炉石传说"}
	elements     = make(map[string]*Template)
	matcherUrl   string
)

type Image struct {
	x, y, width, height int
	bytes               []byte
}

func main() {
	appConfig := fromProperties()

	matcherUrl = appConfig.TemplateMatcherUrl

	for _, t := range appConfig.Templates {
		id := upload(t.Path, matcherUrl)
		if id == "" {
			continue
		}
		t.id = id
		elements[t.Name] = t
	}

	var h *syscall.Handle
	for _, wName := range windowsNames {
		if res, err := FindWindow(wName); err == nil {
			h = &res
			break
		}
	}

	if h == nil {
		_ = fmt.Errorf("no window with title '%s' found", windowsNames)
		return
	}

	window := &Window{hwnd: win.HWND(*h)}

	win.ShowWindow(window.hwnd, 9)
	win.SetActiveWindow(window.hwnd)
	win.SetForegroundWindow(window.hwnd)

	robotgo.MicroSleep(300)

	var rect win.RECT
	win.GetWindowRect(window.hwnd, &rect)

	window.x = int(rect.Left)
	window.y = int(rect.Top)
	window.width = int(rect.Right) - int(rect.Left)
	window.height = int(rect.Bottom) - int(rect.Top)

	fmt.Println(int(rect.Left), int(rect.Top), int(rect.Right), int(rect.Bottom))

	confirm := elements["confirm"]
	choose := elements["choose"]
	turnReady := elements["turnReady"]
	start := elements["start"]
	action := elements["action"]
	opponent := elements["opponent"]
	take := elements["take"]
	unactive := elements["unactive"]
	red := elements["red"]
	green := elements["green"]
	blue := elements["blue"]

	stuckCount := 0
	sleepTime := 2
	for {
		robotgo.Sleep(sleepTime)
		sleepTime = 2

		confirmImage := window.screenElement(confirm)
		if point := detect(confirmImage.bytes, confirm); point != nil {
			robotgo.MoveMouse(confirmImage.x+point.X+20, confirmImage.y+point.Y+20)
			robotgo.MouseClick("left")
			sleepTime = 5
			continue
		}

		chooseImage := window.screenElement(choose)
		if point := detect(chooseImage.bytes, choose); point != nil {
			robotgo.MoveMouse(chooseImage.x+point.X+20, chooseImage.y+point.Y+20)
			robotgo.MouseClick("left")
			sleepTime = 5
			continue
		}

		actionImage := window.screenElement(action)
		if point := detectWithConf(actionImage.bytes, action, 0.70); point != nil {
			robotgo.MoveMouse(actionImage.x+point.X-10, actionImage.y+point.Y-10)
			robotgo.MouseClick("left")

			robotgo.MicroSleep(300)

			opponentImage := window.screenElement(opponent)
			if point := detectWithConf(opponentImage.bytes, opponent, 0.65); point != nil {
				robotgo.MoveMouse(opponentImage.x+point.X-20, opponentImage.y+point.Y+20)
				robotgo.MouseClick("left")
			}
			sleepTime = 1
			continue
		}

		turnReadyImage := window.screenElement(turnReady)
		if point := detectWithConf(turnReadyImage.bytes, turnReady, 0.70); point != nil {
			robotgo.MoveMouse(turnReadyImage.x+point.X+20, turnReadyImage.y+point.Y+20)
			robotgo.MouseClick("left")
			sleepTime = 10
			continue
		}

		startImage := window.screenElement(start)
		if point := detectWithConf(startImage.bytes, start, 70); point != nil {
			robotgo.MoveMouse(startImage.x+point.X+20, startImage.y+point.Y+20)
			robotgo.MouseClick("left")
			sleepTime = 5
			continue
		}

		takeImage := window.screenElement(take)
		if point := detect(takeImage.bytes, take); point != nil {
			robotgo.MoveMouseSmooth(window.x+int(float32(window.width)*0.35), window.y+int(float32(window.height)*0.5), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(300)
			robotgo.MoveMouse(takeImage.x+point.X+20, takeImage.y+point.Y+20)
			robotgo.MouseClick("left")
			sleepTime = 5
			continue
		}

		unactiveImage := window.screenElement(unactive)
		if point := detect(unactiveImage.bytes, unactive); point != nil {
			rn := rand.Intn(3)
			if rn == 0 {
				redImage := window.screenElement(red)
				if point := detect(redImage.bytes, red); point != nil {
					robotgo.MoveMouseSmooth(redImage.x+point.X, redImage.y+point.Y+35, 0.9, 0.9)
					robotgo.MicroSleep(50)
					robotgo.MouseClick("left")
					robotgo.MicroSleep(50)
					robotgo.MouseClick("left")
					continue
				}
			}
			if rn == 1 {
				blueImage := window.screenElement(blue)
				if point := detect(blueImage.bytes, blue); point != nil {
					robotgo.MoveMouseSmooth(blueImage.x+point.X, blueImage.y+point.Y+35, 0.9, 0.9)
					robotgo.MicroSleep(50)
					robotgo.MouseClick("left")
					robotgo.MicroSleep(50)
					robotgo.MouseClick("left")
					continue
				}
			}
			if rn == 2 {
				greenImage := window.screenElement(green)
				if point := detect(greenImage.bytes, green); point != nil {
					robotgo.MoveMouseSmooth(greenImage.x+point.X, greenImage.y+point.Y+35, 0.9, 0.9)
					robotgo.MicroSleep(50)
					robotgo.MouseClick("left")
					robotgo.MicroSleep(50)
					robotgo.MouseClick("left")
					continue
				}
			}
		}

		if stuckCount < 5 {
			r := rand.Intn(150) - 50
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.35), r+window.y+int(float32(window.height)*0.50), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouse(r+window.x+int(float32(window.width)*0.57), r+window.y+int(float32(window.height)*0.5))
			robotgo.MouseClick("left")
			stuckCount++
		} else {
			r := rand.Intn(100) - 50
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.50), r+window.y+int(float32(window.height)*0.20), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.80), r+window.y+int(float32(window.height)*0.33), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.66), r+window.y+int(float32(window.height)*0.75), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.33), r+window.y+int(float32(window.height)*0.33), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.33), r+window.y+int(float32(window.height)*0.75), 0.9, 0.9)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(100)
			robotgo.MoveMouseSmooth(r+window.x+int(float32(window.width)*0.5), r+window.y+int(float32(window.height)*0.80), 0.9, 0.9)
			robotgo.MicroSleep(50)
			robotgo.MouseClick("left")
			robotgo.MicroSleep(50)
			stuckCount = 0
		}
	}
}
