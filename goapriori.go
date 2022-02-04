package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/jimmywmt/goapriori/algorithm"
	log "github.com/sirupsen/logrus"
)

func init() {

	log.SetFormatter(&log.TextFormatter{
		ForceQuote:      true,
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	log.SetLevel(log.ErrorLevel)

}

func main() {
	fmt.Println("Apriori written in GoLang (jimmywmt)")

	var fileName string
	var minsup float64
	var showlog bool

	flag.StringVar(&fileName, "f", "", "DB file")
	flag.Float64Var(&minsup, "m", 0, "Minimal Support (0~1)")
	flag.BoolVar(&showlog, "l", false, "Show log")
	flag.Parse()

	if showlog {
		log.SetLevel(log.InfoLevel)
	}

	if minsup <= 0 || minsup >= 1 {
		log.Errorf("the value of minimal support should be set in 0~1, the input value is %v", minsup)
		os.Exit(1)
	}

	apriori := algorithm.New()

	log.WithFields(log.Fields{
		"fileName": fileName,
	}).Printf("start to read the db file")

	if err := apriori.ReadFile(fileName); err != nil {
		log.Fatal(err)
	}

	log.WithFields(log.Fields{
		"fileName": fileName,
	}).Println("finish to read the db file")

	fmt.Printf("The size of transactions DB:\t%d\n", apriori.GetLenOfDB())
	fmt.Printf("The number of items:\t%d\n", apriori.GetNumItems())

	log.WithFields(log.Fields{
		"minsup": minsup,
	}).Printf("set minsup")
	apriori.SetMinsup(minsup)

	log.Println("run Apriori process")
	start := time.Now()
	apriori.Run()
	elapsed := time.Since(start)
	log.Println("finish Apriori process")
	log.WithFields(log.Fields{
		"elapsed": elapsed,
	}).Printf("process took")

	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	log.WithFields(log.Fields{
		"memory": float64(m.TotalAlloc) / 1048576,
	}).Printf("total memory alloc MiB")

	fmt.Printf("The number of frequent itemsets:\t%d\n", apriori.GetNumFrequentItemsets())
	fmt.Printf("(over or equal to %v)\n", apriori.GetMinsupCount())
	frequentItemsets := apriori.GetFrequentItemsets()

	for _, fis := range frequentItemsets {
		for i := fis.Front(); i != nil; i = i.Next() {
			itemsets := algorithm.ItemsetsToString(i.Value.([]int))
			frequentItemsetsCount := apriori.GetFrequentItemsetsCount()
			fmt.Printf("%s:\t%v\n", itemsets, frequentItemsetsCount[itemsets])
		}
	}
}
