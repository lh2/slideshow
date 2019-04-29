slideshow(1)

# NAME

slideshow - imageviewer controlled by stdin

# SYNOPSIS

*slideshow* [options...] < _filepaths_

# OPTIONS

*-bg*
	Background color as 6 digit hex value.

# DESCRIPTION

On launch, slideshow displays an empty window. Feed file paths to standard input
to load a new image. On every newline character, slideshow loads and displays
the image found at the specified filepath. Error messages will be written to
standard error, as long as the error in question is not fatal, the program will
remain running.

Image files need to be in *farbfeld*(5) format!

# EXAMPLES

```
$ find ~/my_farbfeld_images/ -type f | while read -r file; do echo "$file"; sleep 5s; done | slideshow
```

# SEE ALSO

*farbfeld*(5)

# AUTHORS

Lukas Henkel <lh@entf.net>
