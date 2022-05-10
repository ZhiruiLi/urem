package infocmd

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sync"
	"sync/atomic"

	"github.com/zhiruili/urem/core"
)

type grepPattern struct {
	Name   string
	Raw    string
	Regexp *regexp.Regexp
}

type grepContext struct {
	WaitGroup *sync.WaitGroup
	OutResult chan<- *grepResult

	stopFlag int32
}

func (ctx *grepContext) stop() {
	atomic.AddInt32(&ctx.stopFlag, 1)
}

func (ctx *grepContext) isStopped() bool {
	return atomic.LoadInt32(&ctx.stopFlag) > 0
}

type grepResult struct {
	Error    error
	FileName string
	Pattern  string
	Matched  []string
	LineNo   int
	LineText string
}

func errorResult(filename string, e error) *grepResult {
	return &grepResult{
		Error:    e,
		FileName: filename,
	}
}

func errorResultf(filename string, f string, a ...interface{}) *grepResult {
	return errorResult(filename, fmt.Errorf(f, a...))
}

func grepForPattern(file *os.File, patterns []*grepPattern, ctx *grepContext) {
	fileReader := bufio.NewReader(file)
	lineIdx := 0

	for {
		if ctx.isStopped() {
			core.LogD("early return %s:%d", file.Name(), lineIdx)
			return
		}

		line, err := fileReader.ReadString('\n')
		lineIdx++

		for _, pattern := range patterns {
			matched := pattern.Regexp.FindStringSubmatch(line)
			if matched != nil {
				ctx.OutResult <- &grepResult{
					FileName: file.Name(),
					Pattern:  pattern.Name,
					Matched:  matched,
					LineNo:   lineIdx,
					LineText: line,
				}
				goto CONTINUE_OUT
			}
		}

		if err == io.EOF {
			return
		}

		if err != nil {
			ctx.OutResult <- errorResultf(file.Name(), "read file: %w", err)
		}
	CONTINUE_OUT:
	}
}

func grepOneFile(patterns []*grepPattern, filename string, ctx *grepContext) {
	ext := filepath.Ext(filename)
	if ext != ".h" && ext != ".hpp" {
		core.LogD("not header file, ignore: %s", filename)
		return
	}

	core.LogD("greping file %s", filename)

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		ctx.OutResult <- errorResultf(filename, "read file: %w", err)
		return
	}
	defer file.Close()

	grepForPattern(file, patterns, ctx)
}

func grepOneDir(patterns []*grepPattern, dirname string, ctx *grepContext) {
	defer ctx.WaitGroup.Done()

	core.LogD("greping dir %s", dirname)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		ctx.OutResult <- errorResultf(dirname, "read dir: %w", err)
		return
	}

	for _, file := range files {
		if ctx.isStopped() {
			core.LogD("early return %s:%s", dirname, file)
			return
		}

		fullPath := filepath.Join(dirname, file.Name())
		if !file.IsDir() {
			grepOneFile(patterns, fullPath, ctx)
			continue
		}

		ctx.WaitGroup.Add(1)
		go grepOneDir(patterns, fullPath, ctx)
	}
}

func grepManyDir(patterns []*grepPattern, dirnames []string) []*grepResult {
	ch := make(chan *grepResult, 3)
	ctx := &grepContext{
		WaitGroup: &sync.WaitGroup{},
		OutResult: ch,
	}

	ctx.WaitGroup.Add(len(dirnames))
	for _, dirname := range dirnames {
		go grepOneDir(patterns, dirname, ctx)
	}

	var results []*grepResult
	go func() {
		for {
			result := <-ch
			results = append(results, result)
		}
	}()

	ctx.WaitGroup.Wait()
	close(ch)

	core.LogD("results count %d", len(results))
	return results
}