package imglib

import (
	"io"
	"os"
	"path"
	"strings"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strconv"
)


type ImageConfiguration struct {
	ThumbMaxSize int
	LargeMaxSize int
	JpegQuality int
}

var DefaultImageConfiguration = ImageConfiguration{
	ThumbMaxSize: 300,
	LargeMaxSize: 1200,
	JpegQuality: 95,
}

func (img *Image) SaveToDirectory(dir string, imageStream io.Reader, imgCfg ImageConfiguration) (error) {
	upload, err := writeToTempFile(imageStream);
	if err != nil {
		return err
	}
	//defer os.Remove(upload)
	fmt.Printf("dir: %v\n", dir)
	fmt.Printf("tmpfile: %v\n", upload)
	img.LargeFilename = fmt.Sprintf("%v.jpeg", img.ID)
	fileNamePattern := fmt.Sprintf("%v[Q=%v]", path.Join(dir, img.LargeFilename), imgCfg.JpegQuality)
	err = scale(fileNamePattern, upload, imgCfg.LargeMaxSize)
	if err != nil {
		return err
	}
	img.ThumbFilename = fmt.Sprintf("%vtn.jpeg", img.ID)
	fileNamePattern = fmt.Sprintf("%v[Q=%v]", path.Join(dir, img.ThumbFilename), imgCfg.JpegQuality)
	err = scale(fileNamePattern, upload, imgCfg.ThumbMaxSize)

	return img.updateDimensions(dir)
}

func (img *Image) updateDimensions(dir string) (err error) {
	if img.LargeH, err = getImageIntHeader(path.Join(dir, img.LargeFilename), "height"); err != nil {
		return err
	}
	if img.LargeW, err = getImageIntHeader(path.Join(dir, img.LargeFilename), "width"); err != nil {
		return err
	}
	if img.ThumbH, err = getImageIntHeader(path.Join(dir,img.ThumbFilename), "height"); err != nil {
		return err
	}
	if img.ThumbW, err = getImageIntHeader(path.Join(dir, img.ThumbFilename), "width"); err != nil {
		return err
	}
	return nil
}

func getImageIntHeader(file string, header string) (int, error) {
	cmd := []string{"vipsheader", "-f", header, file}
	//fmt.Printf("cmd: %v\n", cmd)
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(strings.TrimSpace(string(out)))
}

func scale(dst, src string, maxSize int) (error) {
	cmd := []string{"vipsthumbnail", "--rotate", "-s", strconv.Itoa(maxSize), "-o", dst, src}
	//fmt.Printf("cmd: %v\n", cmd)
	out, err := exec.Command(cmd[0], cmd[1:]...).CombinedOutput()
	if len(out) > 0 {
		return fmt.Errorf("vipsthumbnail: %v", string(out))
	}
	return err
}

func writeToTempFile(stream io.Reader) (filename string, err error) {
	file, err := ioutil.TempFile("", "gallery-upload-")
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(file, stream); err != nil {
		return "", err
	}

	if err := file.Close(); err != nil {
		return "", err
	}

	return file.Name(), nil
}

func (img *Image) DeleteFromDirectory(directoryPath string) (error, error) {
	err1 := os.Remove(path.Join(directoryPath, img.LargeFilename))
	err2 := os.Remove(path.Join(directoryPath, img.ThumbFilename))
	return err1, err2
}
