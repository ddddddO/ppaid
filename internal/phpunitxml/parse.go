package phpunitxml

import (
	"encoding/xml"
	"fmt"
	"os"
)

const PHPUnitXMLPath = "phpunit.xml"

type ErrReadPHPUnitXML struct {
	err error
}

func (e *ErrReadPHPUnitXML) Error() string {
	return fmt.Errorf("failed to Open phpunit.xml; %w", e.err).Error()
}

func generatePHPUnitXMLFromExisting(insertData *phpunitXMLData, phpunitMajorVersion int) error {
	phpunitXML, err := parsePHPUnitXML()
	if err != nil {
		return err
	}

	phpunitXML.XmlnsXsi = "http://www.w3.org/2001/XMLSchema-instance"
	phpunitXML.XsiNoNamespaceSchemaLocation = "vendor/phpunit/phpunit/phpunit.xsd"

	// insert testsuite
	files := make([]File, len(insertData.TargetTestFiles))
	for i := range insertData.TargetTestFiles {
		files[i] = File{insertData.TargetTestFiles[i]}
	}
	if phpunitXML.TestSuites == nil {
		phpunitXML.TestSuites = &TestSuites{
			TestSuite: []TestSuite{
				{Name: insertData.TestSuiteName, Files: files},
			},
		}
	} else {
		phpunitXML.TestSuites.TestSuite = append(phpunitXML.TestSuites.TestSuite, TestSuite{Name: insertData.TestSuiteName, Files: files})
	}

	// insert target coverage
	switch phpunitMajorVersion {
	case 9:
		dir := Directory{Suffix: ".php", Content: insertData.TargetCoverageDir}
		if phpunitXML.Coverage == nil {
			phpunitXML.Coverage = &Coverage{Include: IncludeExclude{Directories: []Directory{dir}}}
		} else {
			phpunitXML.Coverage.Include.Directories = append(phpunitXML.Coverage.Include.Directories, dir)
		}
	case 11:
		dir := Directory{Suffix: ".php", Content: insertData.TargetCoverageDir}
		if phpunitXML.Source == nil {
			phpunitXML.Source = &Source{Include: IncludeExclude{Directories: []Directory{dir}}}
		} else {
			phpunitXML.Source.Include.Directories = append(phpunitXML.Source.Include.Directories, dir)
		}
	default:
		if len(insertData.TargetCoverageDir) == 0 {
			return fmt.Errorf("unsupported PHPUnit version: %d", phpunitMajorVersion)
		}
	}

	output, err := xml.MarshalIndent(phpunitXML, "", "	")
	if err != nil {
		return err
	}
	fixedXML := []byte(xml.Header)
	fixedXML = append(fixedXML, output...)
	if err := os.WriteFile(OutputPHPUnitXML, fixedXML, 0644); err != nil {
		return err
	}
	return nil
}

func parsePHPUnitXML() (*PHPUnitXML, error) {
	data, err := os.ReadFile(PHPUnitXMLPath)
	if err != nil {
		// 読み取れなかっただけなら、まっさらに新規作成するからエラー定義して返す
		return nil, &ErrReadPHPUnitXML{err}
	}

	var phpunitXML PHPUnitXML
	if err := xml.Unmarshal(data, &phpunitXML); err != nil {
		return nil, err
	}

	return &phpunitXML, nil
}

// PHPUnitXML - ルート要素 <phpunit ...>
type PHPUnitXML struct {
	XMLName xml.Name `xml:"phpunit"`

	// --- ネームスペース属性（必須） ---
	XmlnsXsi                     string `xml:"xmlns:xsi,attr"`
	XsiNoNamespaceSchemaLocation string `xml:"xsi:noNamespaceSchemaLocation,attr"`

	// --- ルートの属性 ---
	Bootstrap                           string `xml:"bootstrap,attr"`
	CacheResult                         string `xml:"cacheResult,attr"`
	Colors                              string `xml:"colors,attr"`
	ExecutionOrder                      string `xml:"executionOrder,attr"`
	FailOnRisky                         string `xml:"failOnRisky,attr"`
	FailOnWarning                       string `xml:"failOnWarning,attr"`
	FailOnEmptyTestSuite                string `xml:"failOnEmptyTestSuite,attr"`
	BeStrictAboutChangesToGlobalState   string `xml:"beStrictAboutChangesToGlobalState,attr"`
	BeStrictAboutOutputDuringTests      string `xml:"beStrictAboutOutputDuringTests,attr"`
	CacheDirectory                      string `xml:"cacheDirectory,attr"`
	BeStrictAboutCoverageMetadata       string `xml:"beStrictAboutCoverageMetadata,attr"`
	DisplayDetailsOnPhpunitDeprecations string `xml:"displayDetailsOnPhpunitDeprecations,attr"`
	DisplayDetailsOnIncompleteTests     string `xml:"displayDetailsOnIncompleteTests,attr"`
	StopOnFailure                       string `xml:"stopOnFailure,attr,omitempty"`
	// 他の多くの属性を必要に応じてここに追加

	// --- 主要な子要素（ポインタやスライスでオプションとして定義） ---
	Extensions *Extensions `xml:"extensions,omitempty"`
	TestSuites *TestSuites `xml:"testsuites,omitempty"`
	Groups     *Groups     `xml:"groups,omitempty"`

	// ここphpunit versionで変わる<coverage><source>とか. <filter>もどこかある？
	Coverage *Coverage `xml:"coverage,omitempty"`
	Source   *Source   `xml:"source,omitempty"`

	Php       *Php       `xml:"php,omitempty"`
	Listeners *Listeners `xml:"listeners,omitempty"`
	// Loggingは通常空タグだが、子要素を持つことも可能
	Logging *Logging `xml:"logging,omitempty"`
}

// Extensions <extensions>
type Extensions struct {
	Bootstrap []ExtensionBootstrap `xml:"bootstrap"`
}

// ExtensionBootstrap <bootstrap class="...">
type ExtensionBootstrap struct {
	Class string `xml:"class,attr"`
}

// TestSuites <testsuites>
type TestSuites struct {
	TestSuite []TestSuite `xml:"testsuite"`
}

// TestSuite <testsuite name="...">
type TestSuite struct {
	Name        string      `xml:"name,attr"`
	Files       []File      `xml:"file,omitempty"`
	Directories []Directory `xml:"directory,omitempty"`
	Excludes    []Exclude   `xml:"exclude,omitempty"`
}

// Directory <directory suffix="...">
type Directory struct {
	Suffix  string `xml:"suffix,attr,omitempty"`
	Content string `xml:",chardata"`
}

// Exclude <exclude>
type Exclude struct {
	Content string `xml:",chardata"`
}

// Groups <groups>
type Groups struct {
	Exclude []Group `xml:"exclude>group,omitempty"`
	// Include []Group `xml:"include>group,omitempty"` // 挿入対象の定義
}

// Group <group>
type Group struct {
	Name string `xml:",chardata"`
}

// Coverage <coverage>
type Coverage struct {
	Include IncludeExclude `xml:"include,omitempty"`
	Exclude IncludeExclude `xml:"exclude,omitempty"`
	Report  Report         `xml:"report,omitempty"`
}

type Source struct {
	Include IncludeExclude `xml:"include,omitempty"`
	Exclude IncludeExclude `xml:"exclude,omitempty"`
}

// IncludeExclude <include> / <exclude> (coverage用)
type IncludeExclude struct {
	Directories []Directory `xml:"directory,omitempty"`
	Files       []File      `xml:"file,omitempty"`
}

// File <file>
type File struct {
	Content string `xml:",chardata"`
}

// Report <report>
type Report struct {
	Clover *Clover `xml:"clover,omitempty"`
	// ... 他のレポート形式をここに追加
}

// Clover <clover outputFile="...">
type Clover struct {
	OutputFile string `xml:"outputFile,attr"`
}

// Php <php>
type Php struct {
	Ini         []Ini         `xml:"ini,omitempty"`
	Env         []Env         `xml:"env,omitempty"`
	Const       []Const       `xml:"const,omitempty"`
	IncludePath []IncludePath `xml:"includePath,omitempty"`
}

// Ini <ini name="..." value="...">
type Ini struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// Env <env name="..." value="...">
type Env struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
	Force string `xml:"force,attr,omitempty"`
}

// Const <const name="..." value="...">
type Const struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// IncludePath <includePath>
type IncludePath struct {
	Content string `xml:",chardata"`
}

// Listeners <listeners>
type Listeners struct {
	Listener []Listener `xml:"listener,omitempty"`
}

// Listener <listener class="...">
type Listener struct {
	Class     string             `xml:"class,attr"`
	File      string             `xml:"file,attr,omitempty"`
	Arguments []ListenerArgument `xml:"arguments>argument,omitempty"`
}

// ListenerArgument <argument>
type ListenerArgument struct {
	Content string `xml:",chardata"`
	Type    string `xml:"type,attr,omitempty"`
}

// Logging <logging>
type Logging struct {
	// ReportやReportListenerなどの子要素が存在する
}
