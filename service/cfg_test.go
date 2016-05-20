package service

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestCfg(t *testing.T) {

	assert.Equal(t, "/tmp/gallery", Cfg("galleryDir"))

	os.Setenv("galleryDir", "test")
	assert.Equal(t, "test", Cfg("galleryDir"))

	os.Setenv("intValue", "42")
	assert.Equal(t, 42, CfgInt("intValue"))
}

func TestPanicOnWrongKey(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic on missing key")
			t.Fail()
		}
	}()
	Cfg("xyz")
}
