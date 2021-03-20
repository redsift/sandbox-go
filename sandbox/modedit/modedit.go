package modedit

import (
	"io/ioutil"

	"golang.org/x/mod/modfile"
)

func load(fn string) (*modfile.File, error) {
	buf, err := ioutil.ReadFile(fn)
	if err != nil {
		return nil, err
	}
	f, err := modfile.Parse(fn, buf, nil)
	if err != nil {
		return nil, err
	}
	return f, nil
}

// CopyReplace copies all replace directives in fromFile and adds them to toFile, writing the output
// to newFile.
func CopyReplace(fromFile string, toFile string, newFile string) error {
	f, err := load(fromFile)
	if err != nil {
		return err
	}
	t, err := load(toFile)
	if err != nil {
		return err
	}
	for _, r := range f.Replace {
		t.AddReplace(r.Old.Path, r.Old.Version, r.New.Path, r.New.Version)
	}
	t.SetRequire(f.Require)
	t.AddRequire(f.Module.Mod.Path, "v1.0.0")
	t.Cleanup()
	buf, err := t.Format()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(newFile, buf, 0644)
}
