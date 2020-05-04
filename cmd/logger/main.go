package main

import "github.com/spf13/viper"

func main() {
	viper.SetEnvPrefix("LOGGER")

	runner, err := newRunner()
	if err != nil {
		panic(err)
	}

	if err := runner.run(); err != nil {
		panic(err)
	}
}
