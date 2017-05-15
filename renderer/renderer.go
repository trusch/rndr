package renderer

import (
	"io"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

var stdoutPath = "/dev/stdout"

// Renderer is a struct which renders templates or folders with templates
type Renderer struct{}

// Render renders src to dest.
// src can be a file or a folder
// dest must be a folder
func (renderer *Renderer) Render(src, dest string, data map[string]interface{}) error {
	src = filepath.Clean(src)
	dest = filepath.Clean(dest)
	srcStat, err := os.Stat(src)
	if err != nil {
		return nil
	}
	if !srcStat.IsDir() {
		if dest != stdoutPath {
			return renderer.renderFile(src, filepath.Join(dest, src), data)
		}
		return renderer.renderFile(src, dest, data)
	}
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if dest != stdoutPath {
				return renderer.renderFile(path, filepath.Join(dest, path), data)
			}
			return renderer.renderFile(path, dest, data)
		}
		return nil
	})
}

func (renderer *Renderer) renderFile(src, dest string, data map[string]interface{}) error {
	if renderer.isTemplate(src) {
		t, err := template.ParseFiles(src)
		if err != nil {
			return err
		}
		if dest != stdoutPath {
			dest = dest[:len(dest)-5]
			if err = os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
				return err
			}
			f, err := os.Create(dest)
			if err != nil {
				return err
			}
			log.Printf("render %v to %v", src, dest)
			return t.Execute(f, data)
		}
		log.Printf("render %v to stdout", src)
		return t.Execute(os.Stdout, data)
	}
	return renderer.copyFile(src, dest)
}

func (renderer *Renderer) isTemplate(path string) bool {
	return filepath.Ext(path) == ".tmpl"
}

func (renderer *Renderer) copyFile(src, dest string) error {
	if src == dest {
		return nil
	}
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	if err = os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	log.Printf("copied %v to %v", src, dest)
	return err
}
