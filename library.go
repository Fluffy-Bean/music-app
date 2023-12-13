package main

import (
	"errors"
	rl "github.com/gen2brain/raylib-go/raylib"
	"os"
)

type Music struct {
	Path   string
	Name   string
	File   rl.Music
	Artist string
	Liked  bool
}

func loadLibrary(musicFolder string) ([]Music, error) {
	if _, err := os.Stat(musicFolder); os.IsNotExist(err) {
		return nil, errors.New("music folder does not exist")
	}

	files, err := os.ReadDir(musicFolder)
	if err != nil {
		return nil, err
	}

	var musicList []Music
	for _, file := range files {
		musicList = append(musicList, Music{Path: musicFolder, Name: file.Name(), Artist: "John Doe", Liked: false})
	}

	return musicList, nil
}

func loadMusic(music Music) error {
	if musicCurrent != nil {
		rl.StopMusicStream(musicCurrent.File)
		rl.UnloadMusicStream(musicCurrent.File)
	}

	musicPlaying = false
	musicCurrent = &music
	musicCurrent.File = rl.LoadMusicStream(musicCurrent.Path + "/" + musicCurrent.Name)
	if musicCurrent == nil {
		return errors.New("failed to load music")
	}

	rl.PlayMusicStream(musicCurrent.File)
	progressBarMaxValue = rl.GetMusicTimeLength(musicCurrent.File)
	progressBarValue = rl.GetMusicTimePlayed(musicCurrent.File)
	musicPlaying = true

	return nil
}
