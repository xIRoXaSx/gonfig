package gonfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

type GonfType int8

const (
	// GonfJson corresponds to a config of type JSON.
	GonfJson GonfType = iota

	// GonfYAML corresponds to a config of type YAML.
	GonfYAML GonfType = iota

	yamlExtension = ".yaml"
	jsonExtension = ".json"
)

type Gonf struct {
	FullPath string
	Opts     *GonfOptions
	mux      *sync.Mutex
}

type GonfOptions struct {
	OverwriteExisting bool
	DirName           string
	ConfigName        string
	Type              GonfType
	mux               *sync.Mutex
}

func New(opts *GonfOptions) (g *Gonf, err error) {
	if opts == nil {
		return nil, errors.New(ErrExpectedNonNilOrEmpty)
	}
	opts.mux = &sync.Mutex{}
	g = &Gonf{mux: &sync.Mutex{}, Opts: opts}
	return
}

func (g *Gonf) WriteToFile(conf interface{}) (err error) {
	path, err := g.Opts.fullPath()
	if err != nil {
		return
	}

	// Create the folder which contains the config.
	p, err := os.UserConfigDir()
	if err != nil {
		return
	}
	p = filepath.Join(p, g.Opts.DirName)
	_, err = os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(p, 0700)
		if err != nil {
			return
		}
	} else if err != nil {
		return fmt.Errorf("%s: %v", ErrUnexpected, err)
	}

	if g.Opts.OverwriteExisting {
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("%s: %v", ErrOverwrite, err)
		}
	}

	// Check if creating the file would overwrite the file.
	_, err = os.Stat(path)
	if errors.Is(err, os.ErrNotExist) {
		var b []byte
		if g.Opts.Type == GonfJson {
			b, err = json.MarshalIndent(conf, "", "\t")
		} else {
			b, err = yaml.Marshal(conf)
		}
		if err != nil {
			return fmt.Errorf("%s: %v", ErrMarshalling, err)
		}
		err = os.WriteFile(path, b, 0700)
	}
	return
}

func (g *Gonf) LoadFile(val interface{}) (err error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	p, err := g.Opts.fullPath()
	if err != nil {
		return
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return
	}

	if g.Opts.Type == GonfJson {
		return json.Unmarshal(b, val)
	}
	return yaml.Unmarshal(b, val)
}

// FullPath returns the full path of the configuration file.
func (opts *GonfOptions) fullPath() (p string, err error) {
	opts.mux.Lock()
	defer opts.mux.Unlock()

	if opts.DirName == "" || opts.ConfigName == "" {
		return "", errors.New(ErrExpectedNonNilOrEmpty)
	}

	p, err = os.UserConfigDir()
	if err != nil {
		return
	}
	p = filepath.Join(p, opts.DirName, opts.ConfigName)
	if opts.Type == GonfJson {
		p += jsonExtension
	} else {
		p += yamlExtension
	}
	return
}
