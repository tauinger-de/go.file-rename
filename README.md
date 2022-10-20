# go.file-rename

## about

This project aims to create a command-line tool that performs bulk renaming of media 
files (photos, videos) in a folder so that the files can be sorted by filename only to
see them in chronological order.

I want to provide extensive support for EXIF data reading as well as fetching creation date-time
information from other sources to obtain the most accurate creation date-time for a media file.

## usage

Build and install with `go install`.

Run the tool with `gofire -source <source-dir> -topic <your-topic>`. Defaults are
- the current directory for source
- the word "media" for topic

The `topic` is a string that will appear in each filename, such as "2022-10-19 20.35.18 **media**-0001.mov"
