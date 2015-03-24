package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"path"

	"gopkg.in/yaml.v2"
)

// scancmd executes the scan command
func scancmd(project, file, format string, all bool) error {

	// if no project directory is given, default to current directory
	if len(project) == 0 {
		project = "."
	}

	// check if any result format has been specified
	if len(format) == 0 {

		// if no file to write to has been specifed default to text,
		// otherwise attempt to determine the file type by file extension
		if len(file) == 0 {
			format = "txt"
		} else {
			ext := path.Ext(file)
			format = ext[1:]
		}
	}

	// scan the go project provided
	pkgs, err := scan(project)
	if err != nil {
		return err
	}

	if !all {

		// remove standard packages
		pkgs = removestdpkgs(pkgs)
	}

	projectpath, err := importpath(project)
	if err != nil {
		return err
	}

	// filter out packages internal to the project
	pkgs = removeprefix(projectpath, pkgs)

	// create an slice of bytes to print or write results
	var b []byte

	// switch on format
	switch format {

	case "txt": // if text, use a byte.Buffer to format package paths
		var buff bytes.Buffer
		for _, pkg := range pkgs {
			buff.WriteString(pkg + "\n")
		}
		b = buff.Bytes()

	case "xml": // marshal to xml with indentation
		b, err = xml.MarshalIndent(pkgs, "", "  ")
		if err != nil {
			return err
		}

	case "yaml", "yml": // marshal to yml with indentation
		b, err = yaml.Marshal(pkgs)
		if err != nil {
			return err
		}

	case "json": // marshal to json with indentation
		b, err = json.MarshalIndent(pkgs, "", "  ")
		if err != nil {
			return err
		}

	default: // error out on unsupported formats
		return errors.New("unsupported format: " + format)
	}

	// if a file to write to was specified, write to it.
	if len(file) > 0 {
		if err := ioutil.WriteFile(file, b, 0644); err != nil {
			return err
		}
		return nil
	}

	// no file specified, just print results to screen.
	fmt.Printf("%s", b)
	return nil
}