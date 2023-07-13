package components

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"wa-blast/loggers"
	//"./centrifugo"
	"wa-blast/components/rabbitmq"
)

// configItem represents component configuration item
type configItem struct {
	Name          string            `yaml:"name"`
	Type          string            `yaml:"type"`
	Configuration map[string]string `yaml:"configuration"`
}

var log = loggers.Get()

var components = make(map[string]interface{})

func Init() {
	// InitCron()
	// Read component config file
	var c []configItem
	bytes, err := ioutil.ReadFile("components.yml")
	if err != nil {
		fmt.Printf("Unable to read components.yml file. Error: %s\n", err.Error())
		os.Exit(16)
	}
	// Parse component list
	err = yaml.Unmarshal(bytes, &c)
	if err != nil {
		fmt.Printf("Unable to parse components.yml file. Error: %s\n", err.Error())
		os.Exit(17)
	}
	// Initiate and map component
	for _, item := range c {
		log.Debugf("Init component. Name: %s, Type: %s", item.Name, item.Type)
		// Get config
		cfg := item.Configuration
		// Switch type
		switch item.Type {
		case rabbitmq.ComponentType:
			// Get Username, Password, Host and Port
			username := cfg["username"]
			password := cfg["password"]
			host := cfg["host"]
			port := cfg["port"]
			// Init component
			components[item.Name] = rabbitmq.Init(username, password, host, port)
		default:
			fmt.Printf("Failed to parse component %s. Error: unknown component type\n", item.Name)
			os.Exit(18)
		}
	}
}

func GetRabbitMQ(name string) *rabbitmq.RabbitComponent {
	return components[name].(*rabbitmq.RabbitComponent)
}
