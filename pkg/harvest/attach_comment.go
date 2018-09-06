package harvest

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
)

func AttachComment(baseDir string, c *Comment) (err error, dirname string, filename string, createPath bool) {

	c.AuthorEmail = ""

	if strings.HasPrefix(c.AttachTo, "comments/") {
		// comment to a comment; strip comment file basename
		dirname = path.Dir(c.AttachTo)
	} else {
		// it's the first comment to that post or page
		dirname = fmt.Sprintf("comments/%s", c.AttachTo)
	}

	fullDir := fmt.Sprintf("%s/%s", baseDir, dirname)
	fullDirExists := true
	var dirInfo os.FileInfo
	if dirInfo, err = os.Stat(fullDir); err != nil {
		fullDirExists = false
	}
	if fullDirExists && !dirInfo.IsDir() {
		return errors.New(fmt.Sprintf("Comment %v: %s exists, but is no directory", c, fullDir)),
			"", "", false
	}
	createPath = !fullDirExists

	if fullDirExists {
		if err, filename = makeCommentFilename(fullDir, c); err != nil {
			return err, "", "", false
		}
	} else {
		filename = "comment_0_1.json"
		c.Id = 1
	}
	c.IndentLevel = strings.Count(filename, "_") - 2
	return nil, dirname, filename, createPath
}

var commentRegExp = regexp.MustCompile(`_(\d+).json$`) // highest ID is last

func makeCommentFilename(dir string, c *Comment) (error, string) {
	// first determine the highest ID used so far
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return err, ""
	}
	maxId := 0
	for _, dEntry := range files {
		if dEntry.IsDir() {
			continue
		}
		matches := commentRegExp.FindStringSubmatch(dEntry.Name())
		if len(matches) != 2 {
			return errors.New(fmt.Sprintf("garbled file %s in dir %s", dEntry.Name(), dir)), ""
		}
		id, _ := strconv.Atoi(matches[1]) // already matches, ignore error

		if id > maxId {
			maxId = id
		}
	}
	commentId := maxId + 1
	c.Id = commentId
	if strings.HasPrefix(c.AttachTo, "comments/") {
		// suffix last part of c.AttachTo with commentId
		return nil, fmt.Sprintf("%s_%d.json", path.Base(c.AttachTo), commentId)
	} else {
		// it's the first comment to that post or page
		return nil, fmt.Sprintf("comment_0_%d.json", commentId)
	}
}
