package imglib

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"image"
	"image/color"
	"image/jpeg"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestImageScalingAndSaving(t *testing.T) {
	image := &Image{
		ID: 42,
	}
	tmp := tmpDir()

	err := image.SaveToDirectory(tmp, bytes.NewBuffer(createDummyImage()), ImageConfiguration{200, 1200, 95})

	assert.NoError(t, err)

	fmt.Printf("image: %v", image)

	assert.Equal(t, "42.jpeg", image.LargeFilename)
	assert.Equal(t, "42tn.jpeg", image.ThumbFilename)

	assert.Equal(t, 200, image.ThumbW)
	assert.Equal(t, 150, image.ThumbH)
	assert.Equal(t, 1200, image.LargeW)
	assert.Equal(t, 900, image.LargeH)

	assertImage(t, path.Join(tmp, image.LargeFilename), 1200, 900)
	assertImage(t, path.Join(tmp, image.ThumbFilename), 200, 150)
}

func assertImage(t *testing.T, imagePath string, w, h int) {
	file, err := os.Open(imagePath)
	assert.NoError(t, err)
	defer file.Close()

	image, _, err := image.Decode(file)
	assert.NoError(t, err)
	if err != nil {
		return
	}
	assert.Equal(t, w, image.Bounds().Dx())
	assert.Equal(t, h, image.Bounds().Dy())
}

func tmpDir() string {
	name, err := ioutil.TempDir("", "gallery-unittest-")
	if err != nil {
		panic("can not create temp dir")
	}
	return name
}

func createDummyImage() []byte {
	img := image.NewRGBA(image.Rect(0, 0, 2000, 1500))
	col := color.RGBA{255, 0, 0, 255} // Red
	for i := 100; i < 500; i++ {
		for j := 100; j < 500; j++ {
			img.Set(i, j, col)
		}
	}
	buff := new(bytes.Buffer)
	jpeg.Encode(buff, img, nil)
	return buff.Bytes()
}
