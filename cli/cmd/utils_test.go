package cmd

import (
	"testing"
)

func TestOpenBrowser(t *testing.T) {
	url := "https://www.google.com/"
	err := OpenBrowser(url)
	if err != nil {
		t.Errorf("OpenBrowser url: %s, error: %v", url, err)
	}

	url = "dummy"
	err = OpenBrowser(url)
	if err == nil {
		t.Errorf("OpenBrowser url: %s, no expected errors", url)
	}
}
