package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Generate CSV of the format:
// git-repo,License

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: license-csv-gen minio|mc|operator")
	}
	project := os.Args[1]
	resp, err := http.Get(fmt.Sprintf("https://raw.githubusercontent.com/minio/%s/master/CREDITS", project))
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(resp.Body)
	for {
		s, err := r.ReadString('\n')
		if err != nil {
			return
		}
		s = strings.TrimSuffix(s, "\n")

		if strings.HasPrefix(s, "https") {
			fmt.Print(s, ",")
			s, err = r.ReadString('=')
			if err != nil {
				return
			}
			switch {
			case strings.Contains(s, "Apache"):
				fmt.Println("Apache-2.0")
			case strings.Contains(s, "MIT "), strings.Contains(s, "Permission is hereby granted, free of charge, to any person"):
				fmt.Println("MIT")
			case strings.Contains(s, "ISC License"):
				fmt.Println("ISC")
			case strings.Contains(s, "Mozilla Public License Version 2.0"), strings.Contains(s, "Mozilla Public License"):
				fmt.Println("MPL-2.0")
			case strings.Contains(s, "Redistribution and use in source and binary forms, with o"):
				fmt.Println("BSD")
			case strings.Contains(s, "Eclipse"):
				fmt.Println("Eclipse Public License 1.0 or Eclipse Distribution License 1.0 - dual licese")
			case strings.Contains(s, "Pascal S. de Kloe"):
				fmt.Println("Public Domain")
			default:
				fmt.Println("NA")
			}
			_, err = r.ReadString('\n')
			if err != nil {
				return
			}
		}
	}
}
