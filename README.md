# Metronome.
A TUI metronome with support for iterating over declared chord patterns.

# Installation
* `git clone git@github.com:andrewwillette/metronome.git`
* Ensure A yml file is [configured](#song-display-configuration)
* Execute `go run .`

# Song Display Configuration
Inside the `./resources` directory create a yaml file. It will represent a musical song and its chords, per beat.
An example specification for the classic song made famous by Hank Williams, "Lost Highway" with the filename stored at `./resources/LostHighway.yml`.

```yml
song: Lost Highway
sections:
  a:
    - [D,D,D,D]
    - [D,D,G,G]
    - [D,D,D,D]
    - [D,D,D,D]
    - [D,D,D,D]
    - [D,D,G,G]
    - [A,A,A,A]
    - [A,A,A,A]
  b:
    - [G,G,G,G]
    - [G,G,G,G]
    - [D,D,D,D]
    - [D,D,D,D]
    - [D,D,D,D]
    - [D,D,A,A]
    - [D,D,D,D]
    - [D,D,D,D]
```
