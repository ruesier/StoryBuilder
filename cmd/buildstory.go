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

var paths2SB string
var repeat int
var seedstr string
var key string

var paths []string

func init() {
	flag.IntVar(&repeat, "r", 1, "repeat generation")
	flag.StringVar(
		&paths2SB, 
		"sb", 
		"./sb.yaml", 
		"paths to generation configuration, colon ':' separated. the init of the first config is used, default=\"./sb.yaml\"",
	)
	flag.StringVar(&seedstr, "seed", "", "randomization seed for generation")
	flag.StringVar(&key, "k", "", "starting generation key")
	flag.Parse()
	paths = strings.Split(paths2SB, ":")
}

func main() {
	var sb *storybuilder.StoryBuilder
	var err error
	for _, path := range paths {
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
