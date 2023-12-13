package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	screenWidth  = 310
	screenHeight = 575
	musicPath    = "./music"
)

const (
	MUSIC_PANEL int = iota
	SETTINGS_PANEL
)

var (
	musicLibrary      []Music
	musicCurrent      *Music
	musicCurrentIndex int     = -1
	musicVolume       float32 = 50
	musicSeek         float32 = 0
	musicPlaying              = false
	playPauseIcon             = gui.IconText(gui.ICON_PLAYER_PLAY, "")

	panelRec        = rl.NewRectangle(5, 30, 300, 462)
	panelContentRec = rl.NewRectangle(0, 0, 300-14, 0)
	panelScroll     = rl.NewVector2(0, 0)

	statusBarText string
	titleBarText  string

	progressBarMaxValue float32 = 1
	progressBarValue    float32 = 0
	progressBarMinValue float32 = 0

	currentPanel    = 0
	showQuitPopup   = false
	exitApplication = false

	err error
)
