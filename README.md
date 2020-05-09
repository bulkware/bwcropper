# bwCropper

A tool to rotate and crop scanned images automatically.


## Preamble
This tool is intended for rotating and cropping scanned images automatically.
It was specifically made for processing comic magazine pages. It assumes the
scanned input image is surrounded by black and grey colours. It then measures
the distance between the edges and the image by calculating the RGB colour
luminance of each pixel. If the pixel luminance exceeds given value, it assumes
the image begins. It also rotates the pages to different angles, trying to find
the optimal angle which to use for cropping.

This tool was also made because I wanted to learn Go programming language. I am
still learning the very basics, so any issues, improvements or any input really
would be greatly appreciated =)


## Relative luminance
The formula for calculating relative luminance is taken from here:
https://en.wikipedia.org/wiki/Relative_luminance
