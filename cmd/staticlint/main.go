// Package staticlint
package main

import (
	"strings"

	errname "github.com/Antonboom/errname/pkg/analyzer"
	"github.com/lasiar/canonicalheader"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/multichecker"
	"golang.org/x/tools/go/analysis/passes/appends"
	"golang.org/x/tools/go/analysis/passes/assign"
	"golang.org/x/tools/go/analysis/passes/atomic"
	"golang.org/x/tools/go/analysis/passes/atomicalign"
	"golang.org/x/tools/go/analysis/passes/bools"
	"golang.org/x/tools/go/analysis/passes/buildtag"
	"golang.org/x/tools/go/analysis/passes/composite"
	"golang.org/x/tools/go/analysis/passes/copylock"
	"golang.org/x/tools/go/analysis/passes/defers"
	"golang.org/x/tools/go/analysis/passes/directive"
	"golang.org/x/tools/go/analysis/passes/errorsas"
	"golang.org/x/tools/go/analysis/passes/findcall"
	"golang.org/x/tools/go/analysis/passes/httpmux"
	"golang.org/x/tools/go/analysis/passes/httpresponse"
	"golang.org/x/tools/go/analysis/passes/ifaceassert"
	"golang.org/x/tools/go/analysis/passes/loopclosure"
	"golang.org/x/tools/go/analysis/passes/lostcancel"
	"golang.org/x/tools/go/analysis/passes/nilfunc"
	"golang.org/x/tools/go/analysis/passes/nilness"
	"golang.org/x/tools/go/analysis/passes/pkgfact"
	"golang.org/x/tools/go/analysis/passes/printf"
	"golang.org/x/tools/go/analysis/passes/shadow"
	"golang.org/x/tools/go/analysis/passes/shift"
	"golang.org/x/tools/go/analysis/passes/sigchanyzer"
	"golang.org/x/tools/go/analysis/passes/sortslice"
	"golang.org/x/tools/go/analysis/passes/stdmethods"
	"golang.org/x/tools/go/analysis/passes/stdversion"
	"golang.org/x/tools/go/analysis/passes/stringintconv"
	"golang.org/x/tools/go/analysis/passes/structtag"
	"golang.org/x/tools/go/analysis/passes/testinggoroutine"
	"golang.org/x/tools/go/analysis/passes/tests"
	"golang.org/x/tools/go/analysis/passes/timeformat"
	"golang.org/x/tools/go/analysis/passes/unmarshal"
	"golang.org/x/tools/go/analysis/passes/unreachable"
	"golang.org/x/tools/go/analysis/passes/unsafeptr"
	"golang.org/x/tools/go/analysis/passes/unusedresult"
	"golang.org/x/tools/go/analysis/passes/unusedwrite"
	"golang.org/x/tools/go/analysis/passes/usesgenerics"
	"honnef.co/go/tools/staticcheck"

	"github.com/novoseltcev/go-course/pkg/linters/noexit"
)

func main() {
	mychecks := []*analysis.Analyzer{
		appends.Analyzer,          // Check append usage
		assign.Analyzer,           // Check useless assignments
		atomic.Analyzer,           // Check usage of sync/atomic usage
		atomicalign.Analyzer,      // Check alignment for sync/atomic value
		bools.Analyzer,            // Check errors in boolean ops
		buildtag.Analyzer,         // Check go:build
		composite.Analyzer,        // Check struct args
		copylock.Analyzer,         // Check copy of locks
		defers.Analyzer,           // Check defer calls
		directive.Analyzer,        // Check go:build and go:debug
		errorsas.Analyzer,         // Check errors.As
		findcall.Analyzer,         // Check partitional calls
		httpmux.Analyzer,          // Check new http/mux syntax for go1.22
		httpresponse.Analyzer,     // Check defer body.Close()
		ifaceassert.Analyzer,      // Check imposible type asserts
		loopclosure.Analyzer,      // Check usage loop value in closure
		lostcancel.Analyzer,       // Check lost cancel call
		nilfunc.Analyzer,          // Check useless nil comparisons
		nilness.Analyzer,          // Check imposible nil comparisons
		pkgfact.Analyzer,          // Check package facts
		printf.Analyzer,           // Check printf calls
		shadow.Analyzer,           // Check shadowed variables
		shift.Analyzer,            // Check integer shifts
		sigchanyzer.Analyzer,      // Check channel buffer for os.Signal
		sortslice.Analyzer,        // Check sort.Slice argument type
		stdmethods.Analyzer,       // Check std methods naming
		stdversion.Analyzer,       // Check
		stringintconv.Analyzer,    // Check string(int) conversions
		structtag.Analyzer,        // Check struct tags
		testinggoroutine.Analyzer, // Check t.Fatal from goroutines
		tests.Analyzer,            // Check common tests errors
		timeformat.Analyzer,       // Check timeformat
		unmarshal.Analyzer,        // Check usage interface or pointer marshaling
		unreachable.Analyzer,      // Check unreachable code
		unsafeptr.Analyzer,        // Check convertation to unsafe.Pointer
		unusedresult.Analyzer,     // Check unused func result
		unusedwrite.Analyzer,      // Check never readed writer
		usesgenerics.Analyzer,     // Check allow to use generic
		noexit.Analyzer,           // Check os.Exit usage in main function of main packaged
		canonicalheader.Analyzer,  // Check usage http.Header
		errname.New(),             // Check errors naming
	}

	for _, v := range staticcheck.Analyzers {
		if strings.HasPrefix(v.Analyzer.Name, "SA") || v.Analyzer.Name == "ST1000" {
			mychecks = append(mychecks, v.Analyzer)
		}
	}

	multichecker.Main(mychecks...)
}
