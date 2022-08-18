package main

import (
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"local/storybuilder"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

var repeat int
var seedstr string
var key string
var printSeed bool

var paths = make(pathFlag)

func init() {
	flag.IntVar(&repeat, "r", 1, "repeat generation")
	flag.Var(
		&paths,
		"sb",
		"path to generation configuration. This flag can be used multiple times. the init of the first config is used. default=\"./sb.yaml\"",
	)
	flag.StringVar(&seedstr, "seed", "", "randomization seed for generation")
	flag.StringVar(&key, "k", "", "starting generation key")
	flag.BoolVar(&printSeed, "printseed", false, "prints seed value")
	flag.Parse()
	if len(paths) == 0 {
		paths["./sb.yaml"] = struct{}{}
	}
}

func main() {
	var sb *storybuilder.StoryBuilder
	var err error
	for path := range paths {
		var f *os.File
		if f, err = os.Open(path); err != nil {
			panic(err)
		}
		if sb == nil {
			sb = new(storybuilder.StoryBuilder)
			if err = yaml.NewDecoder(f).Decode(sb); err != nil {
				panic(fmt.Errorf("%s: %w", path, err))
			}
		} else {
			temp := new(storybuilder.StoryBuilder)
			if err = yaml.NewDecoder(f).Decode(temp); err != nil {
				panic(fmt.Errorf("%s: %w", path, err))
			}
			sb.Combine(temp)
		}
		f.Close()
	}
	if seedstr != "" {
		var seedGen = fnv.New64()
		if _, err := io.WriteString(seedGen, seedstr); err != nil {
			panic(err)
		}
		sb.SetSeed(int64(seedGen.Sum64()))
	}
	for i := 0; i < repeat; i++ {
		builder := &strings.Builder{}
		if key == "" {
			if _, err = sb.WriteTo(builder); err != nil {
				panic(err)
			}
		} else {
			fmt.Println(key)
			if _, err = sb.WriteKey(builder, key); err != nil {
				panic(err)
			}
		}
		fmt.Println(builder.String())
	}
}

type pathFlag map[string]struct{}

func (pf *pathFlag) String() string {
	builder := &strings.Builder{}
	for k := range *pf {
		builder.WriteString(k)
	}
	return builder.String()
}

func (pf *pathFlag) Set(v string) error {
	(*pf)[v] = struct{}{}
	return nil
}
