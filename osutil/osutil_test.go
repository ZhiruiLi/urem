package osutil

import (
	"path/filepath"
	"testing"
)

// TestFindFileBottomUp 测试 FindFileBottomUp 函数。
func TestFindFileBottomUp(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		exts     []string
		expectOK bool
		expect   string
	}{
		{
			name:     "find uplugin",
			path:     "../testdata/FakePlugin/Source/FakeModule",
			exts:     []string{"*.uplugin"},
			expectOK: true,
			expect:   filepath.Join("..", "testdata", "FakePlugin", "FakePlugin.uplugin"),
		},
		{
			name:     "illegal path",
			path:     "../testdat/FakePlugin/Source/FakeModule",
			exts:     []string{"*.uplugin"},
			expectOK: true,
			expect:   "",
		},
		{
			name:     "no find uproject",
			path:     "../testdata/FakePlugin/Source/FakeModule",
			exts:     []string{"*.uproject"},
			expectOK: true,
			expect:   "",
		},
		{
			name:     "find uplugin with project",
			path:     "../testdata/FakePlugin/Source/FakeModule",
			exts:     []string{"*.uproject", "*.uplugin"},
			expectOK: true,
			expect:   filepath.Join("..", "testdata", "FakePlugin", "FakePlugin.uplugin"),
		},
	}

	for i, c := range cases {
		actual, err := FindFileBottomUp(c.path, c.exts...)
		if c.expectOK != (err == nil) {
			t.Errorf("%d:%s: expect OK = %t, got error: %s", i, c.name, c.expectOK, err)
		}
		if actual != c.expect {
			t.Errorf("%d:%s:\nexpect:\n%s\n\nactual:\n%s", i, c.name, c.expect, actual)
		}
	}
}
