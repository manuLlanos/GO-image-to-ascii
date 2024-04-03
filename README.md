# GO image to ascii converter

Takes a png file and creates ascii art in a text file in the same directory.

Usage:
	
	imageToAscii YourImage.png [width]
	
[width] is optional, if unspecified or 0, it uses all the pixels of the original image.
If it's a positive number it will temporarily resize the image to that width and create a smaller or larget ascii art.