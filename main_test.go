package gonfig

import (
	"os"
	"path/filepath"
	"testing"

	r "github.com/stretchr/testify/require"
)

func TestNewGonfig(t *testing.T) {
	type test struct {
		A string
		B int
		C any
	}

	type config struct {
		String  string
		Int8    int8
		Uint8   uint8
		Integer int
		Float   float32
		Slice   []string
		Struct  test
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

	g, err := New(GonfOptions{})
	r.Error(t, err)
	r.Error(t, g.WriteToFile(v))

	dir := "GonfigTest"
	file := "gonfig"
	g, err = New(GonfOptions{
		DirName:    dir,
		ConfigName: file,
		Type:       GonfJson,
	})
	r.NoError(t, err)

	g, err = New(GonfOptions{
		DirName:    dir,
		ConfigName: file,
		Type:       GonfYAML,
	})
	r.NoError(t, err)
	r.NoError(t, g.WriteToFile(v))
	r.Error(t, g.WriteToFile(v))
	yPath := g.FullPath()
	var c config
	r.NoError(t, g.LoadFile(&c))
	r.Equal(t, v, &c)

	g, err = New(GonfOptions{
		DirName:    dir,
		ConfigName: file,
		Type:       GonfJson,
	})
	r.NoError(t, g.WriteToFile(v))
	jPath := g.FullPath()

	v.Int8 = -1
	g.opts.OverwriteExisting = true
	r.NoError(t, g.WriteToFile(v))
	r.NoError(t, g.LoadFile(&c))
	r.Exactly(t, v.Int8, c.Int8)

	p, err := os.UserConfigDir()
	r.NoError(t, err)
	r.Exactly(t, filepath.Join(p, dir, file)+".json", g.FullPath())
	r.NoError(t, os.Remove(yPath))
	r.NoError(t, os.Remove(jPath))
	r.NoError(t, os.Remove(filepath.Dir(g.FullPath())))
}
