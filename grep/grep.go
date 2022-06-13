package grep

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"
	"unicode"

	"github.com/zhiruili/urem/core"
)

// FilePredicate 用于过滤文件路径。
type FilePredicate func(string) bool

// Any 是一个 FilePredicate，任意输入值都会返回 true。
func Any(p string) bool {
	return true
}

// WithRegexp 是一个 FilePredicate 构造器，仅对符合给定正则表达式的路径返回 true。
func WithRegexp(r string) FilePredicate {
	re := regexp.MustCompile(r)
	return re.MatchString
}

// WithExts 是一个 FilePredicate 构造器，仅对指定后缀的路径返回 true。
func WithExts(exts ...string) FilePredicate {
	return func(p string) bool {
		ext := filepath.Ext(p)
		return core.StrContains(exts, ext)
	}
}

// Pattern 用于表示一个待查找数据。
type Pattern struct {
	Name   string
	Raw    string
	Regexp *regexp.Regexp
}

// Item 用于表示一个查找结果。
type Item struct {
	Error     error
	FileName  string
	Pattern   string
	Matched   []string
	LineNo    int
	LineText  string
	HeadLines []string
}

// Grep 对指定目录进行查找，结果通过 onFound 函数输出。
func Grep(patterns []*Pattern, dirnames []string, needGrep FilePredicate, onFound func(*Item) bool) {
	ch := make(chan *Item, 5)
	ctx := &context{
		waitGroup: &sync.WaitGroup{},
		outItem:   ch,
		stopFlag:  0,
	}

	ctx.waitGroup.Add(len(dirnames))
	for _, dirname := range dirnames {
		go ctx.grepOneDir(patterns, dirname)
	}

	go func() {
		stopped := false
		for {
			item := <-ch
			if item == nil {
				return
			}

			if !stopped && !onFound(item) {
				ctx.stop()
				stopped = true
			}
		}
	}()

	ctx.waitGroup.Wait()
	close(ch)
}

// GrepResult 对指定目录进行查找，直接返回结果。
func GrepResult(patterns []*Pattern, dirnames []string, needGrep FilePredicate) []*Item {
	var items []*Item
	Grep(patterns, dirnames, needGrep, func(i *Item) bool {
		items = append(items, i)
		return true
	})

	return items
}

// context 用于表示一个查找上下文。
type context struct {
	waitGroup *sync.WaitGroup
	outItem   chan<- *Item
	needGrep  FilePredicate
	stopFlag  int32
}

func (ctx *context) stop() {
	atomic.AddInt32(&ctx.stopFlag, 1)
}

func (ctx *context) isStopped() bool {
	return atomic.LoadInt32(&ctx.stopFlag) > 0
}

func errorItem(filename string, e error) *Item {
	return &Item{
		Error:    e,
		FileName: filename,
	}
}

func errorItemf(filename string, f string, a ...interface{}) *Item {
	return errorItem(filename, fmt.Errorf(f, a...))
}

func trimMatchLine(line string) string {
	return strings.TrimSpace(line)
}

func notSpace(r rune) bool {
	return !unicode.IsSpace(r)
}

func trimContextLines(ctxLines []string, line string) []string {
	spaceRuneNum := strings.IndexFunc(line, notSpace)
	return core.StrSliceMap(ctxLines, func(cl string) string {
		spaceRuneNum1 := strings.IndexFunc(cl, notSpace)
		if spaceRuneNum1 > spaceRuneNum {
			spaceRuneNum1 = spaceRuneNum
		}
		cl = cl[spaceRuneNum1:]
		return strings.TrimRightFunc(cl, unicode.IsSpace)
	})
}

var emptyLineRe = regexp.MustCompile(`^\s*$`)

func grepForPattern(name string, reader io.Reader, patterns []*Pattern, ctx *context) {
	fileReader := bufio.NewReader(reader)
	lineIdx := 0
	var headLines []string

	for {
		if ctx.isStopped() {
			core.LogD("early return %s:%d", name, lineIdx)
			return
		}

		line, err := fileReader.ReadString('\n')
		if err == io.EOF {
			return
		}

		if err != nil {
			ctx.outItem <- errorItemf(name, "read file: %w", err)
		}

		lineIdx++
		if emptyLineRe.MatchString(line) {
			headLines = nil
			continue
		}

		for _, pattern := range patterns {
			var matched []string
			if pattern.Regexp == nil {
				if strings.Contains(line, pattern.Raw) {
					matched = []string{pattern.Raw}
				}
			} else {
				matched = pattern.Regexp.FindStringSubmatch(line)
			}

			if matched != nil {
				ctx.outItem <- &Item{
					FileName:  name,
					Pattern:   pattern.Name,
					Matched:   matched,
					LineNo:    lineIdx,
					LineText:  trimMatchLine(line),
					HeadLines: trimContextLines(headLines, line),
				}
				goto CONTINUE_OUT
			}
		}

		headLines = append(headLines, strings.TrimRightFunc(line, unicode.IsSpace))
	CONTINUE_OUT:
	}
}

func (ctx *context) grepOneFile(patterns []*Pattern, filename string) {
	if ctx.needGrep != nil && !ctx.needGrep(filename) {
		core.LogD("ignore: %s", filename)
		return
	}

	core.LogD("greping file %s", filename)

	file, err := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		ctx.outItem <- errorItemf(filename, "read file: %w", err)
		return
	}
	defer file.Close()

	grepForPattern(file.Name(), file, patterns, ctx)
}

func (ctx *context) grepOneDir(patterns []*Pattern, dirname string) {
	defer ctx.waitGroup.Done()

	core.LogD("greping dir %s", dirname)

	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		ctx.outItem <- errorItemf(dirname, "read dir: %w", err)
		return
	}

	for _, file := range files {
		if ctx.isStopped() {
			core.LogD("early return %s:%s", dirname, file)
			return
		}

		fullPath := filepath.Join(dirname, file.Name())
		if !file.IsDir() {
			ctx.grepOneFile(patterns, fullPath)
			continue
		}

		ctx.waitGroup.Add(1)
		go ctx.grepOneDir(patterns, fullPath)
	}
}
