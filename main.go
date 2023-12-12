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
	strInput = CaseCommandInvoker(strInput, "up")
	strInput = CaseCommandInvoker(strInput, "low")
	strInput = CaseCommandInvoker(strInput, "cap")
	strInput = articleACorrect(strInput)
	strInput = punctuationCorrect(strInput)
	strInput = quotationsCorrect(strInput)

	// reUp := regexp.MustCompile(` (\((up|low|bin|hex|cap)\)) | (\((up|low|bin|hex|cap), \d{1,8}\))`)
	// strInput = reUp.ReplaceAllString(strInput, "")

	inputByte = []byte(strInput)

	err = os.WriteFile(args[1], inputByte, 0)
	if err != nil {
		fmt.Println(err)
	}
}

func correctQuotationsMatch(s string, match []int) string {
	// wholeMatchStr := s[match[0]:match[1]] // all match

	// spaceBeforeSentenceStr := s[match[2]:match[3]] // space before sentence

	sentenceStr := s[match[4]:match[5]] // sentence

	// spaceAfterSentenceStr :=  s[match[6]:match[7]] // space after sentence

	strBeforeMatch := s[:match[0]]
	strAfterMatch := s[match[1]:]

	s = strBeforeMatch + "'" + sentenceStr + "'" + strAfterMatch

	return s
}

func quotationsCorrect(s string) string {
	pattern := `'(\s+)([^\s]+.*[^\s]+)(\s+)'`
	comp := regexp.MustCompile(pattern)

	countAllQuotations := len(comp.FindAllString(s, -1))

	for i := 0; i < countAllQuotations; i++ {
		match := comp.FindStringSubmatchIndex(s)

		s = correctQuotationsMatch(s, match)

	}

	// fmt.Println(match)

	// fmt.Println("1", s[match[0]:match[1]]) // all match
	// fmt.Println("2", s[match[2]:match[3]]) // space before sentence
	// fmt.Println("3", s[match[4]:match[5]]) // sentence
	// fmt.Println("4", s[match[6]:match[7]]) // space after sentence
	return s
}

func punctuationCorrect(s string) string {
	pattern := `\s+(\.{3}|!\?|[\?!.,:;])(\s*[a-zA-Z]*)`
	comp := regexp.MustCompile(pattern)
	countMisplacedPunct := len(comp.FindAllString(s, -1))

	for i := 0; i < countMisplacedPunct; i++ {
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

	if string(s[match[3]]) != " " {
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

func CaseCommandInvoker(s string, commandName string) string {
	newS := s
	switch commandName {
	case "up":
		newS = CaseCommand(s, commandName, strings.ToUpper)
	case "low":
		newS = CaseCommand(s, commandName, strings.ToLower)
	case "cap":
		caser := cases.Title(language.English)
		newS = CaseCommand(s, commandName, caser.String)
	}
	return newS
}

func CaseCommand(s string, commandName string, toCase func(string) string) string {
	patternUpMultipule := fmt.Sprintf(`(\(%s, \d{1,8}\))|(\(%s\))`, commandName, commandName)
	compPatUpMultipule := regexp.MustCompile(patternUpMultipule)
	countUpMultipule := len(compPatUpMultipule.FindAllString(s, -1))
	commandNameLen := len(commandName)

	for i := 0; i < countUpMultipule; i++ {
		match := compPatUpMultipule.FindStringSubmatchIndex(s)

		if match[1]-match[0] == 2+commandNameLen {
			s = toCaseMatch(match, s, toCase, 1)
		} else {

			numberToUpStr := s[match[0]+commandNameLen+3 : match[1]-1]
			numberToUpInt, err := strconv.Atoi(numberToUpStr)
			if err != nil {
				fmt.Println(err)
			}

			s = toCaseMatch(match, s, toCase, numberToUpInt)
		}
	}

	return s
}

func toCaseMatch(matches []int, s string, toCase func(string) string, n int) string {
	strBeforeCommand := s[:matches[0]]
	strAfterCommand := s[matches[1]:]
	pattern := `\b((([a-zA-Z]+)))`
	compWords := regexp.MustCompile(pattern)
	matchesWordsBefore := compWords.FindAllStringSubmatchIndex(strBeforeCommand, -1)
	changedWordsStr := ""

	countMatchesWordsBefore := len(matchesWordsBefore)

	if n > countMatchesWordsBefore {
		n = countMatchesWordsBefore
	}

	for i := len(matchesWordsBefore) - 1; i > 0 && i > len(matchesWordsBefore)-1-n; i-- {
		wordLen := len(matchesWordsBefore[i])
		wordToChange := s[matchesWordsBefore[i][0]:matchesWordsBefore[i][wordLen-1]]
		wordToCase := toCase(wordToChange)

		strBeforePrevWord := s[matchesWordsBefore[i-1][len(matchesWordsBefore[i-1])-1]:matchesWordsBefore[i][0]]
		changedWordsStr = strBeforePrevWord + wordToCase + changedWordsStr
	}

	indexFirstWordToChange := countMatchesWordsBefore - 1 - n

	if indexFirstWordToChange <= 0 {
		missedFirstWord := s[matchesWordsBefore[0][0]:matchesWordsBefore[0][1]]
		missedFirstWordToCase := toCase(missedFirstWord)
		s = s[:matchesWordsBefore[0][0]] + missedFirstWordToCase + changedWordsStr + strAfterCommand
		return s
	}

	lenOfFirstWordToChange := len(matchesWordsBefore[indexFirstWordToChange])

	strBeforeChangedWords := s[:matchesWordsBefore[indexFirstWordToChange][lenOfFirstWordToChange-1]]
	// lastChangedWordIndexs := matchesWordsBefore[len(matchesWordsBefore)-1][len(matchesWordsBefore[len(matchesWordsBefore)-1])-1]

	s = strBeforeChangedWords + changedWordsStr + strAfterCommand
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
	pattern := `\b([0-9A-Fa-f]+)\s*[\?!.,:;\(]*\s*\((bin)\)`
	compPat := regexp.MustCompile(pattern)
	countHex := len(compPat.FindAllString(s, -1))

	for i := 0; i < countHex; i++ {
		match := compPat.FindStringSubmatchIndex(s)
		// fmt.Println(match)
		s = toDecimal(match, s, 2)
	}

	// delete all "(hex)" indicatord

	patternHexes := `\((hex)\)`
	patternHexesRX := regexp.MustCompile(patternHexes)
	s = patternHexesRX.ReplaceAllString(s, "")

	return s
}
