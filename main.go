package main

import (
    "github.com/jessevdk/go-flags"
    "os"
)

type Options struct{
    Size int `short:"s" long:"size" description:"icon size" default:"200"`
    Path string `short:"o" long:"out" description:"Output filePath" default:"identile.png"`
    Algo string `short:"a" long:"algo" description:"Hash type" default:"md5" choice:"md5" choice:"sha1" choice:"sha256" choice:"sha512"`
    Salt string `long:"salt" description:"Hash salt"`
    // Dryrun bool `long:"dryrun" description:"Output to console"`
}

var SALT string = "_identile_salt_"

func main() {
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    parser.Name = "identile"
    parser.Usage = "TEXT [OPTIONS]"
    args, _ := parser.Parse()
    if len(args) == 0 {
        parser.WriteHelp(os.Stdout)
        os.Exit(1)
    }

    if opts.Salt == "" {
        opts.Salt = SALT
    }
    code := GetIdentileCodeByAlgo(os.Args[1], opts.Salt,GetIdentileAlgoByString(opts.Algo))
    renderer := NewSimpleRenderer(opts.Size)
    renderer.Render(code, opts.Path)
}
