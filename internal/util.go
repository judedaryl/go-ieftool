package internal

import (
	"log"
	"path/filepath"

	"github.com/spf13/pflag"
)

func remove[T any](slice []T, s int) []T {
	if len(slice) == 1 {
		return []T{}
	}
	newArr := []T{}
	for i := range slice {
		if i != s {
			newArr = append(newArr, slice[i])
		}
	}
	return newArr
}

func MustAbsPathFromFlag(f *pflag.FlagSet, n string) string {
	p, err := f.GetString(n)
	if err != nil {
		log.Fatalln(err)
	}
	if !filepath.IsAbs(p) {
		p, err = filepath.Abs(p)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return p
}
