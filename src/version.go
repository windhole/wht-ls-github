package main

var version = "dev"

func versionString() string {
	if version == "" {
		return "dev"
	}
	return version
}
