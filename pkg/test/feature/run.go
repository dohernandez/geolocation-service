package feature

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"testing"

	"github.com/DATA-DOG/godog"
	"github.com/DATA-DOG/godog/colors"
	"github.com/stretchr/testify/assert"
)

// Used by init()
// nolint:gochecknoglobals
var (
	// RunGoDogTests is true if godog tests need to run
	RunGoDogTests bool
	stopOnFailure bool
	runWithTags   string
	runFeature    string
)

// This has to run on init to define -godog flag, otherwise "undefined flag" error happens
// nolint:gochecknoinits
func init() {
	flag.BoolVar(&RunGoDogTests, "godog", false, "Set this flag is you want to run godog BDD tests")
	flag.BoolVar(&stopOnFailure, "stop-on-failure", false,
		"Stop processing on first failed scenario.. Flag is passed to godog")

	descTagsOption := "Filter scenarios by tags. Expression can be:\n" +
		strings.Repeat(" ", 4) + "- " + colors.Yellow(`"@dev"`) + ": run all scenarios with wip tag\n" +
		strings.Repeat(" ", 4) + "- " + colors.Yellow(`"~@dev"`) +
		": exclude all scenarios with wip tag\n" +
		strings.Repeat(" ", 4) + "- " + colors.Yellow(`"@dev && ~@notImplemented"`) +
		": run wip scenarios, but exclude new\n" +
		strings.Repeat(" ", 4) + "- " + colors.Yellow(`"@dev,@undone"`) + ": run wip or undone scenarios"

	flag.StringVar(&runWithTags, "tag", "", descTagsOption)
	flag.StringVar(&runFeature, "feature", "",
		"Optional feature to run. Filename without the extension .feature")

	flag.Parse()
}

// RunSuite performs feature tests
func RunSuite(path string, featureContext func(t *testing.T, s *godog.Suite), t *testing.T) {
	if !RunGoDogTests {
		t.Skip(`Missing "-godog" flag, skipping integration test.`)
	}

	var paths []string

	if runFeature != "" {
		paths = []string{fmt.Sprintf("%s/%s.feature", path, runFeature)}
	} else {
		files, err := ioutil.ReadDir(path + "/")
		assert.NoError(t, err)

		paths = make([]string, 0, len(files))
		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".feature") {
				paths = append(paths, fmt.Sprintf("%s/%s", path, f.Name()))
			}
		}
	}

	randomSeed := rand.Int63()
	fmt.Println("Running test with random seed:", randomSeed)

	for _, path := range paths {
		path := path // Pinning ranged variable, more info: https://github.com/kyoh86/scopelint
		t.Run(path, func(t *testing.T) {
			status := godog.RunWithOptions(
				"Integration",
				func(s *godog.Suite) {
					featureContext(t, s)
				},
				godog.Options{
					Format:        "pretty",
					Paths:         []string{path},
					Randomize:     randomSeed,
					StopOnFailure: stopOnFailure,
					Tags:          runWithTags,
					Strict:        true,
				},
			)

			if status != 0 {
				assert.Fail(t, "one or more scenarios failed in feature: "+path)
			}
		})
	}
}
