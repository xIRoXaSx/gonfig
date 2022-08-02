package gonfig

import (
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
	g, err := New(nil)
	r.Error(t, err)

	g, err = New(&GonfOptions{})
	r.NoError(t, err)
	r.Error(t, g.WriteToFile(v))

	g, err = New(&GonfOptions{
		DirName:    "GonfigTest",
		ConfigName: "gonfig",
		Type:       GonfJson,
	})
	r.NoError(t, err)

	g, err = New(&GonfOptions{
		DirName:    "GonfigTest",
		ConfigName: "gonfig",
		Type:       GonfYAML,
	})
	r.NoError(t, err)
	r.NoError(t, g.WriteToFile(v))

	var c config
	r.NoError(t, g.LoadFile(&c))
	r.Equal(t, v, &c)

	g.Opts.Type = GonfJson
	r.NoError(t, g.WriteToFile(v))
}
