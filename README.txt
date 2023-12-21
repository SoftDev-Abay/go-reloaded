GO-RELOADED

Text Formatting Tool

This is a simple command-line tool for text formatting in Go. The tool takes two file names as input and formats the content of the first file according to specified rules, then writes the formatted text to the second file.
Usage

bash

go run main.go <input_file> <output_file>

Replace <input_file> with the name of the file you want to format, and <output_file> with the desired name for the formatted output.
Rules

    The input files must exist, and their extensions must be ".txt".

    The tool performs various text formatting operations, including:
        Hexadecimal to decimal conversion: (hex) indicator.
        Binary to decimal conversion: (bin) indicator.
        Capitalization changes: (up), (low), (cap) indicators.
        Article "a" corrections: Ensures proper usage of "a" and "an."
        Punctuation spacing corrections.

    The formatted text is written to the specified output file.

Example

bash

go run main.go input.txt output.txt

Notes

    Make sure to provide the correct file names and extensions.
    The tool provides information about the processed files.

If you wish to test the program, there are already tests you can run them with

go test
