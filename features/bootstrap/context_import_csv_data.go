package bootstrap

import (
	"io/ioutil"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/gherkin"
)

type importCsvDataContext struct {
}

// RegisterFileContext register execute create file steps
func RegisterFileContext(s *godog.Suite) {
	c := importCsvDataContext{}

	s.Step(`^there is a csv file in the path "([^"]*)" with the following content$`, c.thereIsACsvFileInThePathWithTheFollowingContent)
}

func (c *importCsvDataContext) thereIsACsvFileInThePathWithTheFollowingContent(filepath string, content *gherkin.DocString) error {
	return ioutil.WriteFile(filepath, []byte(content.Content), 0644)
}
