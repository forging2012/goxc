package core

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

//just for test-driving the 'test' task
/*
func TestFail(t *testing.T) {
	t.Fatalf("FAIL")
}
*/
func TestGetGoPath(t *testing.T) {
	orig := os.Getenv("GOPATH")
	p1 := path.Join("a", "b")

	os.Setenv("GOPATH", JoinList(p1, "..", "c"))
	gopath := GetGoPathElement(".")
	if gopath != ".." {
		t.Fatalf("Could not load gopath correctly 1 - %s %s", os.Getenv("GOPATH"), gopath)

	}
	os.Setenv("GOPATH", p1)
	gopath = GetGoPathElement(".")
	if gopath != p1 {
		t.Fatalf("Could not load gopath correctly 2 - %s %s", os.Getenv("GOPATH"), gopath)
	}
	os.Setenv("GOPATH", orig)
}

func JoinList(elem ...string) string {
	for i, e := range elem {
		if e != "" {
			return filepath.Clean(strings.Join(elem[i:], string(os.PathListSeparator)))
		}
	}
	return ""
}

func TestSanityCheck(t *testing.T) {
	//goroot := runtime.GOROOT()
	err := SanityCheck("")
	if err == nil {
		t.Fatalf("sanity check failed! Expected to flag missing GOROOT variable")
	}
	tmpDir, err := ioutil.TempDir("", "goxc_test_sanityCheck")
	defer os.RemoveAll(tmpDir)
	err = SanityCheck(tmpDir)
	if err == nil {
		t.Fatalf("sanity check failed! Expected to notice missing src folder")
	}

	srcDir := filepath.Join(tmpDir, "src")
	os.Mkdir(srcDir, 0700)
	scriptname := GetMakeScriptPath(tmpDir)
	ioutil.WriteFile(scriptname, []byte("1"), 0111)
	err = SanityCheck(tmpDir)
	if err != nil {
		t.Fatalf("sanity check failed! Did not find src folder: %v", err)
	}
	//chmod doesnt work in Windows.
	//TODO: verify which OSes support chmod
	if runtime.GOOS == "linux" {
		os.Chmod(srcDir, 0600)
		defer os.Chmod(srcDir, 0700)
		err = SanityCheck(tmpDir)
		if err == nil {
			t.Fatalf("sanity check failed! Expected NOT to be able to open src folder")
		}
	}
}
