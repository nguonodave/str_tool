package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// initialization of global variables
var (
	sample                  = os.Args[1]
	result                  = os.Args[2]
	file_content, read_err  = os.ReadFile(sample)                                                  // reads and stores the contents of the file in a variable, file_content
	sentence                = string(file_content)                                                 // converts the contents from bytes to a string and store the result to a variable, sentence
	words                   = strings.Split(sentence, " ")                                         // splits the sentence into array of words by using a space as a separator
	new_arr_result          = []string{}                                                           // empty array for storing non-instruction words
	output_file, create_err = os.Create(result)                                                    // creates a file for storing the final result, which is the sentence converted from the array
	letters                 = []string{"a", "e", "i", "o", "u", "A", "E", "I", "O", "U", "H", "h"} // array of vowels and "h"
	punctuations            = []string{".", ",", "!", "?", ":", ";"}                               // array of punctuations
)

/*
--> loop through the words of the file
--> check for specific instructions, wheather (up), or (up, 2)... etc
--> words[i-1] this selects the initial item of the current iteration in the array loop
--> in case there is a number use trim to remove all the characters after the number
--> atoi is for converting the number in string form to integer
--> ParseInt converts the hex number to int64
--> if there is an error, meaning it's not a hex value
--> itoa is for converting the int(int64 to int) number to string
--> for capitalizing use the format in the library golang.org/x/text/cases since strings.Title is depricated
*/
func modifications() {
	for i, word := range words {
		// (up) to uppercase
		// when there is no number
		if word == "(up)" {
			words[i-1] = strings.ToUpper(words[i-1])
		} else if word == "(up," { // when there is a number
			n := strings.Trim(words[i+1], words[i+1][1:])
			nmbr, _ := strconv.Atoi(string(n))
			for j := 1; j <= nmbr; j++ {
				words[i-j] = strings.ToUpper(words[i-j])
			}
		}

		// (cap) to titlecase
		// when there is no number
		if word == "(cap)" {
			words[i-1] = cases.Title(language.English, cases.Compact).String(words[i-1])
		} else if word == "(cap," { // when there is a number
			n := strings.Trim(words[i+1], words[i+1][1:])
			nmbr, _ := strconv.Atoi(string(n))
			for j := 1; j <= nmbr; j++ {
				words[i-j] = cases.Title(language.English, cases.Compact).String(words[i-j])
			}
		}

		// (low) to lowercase
		// when there is no number
		if word == "(low)" {
			words[i-1] = strings.ToLower(words[i-1])
		} else if word == "(low," { // when there is a number
			n := strings.Trim(words[i+1], words[i+1][1:])
			nmbr, _ := strconv.Atoi(n)
			for j := 1; j <= nmbr; j++ {
				words[i-j] = strings.ToLower(words[i-j])
			}
		}

		// (hex) to decimal
		if word == "(hex)" {
			hex_num := words[i-1]
			dec_num, _ := strconv.ParseInt(hex_num, 16, 64)
			hex_num = strconv.Itoa(int(dec_num))
			words[i-1] = hex_num
		}

		// (bin) to decimal
		if word == "(bin)" {
			hex_num := words[i-1]
			dec_num, _ := strconv.ParseInt(hex_num, 2, 64)
			hex_num = strconv.Itoa(int(dec_num))
			words[i-1] = hex_num
		}

		// a to an
		for _, letter := range letters {
			if word == "a" && string(words[i+1][0]) == letter {
				words[i] = "an"
			}
		}
	}
}

/*
--> after all modifications have been done, append words that are not instructions, i.e not "(up), (cap, 4) etc", to an empty array, "new_array_results
--> write the sentence to the created file using io.WriteString
--> defer output_file.Close() is for file handling
*/
func main() {
	// check for file reading errors
	if read_err != nil {
		log.Fatalf("unable to read file: %s", read_err)
		fmt.Println()
	}

	modifications()

	for _, item := range words {
		new_arr_result = append(new_arr_result, item)
	}

	output_sentence := strings.Join(new_arr_result, " ")

	instructions := `\((up|low|cap|hex|bin)(,\s*\d+)?\)\s*`
	remove_instructions := regexp.MustCompile(instructions)
	semi_output_sentence := remove_instructions.ReplaceAllString(output_sentence, "")

	// remove spaces before punctuations
	spaces := `( +)([.,!?:;]{1,3})( *)`
	remove_spaces := regexp.MustCompile(spaces)
	final_output_sentence := remove_spaces.ReplaceAllString(semi_output_sentence, "$2 ")

	// check for file creating errors
	if create_err != nil {
		log.Fatalf("unable to create file: %s", create_err)
		fmt.Println()
	}

	if _, write_err := io.WriteString(output_file, final_output_sentence); write_err != nil {
		log.Fatalf("unable to write to file: %s", read_err)
		fmt.Println()
	}
	defer output_file.Close()

	fmt.Println(final_output_sentence)
}
