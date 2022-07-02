package ascii

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"strings"
)

const height = 8

var hash = map[string][]byte{
	"standard":   {0xe1, 0x94, 0xf1, 0x3, 0x34, 0x42, 0x61, 0x7a, 0xb8, 0xa7, 0x8e, 0x1c, 0xa6, 0x3a, 0x20, 0x61, 0xf5, 0xcc, 0x7, 0xa3, 0xf0, 0x5a, 0xc2, 0x26, 0xed, 0x32, 0xeb, 0x9d, 0xfd, 0x22, 0xa6, 0xbf},
	"shadow":     {0x26, 0xb9, 0x4d, 0xb, 0x13, 0x4b, 0x77, 0xe9, 0xfd, 0x23, 0xe0, 0x36, 0xb, 0xfd, 0x81, 0x74, 0xf, 0x80, 0xfb, 0x7f, 0x65, 0x41, 0xd1, 0xd8, 0xc5, 0xd8, 0x5e, 0x73, 0xee, 0x55, 0xf, 0x73},
	"thinkertoy": {0xa5, 0x7b, 0xee, 0xc4, 0x3f, 0xde, 0x67, 0x51, 0xba, 0x1d, 0x30, 0x49, 0x5b, 0x9, 0x26, 0x58, 0xa0, 0x64, 0x45, 0x2f, 0x32, 0x1e, 0x22, 0x1d, 0x8, 0xc3, 0xac, 0x34, 0xa9, 0xdc, 0x12, 0x94},
}

func Art(args ...string) (string, error) {
	fonts := []string{"standard", "thinkertoy", "shadow"}
	var words, filename, currentFont string
	var ascii, fs bool
	for i, arg := range args {
		for index, font := range fonts {
			if (arg == font && i == len(args)-2) || (arg == font && i == len(args)-1) {
				filename = "fonts/" + arg + ".txt"
				fs = true
				currentFont = arg
			} else if i == len(args)-1 && index == len(fonts)-1 && !fs {
				words += arg
				filename = "fonts/standard.txt"
				ascii = true
			}
		}
		if !ascii && !fs {
			words += arg
			if i != len(args)-1 {
				words += " "
			}
		}
	}
	data, err := ioutil.ReadFile(filename)
	h := sha256.New()
	h.Write(data)
	if err != nil {
		return "", err
	} else if checkInputfont(hash[currentFont], h.Sum(nil)) {
		return "", fmt.Errorf("file %s.txt is corrupted", currentFont)
	}

	arrData := strings.Split(string(data), "\n")
	m := make(map[rune][]string)
	space := ' '
	first := 1
	for index, line := range arrData {
		if index >= first && index <= first+height {
			m[space] = append(m[space], line)
			if index == first+height {
				space++
				first += height + 1
			}
		}
	}
	var res string
	for _, set := range strings.Split(words, "\n") {
		if len(set) == 0 {
			res += "\r\n"
			continue
		}
		for line := 0; line < height; line++ {
			for _, r := range set {
				if ascii || fs {
					res += printLine(m[r][line])
				}
			}
			if ascii || fs {
				res += "\n"
			}
		}
	}
	return res, err
}

func printLine(mapStr string) string {
	var str string
	for _, symbol := range mapStr {
		str += string(symbol)
	}
	return str
}

func checkInputfont(arr1, arr2 []byte) bool {
	for i := range arr1 {
		if arr1[i] != arr2[i] {
			return true
		}
	}
	return false
}
