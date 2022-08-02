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
	opts gonfOptions
	path string
	mux  *sync.Mutex
}

type gonfOptions struct {
	overwrite bool
	dirName   string
	fileName  string
	fileType  GonfType
	mux       *sync.Mutex
}

func New(dirName, fileName string, configType GonfType, overwrite bool) (g *Gonf, err error) {
	if dirName == "" || fileName == "" {
		return nil, errors.New(ErrExpectedNonNilOrEmpty)
	}

	opts := gonfOptions{
		overwrite: overwrite,
		dirName:   dirName,
		fileName:  fileName,
		fileType:  configType,
		mux:       &sync.Mutex{},
	}
	opts.mux.Lock()
	defer opts.mux.Unlock()

	p, err := os.UserConfigDir()
	if err != nil {
		return
	}
	p = filepath.Join(p, opts.dirName, opts.fileName)
	ext := jsonExtension
	if opts.fileType == GonfYAML {
		ext = yamlExtension
	}
	p += ext
	opts.fileName += ext
	g = &Gonf{mux: &sync.Mutex{}, opts: opts, path: p}
	return
}

// FullPath returns the full path of the Gonfig file.
func (g *Gonf) FullPath() string {
	return g.path
}

// DirName returns the given directory name of the Gonfig file.
func (g *Gonf) DirName() string {
	return g.opts.dirName
}

// FileName returns the given name of the Gonfig file.
func (g *Gonf) FileName() string {
	return g.opts.fileName
}

// Type returns the given type of the Gonfig file.
func (g *Gonf) Type() GonfType {
	return g.opts.fileType
}

// WriteToFile writes the Gonfig to the Gonfig file.
func (g *Gonf) WriteToFile(conf interface{}) (err error) {
	if g == nil {
		return errors.New(ErrExpectedNonNilOrEmpty)
	}

	// Create the folder which contains the config.
	p, err := os.UserConfigDir()
	if err != nil {
		return
	}
	p = filepath.Join(p, g.opts.dirName)
	_, err = os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(p, 0700)
		if err != nil {
			return
		}
	} else if err != nil {
		return fmt.Errorf("%s: %v", ErrUnexpected, err)
	}

	fp := g.FullPath()
	if g.opts.overwrite {
		err = os.Remove(fp)
		if err != nil {
			return fmt.Errorf("%s: %s: %v", ErrCreatingConfig, ErrOverwrite, err)
		}
	}

	// Check if creating the file would overwrite the file.
	_, err = os.Stat(fp)
	if err == nil && !g.opts.overwrite {
		return fmt.Errorf("%s: %s", ErrCreatingConfig, ErrOverwriteDisabled)
	}
	if errors.Is(err, os.ErrNotExist) {
		var b []byte
		if g.opts.fileType == GonfJson {
			b, err = json.MarshalIndent(conf, "", "\t")
		} else {
			b, err = yaml.Marshal(conf)
		}
		if err != nil {
			return fmt.Errorf("%s: %v", ErrMarshalling, err)
		}
		err = os.WriteFile(fp, b, 0700)
	}
	return
}

// LoadFile unmarshalls the Gonfig to the given interface.
func (g *Gonf) LoadFile(val interface{}) (err error) {
	g.mux.Lock()
	defer g.mux.Unlock()

	p := g.FullPath()
	if err != nil {
		return
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return
	}

	if g.opts.fileType == GonfJson {
		return json.Unmarshal(b, val)
	}
	return yaml.Unmarshal(b, val)
}
