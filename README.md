# cfg_flags

Go package to use an INI configuration file alongside the "flag" package.

Inspired by [iniflags](https://github.com/vharitonsky/iniflags), cfg_flags goal is to be simple, without any advanced features like configuration reloading.
It is only meant to integrate an INI file with golang's flag package, nothing more.

## Usage

Just replace `flag.Parse()` by `cfg_flags.Parse()` in your code.
`cfg_flags.Parse()` will return any error that occurs while parsing the configuration file, it is up to you to handle them.

Code sample:

    if err := cfg_flags.Parse(); err != nil {
        flag.Usage()
        log.Fatal(err)
    }


In order to specify the configuration file to use, you need to set the `-config` flag, e.g. `my_app -config my_config.ini`.
