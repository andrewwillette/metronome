package song

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
	require.Equal(t, 3, len(songs))
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

func Test_getSongsXdg(t *testing.T) {
	songs := GetSongsXdg()
	require.Equal(t, 2, len(songs))
}

// assertSongCheap Assert some values, a for effort
func assertSongCheap(t *testing.T, song Song, title string, sectionAbar1, sectionBbar1 []string) {
	require.Equal(t, title, song.Title)
	require.Equal(t, sectionAbar1, song.Sections.ASection[0])
	require.Equal(t, sectionBbar1, song.Sections.BSection[0])
}

func BenchmarkReadSong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readSong("./../resources/LostHighway.yml")
	}
}

func BenchmarkReadSongsFromDir(b *testing.B) {
	for i := 0; i < b.N; i++ {
		readSongsFromDir("./../resources")
	}
}

func Test_GetLostCowboySongs(t *testing.T) {
	songs := GetLostCowboySongs()
	lhw := songs[0]
	require.Equal(t, "Lost Highway", lhw.Title)
	require.Equal(t, []string{"D", "D", "D", "D"}, lhw.Sections.ASection[0])
	require.Equal(t, []string{"D", "D", "G", "G"}, lhw.Sections.ASection[1])
	require.Equal(t, []string{"A", "A", "A", "A"}, lhw.Sections.ASection[6])
	require.Equal(t, []string{"G", "G", "G", "G"}, lhw.Sections.BSection[1])
	require.Equal(t, []string{"D", "D", "D", "D"}, lhw.Sections.BSection[7])
}

func Test_getSongFrames(t *testing.T) {
	// lcp := musicparse.GetLostCowboySongs()
	// lcps := lcp[0]
	tests := []struct {
		song           Song
		expectedFrames []string
	}{
		{
			song:           TwelveBarBlues,
			expectedFrames: ExptectedTwelveBarBluesFrames(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.song.Title, func(t *testing.T) {
			frames := getSongFrames(tt.song)
			require.Equal(t, tt.expectedFrames, frames)
		})
	}
}
