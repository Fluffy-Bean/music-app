package main

import (
	gui "github.com/gen2brain/raylib-go/raygui"
	rl "github.com/gen2brain/raylib-go/raylib"

	"fmt"
	"math"
	"os"
)

func main() {
	const (
		screenWidth  = 310
		screenHeight = 575
	)
	var (
		musicName     string
		musicFile     rl.Music
		musicVolume   float32 = 50
		musicSeek     float32 = 0
		musicPlaying          = false
		playPauseIcon         = "#131#" // #131# = play, #132# = pause

		listViewExScrollIndex int32 = -1
		listViewExActive      int32 = 2
		listViewExFocus       int32 = 0
		listViewExList        []string

		currentPanel    = 0
		showQuitPopup   = false
		exitApplication = false
	)

	rl.InitWindow(screenWidth, screenHeight, "music")
	rl.InitAudioDevice()
	rl.SetExitKey(0)
	rl.SetTargetFPS(60)

	if _, err := os.Stat("./theme"); !os.IsNotExist(err) {
		gui.LoadStyle("./theme/style_cyber.rgs")
	}

	// Create music folder if it doesn't exist
	if _, err := os.Stat("./music"); os.IsNotExist(err) {
		err := os.Mkdir("./music", 0755)
		if err != nil {
			return
		}
	}

	// Check if there are any music files in the music folder
	if items, _ := os.ReadDir("./music"); len(items) == 0 {
		exitApplication = false
		for !exitApplication {
			rl.BeginDrawing()
			rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

			rl.DrawRectangle(0, 0, int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight()), rl.Fade(rl.DarkGray, 0.8))
			result := gui.MessageBox(rl.Rectangle{float32(rl.GetScreenWidth())/2 - 125, float32(rl.GetScreenHeight())/2 - 50, 250, 120}, gui.IconText(gui.ICON_EXIT, "No Music Files"), "There is no music in the music folder!\n You need to add some to begin.", "Ok")

			if (result == 0) || (result == 1) {
				exitApplication = true
			}

			rl.EndDrawing()
		}

		rl.CloseAudioDevice()
		rl.UnloadMusicStream(musicFile)
		rl.CloseWindow()

		return
	}

	// Load all music files from the music folder
	tracks, _ := os.ReadDir("./music")
	for _, track := range tracks {
		listViewExList = append(listViewExList, track.Name())
	}

	for !exitApplication {
		exitApplication = rl.WindowShouldClose()

		if rl.IsKeyPressed(rl.KeyEscape) {
			showQuitPopup = !showQuitPopup
		}

		if rl.IsKeyPressed(rl.KeySpace) {
			if musicPlaying {
				rl.PauseMusicStream(musicFile)
				rl.SetWindowTitle("Music - Playing: " + musicName)
				playPauseIcon = "#131#"
				musicPlaying = false
			} else {
				rl.ResumeMusicStream(musicFile)
				rl.SetWindowTitle("Music - Playing: " + musicName)
				playPauseIcon = "#132#"
				musicPlaying = true
			}
		}

		if rl.IsKeyPressed(rl.KeyLeft) {
			if rl.GetMusicTimePlayed(musicFile)-10 > 0 {
				rl.SeekMusicStream(musicFile, float32(math.Max(float64(rl.GetMusicTimePlayed(musicFile)-10), 0)))
			}
		}
		if rl.IsKeyPressed(rl.KeyRight) {
			if rl.GetMusicTimePlayed(musicFile)+10 < rl.GetMusicTimeLength(musicFile) {
				rl.SeekMusicStream(musicFile, float32(math.Min(float64(rl.GetMusicTimePlayed(musicFile)+10), float64(rl.GetMusicTimeLength(musicFile)))))
			}
		}

		if rl.IsKeyPressed(rl.KeyUp) {
			if listViewExScrollIndex > 0 {
				listViewExScrollIndex -= 1
			} else {
				listViewExScrollIndex = int32(len(listViewExList)) - 1
			}
		}
		if rl.IsKeyPressed(rl.KeyDown) {
			if listViewExScrollIndex < int32(len(listViewExList))-1 {
				listViewExScrollIndex += 1
			} else {
				listViewExScrollIndex = 0
			}
		}

		rl.BeginDrawing()
		rl.ClearBackground(rl.GetColor(uint(gui.GetStyle(gui.DEFAULT, gui.BACKGROUND_COLOR))))

		// Volume bar
		musicVolume = gui.SliderBar(rl.Rectangle{25, 5, 255, 20}, "#122#", "", musicVolume, 0, 100)
		rl.SetMusicVolume(musicFile, musicVolume/100)

		// Settings button
		if currentPanel == 0 {
			if gui.Button(rl.Rectangle{285, 5, 20, 20}, "#141#") {
				currentPanel = 1
			}
		} else if currentPanel == 1 {
			if gui.Button(rl.Rectangle{285, 5, 20, 20}, "#214#") {
				currentPanel = 0
			}
		}

		if currentPanel == 0 {
			// List all musicFiles
			listViewExActive = gui.ListViewEx(rl.Rectangle{5, 30, 300, 460}, listViewExList, &listViewExFocus, &listViewExScrollIndex, listViewExActive)
		} else if currentPanel == 1 {
			// About
			gui.Panel(rl.NewRectangle(5, 30, 300, 40), "About")
			gui.Label(rl.NewRectangle(10, 55, 290, 20), "Music Player v1.0 by Fluffy")

			// Library
			gui.Panel(rl.NewRectangle(5, 83, 300, 71), "Library")
			gui.Label(rl.NewRectangle(10, 108, 290, 20), "Tracks: "+fmt.Sprintf("%d", len(listViewExList)))
			if gui.Button(rl.NewRectangle(10, 128, 290, 20), "Reload Library") {
				rl.StopMusicStream(musicFile)

				musicPlaying = false
				playPauseIcon = "#131#"
				listViewExList = nil

				tracks, _ := os.ReadDir("./music")
				for _, track := range tracks {
					listViewExList = append(listViewExList, track.Name())
				}
			}

			// Theme
			gui.Panel(rl.NewRectangle(5, 159, 300, 205), "Themes")
			if _, err := os.Stat("./theme"); !os.IsNotExist(err) {
				if gui.Button(rl.NewRectangle(10, 188, 290, 20), "Candy") {
					gui.LoadStyle("./theme/style_candy.rgs")
				}
				if gui.Button(rl.NewRectangle(10, 188+25, 290, 20), "Cherry") {
					gui.LoadStyle("./theme/style_cherry.rgs")
				}
				if gui.Button(rl.NewRectangle(10, 188+50, 290, 20), "Cyber") {
					gui.LoadStyle("./theme/style_cyber.rgs")
				}
				if gui.Button(rl.NewRectangle(10, 188+75, 290, 20), "Dark") {
					gui.LoadStyle("./theme/style_dark.rgs")
				}
				if gui.Button(rl.NewRectangle(10, 188+100, 290, 20), "Jungle") {
					gui.LoadStyle("./theme/style_jungle.rgs")
				}
				if gui.Button(rl.NewRectangle(10, 188+125, 290, 20), "Lavanda") {
					gui.LoadStyle("./theme/style_lavanda.rgs")
				}
				if gui.Button(rl.NewRectangle(10, 188+150, 290, 20), "Default") {
					gui.LoadStyleDefault()
				}
			} else {
				gui.Label(rl.NewRectangle(10, 188, 290, 20), "No themes found")
			}

			// Controls
			gui.Panel(rl.NewRectangle(5, 370, 300, 100), "Controls")
			gui.Label(rl.NewRectangle(10, 395, 290, 20), "Space: Play/Pause")
			gui.Label(rl.NewRectangle(10, 395+17, 290, 20), "Left/Right: Seek")
			gui.Label(rl.NewRectangle(10, 395+35, 290, 20), "Escape: Quit")
			gui.Label(rl.NewRectangle(10, 395+53, 290, 20), "Up/Down: Change Track")
		}

		// Progress bar
		musicSeek = gui.ProgressBar(rl.Rectangle{5, 495, float32(rl.GetScreenWidth()) - 10, 15}, "#122#", "", rl.GetMusicTimePlayed(musicFile), 0, rl.GetMusicTimeLength(musicFile))

		// Previous button
		if gui.Button(rl.Rectangle{5, 515, 90, 30}, "#129#") {
			if listViewExScrollIndex > 0 {
				listViewExScrollIndex -= 1
			} else {
				listViewExScrollIndex = int32(len(listViewExList)) - 1
			}
		}

		// Play/Pause button
		if gui.Button(rl.Rectangle{100, 515, 97 + 14, 30}, playPauseIcon) {
			if musicPlaying {
				rl.PauseMusicStream(musicFile)
				rl.SetWindowTitle("Music")
				playPauseIcon = "#131#"
				musicPlaying = false
			} else {
				rl.ResumeMusicStream(musicFile)
				rl.SetWindowTitle("Music - Playing: " + musicName)
				playPauseIcon = "#132#"
				musicPlaying = true
			}
		}

		// Next button
		if gui.Button(rl.Rectangle{208 + 7, 515, 90, 30}, "#134#") {
			if listViewExScrollIndex < int32(len(listViewExList))-1 {
				listViewExScrollIndex += 1
			} else {
				listViewExScrollIndex = 0
			}
		}

		// convert musicSeek to minutes and seconds, probably easier way to do this, but eh
		minutesPlayed := int(musicSeek / 60)
		secondsPlayed := "00"
		if int(musicSeek-float32(minutesPlayed)*60) < 10 {
			secondsPlayed = fmt.Sprintf("0%d", int(musicSeek-float32(minutesPlayed)*60))
		} else {
			secondsPlayed = fmt.Sprintf("%d", int(musicSeek-float32(minutesPlayed)*60))
		}

		minutesLength := int(rl.GetMusicTimeLength(musicFile) / 60)
		secondsLength := "00"
		if int(rl.GetMusicTimeLength(musicFile)-float32(minutesLength)*60) < 10 {
			secondsLength = fmt.Sprintf("0%d", int(rl.GetMusicTimeLength(musicFile)-float32(minutesLength)*60))
		} else {
			secondsLength = fmt.Sprintf("%d", int(rl.GetMusicTimeLength(musicFile)-float32(minutesLength)*60))
		}

		gui.StatusBar(rl.Rectangle{5, 550, 300, 20}, fmt.Sprintf("Tracks: %d | Time: %d:%s/%d:%s | %s", len(listViewExList), minutesPlayed, secondsPlayed, minutesLength, secondsLength, musicName))

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

		rl.EndDrawing()

		// if end of track, play next track
		// -0.3 as if we do a direct comparison it will never be true, probably due to floating point errors
		if rl.GetMusicTimePlayed(musicFile) >= rl.GetMusicTimeLength(musicFile)-0.3 {
			if listViewExScrollIndex < int32(len(listViewExList))-1 {
				listViewExScrollIndex += 1
			} else {
				listViewExScrollIndex = 0
			}
		}

		if listViewExScrollIndex != -1 {
			if listViewExList[listViewExScrollIndex] != musicName {
				rl.StopMusicStream(musicFile)
				rl.UnloadMusicStream(musicFile)

				musicName = listViewExList[listViewExScrollIndex]
				musicFile = rl.LoadMusicStream("./music/" + musicName)
				musicSeek = 0

				rl.PlayMusicStream(musicFile)
				if !musicPlaying {
					// Stop playback if the user had their music paused
					//rl.PauseMusicStream(musicFile)
				}
			}
			musicName = listViewExList[listViewExScrollIndex]
		}

		rl.UpdateMusicStream(musicFile)
	}

	rl.CloseAudioDevice()
	rl.UnloadMusicStream(musicFile)
	rl.CloseWindow()
}
