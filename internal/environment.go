package internal

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"
)

type Environment struct {
	Name     string                 `yaml:"name"`
	Tenant   string                 `yaml:"tenant"`
	TenantId string                 `yaml:"tenantId"`
	ClientId string                 `yaml:"clientId"`
	Settings map[string]interface{} `yaml:"settings"`
}

func (env Environment) Build(s string, d string) {
	root := s
	_ = filepath.WalkDir(s, func(p string, e fs.DirEntry, err error) error {
		if e.IsDir() {
			return nil
		}
		if filepath.Ext(e.Name()) == ".xml" {
			t := path.Join(d, strings.ReplaceAll(p, root, env.Name))
			c := env.replaceVariables(p)
			err = os.MkdirAll(filepath.Dir(t), os.ModePerm)
			if err != nil {
				return err
			}
			log.Printf("Compiled %s", t)
			err = os.WriteFile(t, c, os.ModePerm)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (env Environment) replaceVariables(p string) []byte {
	content, err := os.ReadFile(p)
	policy := string(content)
	Check(err)

	variables := env.getVariables(policy)
	for _, _var := range variables {
		var val string
		if strings.ToLower(_var) == "tenant" {
			val = env.Tenant
		} else {
			val, err = env.getValue(_var)
			Check(err)
		}
		if val == "" || val == "null" {
			log.Fatalf("Variable %s is not provided in the config file. File: %s\n", _var, p)
		}
		re := regexp.MustCompile(fmt.Sprintf("{Settings:%s}", _var))
		policy = re.ReplaceAllString(policy, val)
	}

	return []byte(policy)
}

func (env Environment) getVariables(c string) []string {
	re := regexp.MustCompile(`{Settings:(.+)}`)
	m := re.FindAllStringSubmatch(c, -1)
	var cm []string
	seen := make(map[string]bool, len(m))
	for _, match := range m {
		if !seen[match[1]] {
			cm = append(cm, match[1])
			seen[match[1]] = true
		}
	}

	return cm
}

func (env Environment) getValue(v string) (string, error) {
	if env.Settings[v] == nil {
		return "", errors.New(fmt.Sprintf("Variable %s is not provided in settings", v))
	}

	return env.Settings[v].(string), nil
}

func (env Environment) Deploy(d string) {
	ps := NewPoliciesFromDir(path.Join(d, env.Name))
	bs := ps.GetBatch()

	g := NewGraphClientFromEnvironment(env)
	for i, b := range bs {
		log.Printf("Processing batch %d", i)
		g.UploadPolicies(b)
	}
}

func (env Environment) ListRemotePolicies() ([]string, error) {
	g := NewGraphClientFromEnvironment(env)
	return g.ListPolicies()
}

func (env Environment) DeleteRemotePolicies() error {
	g := NewGraphClientFromEnvironment(env)
	return g.DeletePolicies()
}

type Environments struct {
	e []Environment
	s string
	d string
}

func NewEnvironmentsFromConfig(p string, n string) Environments {
	var e []Environment

	c, err := os.ReadFile(p)
	if err != nil {
		panic(fmt.Sprintf("Could not read %s: %s", p, err.Error()))
	}

	err = yaml.Unmarshal(c, &e)
	if err != nil {
		panic(fmt.Sprintf("Could not unmarshal config from %s: %s", p, err.Error()))
	}

	es := Environments{
		e: e,
	}
	es.e = e
	es.filter(n)

	return es
}

func (es Environments) Build(s string, d string) {
	es.s = s
	es.d = d

	for _, e := range es.e {
		e.Build(es.s, es.d)
	}
}

func (es Environments) Deploy(d string) {
	es.d = d

	for _, e := range es.e {
		e.Deploy(es.d)
	}
}

func (es Environments) filter(n string) {
	var ne []Environment

	for _, e := range es.e {
		if n == "" || n == e.Name {
			ne = append(ne, e)
		}
	}

	es.e = ne
}

func (es Environments) ListRemotePolicies() (map[string][]string, []error) {
	var errs []error

	r := map[string][]string{}
	for _, e := range es.e {
		ps, err := e.ListRemotePolicies()
		if err != nil {
			errs = append(errs, err)
		}
		r[e.Name] = ps
	}

	return r, errs
}

func (es Environments) DeleteRemotePolicies() []error {
	var errs []error

	for _, e := range es.e {
		err := e.DeleteRemotePolicies()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to delete policies from environment %s: %s", e.Name, err)))
		}
	}

	return errs
}
