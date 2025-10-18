package phpunitxml

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
