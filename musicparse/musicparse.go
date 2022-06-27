package musicparse

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Song struct {
	Title    string `yaml:"song,omitempty"`
	Sections struct {
		ASection [][]string `yaml:"a,flow,omitempty"`
		BSection [][]string `yaml:"b,flow,omitempty"`
		CSection [][]string `yaml:"c,flow,omitempty"`
	}
}

func readSong(fpath string) Song {
	yfile, err := ioutil.ReadFile(fpath)
	if err != nil {
		panic(err)
	}
	var song Song
	err = yaml.Unmarshal(yfile, &song)
	if err != nil {
		panic(err)
	}
	return song
}

func readSongsFromDir(fpath string) []Song {
	f, err := os.Open(fpath)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	files, err := f.Readdir(0)
	if err != nil {
		panic(err)
	}
	songs := []Song{}
	for _, v := range files {
		song := readSong(filepath.Join(fpath, v.Name()))
		songs = append(songs, song)
	}
	return songs
}

func GetDefaultSongs() []Song {
	return readSongsFromDir("./../resources")
}

func GetSongsXdg() []Song {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return readSongsFromDir(filepath.Join(dirname, "/.config/metronome"))
}
