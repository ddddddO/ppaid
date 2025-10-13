package main

import (
	"os"
	"text/template"
)

// PHPUnit version11用
type phpunitXMLv11Data struct {
	TestSuiteName     string
	TargetTestFiles   []string
	TargetCoverageDir string
}

const phpunitXMLv11Template = `<?xml version="1.0" encoding="UTF-8"?>
<phpunit xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
         xsi:noNamespaceSchemaLocation="https://schema.phpunit.de/11.0/phpunit.xsd"
         bootstrap="vendor/autoload.php"
         colors="true"
>
    <source>
        <include>
			<directory suffix=".php">{{.TargetCoverageDir}}</directory>
        </include>
        <exclude>
            <directory>./vendor</directory>
        </exclude>
    </source>

    <testsuites>
        <testsuite name="{{.TestSuiteName}}">
{{range .TargetTestFiles}}            <file>{{.}}</file>
{{end}}        </testsuite>
    </testsuites>
</phpunit>
`

func generatePHPUnitXML(targetTests []string, targetCoverageDir string) error {
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
