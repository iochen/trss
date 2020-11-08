package trss_test

import (
	"fmt"
	"testing"

	"github.com/iochen/trss"
)

func TestGenerateRSS(t *testing.T) {
	rss, err := trss.GenerateRSS("v2fly")
	if err != nil {
		t.Error(err)
	}

	rssB, err := rss.ToRss()
	if err != nil {
		t.Error(err)
	}

	fmt.Println(rssB)
}
