package musicparse

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_musicParse(t *testing.T) {
	song := readSong("./../resources/LostHighway.yml")
	require.NotNil(t, song)
	assertSongCheap(t, song, "Lost Highway", []string{"D", "D", "D", "D"}, []string{"G", "G", "G", "G"})
}

func Test_readSongsFromDir(t *testing.T) {
	songs := readSongsFromDir("./../resources")
	require.Equal(t, 2, len(songs))
	getByTitle := func(ttle string, sngs []Song) Song {
		for _, v := range sngs {
			if v.Title == ttle {
				return v
			}
		}
		return Song{}
	}
	const lhwtl = "Lost Highway"
	lhw := getByTitle(lhwtl, songs)
	assertSongCheap(t, lhw, lhwtl, []string{"D", "D", "D", "D"}, []string{"G", "G", "G", "G"})

	const ccbt = "Carrol County Blues"
	ccb := getByTitle(ccbt, songs)
	assertSongCheap(t, ccb, ccbt, []string{"G", "G", "G", "G"}, []string{"D", "D", "D", "D"})
}

// assertSongCheap Assert some values, a for effort
func assertSongCheap(t *testing.T, song Song, title string, sectionAbar1, sectionBbar1 []string) {
	require.Equal(t, title, song.Title)
	require.Equal(t, sectionAbar1, song.Sections.ASection[0])
	require.Equal(t, sectionBbar1, song.Sections.BSection[0])
}
