package noexit_test

import (
	"testing"

	"golang.org/x/tools/go/analysis/analysistest"

	"github.com/novoseltcev/go-course/pkg/linters/noexit"
)

func Test(t *testing.T) {
	t.Parallel()

	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, noexit.Analyzer, "a", "b")
}
