# Metronome.
A TUI metronome with support for iterating over declared chord patterns.

# Installation
* Execute `git clone git@github.com:andrewwillette/metronome.git`.
* Ensure the `go` executable exists on your `$PATH`.
* Execute `go run .`.

# Configuration
## Content
Files which will be used for metronomes ticking. They represent musical songs and their chords, per beat.
An example specification for the classic song made famous by Hank Williams, "Lost Highway".

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
