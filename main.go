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
	opts GonfOptions
	path string
	mux  *sync.Mutex
}

type GonfOptions struct {
	OverwriteExisting bool
	DirName           string
	ConfigName        string
	Type              GonfType
	mux               *sync.Mutex
}

func New(opts GonfOptions) (g *Gonf, err error) {
	opts.mux = &sync.Mutex{}
	p, err := opts.fullPath()
	if err != nil {
		return
	}
	g = &Gonf{mux: &sync.Mutex{}, opts: opts, path: p}
	return
}

func (g *Gonf) Options() GonfOptions {
	return g.opts
}

func (g *Gonf) FullPath() string {
	return g.path
}

func (g *Gonf) WriteToFile(conf interface{}) (err error) {
	if g == nil {
		return errors.New(ErrExpectedNonNilOrEmpty)
	}
	path, err := g.opts.fullPath()
	if err != nil {
		return
	}

	// Create the folder which contains the config.
	p, err := os.UserConfigDir()
	if err != nil {
		return
	}
	p = filepath.Join(p, g.opts.DirName)
	_, err = os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(p, 0700)
		if err != nil {
			return
		}
	} else if err != nil {
		return fmt.Errorf("%s: %v", ErrUnexpected, err)
	}

	if g.opts.OverwriteExisting {
		err = os.Remove(path)
		if err != nil {
			return fmt.Errorf("%s: %s: %v", ErrCreatingConfig, ErrOverwrite, err)
		}
	}

	// Check if creating the file would overwrite the file.
	_, err = os.Stat(path)
	if err == nil && !g.opts.OverwriteExisting {
		return fmt.Errorf("%s: %s", ErrCreatingConfig, ErrOverwriteDisabled)
	}
	if errors.Is(err, os.ErrNotExist) {
		var b []byte
		if g.opts.Type == GonfJson {
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

	p, err := g.opts.fullPath()
	if err != nil {
		return
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return
	}

	if g.opts.Type == GonfJson {
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
