# Translate SRT
Translates a .srt file via deepl.com (Made for GCI 2017). It only works for small files currently due to DeepL's rate limiting on their private API.

## Installation
Binaries are available in the releases page. Alternatively, if Go is installed, this program can be built from source.
```
$ dep ensure
$ go install
$ translate-srt
```

## Usage
```
$ translate-srt -input input.srt -output output.srt -from EN -to DE
```
