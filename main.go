package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	// "honnef.co/go/tools/pattern"
	// "utils"
)

func main() {

	args := os.Args[1:]
	for i := 0; i < len(args); i++ {
		fmt.Println(args[i])
	}
	inputByte, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
	}

	strInput := string(inputByte[:])

	strInput = hex(strInput)
	strInput = bin(strInput)

	strInput = CaseAllCommand(strInput)

	strInput = articleACorrect(strInput)
	strInput = punctuationCorrect(strInput)
	strInput = quotationsCorrect(strInput)

	inputByte = []byte(strInput)

	err = os.WriteFile(args[1], inputByte, 0)
	if err != nil {
		fmt.Println(err)
	}
}

func correctQuotationsMatch(s string, match []int) (string, int) {
	// Extract the sentence between quotes
	sentenceUnformatted := s[match[2]:match[3]]
	sentenceFormatted := strings.TrimSpace(sentenceUnformatted)

	// Calculate the number of spaces removed
	deletedSpaces := len(sentenceUnformatted) - len(sentenceFormatted)
	fmt.Println("deletedSpaces", deletedSpaces)

	// Replace the original match with the corrected sentence
	s = s[:match[2]] + sentenceFormatted + s[match[3]:]

	return s, deletedSpaces
}

func quotationsCorrect(s string) string {
	pattern := `'([^']*)'`
	comp := regexp.MustCompile(pattern)

	matches := comp.FindAllStringSubmatchIndex(s, -1)

	shiftedLeft := 0

	for _, match := range matches {

		// Update the indices after removing spaces
		fmt.Println("match", match)

		for j := 0; j < len(match); j += 1 {
			match[j] -= shiftedLeft
		}

		// Correct the match
		fmt.Println("match", match)
		spacesRemoved := 0
		s, spacesRemoved = correctQuotationsMatch(s, match)
		shiftedLeft += spacesRemoved

	}

	return s
}

func punctuationCorrect(s string) string {
	pattern := `\s+(\.{3}|!\?|[\?!.,:;])(\s*[a-zA-Z]*)`
	comp := regexp.MustCompile(pattern)
	countMisplacedPunct := len(comp.FindAllString(s, -1))

	for i := 1; i <= countMisplacedPunct; i++ {
		match := comp.FindStringSubmatchIndex(s)

		s = correctPunctuationMatch(s, match)
		// fmt.Println(correctedS)
	}
	return s
}

func correctPunctuationMatch(s string, match []int) string {
	strPunctuation := s[match[2]:match[3]]
	strBeforePunctuation := s[:match[0]]
	separator := ""

	// fmt.Println(string(s[match[3]]))
	if match[3] > len(s)-1 { // check if after shifting charachter after punctuation is the end
	} else if string(s[match[3]]) != " " {
		separator = " "
	}

	strAfterPunctuation := s[match[4]:]
	s = strBeforePunctuation + strPunctuation + separator + strAfterPunctuation

	return s
}

func articleACorrect(s string) string {
	pattern := `\b\s+([a])(\s+[\?!.,:;\(]*\s*[aeiouh][a-zA-Z]+)`
	comp := regexp.MustCompile(pattern)
	countIncorrectArticleMatches := len(comp.FindAllString(s, -1))

	for i := 0; i < countIncorrectArticleMatches; i++ {
		match := comp.FindStringSubmatchIndex(s)

		s = correctArticleMatch(s, match)

	}
	return s
}

func correctArticleMatch(s string, match []int) string {
	strBeforeArticle := s[:match[4]-1]
	strCorrectArticle := "an"
	strAfterArticle := s[match[4]:]
	s = strBeforeArticle + strCorrectArticle + strAfterArticle
	return s
}

func getCaseFunction(commandName string) func(string) string {
	caseTo := strings.ToUpper
	switch commandName {
	case "up":
		caseTo = strings.ToUpper
	case "low":
		caseTo = strings.ToLower
	case "cap":
		caser := cases.Title(language.English)
		caseTo = caser.String
	}
	return caseTo
}

func CaseAllCommand(s string) string {
	patternUpMultipule := `\((up|low|cap)(, (\d{1,8}))?\)`
	compPatUpMultipule := regexp.MustCompile(patternUpMultipule)

	countUpMultipule := len(compPatUpMultipule.FindAllString(s, -1))

	for i := 0; i < countUpMultipule; i++ {
		match := compPatUpMultipule.FindStringSubmatchIndex(s)
		commandName := s[match[2]:match[3]]
		toCase := getCaseFunction(commandName)

		if match[4] == -1 { // doenst have a number

			// run FuncTOCaseMatchOne - ex (cap)
			s = toCaseMatch(match, s, toCase, 1)
		} else { // has a number

			strNum := s[match[6]:match[7]]
			strNumInt, err := strconv.Atoi(strNum)
			if err != nil {
				fmt.Println(err)
			}
			s = toCaseMatch(match, s, toCase, strNumInt)
		}
	}

	return s
}

func toCaseMatch(matches []int, s string, toCase func(string) string, n int) string {
	strBeforeCommand := s[:matches[0]-1]
	strAfterCommand := s[matches[1]+1:]

	fmt.Println("strBeforeCommand", strBeforeCommand)
	fmt.Println("strAfterCommand", strAfterCommand)

	pattern := `[\p{L}\p{M}\d]+` // matches words that may or may not have numbers in it
	compWords := regexp.MustCompile(pattern)
	matchesWordsBefore := compWords.FindAllStringSubmatchIndex(strBeforeCommand, -1)
	changedWordsStr := ""

	countMatchesWordsBefore := len(matchesWordsBefore)

	if countMatchesWordsBefore == 0 {
		s = strBeforeCommand + strAfterCommand
		return s
	}

	if n > countMatchesWordsBefore {
		n = countMatchesWordsBefore
	}

	for i := len(matchesWordsBefore) - 1; i > 0 && i > len(matchesWordsBefore)-1-n; i-- {

		wordLen := len(matchesWordsBefore[i])
		wordToChange := s[matchesWordsBefore[i][0]:matchesWordsBefore[i][wordLen-1]]

		wordToCase := toCase(wordToChange)

		strBeforePrevWord := s[matchesWordsBefore[i-1][len(matchesWordsBefore[i-1])-1]:matchesWordsBefore[i][0]]

		// check if word is fully numeric

		// if IsNumeric(wordToChange) {
		// 	wordToCase = wordToChange
		// }

		changedWordsStr = strBeforePrevWord + wordToCase + changedWordsStr
	}

	indexFirstWordToChange := countMatchesWordsBefore - 1 - n
	strBeforeLastWordAndCommand := s[matchesWordsBefore[countMatchesWordsBefore-1][1]:matches[0]]

	if indexFirstWordToChange < 0 {
		missedFirstWord := s[matchesWordsBefore[0][0]:matchesWordsBefore[0][1]]
		missedFirstWordToCase := toCase(missedFirstWord)
		s = s[:matchesWordsBefore[0][0]] + missedFirstWordToCase + changedWordsStr + strBeforeLastWordAndCommand + strAfterCommand
		return s
	} else if indexFirstWordToChange == 0 {
		s = s[:matchesWordsBefore[indexFirstWordToChange][0]] + changedWordsStr + strBeforeLastWordAndCommand + strAfterCommand
		return s
	}

	lenOfFirstWordToChange := len(matchesWordsBefore[indexFirstWordToChange])

	strBeforeChangedWords := s[:matchesWordsBefore[indexFirstWordToChange][lenOfFirstWordToChange-1]]
	// lastChangedWordIndexs := matchesWordsBefore[len(matchesWordsBefore)-1][len(matchesWordsBefore[len(matchesWordsBefore)-1])-1]

	s = strBeforeChangedWords + changedWordsStr + strBeforeLastWordAndCommand + strAfterCommand
	return s
}

func toDecimal(matches []int, s string, base int) string {
	// [0]...[1] indexed of the whole match
	// [2]...[3] indexed of hex value to transform match
	// [4]...[5] indexed of (hex) match

	valueStrToHex := s[matches[2]:matches[3]]

	decimalValueInt, err := strconv.ParseInt(valueStrToHex, base, 64)
	if err != nil {
		fmt.Println(err)
	}

	decimalValueStr := strconv.FormatInt(decimalValueInt, 10)

	sBeforeHexValue := s[:matches[2]]
	s = sBeforeHexValue + decimalValueStr + s[matches[5]+1:]

	return s
}

func hex(s string) string {
	pattern := `\b([0-9A-Fa-f]+)\s*[\?!.,:;\(]*\s*\((hex)\)`
	compPat := regexp.MustCompile(pattern)
	countHex := len(compPat.FindAllString(s, -1))

	for i := 0; i < countHex; i++ {
		match := compPat.FindStringSubmatchIndex(s)
		// fmt.Println(match)
		s = toDecimal(match, s, 16)
	}

	// delete all "(hex)" indicatord

	patternHexes := `\((hex)\)`
	patternHexesRX := regexp.MustCompile(patternHexes)
	s = patternHexesRX.ReplaceAllString(s, "")

	return s
}

func bin(s string) string {
	pattern := `\b([01]+)\s*[\?!.,:;\(]*\s*\((bin)\)`
	compPat := regexp.MustCompile(pattern)
	countHex := len(compPat.FindAllString(s, -1))

	for i := 0; i < countHex; i++ {
		match := compPat.FindStringSubmatchIndex(s)
		// fmt.Println(match)
		s = toDecimal(match, s, 2)
	}

	// delete all "(hex)" indicatord

	patternHexes := `\((bin)\)`
	patternHexesRX := regexp.MustCompile(patternHexes)
	s = patternHexesRX.ReplaceAllString(s, "")

	return s
}

// => "30 files were added"
// => "It has been 2 years"
// => "Ready, set, GO !"
// => "I should stop shouting"
// => "Welcome to the Brooklyn Bridge"
// => "This is SO EXCITING"
// => "I was sitting over there, and then BAMM!!"
// => "I was thinking... You were right"
// => "I am exactly how they describe me: 'awesome'"
// => "As Elton John said: 'I am the most well-known homosexual in the world'"
// => "There it was. An amazing rock!"
