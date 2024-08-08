package main

import "github.com/seriouspoop/gopush/internal"

func main() {
	root, err := internal.NewRoot()
	if err != nil {
		panic(err)
	}

	err = root.RootCMD().Execute()
	if err != nil {
		panic(err)
	}
}
