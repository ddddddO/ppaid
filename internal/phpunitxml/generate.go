package phpunitxml

import (
	"os"
	"text/template"
)

func Generate(targetTests []string, targetCoverageDir string) error {
	d := &phpunitXMLv11Data{
		TestSuiteName:     "PPAID",
		TargetTestFiles:   targetTests,
		TargetCoverageDir: targetCoverageDir,
	}

	t := template.Must(template.New("xxx").Parse(phpunitXMLv11Template))

	f, err := os.Create("tmp_phpunit.xml")
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, d)
}
