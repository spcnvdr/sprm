// TODO: use os.stat() to make sure source file exists
// TODO: make a samefile function using os.Stat to determine if
// the input file is the same as the output file to avoid an unnecessary copy
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

const sprmVersion = "sprm 0.0.3"

/** cpfile - Copy file src to dst
 *  @param src string the filename/path of the file to copy
 *  @param dst string where to put the copy of src
 *  @returns on success the number of bytes copied,
 *  else 0 and an error message
 */
func cpfile(src, dst string) (int64, error) {
	srcstat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !srcstat.Mode().IsRegular() {
		return 0, err
	}

	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}

	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()

	n, err := io.Copy(out, in)
	return n, err
}

/** Remove every occurence of a rune from s if it occurrs in the rm string
 *  @param s string the string to remove runes from
 *  @param rm a string of characters to remove from s
 *  @returns a new string with every character in rm removed, s is not modified
 *  e.g. s := "This is s"; removeall(s, " s") returns "Thii"
 */
func rmChr(s, rm string) string {
	for _, v := range rm {
		s = strings.ReplaceAll(s, string(v), "")
	}
	return s
}

/** yesno - Get an affirmitive or negative answer to prompt
 *  @param prompt string optional prompt to display to the user
 *  @returns true if user entered an affirmative (y|Y|yes|YES|...)
 *  else a negative
 */
func yesno(prompt string) bool {
	if prompt != "" {
		fmt.Printf("%s: ", prompt)
	}
	reader := bufio.NewReader(os.Stdin)
	s, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}

	up := strings.ToUpper(s)
	if up[0] == 'Y' {
		return true
	}
	return false
}

/*
sprm - Remove all spaces from a filename
@param opath string the path of the file to rename
@param replace string replace spaces with this character
@param strip string if not empty, remove these characters from the filename
@param backup bool copy the file to path with new name instead of renaming
@param ask bool prompt user before renaming/copying file
@returns nil on success, else an error
*/
func sprm(opath string, replace string, strip string, backup, ask bool) error {
	var ext string
	var filename string
	var err error

	dir, fn := path.Split(opath)

	// if it has an extension, save the extension and remove it
	// before modifying filename
	if strings.ContainsAny(fn, ".") {
		ext = path.Ext(fn)
		filename = strings.TrimSuffix(fn, ext)
	} else {
		filename = fn
	}

	// if strip argument used, strip the given chars from filename
	if strip != "" {
		filename = rmChr(filename, strip)
	}

	// replace or remove spaces based on cmd line arguments

	filename = strings.ReplaceAll(filename, " ", replace)

	//add extension back and join it with path
	filename += ext
	newpath := path.Join(dir, filename)

	if backup {
		// before copy, get confirmation
		if ask {
			fmt.Printf("sprm: copy '%s' to '%s'? (y/n): ", opath, newpath)
			if !yesno("") {
				return nil
			}
		}

		n, err := cpfile(opath, newpath)
		if err != nil {
			return err
		}

		if verboseG {
			fmt.Printf("copied file: %s -> %s (%d bytes)\n", opath, newpath, n)
		}

	} else {
		// before rename, get confirmation
		if ask {
			fmt.Printf("sprm: rename '%s' to '%s'? (y/n): ", opath, newpath)
			if !yesno("") {
				return nil
			}
		}

		err = os.Rename(opath, newpath)
		if err != nil {
			return err
		}

		if verboseG {
			fmt.Printf("renamed file: %s -> %s\n", opath, newpath)
		}

	}

	return err
}

/** printUsage - Print a simple usage message
 *
 */
func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: sprm [OPTION...] FILE...\n")
	fmt.Fprintf(os.Stderr, "Try `sprm --help' or `sprm -h' for more information\n")
}

/** printHelp - Print a custom help message
 *
 */
func printHelp() {

	fmt.Fprintf(os.Stderr, "Usage: sprm [OPTION...] FILE...\n")
	fmt.Fprintf(os.Stderr, "Remove spaces and other characters from FILE name(s)\n\n")
	fmt.Fprintf(os.Stderr, "  -b, --backup                Make a copy instead of renaming in place\n")
	fmt.Fprintf(os.Stderr, "  -d, --dash                  Replace spaces with dashes/hyphens\n")
	fmt.Fprintf(os.Stderr, "  -i, --interactive           Prompt before renaming/copying file\n")
	fmt.Fprintf(os.Stderr, "  -s, --strip=CHARS           Remove the given characters from the filename\n")
	fmt.Fprintf(os.Stderr, "  -u, --underscore            Replace spaces with underscores\n")
	fmt.Fprintf(os.Stderr, "  -v, --verbose               Verbosely list files processed\n")
	fmt.Fprintf(os.Stderr, "  -?, -h, --help              Show this help message\n")
	fmt.Fprintf(os.Stderr, "  -V, --version               Print program version\n")
	fmt.Fprintf(os.Stderr, "\n")
}

// global variables for command line arguments, only used in sprm()
var (
	backupG      bool
	dashG        bool
	underscoreG  bool
	versionG     bool
	verboseG     bool
	interactiveG bool
	questionG    bool
	stripG       string
)

// init is automatically called at start, setup cmd line args
func init() {
	// verbose mode and shortcut
	flag.BoolVar(&verboseG, "verbose", false, "verbose output")
	flag.BoolVar(&verboseG, "v", false, "verbose shortcut")

	// strip
	flag.StringVar(&stripG, "strip", "", "Remove the given characters from filename")
	flag.StringVar(&stripG, "s", "", "Strip shortcut")

	// version
	flag.BoolVar(&versionG, "version", false, "Print program version")
	flag.BoolVar(&versionG, "V", false, "Print program version")

	// backup
	flag.BoolVar(&backupG, "backup", false, "Leave the original file unchanged")
	flag.BoolVar(&backupG, "b", false, "Leave the original file unchanged")

	// interactive
	flag.BoolVar(&interactiveG, "interactive", false, "Prompt before renaming/copying file")
	flag.BoolVar(&interactiveG, "i", false, "Interactive shortuct")

	//dash mode
	flag.BoolVar(&dashG, "dash", false, "Replace spaces with dashes/hyphens")
	flag.BoolVar(&dashG, "d", false, "Dash shortcut")

	//underscore mode
	flag.BoolVar(&underscoreG, "underscore", false, "Replace spaces with underscores")
	flag.BoolVar(&underscoreG, "u", false, "Underscore shortcut")

	flag.BoolVar(&questionG, "?", false, "Underscore shortcut")
}

func main() {
	flag.Usage = printHelp

	flag.Parse()
	if dashG && underscoreG {
		fmt.Fprintf(os.Stderr,
			"Error: cannot replace spaces with dashes and underscores!\n")
		os.Exit(1)
	}

	if questionG {
		printHelp()
		os.Exit(0)
	}

	var spaceReplace string
	if dashG {
		spaceReplace = "-"
	} else if underscoreG {
		spaceReplace = "_"
	}

	if versionG {
		fmt.Println(sprmVersion)
		os.Exit(0)
	}

	if len(flag.Args()) == 0 {
		printUsage()
		os.Exit(1)
	}

	// loop through each file and fix the filename
	for _, v := range flag.Args() {
		err := sprm(v, spaceReplace, stripG, backupG, interactiveG)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}

	}
}
