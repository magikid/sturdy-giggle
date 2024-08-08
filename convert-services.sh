#!/usr/bin/env bash

set -euo pipefail
IFS=$'\n\t'

#ffmpeg -i input.wav -vn -ar 44100 -ac 2 -q:a 2 output.mp3
INPUT_FOLDER='./Raw Services'
OUTPUT_FOLDER='./Converted Services'
find "$INPUT_FOLDER" -type f -iname '*.wav' -print0 | while read -d $'\0' file
do
	FILE_NAME=$(basename -s.WAV "$file")
	echo "Converting $file to $OUTPUT_FOLDER/$FILE_NAME.mp3"
	< /dev/null ffmpeg -y -hide_banner -loglevel warning -i "$file" -vn -ar 44100 -ac 2 -channel_layout stereo -q:a 2 "$OUTPUT_FOLDER/$FILE_NAME.mp3"
done

