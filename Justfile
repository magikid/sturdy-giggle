new:
    go run convert-services.go

all:
    go run convert-services.go -all

one input output:
    ffmpeg -y -hide_banner -loglevel warning -i '{{input}}' -vn -ar 44100 -ac 2 -channel_layout stereo -q:a 2 '{{output}}'
