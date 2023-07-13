package util

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bwmarrin/snowflake"

	"wa-blast/flags"
)

// Constants
const envKeyNodeNo = flags.EnvPrefix + "_NODE_NO"

// Snowflake Generator node
var node *snowflake.Node

// initSnowflake initiate Snowflake generator node in singleton pattern. App will exit if an error occurred
func initSnowflake() {
	// Get node number from env
	// _, ok := os.LookupEnv(envKeyNodeNo)
	// if !ok {
	// 	fmt.Printf("%s is not set in system environment\n", envKeyNodeNo)
	// 	os.Exit(7)
	// }
	// Parse node number
	_, err := strconv.ParseInt("1", 10, 64)
	if err != nil {
		fmt.Printf("unable to parse node number. Error: %s\n", err.Error())
		os.Exit(8)
	}
	// Create snowflake node
	n, err := snowflake.NewNode(1)
	if err != nil {
		fmt.Printf("unable to initiate Snowflake Node. Error: %s\n", err.Error())
		os.Exit(9)
	}
	// Set node
	node = n
}

// NewId generate Twitter Snowflake ID
func NewId() string {
	return node.Generate().String()
}
