package storybuilder

import (
	"bytes"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

var injectionParser *regexp.Regexp

func init() {
	var amount = `(?:(?:[1-9][0-9]*)|(?:[0-9]+-[1-9][0-9]*))`
	var repeatList = `(?P<capital>!)?(?:-(?P<amounts>` + amount + `(?:,` + amount + `)*)(?:'(?P<separator>.*)')?)?`

	injectionParser = regexp.MustCompile(`(?:<(?P<key>[a-zA-Z0-9.]+)` + repeatList + `>)`)
}

type chunk interface {
	io.WriterTo
	String() string
}

func buildChunks(str string, sb StoryBuilder) (chunks []chunk) {
	var remaining = str
	chunks = make([]chunk, 0, 1)
	for len(remaining) > 0 {
		indexes := injectionParser.FindStringSubmatchIndex(remaining)
		if indexes == nil {
			chunks = append(chunks, basicChunk(remaining))
			return
		}
		temp := &randomChunk{
			sb:        sb,
			key:       extractKey(remaining, indexes),
			amounts:   extractAmounts(remaining, indexes),
			separator: extractSeparator(remaining, indexes),
			capital:   indexes[4] > 0,
		}
		chunks = append(chunks, basicChunk(remaining[:indexes[0]]), temp)
		remaining = remaining[indexes[1]:]
	}
	return
}

func extractKey(remaining string, indexes []int) string {
	if indexes[2] < 0 {
		return "_"
	}
	return remaining[indexes[2]:indexes[3]]
}

func extractAmounts(remaining string, indexes []int) (amounts []int) {
	if indexes[6] < 0 {
		return []int{1}
	}
	txt := remaining[indexes[6]:indexes[7]]
	ranges := strings.Split(txt, ",")
	for _, r := range ranges {
		if strings.Contains(r, "-") {
			extremes := strings.Split(r, "-")
			start, e := strconv.Atoi(extremes[0])
			if e != nil {
				panic(e)
			}
			end, e := strconv.Atoi(extremes[1])
			if e != nil {
				panic(e)
			}
			for i := start; i < end; i++ {
				amounts = append(amounts, i)
			}
		} else {
			if amount, e := strconv.Atoi(r); e == nil {
				amounts = append(amounts, amount)
			} else {
				panic(e)
			}
		}
	}
	return
}

func extractSeparator(remaining string, indexes []int) string {
	if indexes[8] < 0 {
		return " "
	}
	return remaining[indexes[8]:indexes[9]]
}

type basicChunk string

func (bc basicChunk) String() string {
	return string(bc)
}

func (bc basicChunk) WriteTo(w io.Writer) (int64, error) {
	i, e := io.WriteString(w, string(bc))
	return int64(i), e
}

type randomChunk struct {
	sb        StoryBuilder
	key       string
	amounts   []int
	separator string
	capital   bool
}

func (rc *randomChunk) String() string {
	sb := new(strings.Builder)
	if _, err := rc.WriteTo(sb); err != nil {
		panic(err)
	}
	return sb.String()
}

func (rc *randomChunk) WriteTo(w io.Writer) (size int64, e error) {
	repeat := rc.amounts[rc.sb.random.Intn(len(rc.amounts))]
	for i := 0; i < repeat; i++ {
		var tempSize int64
		if i > 0 {
			var tempSizei int
			tempSizei, e = io.WriteString(w, rc.separator)
			size += int64(tempSizei)
			if e != nil {
				return
			}
		}
		if i == 0 && rc.capital {
			buffer := new(bytes.Buffer)
			_, e = rc.sb.WriteKey(buffer, rc.key)
			if e != nil {
				return 0, e
			}
			data := buffer.Bytes()
			data[0] = byte(unicode.ToUpper(rune(data[0])))
			tempSizei, err := w.Write(data)
			size += int64(tempSizei)
			if err != nil {
				return size, err
			}
		} else {
			tempSize, e = rc.sb.WriteKey(w, rc.key)
			size += tempSize
			if e != nil {
				return
			}
		}
	}
	return
}
