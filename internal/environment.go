package internal

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"com.schumann-it.go-ieftool/internal/msgraph/trustframework"
	"github.com/spf13/pflag"
	"gopkg.in/yaml.v3"
)

//go:embed msgraph/trustframework/ApplicationPatchIdentityFramework.json
var iefApplicationPatch []byte

//go:embed msgraph/trustframework/ApplicationPatchSaml.json
var samlApplicationPatch []byte

type EnvironmentSaml struct {
	AppObjectId *string `yaml:"appObjectId,omitempty"`
	MetadataUrl *string `yaml:"metadataUrl,omitempty"`
	CertPath    *string `yaml:"certPath,omitempty"`
	Cert        []byte
}

type Environment struct {
	Name                                   string                 `yaml:"name"`
	IsProduction                           bool                   `yaml:"isProduction"`
	Tenant                                 string                 `yaml:"tenant"`
	TenantId                               string                 `yaml:"tenantId"`
	ClientId                               string                 `yaml:"clientId"`
	IdentityExperienceFrameworkAppObjectId *string                `yaml:"identityExperienceFrameworkAppObjectId,omitempty"`
	Saml                                   *EnvironmentSaml       `yaml:"saml,omitempty"`
	Settings                               map[string]interface{} `yaml:"settings"`
}

func (env Environment) Build(s string, d string) error {
	var errs Errors
	root := s
	err := filepath.WalkDir(s, func(p string, e fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if e.IsDir() {
			return nil
		}
		if filepath.Ext(e.Name()) == ".xml" {
			t := path.Join(d, strings.ReplaceAll(p, root, env.Name))
			c, ve := env.replaceVariables(p)
			if ve != nil {
				errs = append(errs, ve)
				return nil
			}
			ve = os.MkdirAll(filepath.Dir(t), os.ModePerm)
			if ve != nil {
				errs = append(errs, ve)
				return nil
			}
			log.Printf("Compiled %s", t)
			if env.IsProduction {
				// @TODO remove debug code
				log.Print("Removed debug parameters as this is a prod environment.")
			}
			ve = os.WriteFile(t, c, os.ModePerm)
			if ve != nil {
				errs = append(errs, ve)
				return nil
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (env Environment) replaceVariables(p string) ([]byte, error) {
	content, err := os.ReadFile(p)
	policy := string(content)
	if err != nil {
		return nil, err
	}

	var errs Errors
	for _, v := range env.variables(policy) {
		val, ve := env.value(v)
		if ve != nil {
			errs = append(errs, fmt.Errorf("%s: but required in source %s", ve.Error(), p))
			continue
		}
		if val == "" || val == "null" {
			errs = append(errs, fmt.Errorf("variable '%s' must not be empty: source %s", v, p))
			continue
		}
		re := regexp.MustCompile(fmt.Sprintf("{Settings:%s}", v))
		policy = re.ReplaceAllString(policy, val)
	}

	if errs.HasErrors() {
		return nil, errs.Format()
	}

	return []byte(policy), nil
}

func (env Environment) variables(c string) []string {
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

func (env Environment) value(n string) (string, error) {
	switch n {
	case "Tenant":
		return env.Tenant, nil
	default:
		if env.Settings[n] == nil {
			return "", fmt.Errorf("variable '%s' is not provided in settings", n)
		}

		return env.Settings[n].(string), nil
	}
}

func (env Environment) Deploy(d string) error {
	ps, err := trustframework.NewPoliciesFromDir(path.Join(d, env.Name))
	if err != nil {
		return err
	}
	bs := ps.GetBatch()

	g, err := NewGraphClientFromEnvironment(env)
	if err != nil {
		return err
	}

	for i, b := range bs {
		log.Printf("Processing batch %d", i)
		g.UploadPolicies(b)
	}

	return nil
}

func (env Environment) ListRemotePolicies() ([]string, error) {
	g, err := NewGraphClientFromEnvironment(env)
	if err != nil {
		return nil, err
	}

	return g.ListPolicies()
}

func (env Environment) DeleteRemotePolicies() error {
	g, err := NewGraphClientFromEnvironment(env)
	if err != nil {
		return err
	}

	return g.DeletePolicies()
}

func (env Environment) FixAppRegistrations() error {
	g, err := NewGraphClientFromEnvironment(env)
	if err != nil {
		return err
	}

	if env.IdentityExperienceFrameworkAppObjectId == nil {
		return fmt.Errorf("please specify identityExperienceFrameworkObjectId in envirnment")
	}
	err = g.FixAppRegistration(*env.IdentityExperienceFrameworkAppObjectId, iefApplicationPatch)
	if err != nil {
		log.Fatalln(err)
	}

	if env.Saml != nil && env.Saml.AppObjectId != nil {
		var p map[string]interface{}
		err = json.Unmarshal(samlApplicationPatch, &p)
		if env.Saml.MetadataUrl != nil {
			p["samlMetadataUrl"] = env.Saml.MetadataUrl
		}
		patch, err := json.Marshal(p)
		if err != nil {
			log.Fatalln(err)
		}
		err = g.FixAppRegistration(*env.Saml.AppObjectId, patch)
		if err != nil {
			log.Fatalln(err)
		}
	}

	return nil
}

func (env Environment) CreateKeySets() error {
	g, err := NewGraphClientFromEnvironment(env)
	if err != nil {
		return err
	}

	es := strings.ReplaceAll(fmt.Sprintf("B2C_SAML_CERT_PW_%s", strings.ToUpper(env.Name)), "-", "_")

	return g.CreateKeySets(env.Saml.Cert, es)
}

func (env Environment) DeleteKeySets() interface{} {
	g, err := NewGraphClientFromEnvironment(env)
	if err != nil {
		return err
	}

	return g.DeleteKeySets()
}

type Environments struct {
	e []Environment
	s string
	d string
}

func MustNewEnvironmentsFromFlags(f *pflag.FlagSet) *Environments {
	cf, err := f.GetString("config")
	if err != nil {
		log.Fatalf("could not parse flag 'config': \n%s", err.Error())
	}

	en, err := f.GetString("environment")
	if err != nil {
		log.Fatalf("could not parse flag 'environment': \n%s", err.Error())
	}

	e, err := NewEnvironmentsFromConfig(cf, en)
	if err != nil {
		log.Fatalf("could not read environments config: \n%s", err.Error())
	}

	return e
}

func NewEnvironmentsFromConfig(p string, n string) (*Environments, error) {
	var e []Environment

	c, err := os.ReadFile(p)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not read %s: %s", p, err.Error()))
	}

	err = yaml.Unmarshal(c, &e)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Could not unmarshal config from %s: %s", p, err.Error()))
	}

	es := Environments{
		e: e,
	}
	es.e = e
	es.filter(n)

	for i, _ := range es.e {
		if es.e[i].Saml.CertPath != nil {
			sp, err := filepath.Abs(*es.e[i].Saml.CertPath)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Could not find saml cert %s: %s", *es.e[i].Saml.CertPath, err.Error()))
			}
			b, err := os.ReadFile(sp)
			if err != nil {
				return nil, errors.New(fmt.Sprintf("Could not read saml cert %s: %s", p, err.Error()))
			}
			es.e[i].Saml.Cert = b
		}
	}

	return &es, nil
}

func (es *Environments) Len() int {
	return len(es.e)
}

func (es *Environments) Build(s string, d string) error {
	es.s = s
	es.d = d

	for _, e := range es.e {
		err := e.Build(es.s, es.d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *Environments) Deploy(d string) error {
	es.d = d

	for _, e := range es.e {
		err := e.Deploy(es.d)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *Environments) filter(n string) {
	var ne []Environment

	for _, e := range es.e {
		if n == "" || n == e.Name {
			ne = append(ne, e)
		}
	}

	es.e = ne
}

func (es *Environments) ListRemotePolicies() (map[string][]string, error) {
	var errs Errors

	r := map[string][]string{}
	for _, e := range es.e {
		ps, err := e.ListRemotePolicies()
		if err != nil {
			errs = append(errs, err)
		}
		r[e.Name] = ps
	}

	if errs.HasErrors() {
		return nil, errs.Format()
	}

	return r, nil
}

func (es *Environments) DeleteRemotePolicies() error {
	var errs Errors

	for _, e := range es.e {
		err := e.DeleteRemotePolicies()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to delete policies from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (es *Environments) FixAppRegistrations() error {
	var errs Errors

	for _, e := range es.e {
		err := e.FixAppRegistrations()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to fix app registrations from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (es *Environments) CreateKeySets() error {
	var errs Errors

	for _, e := range es.e {
		err := e.CreateKeySets()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to key sets from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}

func (es *Environments) DeleteKeySets() error {
	var errs Errors

	for _, e := range es.e {
		err := e.DeleteKeySets()
		if err != nil {
			errs = append(errs, errors.New(fmt.Sprintf("Failed to key sets from environment %s: %s", e.Name, err)))
		}
	}

	if errs.HasErrors() {
		return errs.Format()
	}

	return nil
}
