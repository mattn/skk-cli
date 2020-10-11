package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"unicode"

	"github.com/hermanschaaf/kana"
	"github.com/mattn/go-skkdic"
)

const name = "skk-cli"

const version = "0.0.1"

var revision = "HEAD"

var (
	reTrim   = regexp.MustCompile("[aiueo]$")
	replacer = strings.NewReplacer(
		".", "。",
		",", "、",
	)
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "Path to SKK-JISYO"
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

func roma2hira(s string) string {
	return replacer.Replace(kana.RomajiToHiragana(s))
}

func loadDict(d *skkdic.Dict, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return d.Load(f)
}

type Request struct {
	Method string `json:"method"`
	Text   string `json:"text"`
}

type Response struct {
	Status string      `json:"status"`
	Result interface{} `json:"result"`
}

func main() {
	var jm bool
	var paths arrayFlags
	var showVersion bool
	flag.BoolVar(&jm, "json", false, "JSON mode")
	flag.Var(&paths, "d", "Path to SKK-JISYO.L")
	flag.BoolVar(&showVersion, "V", false, "Print the version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}

	if len(paths) == 0 {
		paths = []string{defaultDict()}
	}

	dic := skkdic.New()
	for _, p := range paths {
		if err := loadDict(dic, p); err != nil {
			log.Fatal(err)
		}
	}

	var enc *json.Encoder
	if jm {
		enc = json.NewEncoder(os.Stdout)
	}
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if enc == nil {
			fmt.Print("> ")
		}
		if !scanner.Scan() {
			break
		}
		var req Request
		if enc != nil {
			s := scanner.Text()
			err := json.Unmarshal([]byte(s), &req)
			if err != nil {
				if err := enc.Encode(&Response{Status: "NG", Result: err.Error()}); err != nil {
					log.Fatal(err)
				}
				continue
			}
		} else {
			s := scanner.Text()
			req.Text = s
		}
		words := split(req.Text)
		result := []string{}
		prefix := ""
		for len(words) > 0 {
			rs := []rune(words[0])
			if len(rs) > 0 && unicode.IsUpper(rs[0]) {
				break
			}
			prefix += roma2hira(words[0])
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
				ss += roma2hira(strings.ToLower(word))
			}
			suf = roma2hira(strings.ToLower(suf))
			for _, e := range dic.SearchOkuriAri(ss) {
				for _, word := range e.Words {
					result = append(result, prefix+word.Text+suf)
				}
			}
		} else if len(words) == 1 {
			ss := roma2hira(strings.ToLower(req.Text))
			for _, e := range dic.SearchOkuriNasi(ss) {
				for _, word := range e.Words {
					result = append(result, word.Text)
				}
			}
		}
		if unicode.IsUpper([]rune(req.Text)[0]) {
			for _, e := range dic.SearchOkuriNasiPrefix(req.Text) {
				for _, word := range e.Words {
					result = append(result, word.Text)
				}
			}
			result = append(result, roma2hira(strings.ToLower(req.Text)))
		} else {
			hira := roma2hira(strings.ToLower(req.Text))
			if strings.IndexFunc(hira, func(r rune) bool { return unicode.IsTitle(r) }) == -1 {
				result = append(result, hira)
			}
			for _, e := range dic.SearchOkuriNasiPrefix(req.Text) {
				for _, word := range e.Words {
					result = append(result, word.Text)
				}
			}
		}
		if enc != nil {
			if err := enc.Encode(&Response{Status: "OK", Result: result}); err != nil {
				log.Fatal(err)
			}
		} else {
			for _, r := range result {
				fmt.Println(r)
			}
		}
	}
}
