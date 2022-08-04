package gonfig

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"sync"

	"gopkg.in/yaml.v3"
)

type GonfType int8

const (
	// GonfJson corresponds to a config of type JSON.
	GonfJson GonfType = iota

	// GonfYaml corresponds to a config of type YAML.
	GonfYaml GonfType = iota

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
	dir       string
	fileName  string
	fileType  GonfType
	mux       *sync.Mutex
}

func New(dirName, fileName string, configType GonfType, overwrite bool) (g *Gonf, err error) {
	if dirName == "" || fileName == "" {
		return nil, ErrExpectedNonNilOrEmpty
	}

	p, err := os.UserConfigDir()
	if err != nil {
		return
	}

	ext := jsonExtension
	if configType == GonfYaml {
		ext = yamlExtension
	}

	opts := gonfOptions{
		overwrite: overwrite,
		dir:       filepath.Join(p, dirName),
		fileName:  fileName + ext,
		fileType:  configType,
		mux:       &sync.Mutex{},
	}
	p = filepath.Join(opts.dir, opts.fileName)

	opts.mux.Lock()
	g = &Gonf{mux: &sync.Mutex{}, opts: opts, path: p}
	opts.mux.Unlock()
	return
}

func NewWithPath(fullPath string, configType GonfType, overwrite bool) (g *Gonf, err error) {
	if fullPath == "" {
		return nil, ErrExpectedNonNilOrEmpty
	}

	opts := gonfOptions{
		overwrite: overwrite,
		dir:       filepath.Dir(fullPath),
		fileName:  filepath.Base(fullPath),
		fileType:  configType,
		mux:       &sync.Mutex{},
	}

	opts.mux.Lock()
	g = &Gonf{mux: &sync.Mutex{}, opts: opts, path: fullPath}
	opts.mux.Unlock()
	return
}

// FullPath returns the full path of the Gonfig file.
func (g *Gonf) FullPath() string {
	return g.path
}

// Dir returns the given directory name of the Gonfig file.
func (g *Gonf) Dir() string {
	return g.opts.dir
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
		return ErrExpectedNonNilOrEmpty
	}

	// Create the folder which contains the config.
	p := filepath.Dir(g.FullPath())
	_, err = os.Stat(p)
	if errors.Is(err, os.ErrNotExist) {
		err = os.Mkdir(p, 0700)
		if err != nil {
			return
		}
	} else if err != nil {
		return fmt.Errorf("%v: %w", ErrUnexpected, err)
	}

	fp := g.FullPath()
	_, err = os.Stat(fp)
	if g.opts.overwrite && err == nil {
		err = os.Remove(fp)
		if err != nil {
			return fmt.Errorf("%v: %w", ErrOverwriteRemove, err)
		}
	}

	// Check if creating the file would overwrite the file.
	_, err = os.Stat(fp)
	if err == nil && !g.opts.overwrite {
		return ErrOverwriteDisabled
	}
	if errors.Is(err, os.ErrNotExist) {
		var b []byte
		if g.opts.fileType == GonfJson {
			b, err = json.MarshalIndent(conf, "", "\t")
		} else {
			b, err = yaml.Marshal(conf)
		}
		if err != nil {
			return fmt.Errorf("%v: %w", ErrMarshalling, err)
		}
		err = os.WriteFile(fp, b, 0700)
	}
	return
}

// Load unmarshalls the Gonfig to the given interface.
func (g *Gonf) Load(val interface{}) (err error) {
	if reflect.ValueOf(val).Kind() != reflect.Ptr {
		return ErrMustBeAddressable
	}

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
