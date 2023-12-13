package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	"fmt"
	"math"
	"os"
)

func run() int {
	musicLibrary, err = loadLibrary(musicPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load music library: %v\n", err)
		return 1
	}

	//rl.SetConfigFlags(rl.FlagWindowResizable)
	//rl.SetConfigFlags(rl.FlagWindowHighdpi)

	rl.InitWindow(screenWidth, screenHeight, "music")
	rl.InitAudioDevice()
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	if _, err := os.Stat("./theme"); !os.IsNotExist(err) {
		gui.LoadStyle("./theme/style_cyber.rgs")
	}

	for !exitApplication {
		exitApplication = rl.WindowShouldClose()

		// Take user input
		input()

		// Update logic
		update()

		// Draw
		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))
		draw()
		rl.EndDrawing()

		// Update music stream
		if musicCurrent != nil && musicPlaying {
			rl.SetMusicVolume(musicCurrent.File, musicVolume/100)

			// if end of track, play next track
			// -0.3 as if we do a direct comparison it will never be true, probably due to floating point errors
			if rl.GetMusicTimePlayed(musicCurrent.File) >= rl.GetMusicTimeLength(musicCurrent.File)-0.3 {
				if musicCurrentIndex < len(musicLibrary)-1 {
					musicCurrentIndex += 1
				} else {
					musicCurrentIndex = 0
				}

				loadMusic(musicLibrary[musicCurrentIndex])
			}

			rl.UpdateMusicStream(musicCurrent.File)
		}

		rl.SetWindowTitle(titleBarText)
	}

	if musicCurrent != nil {
		rl.UnloadMusicStream(musicCurrent.File)
	}

	rl.CloseAudioDevice()
	rl.CloseWindow()

	return 0
}

func input() {
	if rl.IsKeyPressed(rl.KeyEscape) {
		showQuitPopup = !showQuitPopup
	}

	if rl.IsKeyPressed(rl.KeySpace) {
		if musicPlaying {
			rl.PauseMusicStream(musicCurrent.File)
			musicPlaying = false
		} else {
			rl.ResumeMusicStream(musicCurrent.File)
			musicPlaying = true
		}
	}

	if rl.IsKeyPressed(rl.KeyLeft) {
		if rl.GetMusicTimePlayed(musicCurrent.File)-10 > 0 {
			rl.SeekMusicStream(musicCurrent.File, float32(math.Max(float64(rl.GetMusicTimePlayed(musicCurrent.File)-10), 0)))
		}
	}
	if rl.IsKeyPressed(rl.KeyRight) {
		if rl.GetMusicTimePlayed(musicCurrent.File)+10 < rl.GetMusicTimeLength(musicCurrent.File) {
			rl.SeekMusicStream(musicCurrent.File, float32(math.Min(float64(rl.GetMusicTimePlayed(musicCurrent.File)+10), float64(rl.GetMusicTimeLength(musicCurrent.File)))))
		}
	}

	if rl.IsKeyPressed(rl.KeyUp) {
		if musicCurrentIndex > 0 {
			musicCurrentIndex -= 1
		} else {
			musicCurrentIndex = len(musicLibrary) - 1
		}

		loadMusic(musicLibrary[musicCurrentIndex])
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		if musicCurrentIndex < len(musicLibrary)-1 {
			musicCurrentIndex += 1
		} else {
			musicCurrentIndex = 0
		}

		loadMusic(musicLibrary[musicCurrentIndex])
	}
}

func update() {
	if musicCurrent == nil {
		progressBarValue = 0
	} else {
		progressBarValue = rl.GetMusicTimePlayed(musicCurrent.File)
	}

	if musicCurrent == nil {
		statusBarText = "Time: 0:00/0:00 | Nothing Playing"
	} else {
		minutesPlayed := int(musicSeek / 60)
		secondsPlayed := "00"
		if int(musicSeek-float32(minutesPlayed)*60) < 10 {
			secondsPlayed = fmt.Sprintf("0%d", int(musicSeek-float32(minutesPlayed)*60))
		} else {
			secondsPlayed = fmt.Sprintf("%d", int(musicSeek-float32(minutesPlayed)*60))
		}

		minutesLength := int(rl.GetMusicTimeLength(musicCurrent.File) / 60)
		secondsLength := "00"
		if int(rl.GetMusicTimeLength(musicCurrent.File)-float32(minutesLength)*60) < 10 {
			secondsLength = fmt.Sprintf("0%d", int(rl.GetMusicTimeLength(musicCurrent.File)-float32(minutesLength)*60))
		} else {
			secondsLength = fmt.Sprintf("%d", int(rl.GetMusicTimeLength(musicCurrent.File)-float32(minutesLength)*60))
		}
		statusBarText = fmt.Sprintf("Time: %d:%s/%d:%s | (%d) %s", minutesPlayed, secondsPlayed, minutesLength, secondsLength, musicCurrentIndex, musicCurrent.Name)
	}

	if musicCurrent == nil {
		titleBarText = "Music"
	} else {
		titleBarText = fmt.Sprintf("Music - Playing: %s", musicCurrent.Name)
	}

	if musicPlaying {
		playPauseIcon = gui.IconText(gui.ICON_PLAYER_PAUSE, "")
	} else {
		playPauseIcon = gui.IconText(gui.ICON_PLAYER_PLAY, "")
	}
}

func draw() {
	// Volume bar
	musicVolume = gui.SliderBar(rl.Rectangle{25, 5, 255, 20}, gui.IconText(gui.ICON_AUDIO, ""), "", musicVolume, 0, 100)

	// Settings button
	if currentPanel == MUSIC_PANEL {
		if gui.Button(rl.Rectangle{285, 5, 20, 20}, gui.IconText(gui.ICON_GEAR, "")) {
			currentPanel = 1
		}
	} else if currentPanel == SETTINGS_PANEL {
		if gui.Button(rl.Rectangle{285, 5, 20, 20}, gui.IconText(gui.ICON_BURGER_MENU, "")) {
			currentPanel = 0
		}
	}

	if currentPanel == MUSIC_PANEL {
		// This is horrible, I hate it I hate it I hate it
		panelContentRec.Height = float32(5 + (25 * len(musicLibrary)))
		// 299 as 300 causes side scrolling
		width := float32(299)
		if panelContentRec.Height > 462 {
			// 13 is scrollbar width
			width = float32(299 - 13)
		}

		gui.ScrollPanel(panelRec, fmt.Sprintf("Library (%d)", len(musicLibrary)), panelContentRec, &panelScroll, &rl.Rectangle{0, 0, 0, 0})

		// 4 magic numbers that I dont know what do but they work
		rl.BeginScissorMode(int32(6), int32(54), int32(width), int32(433))
		for i := 0; i < len(musicLibrary); i++ {
			musicFile := musicLibrary[i]

			// check if button y is in view
			if float32(4+(25*i))+55+panelScroll.Y < 54-20 || float32(4+(25*i))+55+panelScroll.Y > 487 {
				continue
			}

			// WHY 40, THE BUTTON IS 20 PIXELS WIDE
			gui.Label(rl.NewRectangle(10, float32(4+(25*i))+55+panelScroll.Y, width-40, 20), fmt.Sprintf("%s", musicFile.Name))

			if musicCurrentIndex == i {
				gui.SetState(gui.STATE_PRESSED)
			}

			if gui.Button(rl.Rectangle{width - 20, float32(4+(25*i)) + 55 + panelScroll.Y, 20, 20}, gui.IconText(gui.ICON_PLAYER_PLAY, "")) {
				loadMusic(musicFile)
				musicCurrentIndex = i
				fmt.Println(musicCurrentIndex)
				fmt.Println(musicFile.Name)
			}

			gui.SetState(gui.STATE_NORMAL)
		}
		rl.EndScissorMode()
	} else if currentPanel == SETTINGS_PANEL {
		gui.Panel(rl.NewRectangle(5, 30, 300, 460), "Settings")

		// About
		gui.GroupBox(rl.NewRectangle(10, 60, 290, 20), "About")
		{
			gui.Label(rl.NewRectangle(15, 60, 280, 20), "Music Player v1.0 by Fluffy")
		}

		// Library
		gui.GroupBox(rl.NewRectangle(10, 90, 290, 45), "Library")
		{
			gui.Label(rl.NewRectangle(15, 90, 280, 20), "Tracks: "+fmt.Sprintf("%d", len(musicLibrary)))
			if gui.Button(rl.NewRectangle(15, 110, 280, 20), "Reload Library") {
				if musicCurrent != nil {
					rl.StopMusicStream(musicCurrent.File)
					rl.UnloadMusicStream(musicCurrent.File)
					musicCurrent = nil
				}

				musicPlaying = false

				musicLibrary, err = loadLibrary(musicPath)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Failed to load music library: %v\n", err)
				}
			}
		}

		// Theme
		gui.GroupBox(rl.NewRectangle(10, 145, 290, 185), "Themes")
		{
			if gui.Button(rl.NewRectangle(15, 155, 280, 20), "Candy") {
				gui.LoadStyle("./theme/style_candy.rgs")
			}
			if gui.Button(rl.NewRectangle(15, 155+25, 280, 20), "Cherry") {
				gui.LoadStyle("./theme/style_cherry.rgs")
			}
			if gui.Button(rl.NewRectangle(15, 155+50, 280, 20), "Cyber") {
				gui.LoadStyle("./theme/style_cyber.rgs")
			}
			if gui.Button(rl.NewRectangle(15, 155+75, 280, 20), "Dark") {
				gui.LoadStyle("./theme/style_dark.rgs")
			}
			if gui.Button(rl.NewRectangle(15, 155+100, 280, 20), "Jungle") {
				gui.LoadStyle("./theme/style_jungle.rgs")
			}
			if gui.Button(rl.NewRectangle(15, 155+125, 280, 20), "Lavanda") {
				gui.LoadStyle("./theme/style_lavanda.rgs")
			}
			if gui.Button(rl.NewRectangle(15, 155+150, 280, 20), "Default") {
				gui.LoadStyleDefault()
			}
		}

		// Controls
		gui.GroupBox(rl.NewRectangle(10, 340, 290, 75), "Controls")
		{
			gui.Label(rl.NewRectangle(15, 340, 280, 20), "Space: Play/Pause")
			gui.Label(rl.NewRectangle(15, 340+17, 280, 20), "Left/Right: Seek")
			gui.Label(rl.NewRectangle(15, 340+35, 280, 20), "Escape: Quit")
			gui.Label(rl.NewRectangle(15, 340+53, 280, 20), "Up/Down: Change Track")
		}
	}

	// Progress bar
	musicSeek = gui.ProgressBar(rl.Rectangle{5, 495, float32(rl.GetScreenWidth()) - 10, 15}, "", "", progressBarValue, progressBarMinValue, progressBarMaxValue)

	if musicCurrent == nil {
		gui.SetState(gui.STATE_DISABLED)
	}

	// Previous button
	if gui.Button(rl.Rectangle{5, 515, 90, 30}, gui.IconText(gui.ICON_PLAYER_PREVIOUS, "")) {
		if musicCurrentIndex > 0 {
			musicCurrentIndex -= 1
		} else {
			musicCurrentIndex = len(musicLibrary) - 1
		}

		loadMusic(musicLibrary[musicCurrentIndex])
	}

	// Play/Pause button
	if gui.Button(rl.Rectangle{100, 515, 97 + 14, 30}, playPauseIcon) {
		if musicPlaying {
			rl.PauseMusicStream(musicCurrent.File)
			musicPlaying = false
		} else {
			rl.ResumeMusicStream(musicCurrent.File)
			musicPlaying = true
		}
	}

	// Next button
	if gui.Button(rl.Rectangle{208 + 7, 515, 90, 30}, gui.IconText(gui.ICON_PLAYER_NEXT, "")) {
		fmt.Println(musicCurrentIndex)

		if musicCurrentIndex < len(musicLibrary)-1 {
			musicCurrentIndex += 1
		} else {
			musicCurrentIndex = 0
		}

		loadMusic(musicLibrary[musicCurrentIndex])

		fmt.Println(musicCurrentIndex)
	}

	gui.SetState(gui.STATE_NORMAL)

	gui.StatusBar(rl.Rectangle{5, 550, 300, 20}, statusBarText)

	// Close window message box
	if showQuitPopup {
		rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.DarkGray, 0.8))
		result := gui.MessageBox(rl.Rectangle{float32(rl.GetScreenWidth())/2 - 125, float32(rl.GetScreenHeight())/2 - 50, 250, 100}, gui.IconText(gui.ICON_EXIT, "Close Window"), "Do you really want to exit?", "Yes;No")

		if (result == 0) || (result == 2) {
			showQuitPopup = false
		} else if result == 1 {
			exitApplication = true
		}
	}
}

func main() {
	os.Exit(run())
}
