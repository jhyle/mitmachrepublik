package mmr

import (
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type (
	Templates struct {
		sync.RWMutex
		pattern  string
		files    []string
		modTimes []time.Time
		tpls     *template.Template
		funcs    map[string]interface{}
	}
)

func NewTemplates(pattern string, funcs map[string]interface{}) (*Templates, error) {

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading templates from %s", pattern)
	}

	modTimes := make([]time.Time, len(files))
	for i, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, errors.Wrapf(err, "error getting file stats from %s", file)
		}
		modTimes[i] = fileInfo.ModTime()
	}

	tpls, err := template.New("/").Funcs(funcs).ParseFiles(files...)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing templates")
	}

	return &Templates{pattern: pattern, files: files, modTimes: modTimes, tpls: tpls, funcs: funcs}, nil
}

func (templates *Templates) reloadIfChanged() error {

	var err error
	var fileInfo os.FileInfo
	var newTemplates *Templates

	templates.Lock()

	for i, file := range templates.files {
		fileInfo, err = os.Stat(file)
		if err != nil {
			err = errors.Wrapf(err, "error getting file stats from %s", file)
			break
		}
		if fileInfo.ModTime() != templates.modTimes[i] {
			newTemplates, err = NewTemplates(templates.pattern, templates.funcs)
			if err == nil {
				templates.files = newTemplates.files
				templates.modTimes = newTemplates.modTimes
				templates.tpls = newTemplates.tpls
			}
			break
		}
	}

	templates.Unlock()
	return err
}

func (templates *Templates) Find(name string) (*template.Template, error) {

	err := templates.reloadIfChanged()
	if err != nil {
		return nil, errors.Wrap(err, "error reloading changed templates")
	}

	tpl := templates.tpls.Lookup(name)
	if tpl == nil {
		return nil, errors.New("Could not find template " + name + ".")
	}

	return tpl, nil
}

func (templates *Templates) Execute(name string, wr io.Writer, data map[string]interface{}) error {

	tpl, err := templates.Find(name)
	if err != nil {
		return errors.Wrapf(err, "error finding template %s", name)
	}

	err = tpl.Execute(wr, data)
	if err != nil {
		return errors.Wrapf(err, "error executing template %s", name)
	}

	return nil
}
