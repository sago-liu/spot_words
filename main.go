package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

//词库类型参数
var lexicon = flag.String("l", "", "lexicon: sight/tricky")
//反推验证参数
var reverseVerifyFlag = flag.Bool("r", false, "reverseVerify")

func main() {
	//获取参数
	flag.Parse()
	if *lexicon == "" || (*lexicon != "sight" && *lexicon != "tricky") {
		fmt.Println("./spot_words -l sight or ./spot_words -l tricky")
		return
	}
	//fmt.Println(*lexicon)

	//加载词库文件
	lexiconFileName := fmt.Sprintf("./%s words.csv", *lexicon)
	lexiconWords := loadLexiconCSV(lexiconFileName)
	if (*reverseVerifyFlag) {
		//反推验证
		reverseVerify(lexiconWords)
		return
	}

	//加载输入文件
	InputFileName := "./input.txt"
	inputWords := loadInputTxt(InputFileName)

	//匹配和打印
	spot(inputWords, lexiconWords)
	fmt.Println("===============Author: Donovan Liu===============")
}

/*
加载词库文件
 */
func loadLexiconCSV(fileName string) (words []string) {
	words = make([]string, 0, 500)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("open %s error: %v", fileName, err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("read %s error: %v", fileName, err)
			return
		}

		//fmt.Println(record) // record has the type []string
		words = append(words, strings.TrimSpace(record[0]))
	}

	//fmt.Println(words)
	return
}

/*
加载输入文件
 */
func loadInputTxt(fileName string) (words []string) {
	words = make([]string, 0, 500)

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("open %s error: %v", fileName, err)
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("read %s error: %v", fileName, err)
			return
		}

		//fmt.Println(string(line))
		//从句子中提取单词
		words = append(words, extractWords(string(line))...)
	}

	//fmt.Println(words)
	return
}

/*
从句子中提取单词
 */
func extractWords(sentence string) (words []string) {
	sentence = strings.Replace(sentence, ",", "", -1)
	sentence = strings.Replace(sentence, ".", "", -1)
	sentence = strings.Replace(sentence, "?", "", -1)
	sentence = strings.Replace(sentence, "!", "", -1)
	sentence = strings.Replace(sentence, "\"", "", -1)
	words = strings.Fields(strings.TrimSpace(sentence))
	return
}

/*
匹配和打印
 */
func spot(inputWords []string, lexiconWords []string) {
	hits := make(map[string]bool)

	for _, iW := range(inputWords) {
		for _, lW := range(lexiconWords) {
			if strings.ToLower(iW) == strings.ToLower(lW) {
				hits[lW] = true
			} else {
				//如果匹配不上，还原单词词形再匹配
				basicForms := lemmatize(iW)
				for _, bF := range(basicForms) {
					if strings.ToLower(bF) == strings.ToLower(lW) {
						hits[lW] = true
					}
				}
			}
		}
	}

	//fmt.Println(hits)
	for k, _ := range(hits) {
		fmt.Println(k)
	}
}

/*
还原单词词形
 */
func lemmatize(word string) (basicForms []string) {
	//黑名单
	blackLists := []string{"is", "bees", "hers", "toes", "yours"}
	for _, v := range(blackLists) {
		if strings.ToLower(word) == v {
			return
		}
	}

	//不规则名词表x
	/*名词复数形式->原形
	规则1.  *ves --> *f/*fe
	规则2.  *ies --> *y
	规则3.  *es  --> *
	规则4.  *s   --> *
	 */

	//不规则动词表x
	/*动词第三人称单数->原形
	规则2.  *ies --> *y
	规则3.  *es --> *
	规则4.  *s   --> *
	 */
	//动词现在进行时->原形x
	//动词过去时->原形x
	//动词过去分词->原形x

	basicForms = make([]string, 0, 2)

	if len(word) > 3 && "ves" == word[len(word)-3:] {
		basicForms = append(basicForms, word[0:len(word)-3] + "f")
		basicForms = append(basicForms, word[0:len(word)-3] + "fe")
	}

	if len(word) > 3 && "ies" == word[len(word)-3:] {
		basicForms = append(basicForms, word[0:len(word)-3] + "y")
	}

	if len(word) > 2 && "es" == word[len(word)-2:] {
		blackLists := []string{"bees", "toes", "topes", "heres", "uses", "notes", "sites", "themes"}
		for _, v := range(blackLists) {
			if strings.ToLower(word) != v {
				basicForms = append(basicForms, word[0:len(word)-2])
			}
		}
	}

	if len(word) > 1 && "s" == word[len(word)-1:] {
		blackLists := []string{"is", "hers", "ours", "yours", "ass", "cans", "goods", "hiss", "its", "lives", "news"}
		for _, v := range(blackLists) {
			if strings.ToLower(word) != v {
				basicForms = append(basicForms, word[0:len(word)-1])
			}
		}
	}

	return
}

/*
反推验证
 */
func reverseVerify(lexiconWords []string) {
	var words []string
	for _, v := range(lexiconWords) {
		words = make([]string, 0, 10)
		if len(v) > 1 && "f" == v[len(v)-1:] {
			words = append(words, v[0:len(v)-1] + "ves")
		}
		if len(v) > 2 && "fe" == v[len(v)-2:] {
			words = append(words, v[0:len(v)-2] + "ves")
		}

		if len(v) > 1 && "y" == v[len(v)-1:] {
			words = append(words, v[0:len(v)-1] + "ies")
		}

		words = append(words, v + "es")

		words = append(words, v + "s")

		fmt.Println(v, words)
	}
}

