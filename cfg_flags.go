// The MIT License (MIT)
//
// Copyright (c) 2015 Arnaud Vazard
//
// See LICENSE file.

// Package to use an INI configuration file alongside the "flag" package
package cfg_flags

import (
	"flag"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
)

var configFlag *string

func init() {
	// Flag to set the configuration file path via command line
	configFlag = flag.String("config", "", "Configuration file path")
}

func Parse(configFile string) error {
	// Parse the command line flags
	flag.Parse()

	var configFilePath string
	if *configFlag != "" {
		configFilePath = *configFlag
	} else if configFile == "" {
		return fmt.Errorf("cfg_flags: No config file path was given")
	} else {
		configFilePath = configFile
	}

	// Get the flags from the configuration file
	valuesFromFile, err := getValuesFromFile(configFilePath)
	if err != nil {
		return err
	}
	// Get the missing flags (That is the flags that have not been set via the command line)
	missingFlags := getMissingFlags()

	for key, value := range valuesFromFile {
		// Look up the flag from the configuration file in the flag list
		f := flag.Lookup(key)
		// If no flag matching "key" was found, return false
		if f == nil {
			return fmt.Errorf("cfg_flags: Unknow flag found in the configuration file (%q)\n", key)
		}
		// Iterate over the list of flags that are not yet set
		for _, v := range missingFlags {
			// If the flag from the file is found in the "missing" slice
			if f.Name == v {
				// If the value from the file is different from the default value for the flag, we set the value for this flag
				if f.Value.String() != value {
					// If an error happen, return false
					if err := f.Value.Set(value); err != nil {
						return fmt.Errorf("cfg_flags: Error while parsing flag %q (error: %q)\n", key, err)
					}
				}
				// if we found the flag there is no need to continue to loop
				break
			}
		}
	}
	return nil
}

func getValuesFromFile(configFile string) (map[string]string, error) {
	// Read the file to a byte slice
	content, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("cfg_flags: Error while reading configuration file %q: %s\n", configFile, err)
	}

	valuesFromFile := make(map[string]string)
	// Iterate over the file content (split on "\n")
	for i, line := range strings.Split(string(content), "\n") {
		line = strings.TrimSpace(line)
		// If the line is empty, start with a comment or is a section, we do nothing
		if line == "" || line[0] == ';' || line[0] == '#' || line[0] == '[' {
			continue
		}
		// Split the string on the "=" character
		fields := strings.Split(line, "=")
		if len(fields) > 2 {
			return nil, fmt.Errorf("cfg_flags: There is more than one \"=\" in the following line from the configuration file: %q (line %d)\n", line, i)
		} else if len(fields) < 2 {
			return nil, fmt.Errorf("cfg_flags: There is no \"=\" in the following line from the configuration file: %q (line %d)\n", line, i)
		}
		err := cleanString(&fields[1])
		if err != nil {
			return nil, fmt.Errorf("cfg_flags: Error while processing line %d from the configuration file: %s", i, err)
		}
		// Return map: The key is the first field with leading and trailing spaces removed, the value is the second field "cleaned"
		valuesFromFile[strings.TrimSpace(fields[0])] = fields[1]
	}
	return valuesFromFile, nil
}

func getMissingFlags() []string {
	var (
		set, missing []string
		found        bool = false
	)
	// Visit only the flags that have been set
	flag.Visit(func(f *flag.Flag) {
		set = append(set, f.Name)
	})
	// Visit all the flags, even those not set
	flag.VisitAll(func(f *flag.Flag) {
		for _, v := range set {
			if v == f.Name {
				found = true
				break
			}
		}
		// If we don't find the flag in the slice of already set flags, we add it to the missing slice
		if !found {
			missing = append(missing, f.Name)
		}
		found = false
	})
	return missing
}

func cleanString(str *string) error {
	// Trim the spaces from the string and return if the resulting string is empty
	tmp := strings.TrimSpace(*str)
	if len(tmp) == 0 {
		return nil
	}
	// If the string is not quoted, we remove the trailing comments (Beginning with # or ;)
	// If the string is quoted, we unquote it
	if tmp[0] != '"' {
		tmp = strings.Split(strings.Split(tmp, "#")[0], ";")[0]
	} else {
		result, err := strconv.Unquote(tmp)
		if err != nil {
			return err
		}
		tmp = result
	}
	*str = strings.TrimSpace(tmp)
	return nil
}
