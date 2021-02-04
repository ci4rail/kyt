package fwinfo

import (
	"testing"
)

func TestReadFile(t *testing.T) {
	var tt = []struct {
		f    string
		info string
		err  string
	}{
		{"testdata/ok.txt", "Verdin-iMX8MM_Reference-Minimal-Image_CI.OS.LMP-0.1.0-12.Branch.feature-kp-persistent-loggin.Sha.3e3589c2950bd6316b26385f4f8d5bdb3503bbd9-dev", ""},
		{"testdata/nok.txt", "", "No fwinfo line found"},
		{"testdata/nok2.txt", "", "No fwinfo line found"},
		{"testdata/notexist.txt", "", "open testdata/notexist.txt: no such file or directory"},
	}

	for _, w := range tt {
		t.Logf("File=%s", w.f)

		info, err := readFile(w.f)
		t.Logf("Got err=%v, %s", err, info)

		if info != w.info {
			t.Errorf("info wrong. Got %s, want %s", info, w.info)
		}
		if err == nil && w.err != "" {
			t.Errorf("Should return err")
		}
		if err != nil && w.err == "" {
			t.Errorf("Should not return err")
		}
	}
}
