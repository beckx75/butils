package fileutils

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type TheFiles []*TheFile

type TheFile struct {
	Path     string
	Name     string
	Ext      string
	Dir      string
	Selected bool
}

func (tfs TheFiles) Names() []string {
	rec := make([]string, len(tfs))
	for i, tf := range tfs {
		rec[i] = tf.Name
	}
	return rec
}

func GetSelectedFilepathes(tfs []*TheFile) []string {
	rec := []string{}
	for _, tf := range tfs {
		if tf.Selected {
			rec = append(rec, tf.Path)
		}
	}
	return rec
}

func GetFilenames(tfs []*TheFile) []string {
	rec := []string{}
	for _, tf := range tfs {
		rec = append(rec, tf.Name)
	}
	return rec
}

func (tfs *TheFiles) UnselectAll() {
	for _, tf := range *tfs {
		tf.Selected = false
	}
}

func (tfs TheFiles) Select(idx int) {
	tfs[idx].Selected = true
}

func (tfs TheFiles) GetPathFromSelected() []string {
	lst := []string{}
	for _, tf := range tfs {
		if tf.Selected == true {
			lst = append(lst, tf.Path)
		}
	}
	return lst
}

func (tfs TheFiles) GetSelected() []*TheFile {
	selected := []*TheFile{}
	for _, tf := range tfs {
		if tf.Selected == true {
			selected = append(selected, tf)
		}
	}
	return selected
}

func NewTheFile(fp string) *TheFile {
	return &TheFile{
		Path:     fp,
		Name:     filepath.Base(fp),
		Ext:      filepath.Ext(fp),
		Dir:      filepath.Dir(fp),
		Selected: false,
	}
}

func NewTheFileList(files []string) []*TheFile {
	tfs := []*TheFile{}
	for _, file := range files {
		tfs = append(tfs, NewTheFile(file))
	}
	return tfs
}

func AddFile(tfs *[]*TheFile, fp string) {
	if !PathInFc(*tfs, fp) {
		*tfs = append(*tfs,
			NewTheFile(fp),
		)
	}
}

func (tf *TheFile) AddPath(fp string) {
	tf.Path = fp
	tf.Dir = filepath.Dir(fp)
	tf.Name = filepath.Base(fp)
	tf.Ext = filepath.Ext(fp)
	tf.Selected = false
}

func GetFiles(args []string, pattern []string) ([]*TheFile, []string, error) {
	files := []string{}
	for _, arg := range args {
		err := filepath.WalkDir(arg, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			ext := strings.ToLower(filepath.Ext(path))
			if sliceContains(pattern, ext) {
				abspath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				if !sliceContains(files, abspath) {
					files = append(files, abspath)
				}
			}
			return nil
		})
		if err != nil {
			return nil, nil, err
		}
	}
	tfs := make([]*TheFile, len(files))
	for i, file := range files {
		tf := &TheFile{Path: file}
		tf.Ext = filepath.Ext(file)
		tf.Name = filepath.Base(file)
		tf.Dir = filepath.Dir(file)
		tfs[i] = tf
	}
	return tfs, files, nil
}

func AddSuffixToFile(fp string, val string, delExtension bool) string {
	dir, file := filepath.Split(fp)
	ext := filepath.Ext(file)
	file, _ = strings.CutSuffix(file, ext)
	if delExtension {
		file = file + val
	} else {
		file = file + val + ext
	}
	return filepath.Join(dir, file)
}

func sliceContains(list []string, val string) bool {
	for _, e := range list {
		if e == val {
			return true
		}
	}
	return false
}

func PathInFc(tfs []*TheFile, path string) bool {
	for _, fc := range tfs {
		if fc.Path == path {
			return true
		}
	}
	return false
}
