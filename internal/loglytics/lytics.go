package loglytics

import (
	"bufio"
	"io"
	"io/fs"
	"log"
	"os"

	"github.com/odin-software/metro/control"
)

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetOrderedLogFiles() []fs.DirEntry {
	dir := control.DefaultConfig.LogsDirectory
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		log.Fatalf("logs: %s does not exist", dir)
	}

	files, err := os.ReadDir(dir)
	Check(err)
	return files
}

func GetLastLines(file fs.DirEntry, amt int) []string {
	dir := control.DefaultConfig.LogsDirectory
	f, err := os.Open(dir + file.Name())
	Check(err)

	rr := &ReverseReader{file: f}
	rr.SeekEnd()

	scanner := bufio.NewScanner(rr)
	lines := make([]string, 0)

	scanner.Scan()
	for i := 0; i < amt; i++ {
		scanner.Scan()
		text := []rune(scanner.Text())
		lines = append(lines, reverse(text))
	}

	return lines
}

type ReverseReader struct {
	file *os.File
}

// Seek to the final byte of the file
func (r *ReverseReader) SeekEnd() {
	_, err := r.file.Seek(0, io.SeekEnd)
	Check(err)
}

// Read the file backwards
func (r *ReverseReader) Read(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	// This no-op gives us the current offset value
	offset, err := r.file.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	var m int
	for i := 0; i < len(b); i++ {
		if offset == 0 {
			return m, io.EOF
		}
		// Seek in case someone else is relying on seek too
		offset, err = r.file.Seek(-1, io.SeekCurrent)
		if err != nil {
			return m, err // Should never happen
		}

		// Just read one byte at a time
		n, err := r.file.ReadAt(b[i:i+1], offset)
		if err != nil {
			return m + n, err // Should never happen
		}
		m += n
	}
	return m, nil
}

func reverse(r []rune) string {
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}
