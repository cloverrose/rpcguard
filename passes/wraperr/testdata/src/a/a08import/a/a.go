package a

import (
	"a/a08import/excludedpkg"
	"a/a08import/includedpkg"
)

func CallIncludedOKFunc() error { // want CallIncludedOKFunc:"okFunc"
	return includedpkg.OKFunc()
}

func CallIncludedBadFunc() error { // want CallIncludedBadFunc:"badFunc"
	return includedpkg.BadFunc()
}

func CallExcludedOKFunc() error { // want CallExcludedOKFunc:"badFunc"
	return excludedpkg.OKFunc()
}
