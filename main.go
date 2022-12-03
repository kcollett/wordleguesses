package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"
	"unicode"

	mapset "github.com/deckarep/golang-set"
	"golang.org/x/exp/slices"
	"golang.org/x/exp/slog"
)

const program_name string = "wordle_guesses"
const long_usage_preamble string = "usage: " + program_name +
	" [-h] [-e excluded_letters | -i included_letters] template"
const long_usage_para0 string = "When playing Wordle, sometimes it can be " +
	"helpful to write out a list of candidate guesses.  For example, you " +
	"might be considering all the possibilities that arise from changing " +
	"the third character in the sequence '_A?AM'. Assuming through previous " +
	"play you've already ruled out 'R', 'I', 'S', 'E', 'N', 'G', 'Y', 'C', " +
	"'U', and 'K', you would end up generating this list (assuming you are " +
	"being exhaustive and not skipping improbable candidates): _AAAM, " +
	"_ABAM, _ADAM, _AFAM, _AHAM, _AJAM, _ALAM, _AMAM, _AOAM, _APAM, _AQAM, " +
	"_ATAM, _AVAM, _AWAM, _AXAM, and _AZAM."
const long_usage_para1 string = "Writing lists like these out " +
	"can be quite laborious, and can create a significant hindrance to " +
	"those with diminished dexterity. wordle_guesses is a program you can " +
	"use to alleviate this burden by printing out candidate " +
	"Wordle guesses. You specify the pattern for the candidate guesses using " +
	"a 5-letter template composed of alphabetical letters, any number of the " +
	"character '" + blank_char + "', and a single occurrence of the the " +
	"character '" + change_char + "'. The '" + change_char + "' character " +
	"indicates the letter to be changed to generate the candidate guesses. " +
	"('" + change_char + "' is used instead of '?' to avoid issues with " +
	"command-line processors that try to perform substitution using '?'.)"
const long_usage_para2 string = "The program will iterate through the alphabet, " +
	"substituting the '" + change_char + "' with candidate letters to " +
	"generate a guess. You can specify a list of letters to exclude when " +
	"generating the candidate guesses; typically, you would do this for the " +
	"letters which Wordle has indicated aren't in the answer. Alternatively, " +
	"instead of iterating through the alphabet, you can specify the set of " +
	"letters to include when making guesses."
const template_argument string = "  template\ttemplate is a 5-character " +
	"sequence composed of letters, any number of the character '" + blank_char +
	"', and a single instance of the character '" + change_char + "' " +
	"('" + blank_char + "a" + change_char + "am', for example)."

var long_usage string = long_usage_preamble + "\n\n" +
	insert_newlines(long_usage_para0, 70) + "\n\n" +
	insert_newlines(long_usage_para1, 70) + "\n\n" +
	insert_newlines(long_usage_para2, 70) + "\n\n" +
	"positional arguments:\n" +
	insert_newlines_with_prefix(template_argument, 70, "\t\t") + "\n\n" +
	"optional arguments:\n" +
	"  -e excluded_letters\n" +
	"\t\tspecify list of letters to exclude when generating candidate guesses\n" +
	"  -i included_letters\n" +
	"\t\tspecify explicit list of letters to include when generating candidate guesses\n" +
	"  -h\t\tshow a short usage message and exit\n" +
	"  -d\t\tprint out this description and exit"

const short_usage string = "usage: " + program_name + " [-h] [-d] [-e excluded_letters | -i included_letters] template"

func long_usage_message() {
	fmt.Fprintln(os.Stderr, long_usage)
}

func short_usage_message() {
	fmt.Fprintln(os.Stderr, short_usage)
}

const blank_char string = "_"
const change_char string = "."

func insert_newlines(s string, max_length int) string {
	return insert_newlines_with_prefix(s, max_length, "")
}

func insert_newlines_with_prefix(s string, max_length int, prefix string) string {
	var buffer bytes.Buffer
	var last_newline_index int

	for i, rune := range s {
		if unicode.IsSpace(rune) && i-last_newline_index > max_length {
			buffer.WriteString("\n")
			buffer.WriteString(prefix)
			last_newline_index = i
		} else {
			buffer.WriteRune(rune)
		}
	}

	return buffer.String()
}

var all_letters = mapset.NewThreadUnsafeSet()

func init() {
	for letter := 'A'; letter <= 'Z'; letter++ {
		all_letters.Add(byte(letter))
	}
}

func make_letter_set(letters string) mapset.Set {
	inc_letters := mapset.NewThreadUnsafeSet()
	for i := 0; i < len(letters); i++ {
		inc_letters.Add(letters[i])
	}
	return inc_letters
}

func list_guesses(prefix, suffix string, included_letters, excluded_letters mapset.Set) ([]string, error) {

	var letters mapset.Set
	if included_letters.Cardinality() > 0 {
		letters = included_letters
	} else {
		letters = all_letters.Difference(excluded_letters)
	}

	byte_letters := make([]byte, 0, letters.Cardinality())
	for letter := range letters.Iter() {
		byte_letter, ok := letter.(byte)
		if !ok {
			return nil, errors.New("type mismatch")
		}

		byte_letters = append(byte_letters, byte(byte_letter))
	}
	slices.Sort(byte_letters)

	guesses := make([]string, 0, len(byte_letters))
	for _, letter := range byte_letters {
		guess := prefix + string(letter) + suffix
		guesses = append(guesses, guess)
	}
	return guesses, nil
}

func case_strings(strs []string) []string {
	blank_char_byte := blank_char[0]
	result := make([]string, 0, len(strs))
	for _, str := range strs {
		str = strings.ToLower(str)
		if str[0] != blank_char_byte {
			upper := strings.ToUpper(string(str[0]))
			str = upper + str[1:]
		}
		result = append(result, str)
	}
	return result
}

func print_guesses(guesses []string, guesses_per_line int) {
	for i, guess := range guesses {
		if i > 0 {
			if i%guesses_per_line == 0 {
				fmt.Println()
			} else {
				fmt.Printf("\t")
			}
		}
		fmt.Print(guess)
	}

	if len(guesses) > 0 {
		fmt.Println()
	}
}

func main() {
	handler_options := slog.HandlerOptions{AddSource: true, Level: slog.ErrorLevel}
	logger := slog.New(handler_options.NewTextHandler(os.Stderr))

	if len(os.Args) == 1 {
		long_usage_message()
		os.Exit(0)
	}

	var included_letters_arg string
	var excluded_letters_arg string
	var description_arg bool

	flag.StringVar(&included_letters_arg, "i", "",
		"specify list of letters to include when generating candidate guesses")
	flag.StringVar(&excluded_letters_arg, "e", "",
		"specify list of letters to exclude when generating candidate guesses")
	flag.BoolVar(&description_arg, "d", false, "output a long description")

	flag.Usage = short_usage_message

	flag.Parse()
	logger.Info("", "included_letters_arg", included_letters_arg)
	logger.Info("", "excluded_letters_arg", excluded_letters_arg)
	logger.Info("", "description_arg", description_arg)

	if description_arg {
		long_usage_message()
		os.Exit(0)
	}

	if len(included_letters_arg) > 0 && len(excluded_letters_arg) > 0 {
		fmt.Fprintln(os.Stderr, "Error: cannot specify both -e and -i")
		os.Exit(1)
	}

	var remaining_args []string = flag.Args()
	if len(remaining_args) != 1 {
		short_usage_message()
		os.Exit(1)
	}

	template := remaining_args[0]
	if len(template) != 5 {
		fmt.Fprintln(os.Stderr, "Error: template is not 5 letters")
		os.Exit(1)
	}
	logger.Info("", "template", template)

	split_re := regexp.MustCompile("\\" + change_char)
	parts := split_re.Split(template, -1)

	num_parts := len(parts)
	if num_parts != 2 {
		fmt.Fprintf(os.Stderr, "Error: template must have one (and only one) '%s' character\n", change_char)
		os.Exit(1)
	}

	prefix := strings.ToUpper(parts[0])
	suffix := strings.ToUpper(parts[1])
	logger.Info("", "prefix", prefix)
	logger.Info("", "suffix", suffix)

	included_letters_arg = strings.ToUpper(included_letters_arg)
	excluded_letters_arg = strings.ToUpper(excluded_letters_arg)

	inc_letters := make_letter_set(included_letters_arg)
	exc_letters := make_letter_set(excluded_letters_arg)

	// fmt.Print("inc_letters: ")
	// inc_letters.Printf("'%c'")
	// fmt.Println()
	// fmt.Print("exc_letters: ")
	// exc_letters.Printf("'%c'")
	// fmt.Println()

	guesses, e := list_guesses(prefix, suffix, inc_letters, exc_letters)
	if e != nil {
		fmt.Fprintf(os.Stderr, "list_guesses Error: %v\n", e)
		os.Exit(1)
	}

	guesses = case_strings(guesses)
	print_guesses(guesses, 5)
}
