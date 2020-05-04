package main

func main() {
	runner, err := newRunner()
	if err != nil {
		panic(err)
	}

	if err := runner.run(); err != nil {
		panic(err)
	}
}
