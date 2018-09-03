package harvest

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

var input = make(map[string]Comment)
var baseDir = "../../testdata"

// set up globals and run tests
func TestMain(m *testing.M) {
	if err := loadInput(); err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func loadInput() error {
	inputDir := "../../testdata/input"
	dirInfo, err := os.Stat(inputDir)
	if os.IsNotExist(err) {
		l.Printf("inputDir %s does not exist", inputDir)
		return err
	}
	if !dirInfo.IsDir() {
		l.Printf("inputDir %s is no directory", inputDir)
		return err
	}

	files, err := ioutil.ReadDir(inputDir)
	if err != nil {
		return err
	}

	for _, dEntry := range files {
		if dEntry.IsDir() {
			continue
		}
		var f *os.File
		if f, err = os.Open(fmt.Sprintf("%s/%s", inputDir, dEntry.Name())); err != nil {
			return err
		}

		var c Comment
		bytesRead, err := ioutil.ReadAll(f)
		if err != nil {
			l.Printf("error reading from input %s", f.Name())
			return err
		}
		if err = json.Unmarshal(bytesRead, &c); err != nil {
			l.Printf("error unmarshalling %s: %v", f.Name(), err)
			return err
		}
		input[dEntry.Name()] = c

		f.Close()
	}
	return nil
}

func TestAttachComment_PostFirst(t *testing.T) {
	inputFileNAme := "p_2018_08_4333-xxx.json"
	dirNameExpected := fmt.Sprintf("comments/%s", "post/2018/08/4333-wind-driving-the-spray")
	fileNameExpected := "comment_0_1.json"
	doAttachComment(t, inputFileNAme, dirNameExpected, fileNameExpected,
		true, 1, 0)
}

func TestAttachComment_PostThird(t *testing.T) {
	inputFileName := "p_2018_08_4314-xxx.json"
	dirNameExpected := fmt.Sprintf("comments/%s", "post/2018/08/4314-the-tv-tower")
	fileNameExpected := "comment_0_3.json"
	doAttachComment(t, inputFileName, dirNameExpected, fileNameExpected,
		false, 3, 0)
}

func TestAttachComment_CommentSecond(t *testing.T) {
	inputFileName := "c_p_2018_03_4164-xxx_c_0_68280.json"
	dirNameExpected := "comments/post/2018/03/4164-inside-saint-martins-i"
	fileNameExpected := "comment_0_68280_68281.json"
	doAttachComment(t, inputFileName, dirNameExpected, fileNameExpected,
		false, 68281, 1)
}

func TestAttachComment_CommentThird(t *testing.T) {
	inputFileName := "c_p_2018_08_4314-xxx_c_0_1_2.json"
	dirNameExpected := "comments/post/2018/08/4314-the-tv-tower"
	fileNameExpected := "comment_0_1_2_3.json"
	doAttachComment(t, inputFileName, dirNameExpected, fileNameExpected,
		false, 3, 2)
}

func TestAttachComment_PageFirst(t *testing.T) {
	inputFileNAme := "page_about.json"
	dirNameExpected := "comments/page/about"
	fileNameExpected := "comment_0_1.json"
	doAttachComment(t, inputFileNAme, dirNameExpected, fileNameExpected,
		true, 1, 0)
}

func TestAttachComment_PageSecond(t *testing.T) {
	inputFileNAme := "page_site-notice_c_0_17.json"
	dirNameExpected := "comments/page/site-notice"
	fileNameExpected := "comment_0_17_18.json"
	doAttachComment(t, inputFileNAme, dirNameExpected, fileNameExpected,
		false, 18, 1)
}

func doAttachComment(t *testing.T, inputFileName string, dirNameExpected string, fileNameExpected string,
	createPathExpected bool, idExpected int, indentLevelExpected int) {
	c := input[inputFileName]
	if err, commentDirName, commentFileName, createPath := AttachComment(baseDir, &c); err == nil {
		if commentDirName != dirNameExpected {
			t.Errorf("%s: dirname %s != expected %s", inputFileName, commentDirName, dirNameExpected)
		}
		if commentFileName != fileNameExpected {
			t.Errorf("%s: filename %s != expected %s", inputFileName, commentFileName, fileNameExpected)
		}
		if createPath != createPathExpected {
			t.Errorf("%s: %t != expected %t", inputFileName, createPath, createPathExpected)
		}
		if c.Id != idExpected {
			t.Errorf("%s: ID %d != expected %d", inputFileName, c.Id, idExpected)
		}
		if c.IndentLevel != indentLevelExpected {
			t.Errorf("%s: IndentLevel %d != expected %d", inputFileName, c.IndentLevel, indentLevelExpected)
		}
	} else {
		t.Error(err)
	}
}
