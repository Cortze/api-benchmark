package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	cli "github.com/urfave/cli/v2"
)

var RewardsCommand = &cli.Command{
	Name:   "rewards",
	Usage:  "calculate rewards for a validator in a given epoch list",
	Action: LaunchRewardCalculator,
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "outfile",
			Usage: "output csv file",
		},
		&cli.StringFlag{
			Name:  "init-slot",
			Usage: "init slot from where to start",
		},
		&cli.StringFlag{
			Name:  "final-slot",
			Usage: "init slot from where to start",
		},
		&cli.StringFlag{
			Name:  "validator-indexes",
			Usage: "init slot from where to start",
		}},
}

var logRewardsRewards = logrus.WithField(
	"module", "RewardsCommand",
)

var timeout = 30
var query = "http://localhost:5052/eth/v1/beacon/states/{slot}/validator_balances?id={indexes}"

// CrawlAction is the function that is called when running `eth2`.
func LaunchRewardCalculator(c *cli.Context) error {
	logRewardsRewards.Info("parsing flags")
	// check if a config file is set
	if !c.IsSet("init-slot") {
		return errors.New("final slot not provided")
	}
	if !c.IsSet("final-slot") {
		return errors.New("final slot not provided")
	}
	if !c.IsSet("validator-indexes") {
		return errors.New("validator indexes not provided")
	}
	if !c.IsSet("outfile") {
		return errors.New("outputfile no provided")
	}
	outputFile := c.String("outfile")
	initSlot := c.String("init-slot")
	finalSlot := c.String("final-slot")
	validatorIndexFile := c.String("validator-indexes")
	validatorIndexInt := make([]int, 0)
	validatorIndex := make([]string, 0)

	// open file and read all the indexes
	fbytes, err := ioutil.ReadFile(validatorIndexFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(fbytes, &validatorIndexInt)

	logRewardsRewards.Infof("Requesting balances of %d validators", len(validatorIndexInt))

	for _, integer := range validatorIndexInt {
		validatorIndex = append(validatorIndex, fmt.Sprintf("%d", integer))
	}

	initSlotInt, err := strconv.Atoi(initSlot)
	if err != nil {
		return err
	}
	finalSlotInt, err := strconv.Atoi(finalSlot)
	if err != nil {
		return err
	}

	csvFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	// compose lines
	headers := "slot,total"
	for i := 0; i < len(validatorIndex); i++ {
		headers = headers + "," + validatorIndex[i]
	}

	_, err = csvFile.WriteString(headers + "\n")
	if err != nil {
		return err
	}

	// to calculate rewards
	prevBalance := make([]int, len(validatorIndex))
	var i int = 0
	for s := initSlotInt; s < finalSlotInt; s = s + 32 {
		balanceArray := make([]int, len(validatorIndex))

		// generate a http.Client to manage timeouts
		httpCli := http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		}

		composed_query := GenerateQuery(query, fmt.Sprintf("%d", s), validatorIndex)
		logRewardsRewards.Infof("Requesting query %s", composed_query)

		resp, err := httpCli.Get(composed_query)
		if err != nil {
			return err
		}

		// unmarshal json response
		var valResp map[string]interface{}
		fmt.Println(resp.Body)
		respBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		err = json.Unmarshal(respBytes, &valResp)
		if err != nil {
			return err
		}

		balances := ParseBalancesResp(valResp)

		fmt.Println(balances)
		if i == 0 {
			// the prev values are empty
			for idx, bal := range balances {
				fbal, err := strconv.Atoi(bal.balance)
				if err != nil {
					return err
				}
				prevBalance[idx] = fbal
				balanceArray[idx] = 0
			}
		} else {
			// calculate the difference with the previous epoch
			for idx, bal := range balances {
				fbal, err := strconv.Atoi(bal.balance)
				if err != nil {
					return err
				}
				if (fbal - prevBalance[idx]) < 0 {
					balanceArray[idx] = 0
				} else {
					balanceArray[idx] = fbal - prevBalance[idx]
				}
				prevBalance[idx] = fbal
			}
		}
		i++
		fmt.Println(balanceArray)
		total := int64(0)
		row := fmt.Sprintf("%d", s) // slot
		for _, item := range balanceArray {
			total = total + int64(item)
		}
		row = row + "," + fmt.Sprintf("%d", total)

		for _, item := range balanceArray {
			row = row + "," + fmt.Sprintf("%d", item)
		}
		_, err = csvFile.WriteString(row + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

func GenerateQuery(base_query string, slot string, indexes []string) string {
	var composed_query string

	if len(indexes) == 0 {
		return ""
	}

	composed_query = strings.Replace(base_query, "{slot}", slot, 1)
	idxs := indexes[0]
	for i := 1; i < len(indexes); i++ {
		idxs = idxs + "," + indexes[i]
	}

	composed_query = strings.Replace(composed_query, "{indexes}", idxs, 1)
	return composed_query
}

type Balances struct {
	index   string `json:"index"`
	balance string `json:"balances"`
}

func ParseBalancesResp(resp map[string]interface{}) []Balances {
	balances := make([]Balances, 0)

	data := resp["data"].([]interface{})
	for _, item := range data {
		bal := item.(map[string]interface{})
		b := Balances{
			index:   bal["index"].(string),
			balance: bal["balance"].(string),
		}
		balances = append(balances, b)
	}

	return balances
}
