package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	args := os.Args[1:]

	if len(args) != 2 {
		fmt.Println("Two files accepted, ex: sample.txt result.txt")
		return
	}

	if !fileExists(args[0]) || !fileExists(args[1]) {
		fmt.Println("Files should exist!")
		return
	}

	if !fileIsTxt(args[0]) || !fileIsTxt(args[1]) {
		fmt.Println("Files should be of .txt ext!")
		return
	}

	inputByte, err := os.ReadFile(args[0])
	if err != nil {
		fmt.Println(err)
		return
	}

	strInput := string(inputByte[:])

	strInput = FormatText(strInput)

	inputByte = []byte(strInput)

	err = os.WriteFile(args[1], inputByte, 0)

	if err != nil {
		fmt.Println(err)
		return
	}

	for i := 0; i < len(args); i++ {
		fmt.Println(args[i])
	}
}

func fileIsTxt(filename string) bool {
	extension := filename[len(filename)-4:]
	if extension == ".txt" {
		return true
	}
	return false
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func FormatText(s string) string {
	s = hex(s)
	s = bin(s)
	s = CaseAllCommand(s)
	s = articleACorrect(s)
	s = punctuationCorrect(s)
	s = quotationsCorrect(s)
	return s
}

func correctQuotationsMatch(s string, match []int) (string, int) {
	// Extract the sentence between quotes
	sentenceUnformatted := s[match[2]:match[3]]
	sentenceFormatted := strings.TrimSpace(sentenceUnformatted)

	// Calculate the number of spaces removed
	deletedSpaces := len(sentenceUnformatted) - len(sentenceFormatted)

	// Replace the original match with the corrected sentence
	s = s[:match[2]] + sentenceFormatted + s[match[3]:]

	return s, deletedSpaces
}

func quotationsCorrect(s string) string {
	appostReplacement := "zAdWtf6wqT"

	patternAppost := `(\s*)(')(\s*(t|m|ll|ve|re|s|d)\s+)`

	compAppost := regexp.MustCompile(patternAppost)

	s = compAppost.ReplaceAllStringFunc(s, func(match string) string {
		matchAppostComp := regexp.MustCompile(`'`)
		result := matchAppostComp.ReplaceAllString(match, appostReplacement)
		return result
	})

	pattern := `'([^']*)'`
	comp := regexp.MustCompile(pattern)

	matches := comp.FindAllStringSubmatchIndex(s, -1)

	shiftedLeft := 0

	for _, match := range matches {

		// Update the indices after removing spaces

		for j := 0; j < len(match); j += 1 {
			match[j] -= shiftedLeft
		}

		// Correct the match
		spacesRemoved := 0
		s, spacesRemoved = correctQuotationsMatch(s, match)
		shiftedLeft += spacesRemoved

	}

	compAppost = regexp.MustCompile(appostReplacement)

	s = compAppost.ReplaceAllString(s, "'")

	return s
}

func punctuationCorrect(s string) string {
	pattern := `\s+(\.{3}|!\?|[\?!.,:;])(\s*[a-zA-Z]*)`
	comp := regexp.MustCompile(pattern)

	match := comp.FindStringSubmatchIndex(s)

	for len(match) > 1 {

		s = correctPunctuationMatch(s, match)

		match = comp.FindStringSubmatchIndex(s)
	}

	patternBetweenWords := `([a-zA-Z]+)(\.{3}|!\?|[\?!.,:;])([a-zA-Z]+)`
	comp = regexp.MustCompile(patternBetweenWords)

	match = comp.FindStringSubmatchIndex(s)

	for len(match) > 1 {

		s = correctPunctuationMatchBetweenWords(s, match)

		match = comp.FindStringSubmatchIndex(s)
	}

	return s
}

func correctPunctuationMatchBetweenWords(s string, match []int) string {
	strPunctuation := s[match[4]:match[5]]
	strBeforePunctuation := s[:match[3]]
	separator := " "

	strAfterPunctuation := s[match[6]:]
	s = strBeforePunctuation + strPunctuation + separator + strAfterPunctuation

	return s
}

func correctPunctuationMatch(s string, match []int) string {
	strPunctuation := s[match[2]:match[3]]
	strBeforePunctuation := s[:match[0]]
	separator := ""

	if match[3] > len(s)-1 { // check if after shifting charachter after punctuation is the end
	} else if string(s[match[3]]) != " " {
		separator = " "
	}

	strAfterPunctuation := s[match[4]:]
	s = strBeforePunctuation + strPunctuation + separator + strAfterPunctuation

	return s
}

func articleACorrect(s string) string {
	patternA := `\b([aA])(\s+[\?!.,:;\(]*\s*[AEIOUHaeiouh][a-zA-Z]*)` // to include not only words but also sing charachters change + to *
	comp := regexp.MustCompile(patternA)

	match := comp.FindStringSubmatchIndex(s)

	for len(match) != 0 {

		s = correctArticleMatch(s, match)

		match = comp.FindStringSubmatchIndex(s)

	}

	patternAn := `\b([aA][nN])(\s+[\?!.,:;\(]*\s*[bcdfghjklmnpqrstvwxyzBCDFGHJKLMNPQRSTVWXYZ][a-zA-Z]*)` // to include not only words but also sing charachters change + to *
	comp = regexp.MustCompile(patternAn)

	match = comp.FindStringSubmatchIndex(s)

	for len(match) != 0 {

		s = correctArticleMatch(s, match)
		match = comp.FindStringSubmatchIndex(s)

	}

	return s
}

func correctArticleMatch(s string, match []int) string {
	strBeforeArticle := s[:match[4]-1]

	sizeArticle := match[3] - match[2] - 1

	article := s[match[2] : match[3]+sizeArticle]

	strCorrectArticle := ""

	switch article {
	case "A":
		strCorrectArticle = "An"
	case "a":
		strCorrectArticle = "an"
	case "An":
		strCorrectArticle = "A"
	case "AN":
		strCorrectArticle = "A"
	case "an":
		strCorrectArticle = "a"
	}

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
		caseTo = strings.Title
	}
	return caseTo
}

func CaseAllCommand(s string) string {
	patternUpMultipule := `\((up|low|cap)(, (\d{1,20}))?\)`
	compPatUpMultipule := regexp.MustCompile(patternUpMultipule)

	countUpMultipule := len(compPatUpMultipule.FindAllString(s, -1))

	for i := 0; i < countUpMultipule; i++ {
		match := compPatUpMultipule.FindStringSubmatchIndex(s)
		commandName := s[match[2]:match[3]]
		toCase := getCaseFunction(commandName)

		if match[4] == -1 { // doenst have a number

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

	patternToDelete := `\((up|low|cap)(, ([-+]*\d+))?\)`
	compPatToDelete := regexp.MustCompile(patternToDelete)
	s = compPatToDelete.ReplaceAllString(s, "")

	return s
}

func toCaseMatch(matches []int, s string, toCase func(string) string, n int) string {
	if n == 0 {
		return s
	}

	strBeforeCommand := s[:matches[0]]
	strAfterCommand := s[matches[1]:]

	pattern := `[\p{L}\p{M}\d]+` // matches words that may or may not have numbers in it
	compWords := regexp.MustCompile(pattern)
	matchesWordsBefore := compWords.FindAllStringSubmatchIndex(strBeforeCommand, -1)

	countMatchesWordsBefore := len(matchesWordsBefore)

	if n > countMatchesWordsBefore {
		n = countMatchesWordsBefore
	}

	firstIndexOfFirstWord := matchesWordsBefore[countMatchesWordsBefore-n][0]
	lastIndexOfLastWord := matchesWordsBefore[countMatchesWordsBefore-1][1]

	strToChange := s[firstIndexOfFirstWord:lastIndexOfLastWord]
	strToCase := toCase(strings.ToLower(strToChange))

	strBeforeCommandAndLastWord := s[lastIndexOfLastWord:matches[0]]

	if strBeforeCommandAndLastWord == "" {
		strBeforeCommandAndLastWord = " "
	}

	s = s[:firstIndexOfFirstWord] + strToCase + strBeforeCommandAndLastWord + strAfterCommand

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
		s = toDecimal(match, s, 2)
	}

	// delete all "(hex)" indicatord

	patternHexes := `\((bin)\)`
	patternHexesRX := regexp.MustCompile(patternHexes)
	s = patternHexesRX.ReplaceAllString(s, "")

	return s
}
