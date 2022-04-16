package cmd

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"

	"github.com/cortze/api-benchmark/pkg/config"
	"github.com/cortze/api-benchmark/pkg/requester"
)

var RunCommand = &cli.Command{
	Name:   "run",
	Usage:  "perform the api-benchmark on the given endpoint",
	Action: LaunchBenchmark,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:        "config-file",
			Usage:       "path to the <config.json> file used to configure the benchmark",
			EnvVars:     []string{"BENCHMARK_CONFIG_FILE_NAME"},
			DefaultText: config.ConfigFile,
			Value:       config.ConfigFile,
		}},
}

var log = logrus.WithField(
	"module", "RunCommand",
)

// CrawlAction is the function that is called when running `eth2`.
func LaunchBenchmark(c *cli.Context) error {
	log.Info("parsing flags")
	// check if a config file is set
	if !c.IsSet("config-file") {
		return errors.New("config-file with the benchmark configuration not provided")
	}

	log.Info("parsing configuration")
	// compose configuration form config .json
	conf, err := config.NewConfigFromJson(c.String("config-file"))
	if err != nil {
		return errors.Wrap(err, "unable to compose configuration for the benchmark")
	}

	log.Info("generating benchmark")
	// generate the benchmark with the given configuration
	benchmark := requester.NewBenchmark(c.Context, conf)

	// Check if we have to import or export the queries for the test
	switch conf.QueryBackup {
	case "import":
		log.Infof("importing queries from file %s", conf.QueryFile)
		benchmark.ImportQueryListFromFile(conf.QueryFile)
	default:
		// export by default
		log.Infof("exporting queries to file %s", conf.QueryFile)
		benchmark.ComposeQueryList()
		benchmark.ExportQueryList(conf.QueryFile)
	}

	log.Info("running the benchmark")
	// run the benchmark
	benchmark.Run()

	return nil
}
