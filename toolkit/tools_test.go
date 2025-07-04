package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomStringUsingRandPrime(t *testing.T) {
	var testtools Tools

	s := testtools.RandomStringUsingRandPrime(7)
	if len(s) != 7 {
		t.Error("wrong length of random string returned")
	}
}

func TestTools_RandomStringUsingRandInt(t *testing.T) {
	var testtools Tools
	s := testtools.RandomStringUsingRandInt(8)
	if len(s) != 8 {
		t.Error("Wrong length of random string returned")
	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{
		name:          "allowd no rename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    false,
		errorExpected: false,
	},
	{
		name:          "allowd no ename",
		allowedTypes:  []string{"image/jpeg", "image/png"},
		renameFile:    true,
		errorExpected: false,
	},
	{
		name:          "not allowd",
		allowedTypes:  []string{"image/jpeg"},
		renameFile:    false,
		errorExpected: true,
	},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		// set up a pipe to avoid buffering
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			// Create the form data field 'file'
			part, err := writer.CreateFormFile("file", "./testdata/img.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/img.png")
			if err != nil {
				t.Error(err)
			}
			defer f.Close()

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("error decoding image", err)
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}

		}()

		// read from pipe which receives data

		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes

		uploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploadFiles/", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploadFiles/%s", uploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}

			//clean up
			_ = os.Remove(fmt.Sprintf("./testdata/uploadFiles/%s", uploadedFiles[0].NewFileName))
		}
		if !e.errorExpected && err != nil {
			t.Errorf("%s: error expected but none received", e.name)
		}

		wg.Wait()
	}
}

func TestTools_UploadOneFile(t *testing.T) {
	//Setup a pipe to avoid buffering
	pr, pw := io.Pipe()
	writer := multipart.NewWriter(pw)

	go func() {
		defer writer.Close()

		//Create the form data from field 'file'
		part, err := writer.CreateFormFile("file", "./testdata/img/png")
		if err != nil {
			t.Error(err)
		}

		f, err := os.Open("./testdata/img.png")
		if err != nil {
			t.Error(err)
		}
		defer f.Close()

		img, _, err := image.Decode(f)
		if err != nil {
			t.Errorf("error decoding image: %s", err)
		}

		err = png.Encode(part, img)
		if err != nil {
			t.Error(err)
		}

	}()

	request := httptest.NewRequest("POST", "/", pr)
	request.Header.Add("Content-Type", writer.FormDataContentType())

	var testTools Tools

	uploadFile, err := testTools.UploadOneFile(request, "./testdata/uploadFiles/")
	if err != nil {
		t.Error(err)
	}

	if _, err := os.Stat(fmt.Sprintf("./testdata/uploadFiles/%s", uploadFile.NewFileName)); os.IsNotExist(err) {
		t.Errorf("expected file to exists: %s", uploadFile.NewFileName)
	}

	_ = os.Remove(fmt.Sprintf("./testdata/uploadFiles/%s", uploadFile.NewFileName))
}

func TestTools_CreateDirIfNotExists(t *testing.T) {
	var testtools Tools

	dirName := "./testdir/dir1"

	err := testtools.CreateDirIfNotExist(dirName)
	if err != nil {
		t.Error(err)
	}

	err = testtools.CreateDirIfNotExist(dirName)
	if err != nil {
		t.Error(err)
	}

	err = os.RemoveAll(dirName)
	if err != nil {
		t.Error(err)
	}
}

var slugTest = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complext string", s: "Now is the time for all GOOD men! + fish & such &^123", expected: "now-is-the-time-for-all-good-men-fish-such-123", errorExpected: false},
}

func TestTool_SlugifyText(t *testing.T) {
	var testTool Tools

	for _, e := range slugTest {
		result, err := testTool.Slugify(e.s)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error received when none expected: %s", e.name, err.Error())
		}

		if !e.errorExpected && result != e.expected {
			t.Errorf("%s: wrong slug returned; expeted %s but got %s", e.name, e.expected, result)
		}
	}

}

func TestTools_DownloadStaticFile(t *testing.T) {
	var tools Tools

	rr := httptest.NewRecorder()

	req, _ := http.NewRequest("GET", "/", nil)

	tools.DownloadStaticFile(rr, req, "./testdata", "img.png", "cloud.png")

	res := rr.Result()
	defer res.Body.Close()

	if res.Header["Content-Length"][0] != "534283" {
		t.Error("wrong content length of", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Disposition"][0] != "attachment; filename=\"cloud.png\"" {
		t.Error("wrong content disposition")
	}

	_, err := io.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}
