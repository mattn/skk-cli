package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"

	"github.com/hermanschaaf/kana"
	"github.com/mattn/go-skkdic"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "path to SKK-JISYO"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func defaultDict() string {
	p, err := os.UserConfigDir()
	if err != nil {
		return ""
	}
	return filepath.Join(p, "skk-cli", "SKK-JISYO.L")
}

var (
	reTrim = regexp.MustCompile("[aiueo]$")
)

func split(s string) []string {
	result := []string{}
	rs := []rune(s)
	j := 0
	for i := 0; i < len(rs)-1; i++ {
		if unicode.IsLower(rs[i]) && unicode.IsUpper(rs[i+1]) {
			result = append(result, string(rs[:i+1]))
			j = i + 1
		}
	}
	return append(result, string(rs[j:]))
}

func loadDict(d *skkdic.Dict, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return d.Load(f)
}

func main() {
	var paths arrayFlags
	flag.Var(&paths, "d", "Some description for this param.")
	flag.Parse()

	if len(paths) == 0 {
		paths = []string{defaultDict()}
	}

	dic := skkdic.New()
	for _, p := range paths {
		if err := loadDict(dic, p); err != nil {
			log.Fatal(err)
		}
	}

	enc := json.NewEncoder(os.Stdout)
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		s := scanner.Text()
		words := split(s)
		result := []string{}
		prefix := ""
		for len(words) > 0 {
			rs := []rune(words[0])
			if len(rs) > 0 && unicode.IsUpper(rs[0]) {
				break
			}
			prefix += kana.RomajiToHiragana(words[0])
			words = words[1:]
		}
		if len(words) > 1 {
			// OkuRu => OkuR => おくr
			// v
			// おく + る
			ss := ""
			suf := words[len(words)-1]
			words[len(words)-1] = reTrim.ReplaceAllString(suf, "")
			for _, word := range words {
				ss += kana.RomajiToHiragana(strings.ToLower(word))
			}
			suf = kana.RomajiToHiragana(strings.ToLower(suf))
			for _, e := range dic.SearchOkuriAri(ss) {
				for _, word := range e.Words {
					result = append(result, prefix+word.Text+suf)
				}
			}
		} else if len(words) == 1 {
			ss := kana.RomajiToHiragana(strings.ToLower(s))
			for _, e := range dic.SearchOkuriNasi(ss) {
				for _, word := range e.Words {
					result = append(result, word.Text)
				}
			}
		}
		if len(result) == 0 {
			result = append(result, kana.RomajiToHiragana(strings.ToLower(s)))
		}
		enc.Encode(result)
	}
}
