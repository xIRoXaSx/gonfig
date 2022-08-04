package gonfig

import (
	"math"
	"os"
	"path/filepath"
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestNewGonfig(t *testing.T) {
	t.Parallel()

	type test struct {
		A string `json:"A" yaml:"A"`
		B int    `json:"B" yaml:"B"`
		C any    `json:"C" yaml:"C"`
	}

	type config struct {
		String  string   `json:"String" yaml:"String"`
		Int8    int8     `json:"Int8" yaml:"Int8"`
		Uint8   uint8    `json:"Uint8" yaml:"Uint8"`
		Integer int      `json:"Integer" yaml:"Integer"`
		Float   float32  `json:"Float" yaml:"Float"`
		Slice   []string `json:"Slice" yaml:"Slice"`
		Struct  test     `json:"Struct" yaml:"Struct"`
	}

	v := &config{
		String:  "Hello world!",
		Int8:    -123,
		Uint8:   123,
		Integer: 321,
		Float:   1.23,
		Slice:   []string{"Hello", " ", "world!"},
		Struct:  test{},
	}

	dir := "GonfigTest"
	file := "gonfig"
	gType := GonfJson
	p, err := os.UserConfigDir()
	r.NoError(t, err)
	baseDir := filepath.Join(p, dir)
	r.NoError(t, cleanup(baseDir))
	p = filepath.Join(baseDir, file)

	// Test edge cases upon creation.
	g, err := New("", file, GonfJson, true)
	r.ErrorIs(t, err, ErrExpectedNonNilOrEmpty)
	g, err = New(dir, "", GonfJson, true)
	r.ErrorIs(t, err, ErrExpectedNonNilOrEmpty)
	g, err = New(dir, file, gType, false)
	r.NoError(t, err)
	err = g.WriteToFile(math.Inf(1))
	r.Error(t, err)

	g, err = New(dir, file, gType, false)
	r.NoError(t, err)
	r.Exactly(t, p+jsonExtension, g.FullPath())
	r.Exactly(t, g.FileName(), file+jsonExtension)
	r.Exactly(t, g.Dir(), baseDir)
	r.Exactly(t, g.Type(), gType)

	g, err = New(dir, file, GonfYaml, false)
	r.NoError(t, err)
	r.Exactly(t, p+yamlExtension, g.FullPath())
	r.Exactly(t, g.FileName(), file+yamlExtension)
	r.NoError(t, g.WriteToFile(v))
	r.ErrorIs(t, g.WriteToFile(v), ErrOverwriteDisabled)
	g, err = New(dir, file, GonfJson, false)
	r.NoError(t, g.WriteToFile(v))

	var c config
	r.ErrorIs(t, g.Load(c), ErrMustBeAddressable)
	r.NoError(t, g.Load(&c))
	r.Equal(t, v, &c)

	v.Int8 = -1
	g, err = New(dir, file, GonfJson, true)
	r.NoError(t, g.WriteToFile(v))
	r.NoError(t, g.Load(&c))
	r.Exactly(t, v.Int8, c.Int8)

	g, err = NewWithPath(p+"_1"+jsonExtension, GonfJson, true)
	r.NoError(t, g.WriteToFile(v))
	r.NoError(t, cleanup(baseDir))
}

func cleanup(p string) (err error) {
	_, err = os.Stat(p)
	if err == os.ErrNotExist {
		return
	}
	return os.RemoveAll(p)
}
