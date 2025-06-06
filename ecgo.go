package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

var (
	noNewline     bool
	enableEscapes bool
	colorName     string
)

func init() {
	flag.BoolVar(&noNewline, "n", false, "Do not output the trailing newline.")
	flag.BoolVar(&enableEscapes, "e", false, "Enable interpretation of backslash escapes (e.g., \\n, \\t).")
	flag.StringVar(&colorName, "color", "", "Specify output color (e.g., red, green, blue, yellow, magenta, cyan, white, black).")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [OPTIONS] [STRING...]:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("\nAvailable colors: red, green, blue, yellow, magenta, cyan, white, black.")
	}
}

func applyEscapes(s string) string {
	s = strings.ReplaceAll(s, "\\n", "\n")
	s = strings.ReplaceAll(s, "\\t", "\t")
	return s
}

func main() {
	flag.Parse()

	args := flag.Args()

	if len(args) == 0 {
		if !noNewline {
			fmt.Println()
		}
		return
	}

	var output strings.Builder
	for i, arg := range args {
		if enableEscapes {
			output.WriteString(applyEscapes(arg))
		} else {
			output.WriteString(arg)
		}
		if i < len(args)-1 {
			output.WriteString(" ")
		}
	}

	outputString := output.String()

	var printer *color.Color

	if colorName != "" {
		switch strings.ToLower(colorName) {
		case "red":
			printer = color.New(color.FgRed)
		case "green":
			printer = color.New(color.FgGreen)
		case "blue":
			printer = color.New(color.FgBlue)
		case "yellow":
			printer = color.New(color.FgYellow)
		case "magenta":
			printer = color.New(color.FgMagenta)
		case "cyan":
			printer = color.New(color.FgCyan)
		case "white":
			printer = color.New(color.FgWhite)
		case "black":
			printer = color.New(color.FgBlack)
		default:
			fmt.Fprintf(os.Stderr, "Warning: Unknown color '%s'. Printing without color.\n", colorName)
			printer = nil
		}
	}

	if printer != nil {
		if noNewline {
			printer.Print(outputString)
		} else {
			printer.Println(outputString)
		}
	} else {
		if noNewline {
			fmt.Print(outputString)
		} else {
			fmt.Println(outputString)
		}
	}
}
