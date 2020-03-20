/*
modification history
--------------------
2020/03/02 16:25:05, by NKztq, create
*/

// Package main is special.  It defines a
// standalone executable program, not a library.
// Within package main the function main is also
// special-it's where execution of the program begins.
// Whatever main does is what the program does.
package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/baidu/go-lib/log"
	log4go "github.com/baidu/go-lib/log/log4go"

	"github.com/NKztq/spider/conf"
	"github.com/NKztq/spider/crawler"
	"github.com/NKztq/spider/fetcher"
	"github.com/NKztq/spider/outputer"
	"github.com/NKztq/spider/seed"
)

var (
	help     *bool   = flag.Bool("h", false, "show help")
	confRoot *string = flag.String("c", "../conf", "root path of config file")
	logPath  *string = flag.String("l", "../log", "dir path of log")
	stdOut   *bool   = flag.Bool("s", false, "show log in stdout")
	showVer  *bool   = flag.Bool("v", false, "show version")
	debugLog *bool   = flag.Bool("d", false, "show debug level log msg")
)

var (
	confFileName = "spider.conf"
)

// main the function where execution of the program begins
func main() {
	var err error
	var logSwitch string

	flag.Parse()
	if *help {
		flag.PrintDefaults()
		return
	}
	if *showVer {
		fmt.Printf("mini_spider: version %s\n", version)
		return
	}

	// debug switch
	if *debugLog {
		logSwitch = "DEBUG"
	} else {
		logSwitch = "INFO"
	}

	err = initLog(logSwitch, logPath, stdOut)
	if err != nil {
		log.Logger.Error("main(): initLog(): %v", err)
		gracefullyExit(-1)
	}

	cfg, err := conf.LoadAndCheck(path.Join(*confRoot, confFileName))
	if err != nil {
		log.Logger.Error("main(): conf.LoadAndCheck(): %v", err)
		gracefullyExit(-2)
	}

	err = createOutputDirectory(cfg.Outputer.OutputDirectory)
	if err != nil {
		log.Logger.Error("main(): createOutputDirectory(): %v", err)
		gracefullyExit(-3)
	}

	seeds, err := seed.Load(cfg.Basic.UrlListFile)
	if err != nil {
		log.Logger.Error("main(): seed.Load(): %v", err)
		gracefullyExit(-4)
	}

	// create fetcher
	fetcher := fetcher.NewFetcher(cfg.Fetcher)

	// create outputer
	outputer, err := outputer.NewOutputer(cfg.Outputer)
	if err != nil {
		log.Logger.Error("main(): outputer.NewOutputer(): %v", err)
		gracefullyExit(-5)
	}

	// create crawler
	crawler := crawler.NewCrawler(cfg.Crawler, seeds, fetcher, outputer)

	// run crawler
	err = crawler.RunOnce()
	if err != nil {
		log.Logger.Error("main(): crawler.RunOnce(): %v", err)
		gracefullyExit(-6)
	}

	gracefullyExit(0)
}

func initLog(logSwitch string, logPath *string, stdOut *bool) error {
	/* initialize log   */
	/* set log buffer size  */
	log4go.SetLogBufferLength(10000)
	/* if blocking, log will be dropped */
	log4go.SetLogWithBlocking(false)

	err := log.Init("mini_spider", logSwitch, *logPath, *stdOut, "midnight", 5)
	if err != nil {
		return fmt.Errorf("gtc_api(): err in log.Init():%s\n", err.Error())
	}

	return nil
}

func createOutputDirectory(outputDirectory string) error {
	var err error

	_, err = os.Stat(outputDirectory)

	// mkdir if not exist
	if os.IsNotExist(err) {
		err = os.Mkdir(outputDirectory, os.ModePerm)
	}

	return err
}

func gracefullyExit(code int) {
	log.Logger.Close()

	// wait until log closed
	time.Sleep(1 * time.Second)

	os.Exit(code)
}
