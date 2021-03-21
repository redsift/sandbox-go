package modedit

import (
	"io/ioutil"
	"os"
	"sort"
	"strings"

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

func alreadyPresent(r *modfile.Require, list []*modfile.Require) bool {
	for _, l := range list {
		if r.Mod == l.Mod {
			return true
		}
	}
	return false
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
	for _, r := range f.Require {
		if alreadyPresent(r, t.Require) {
			continue
		}
		t.AddNewRequire(r.Mod.Path, r.Mod.Version, r.Indirect)
	}
	t.AddRequire(f.Module.Mod.Path, "v1.0.0")
	t.Cleanup()
	buf, err := t.Format()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(newFile, buf, 0644)
}

func parseSum(data []byte) map[string]string {
	lines := strings.Split(string(data), "\n")

	m := make(map[string]string, len(lines))
	for _, l := range lines {
		w := strings.Fields(l)
		if len(w) != 3 {
			continue
		}
		m[w[0]+" "+w[1]] = l
	}
	return m
}

// CopySum copies all replace sum lines in f1 and adds them to f2, writing the output
// to newFile.
func CopySum(f1 string, f2 string, newFile string) error {
	f, err := os.ReadFile(f1)
	if err != nil {
		return err
	}
	t, err := os.ReadFile(f2)
	if err != nil {
		return err
	}
	fl := parseSum(f)
	tl := parseSum(t)

	out := strings.Split(string(t), "\n")
	out = out[:len(out)-1]
	for m, l := range fl {
		_, alreadyThere := tl[m]
		if alreadyThere || l == "" {
			continue
		}
		out = append(out, l)
	}
	sort.Strings(out)
	o, err := os.OpenFile(newFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	buf := strings.Join(out, "\n") + "\n"
	_, err = o.Write([]byte(buf))
	if err != nil {
		return err
	}
	return o.Close()
}
