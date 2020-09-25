package plugins

import (
	"fmt"
	"github.com/dustin/go-humanize"
	"strings"
)

type WriteCounter struct {
	current uint64
	total   uint64
}

func NewWriteCounter(total uint64) *WriteCounter {
	if total <= 0 {
		total = 0
	}
	return &WriteCounter{
		total: total,
	}
}

func (c *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	c.current += uint64(n)
	c.PrintProgress()
	return n, nil
}

func (c *WriteCounter) PrintProgress() {
	fmt.Printf("\r%s", strings.Repeat(" ", 35))

	var percent string
	if c.total > 0 {
		percent = fmt.Sprintf("%.2f%%", float64(c.current)/float64(c.total)*100)
	} else {
		percent = "N/A"
	}

	fmt.Printf("\rDownloading... %s(%s) complete", humanize.Bytes(c.current), percent)
}
