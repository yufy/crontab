package main

import "flag"

var configFile = flag.String("f", "configs/config.yaml", "set config file which viper will loading")

func main() {
	flag.Parse()

	w, err := CreateWorker(*configFile)
	if err != nil {
		panic(err)
	}

	if err := w.Run(); err != nil {
		panic(err)
	}

	w.AwaitSinal()
}
