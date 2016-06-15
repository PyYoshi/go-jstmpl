package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	"github.com/ghodss/yaml"
	"github.com/jessevdk/go-flags"
	"github.com/lestrrat/go-jshschema"
	"github.com/minodisk/go-jstmpl"
)

func main() {
	os.Exit(_main())
}

type options struct {
	Schema   string `short:"s" long:"schema" description:"the source JSON Schema file"`
	OutDir   string `short:"o" long:"outfile" description:"output file to generate"`
	Template string `short:"t" log:"tmpl" description:"template file to generate document"`
}

func _main() int {
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		log.Printf("fail to parse flags: %s", err)
		return 1
	}

	f, err := os.Open(opts.Schema)
	if err != nil {
		log.Printf("fail to open the source JSON Schema file: %s", err)
		return 1
	}
	defer f.Close()

	var m map[string]interface{}
	switch ext := filepath.Ext(opts.Schema); ext {
	case ".json":
		if err := json.NewDecoder(f).Decode(&m); err != nil {
			log.Printf("fail to decode JSON: %s", err)
			return 1
		}
	case ".yml", ".yaml":
		b, err := ioutil.ReadFile(opts.Schema)
		if err != nil {
			log.Printf("fail to read the source JSON Schema file: %s", err)
			return 1
		}
		if err := yaml.Unmarshal(b, &m); err != nil {
			log.Printf("fail to unmarshal YAML: %s", err)
			return 1
		}
	default:
		log.Printf("undefined extension: %s", ext)
		return 1
	}

	hs := hschema.New()
	if err := hs.Extract(m); err != nil {
		log.Printf("fail to extract JSON Schema: %s", err)
		return 1
	}

	b := jstmpl.NewBuilder()
	ts, err := b.Build(hs)
	if err != nil {
		log.Printf("fail to build TSModel: %s", err)
		return 1
	}

	err = filepath.Walk(opts.Template, func(i string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		r, err := filepath.Rel(opts.Template, i)
		if err != nil {
			return err
		}
		o := filepath.Join(opts.OutDir, r)

		var tmpl []byte
		if i != "" {
			f, err := os.Open(i)
			if err != nil {
				return err
			}
			defer f.Close()
			tmpl, err = ioutil.ReadAll(f)
			if err != nil {
				return err
			}
		}

		var out io.Writer
		if o != "" {
			d := filepath.Dir(o)
			_, err := os.Stat(d)
			if err != nil {
				err := os.MkdirAll(d, 0755)
				if err != nil {
					return err
				}
			}
			f, err := os.Create(o)
			if err != nil {
				return err
			}
			defer f.Close()
			out = f
		}

		g := jstmpl.NewGenerator()
		if err := g.Process(out, ts, tmpl); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Printf("fail to walk: %s", err)
		return 1
	}

	return 0
}