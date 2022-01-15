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

var path2SB string
var repeat int
var seedstr string
var key string

func init() {
	flag.IntVar(&repeat, "r", 1, "repeat generation")
	flag.StringVar(&path2SB, "sb", "./sb.yaml", "path to generation configuration, default=\"./sb.yaml\"")
	flag.StringVar(&seedstr, "seed", "", "randomization seed for generation")
	flag.StringVar(&key, "k", "", "starting generation key")
	flag.Parse()
}

func main() {
	var f *os.File
	var err error
	if f, err = os.Open(path2SB); err != nil {
		panic(err)
	}
	var sb storybuilder.StoryBuilder
	if err = yaml.NewDecoder(f).Decode(&sb); err != nil {
		panic(err)
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
