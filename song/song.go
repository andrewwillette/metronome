package song

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

type Sections struct {
	ASection [][]string `yaml:"a,flow,omitempty"`
	BSection [][]string `yaml:"b,flow,omitempty"`
	CSection [][]string `yaml:"c,flow,omitempty"`
}

type Song struct {
	Title    string `yaml:"song,omitempty"`
	Sections Sections
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
func GetLostCowboySongs() []Song {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return []Song{readSong(filepath.Join(dirname, "/.config/metronome/LostHighway.yml"))}
}

func GetSongsXdg() []Song {
	dirname, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return readSongsFromDir(filepath.Join(dirname, "/.config/metronome"))
}

func getFramesAsSpaces(section [][]string) []string {
	spaceFrames := []string{}
	return spaceFrames
}

func appendSectionFrames(section [][]string, frames []string) []string {
	sectionFrame := 0
	// spaceFrames := getFramesAsSpaces(song.Sections.ASection)
	for _, bar := range section {
		for range bar {
			frame := getSectionFrame(section, sectionFrame)
			sectionFrame++
			frames = append(frames, frame)
		}
	}
	return frames
}

func GetSongFrames(song Song) []string {
	frames := []string{}
	if len(song.Sections.ASection) > 0 {
		frames = appendSectionFrames(song.Sections.ASection, frames)
	}
	if len(song.Sections.BSection) > 0 {
		frames = appendSectionFrames(song.Sections.BSection, frames)
	}
	if len(song.Sections.CSection) > 0 {
		frames = appendSectionFrames(song.Sections.CSection, frames)
	}
	return frames
}

// getSectionFrame return string representation of the song
// when the given sectionIndex is active. It ends up being spaces
// everywhere but the index value.
func getSectionFrame(section [][]string, sectionIndex int) string {
	var sb strings.Builder
	i := 0
	for barIndex, bar := range section {
		for _, beat := range bar {
			if i == sectionIndex {
				sb.WriteString(beat)
			} else {
				sb.WriteString(" ")
			}
			i++
		}
		if barIndex != (len(section) - 1) {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
