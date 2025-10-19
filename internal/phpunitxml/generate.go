package phpunitxml

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/ddddddO/ppaid/internal/command"
)

const OutputPHPUnitXML = "phpunitxml_generated_by_ppaid.xml"

type phpunitXMLData struct {
	TestSuiteName     string
	TargetTestFiles   []string
	TargetCoverageDir string
}

func Generate(commandToSpecifyBeforePHPCommand string, targetTests []string, targetCoverageDir string) error {
	insertData := &phpunitXMLData{
		TestSuiteName:     "PPAID",
		TargetTestFiles:   targetTests,
		TargetCoverageDir: targetCoverageDir,
	}

	majorVersion, err := command.ParsePHPUnitVersion(commandToSpecifyBeforePHPCommand)
	if err != nil {
		return err
	}

	if err := generatePHPUnitXMLFromExisting(insertData, majorVersion); err == nil {
		return nil
	} else if !errors.Is(err, &ErrReadPHPUnitXML{}) {
		return err
	}

	var t *template.Template
	switch majorVersion {
	case 9:
		t = template.Must(template.New("xxx").Parse(phpunitXMLv9Template))
	case 11:
		t = template.Must(template.New("xxx").Parse(phpunitXMLv11Template))
	default:
		return fmt.Errorf("unsupported PHPUnit version: %d", majorVersion)
	}

	f, err := os.Create(OutputPHPUnitXML)
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, insertData)
}
