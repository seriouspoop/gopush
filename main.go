package main

import "github.com/seriouspoop/gopush/internal"

func main() {
	root, err := internal.NewRoot()
	if err != nil {
		return
	}

	err = root.RootCMD().Execute()
	if err != nil {
		return
	}
}
