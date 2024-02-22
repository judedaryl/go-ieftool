package internal

import (
	"log"
	"path/filepath"

	"github.com/spf13/pflag"
)

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
