package goReloaded

import (
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
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
)

/*
--> loop through the words of the file
--> check for specific instructions, wheather (up), or (up, 2)... etc
--> words[i-1] this selects the initial item of the current iteration in the array loop
--> in case there is a number use trim to remove all the characters after the number
--> atoi is for converting the number in string form to integer
--> loop through the number to modify the initial words per the specific number
--> ParseInt converts the hex number to int64
--> output error if there is one, meaning it's not a hex value
--> itoa is for converting the int(int64 to int) number to string
*/
func Instructions() {
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
			words[i-1] = strings.Title(words[i-1])
		} else if word == "(cap," { // when there is a number
			n := strings.Trim(words[i+1], words[i+1][1:])
			nmbr, _ := strconv.Atoi(string(n))
			for j := 1; j <= nmbr; j++ {
				words[i-j] = strings.Title(words[i-j])
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
			if dec_num, hex_err := strconv.ParseInt(hex_num, 16, 64); hex_err != nil {
				log.Fatalf("the value might not be a hex number: %s", hex_err)
			} else {
				hex_num = strconv.Itoa(int(dec_num))
				words[i-1] = hex_num
			}
		}

		// (bin) to decimal
		if word == "(bin)" {
			bin_num := words[i-1]
			if dec_num, bin_err := strconv.ParseInt(bin_num, 2, 64); bin_err != nil {
				log.Fatalf("the value might not be a binary number: %s", bin_err)
			} else {
				bin_num = strconv.Itoa(int(dec_num))
				words[i-1] = bin_num
			}
		}

		// a to an
		for _, letter := range letters {
			if word == "a" && string(words[i+1][0]) == letter {
				words[i] = "an"
			}
		}
	}
}

func FinalTool() {
	// check for file reading errors
	if read_err != nil {
		log.Fatalf("unable to read file: %s", read_err)
		fmt.Println()
	}

	Instructions()

	for _, item := range words {
		new_arr_result = append(new_arr_result, item)
	}

	output_sentence := strings.Join(new_arr_result, " ")

	instructions := `\((up|low|cap|hex|bin)(,\s*\d+)?\)\s*`
	remove_instructions := regexp.MustCompile(instructions)
	output_sentence_1 := remove_instructions.ReplaceAllString(output_sentence, "")

	/*
		--> remove spaces before punctuations and add space after punctuation if none
		--> this metthod introduces a space after a punctuation incase they are groups of punctuations
		--> solved it with the next regex pattern
	*/
	spaces := `(\s+)([.,!?:;])`
	remove_spaces := regexp.MustCompile(spaces)
	output_sentence_2 := remove_spaces.ReplaceAllString(output_sentence_1, "$2 ")

	// remove spaces between punctuations
	spaced_punct := `(\s+)([.,!?:;])`
	remove_spaced_punct := regexp.MustCompile(spaced_punct)
	output_sentence_3 := remove_spaced_punct.ReplaceAllString(output_sentence_2, "$2")

	// remove spaces immediately after and before apostrophe
	space_apostrophe := `(\')(\s*)(.*?)(\s*)(\')`
	remove_apostrophe := regexp.MustCompile(space_apostrophe)
	output_sentence_4 := remove_apostrophe.ReplaceAllString(output_sentence_3, "$1$3$5")

	// removing multiple spaces
	multi_spaces := `(\s)(\s+)`
	remove_multi_spaces := regexp.MustCompile(multi_spaces)
	final_output_sentence := remove_multi_spaces.ReplaceAllString(output_sentence_4, "$1")

	// removed spaces after end of sentence
	final_output_sentence = strings.TrimSpace(final_output_sentence)

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

	// fmt.Println(final_output_sentence)
}
