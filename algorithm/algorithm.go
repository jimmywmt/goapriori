package algorithm

import (
	"bufio"
	"container/list"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

const defaultNumFrequentItemset = 128
const defaultNumFrequentItemsetMaxLen = 4

type Apriori struct {
	relatedTransactions  map[int]*list.List
	frequentItemsetCount map[string]int
	frequentItemsets     []*list.List
	numFrequentItemsets  int
	minsup               float64
	minsupCount          int
	valuOfMinsup         float64
	fileName             string
	maxItemID            int
	lenOfDB              int
}

func New() *Apriori {
	return &Apriori{}
}

func ItemsetsToString(itemsets []int) string {
	return strings.Trim(strings.Replace(fmt.Sprint(itemsets), " ", ",", -1), "[]")
}

func (apriori *Apriori) ReadFile(file string) error {
	apriori.fileName = file
	apriori.relatedTransactions = make(map[int]*list.List)
	apriori.maxItemID = 0
	apriori.lenOfDB = 0
	content, err := os.Open(apriori.fileName)

	if err != nil {
		return err
	}

	defer content.Close()

	scanner := bufio.NewScanner(content)

	for scanner.Scan() {
		stringArray := strings.Split(scanner.Text(), " ")

		for _, value := range stringArray {
			intValue, err := strconv.Atoi(value)

			if intValue > apriori.maxItemID {
				apriori.maxItemID = intValue
			}

			if err != nil {
				return err
			}

			transactionList := apriori.relatedTransactions[intValue]

			if transactionList == nil {
				transactionList = &list.List{}
			}

			transactionList.PushBack(apriori.lenOfDB)
			apriori.relatedTransactions[intValue] = transactionList
		}
		apriori.lenOfDB++
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func (apriori *Apriori) SetMinsup(minsup float64) {
	if apriori.relatedTransactions != nil {
		apriori.minsup = minsup
		apriori.minsupCount = int(math.Ceil(float64(apriori.lenOfDB) * minsup))
	}
}

func (apriori *Apriori) GetMinsup() float64 {
	return apriori.minsup
}

func (apriori *Apriori) GetFrequentItemsets() []*list.List {
	return apriori.frequentItemsets
}

func (apriori *Apriori) GetFrequentItemsetsCount() map[string]int {
	return apriori.frequentItemsetCount
}

func (apriori *Apriori) GetLenOfDB() int {
	return apriori.lenOfDB
}

func (apriori *Apriori) GetNumItems() int {
	return len(apriori.relatedTransactions)
}

func (apriori *Apriori) GetNumFrequentItemsets() int {
	return apriori.numFrequentItemsets
}

func (apriori *Apriori) GetMinsupCount() int {
	return apriori.minsupCount
}

func (apriori *Apriori) Run() {
	if apriori.lenOfDB != 0 && apriori.minsup != 0.0 {
		apriori.frequentItemsetCount = make(map[string]int)
		apriori.frequentItemsets = make([]*list.List, 0, defaultNumFrequentItemsetMaxLen)
		apriori.numFrequentItemsets = 0

		// find large-1 itemsets
		largeItemsets := &list.List{}

		for i := 1; i <= apriori.maxItemID; i++ {
			value := apriori.relatedTransactions[i]
			if value.Len() >= apriori.minsupCount {
				apriori.numFrequentItemsets++
				tempSlice := []int{i}
				largeItemsets.PushBack(tempSlice)
				apriori.frequentItemsetCount[strconv.Itoa(i)] = value.Len()
			}
		}

		if largeItemsets.Len() != 0 {
			apriori.frequentItemsets = append(apriori.frequentItemsets, largeItemsets)
			apriori.nextLevelLargeItemsets(largeItemsets)
		}

	}
}

func (apriori *Apriori) nextLevelLargeItemsets(largeItemsets *list.List) {
	newLargeItemsets := &list.List{}
	i := largeItemsets.Front()
	candidateItemsets := make(map[string]int)

	for i.Next() != nil {
		oneItemset := i.Value.([]int)
		j := i.Next()

		for j != nil {
			twoItemset := j.Value.([]int)
			threeItemset := make([]int, len(oneItemset)+1)

			if mergeItemsets(oneItemset, twoItemset, threeItemset) {
				itemsetsString := ItemsetsToString(threeItemset)
				count := candidateItemsets[itemsetsString]
				count++

				if count == len(oneItemset) {
					freqCount := apriori.itemsetCount(threeItemset)

					if freqCount >= apriori.minsupCount {
						apriori.numFrequentItemsets++
						newLargeItemsets.PushBack(threeItemset)
						apriori.frequentItemsetCount[ItemsetsToString(threeItemset)] = freqCount
					}
				} else {
					candidateItemsets[itemsetsString] = count
				}
			}

			j = j.Next()
		}

		i = i.Next()
	}

	if newLargeItemsets.Len() != 0 {
		apriori.frequentItemsets = append(apriori.frequentItemsets, newLargeItemsets)
		if newLargeItemsets.Len() >= len(newLargeItemsets.Front().Value.([]int))+1 {
			apriori.nextLevelLargeItemsets(newLargeItemsets)
		}
	}
}

func mergeItemsets(one []int, two []int, newItemset []int) bool {
	var i, j, k, d int
	length := len(one)

	for i < length && j < length {
		if one[i] == two[j] {
			newItemset[k] = one[i]
			i++
			j++
		} else if one[i] < two[j] {
			d++

			if d < 2 {
				newItemset[k] = one[i]
				i++
			} else {
				return false
			}
		} else {
			return false
		}
		k++
	}
	newItemset[k] = two[j]
	return true
}

func (apriori *Apriori) itemsetCount(itemset []int) int {
	firstTransactions := apriori.relatedTransactions[itemset[0]]

	if len(itemset) == 1 {
		return firstTransactions.Len()
	}

	var itemsetTransactions *list.List

	for i := 1; i < len(itemset); i++ {
		itemsetTransactions = &list.List{}
		secondTransactions := apriori.relatedTransactions[itemset[i]]

		it1 := firstTransactions.Front()
		it2 := secondTransactions.Front()

		if it1 != nil && it2 != nil {
			oneID := it1.Value.(int)
			twoID := it2.Value.(int)

			for {
				if oneID == twoID {
					itemsetTransactions.PushBack(oneID)
					it1 = it1.Next()
					it2 = it2.Next()

					if it1 != nil && it2 != nil {
						oneID = it1.Value.(int)
						twoID = it2.Value.(int)
					} else {
						break
					}
				} else if oneID < twoID {
					it1 = it1.Next()
					if it1 != nil {
						oneID = it1.Value.(int)
					} else {
						break
					}
				} else {
					it2 = it2.Next()
					if it2 != nil {
						twoID = it2.Value.(int)
					} else {
						break
					}
				}
			}

		}

		firstTransactions = itemsetTransactions
	}

	return itemsetTransactions.Len()
}
