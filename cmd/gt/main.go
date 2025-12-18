package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

var (
	flagHelp         bool
	flagHelpFuncs    bool
	flagTemplateFile string
	flagTemplateData string
	flagOutputFile   string
	flagHeader       string
)

func init() {
	flag.BoolVar(&flagHelp, "h", false, "display help")
	flag.BoolVar(&flagHelp, "help", false, "display help")
	flag.BoolVar(&flagHelpFuncs, "funcs", false, "display a list of included template functions")
	flag.StringVar(&flagTemplateFile, "t", "", "input template file")
	flag.StringVar(&flagTemplateFile, "template", "", "input template file")
	flag.StringVar(&flagTemplateData, "d", "", "input template data file (optional, yaml format)")
	flag.StringVar(&flagTemplateData, "data", "", "input template data file (optional)")
	flag.StringVar(&flagOutputFile, "o", "", "output file")
	flag.StringVar(&flagOutputFile, "output", "", "output file")
	flag.StringVar(&flagHeader, "header", "", "add header to top of file before rendering template")
}

func checkFlags() error {
	if flagTemplateFile == "" {
		return errors.New("no template file specified")
	}
	if flagOutputFile == "" {
		return errors.New("no output file specified")
	}
	return nil
}

func main() {
	flag.Parse()
	if flagHelp {
		fmt.Println("gt template processor")
		flag.PrintDefaults()
		return
	}
	if err := checkFlags(); err != nil {
		fmt.Println("gt template processor")
		flag.PrintDefaults()
		log.Fatal(err)
	}
	if err := doExec(flagTemplateFile, flagTemplateData, flagOutputFile); err != nil {
		log.Fatal(err)
	}
}

func includeFunc(path string) (string, error) {
	data, err := os.ReadFile(path)
	return string(data), err
}

func newTemplate() *template.Template {
	tmpl := template.New("gt")
	funcMap := template.FuncMap{
		"include": includeFunc,
	}
	return tmpl.Funcs(funcMap)
}

func doExec(tmplPath, dataPath, outputPath string) error {
	tmplBytes, err := os.ReadFile(tmplPath)
	if err != nil {
		return fmt.Errorf("can't read %s: %s", tmplPath, err)
	}
	data := make(map[string]any)
	if dataPath != "" {
		dataBytes, err := os.ReadFile(dataPath)
		if err != nil {
			return fmt.Errorf("can't read %s: %s", dataPath, err)
		}
		if err := yaml.Unmarshal(dataBytes, &data); err != nil {
			return fmt.Errorf("couldn't parse yaml data: %s", err)
		}
	}
	outfile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("can't open %s: %s", outputPath, err)
	}
	if flagHeader != "" {
		if _, err := fmt.Fprintln(outfile, flagHeader); err != nil {
			return fmt.Errorf("error writing template header: %s", err)
		}
	}
	tmpl := newTemplate()
	if tmpl, err = tmpl.Parse(string(tmplBytes)); err != nil {
		return fmt.Errorf("can't parse template: %s", err)
	}
	if err := tmpl.Execute(outfile, data); err != nil {
		return fmt.Errorf("error executing template: %s", err)
	}
	return nil
}
