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

var paths = pathFlag{set: make(map[string]struct{})}

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
	if len(paths.set) == 0 {
		paths.first = "./sb.yaml"
		paths.set["./sb.yaml"] = struct{}{}
	}
}

func main() {
	sb := new(storybuilder.StoryBuilder)
	var err error
	{
		path := paths.first
		var f *os.File
		if f, err = os.Open(path); err != nil {
			panic(err)
		}
		if err = yaml.NewDecoder(f).Decode(sb); err != nil {
			panic(fmt.Errorf("%s: %w", path, err))
		}
		f.Close()
	}
	for path := range paths.set {
		if path == paths.first {
			continue
		}
		var f *os.File
		if f, err = os.Open(path); err != nil {
			panic(err)
		} 
		temp := new(storybuilder.StoryBuilder)
		if err = yaml.NewDecoder(f).Decode(temp); err != nil {
			panic(fmt.Errorf("%s: %w", path, err))
		}
		sb.Combine(temp)
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

type pathFlag struct {
	first string
	set map[string]struct{}
}

func (pf *pathFlag) String() string {
	builder := &strings.Builder{}
	for k := range (*pf).set {
		builder.WriteString(k)
	}
	return builder.String()
}

func (pf *pathFlag) Set(v string) error {
	if len(pf.first) == 0 {
		pf.first = v
	}
	(*pf).set[v] = struct{}{}
	return nil
}
