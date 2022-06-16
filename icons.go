package main

import "io/ioutil"

func GetIconData(filename string) []byte {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	return data
}
