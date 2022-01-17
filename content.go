package storybuilder

import (
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"
)

type StoryBuilder struct {
	Module			   string              `yaml:"_module_"`
	Init               string              `yaml:"_init_"`
	Fill               map[string][]string `yaml:"-,inline"`
	random             *rand.Rand
	addedModules	   map[string]bool
}

func (sb *StoryBuilder) init() {
	if sb.random == nil {
		sb.random = rand.New(rand.NewSource(time.Now().UnixNano()))
	}
}

func (sb *StoryBuilder) SetSeed(seed int64) {
	sb.random = rand.New(rand.NewSource(seed))
}

func (sb StoryBuilder) Generate(key string) string {
	build := new(strings.Builder)
	if _, e := sb.WriteKey(build, key); e != nil {
		panic(e)
	}
	return build.String()
}

func (sb *StoryBuilder) WriteTo(w io.Writer) (size int64, err error) {
	sb.init()
	chunks := buildChunks(sb.Init, sb)
	for _, ch := range chunks {
		temp, e := ch.WriteTo(w)
		size += temp
		if e != nil {
			return size, e
		}
	}
	return
}

func (sb *StoryBuilder) WriteKey(w io.Writer, key string) (size int64, err error) {
	sb.init()
	if sb.Fill[key] == nil || len(sb.Fill[key]) < 1 {
		var written int
		written, err = w.Write([]byte{'<'})
		size += int64(written)
		if err != nil {
			return
		}

		written, err = io.WriteString(w, key)
		size += int64(written)
		if err != nil {
			return
		}

		written, err = w.Write([]byte{'>'})
		size += int64(written)
		if err != nil {
			return
		}
	} else {
		selection := sb.Fill[key][sb.random.Intn(len(sb.Fill[key]))]
		chunks := buildChunks(selection, sb)
		for _, ch := range chunks {
			temp, e := ch.WriteTo(w)
			size += temp
			if e != nil {
				return size, e
			}
		}
	}
	return
}

func (sb *StoryBuilder) Combine(other *StoryBuilder) {
	if len(sb.Module) == 0 {
		panic(fmt.Errorf("primary storybuilder is missing '_module_' field and so cannot be combined"))
	}
	if len(other.Module) == 0 {
		panic(fmt.Errorf("failed to combine storybuilder due to missing '_module_' field"))
	}
	if sb.addedModules == nil {
		sb.addedModules = map[string]bool{
			sb.Module: true,
		}
	}
	if added := sb.addedModules[other.Module]; !added {
		sb.addedModules[other.Module] = true
		for key, list := range other.Fill {
			if current, ok := sb.Fill[key]; ok {
				sb.Fill[key] = append(current, list...)
			} else {
				sb.Fill[key] = list
			}
		}
	}
}
