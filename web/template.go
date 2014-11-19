package mmr

import (
	"errors"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type (
	Templates struct {
		sync.RWMutex
		pattern  string
		files    []string
		modTimes []time.Time
		tpls     *template.Template
	}
)

func NewTemplates(pattern string) (*Templates, error) {

	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	modTimes := make([]time.Time, len(files))
	for i, file := range files {
		fileInfo, err := os.Stat(file)
		if err != nil {
			return nil, err
		}
		modTimes[i] = fileInfo.ModTime()
	}

	tpls, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}

	return &Templates{pattern: pattern, files: files, modTimes: modTimes, tpls: tpls}, nil
}

func (templates *Templates) reloadIfChanged() error {

	var err error = nil
	var fileInfo os.FileInfo
	var newTemplates *Templates
	
	templates.Lock()

	for i, file := range templates.files {
		fileInfo, err = os.Stat(file)
		if err != nil {
			break
		}
		if fileInfo.ModTime() != templates.modTimes[i] {
			newTemplates, err = NewTemplates(templates.pattern)
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
		return nil, err
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
		return err
	}
	return tpl.Execute(wr, data)
}
