package main

import (
	"bufio"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
)

type Banner struct {
	Image string
	Name  string
}

type FullBlockStruct struct {
	Image   string
	Name    string
	Company string
	Price   int
	MRP     int
}

type HalfBlockStruct struct {
	Image1 string
	Name1  string
	Image2 string
	Name2  string
}

type SuggestionResult struct {
	Result1 string
	Result2 string
	Result3 string
}

type StockItemStruct struct {
	ItemName string
	Barcode  int
	Company  string
	Stock    int
	Images   []byte
	Price    int
	MRP      int
}

type StockEntryStruct struct {
	Barcode int
	Rate    int
	Price   int
	ExpDate string
	BatchNo string
	Mfd     string
}

// Node represents each node in the trie
type Node struct {
	children [AlphabetSize]*Node
	weight   [AlphabetSize]int
	isEnd    bool
}

// Trie represents a trie and has a pointer to the root node
type Trie struct {
	root *Node
}

//AlphabetSize is the number of possible characters in the trie
const AlphabetSize = 26

var database *sql.DB

var testTrie *Trie

//InitTrie will create new Trie
func InitTrie() *Trie {
	result := &Trie{root: &Node{}}
	return result
}

// Insert will take in a word and add it to the trie
func (t *Trie) Insert(w string) {
	wordLength := len(w)
	currentNode := t.root
	for i := 0; i < wordLength; i++ {
		charIndex := w[i] - 'a'
		currentNode.weight[charIndex]++
		if currentNode.children[charIndex] == nil {
			currentNode.children[charIndex] = &Node{}
		}
		currentNode = currentNode.children[charIndex]
	}
	currentNode.isEnd = true
}

// Search will take in a word and RETURN true if that word is included in the trie
func (t *Trie) Search(w string) bool {
	wordLength := len(w)
	currentNode := t.root
	for i := 0; i < wordLength; i++ {
		charIndex := w[i] - 'a'
		currentNode.weight[charIndex]++
		if currentNode.children[charIndex] == nil {
			return false
		}
		currentNode = currentNode.children[charIndex]
	}
	return currentNode.isEnd
}

func (t *Trie) GetSuggestion(w string) SuggestionResult {
	wordLength := len(w)
	currentNode := t.root
	for i := 0; i < wordLength; i++ {
		charIndex := w[i] - 'a'
		//	fmt.Println(currentNode.weight[charIndex])
		//	currentNode.weight[charIndex]++
		if currentNode.children[charIndex] == nil {
			return SuggestionResult{"no suggestion", "no suggestion", "no suggestion"}
		}
		currentNode = currentNode.children[charIndex]
	}

	indexMax1 := 0
	max1 := 0

	indexMax2 := 0
	max2 := 0

	indexMax3 := 0
	max3 := 0

	for i := 0; i < 26; i++ {
		if currentNode.weight[i] > max1 {
			max3 = max2
			indexMax3 = indexMax2
			max2 = max1
			indexMax2 = indexMax1
			max1 = currentNode.weight[i]
			indexMax1 = i
		} else if currentNode.weight[i] > max2 {
			max3 = max2
			indexMax3 = indexMax2
			max2 = currentNode.weight[i]
			indexMax2 = i
		} else if currentNode.weight[i] > max3 {
			max3 = currentNode.weight[i]
			indexMax3 = i
		}

	}
	searchResult := SuggestionResult{}
	if max1 == 0 {
		searchResult.Result1 = ""
	} else {
		searchResult.Result1 = w + string(rune('a'+indexMax1)) + findMax(currentNode.children[indexMax1])
	}
	if max2 == 0 {
		searchResult.Result2 = ""
	} else {
		searchResult.Result2 = w + string(rune('a'+indexMax2)) + findMax(currentNode.children[indexMax2])
	}
	if max3 == 0 {
		searchResult.Result3 = ""
	} else {
		searchResult.Result3 = w + string(rune('a'+indexMax3)) + findMax(currentNode.children[indexMax3])
	}

	return searchResult
}

func findMax(node *Node) string {
	index_max := 0
	max := 0
	foundmax := false
	for i := 0; i < 26; i++ {
		if node.weight[i] > max {
			max = node.weight[i]
			index_max = i
			foundmax = true
		}
	}
	if !foundmax {
		return ""
	}
	return string(rune('a'+index_max)) + findMax(node.children[index_max])
}

func putItem(item StockItemStruct) {
	_, _ = database.Query("INSERT INTO StockItem VALUES ('" + item.ItemName + "','" + string(rune(item.Barcode)) + "','" + item.Company + "','" + string(rune(item.Stock)) + "','" + string(item.Images) + "'" + string(rune(item.Price)) + "'" + string(rune(item.MRP)) + "'" + ")")
}

func main() {
	// Init Router
	database, _ = sql.Open("mysql", "shreyanshumalviya:homnhomnhomn@(127.0.0.1:3306)/shreyanshumalviya")
	testTrie = InitTrie()
	testTrie.Insert("abhishek")
	testTrie.Insert("aman")
	testTrie.Insert("anil")
	testTrie.Insert("ankit")
	testTrie.Insert("anshul")
	testTrie.Insert("kanti")
	testTrie.Insert("oats")
	testTrie.Insert("dalia")
	testTrie.Insert("ghrit")
	testTrie.Insert("anu")
	testTrie.Insert("amlapickle")
	testTrie.Insert("ashokaristh")
	testTrie.Insert("giloy")
	testTrie.Insert("gulabjamun")
	testTrie.Insert("abhyaristh")
	testTrie.Insert("amlamurabba")
	testTrie.Insert("arjunaristh")
	testTrie.Insert("chyawanprash")

	fmt.Println("Creating Server")
	r := mux.NewRouter()

	r.HandleFunc("/", suggestionRequest).Methods(http.MethodGet)
	r.HandleFunc("/search", search).Methods(http.MethodGet)
	r.HandleFunc("/banners", getBanners).Methods(http.MethodGet)
	r.HandleFunc("/homeitems", getHomeItems).Methods(http.MethodGet)
	r.HandleFunc("/addnewstock", addNewStock).Methods(http.MethodGet)
	log.Fatal(http.ListenAndServe(":8080", r))
}

func testPutting() {

	/*buf := new(bytes.Buffer)
	err := jpeg.Encode(buf, new_image, nil)
	if err != nil {
		send_s3 := buf.Bytes()

		putItem(StockItemStruct{"coca cola", 192021333, "coke", 214, send_s3, 34, 23})
	}*/
}

func addNewStock(writer http.ResponseWriter, request *http.Request) {

}

func getHomeItems(writer http.ResponseWriter, request *http.Request) {
	f, _ := os.Open("homeItems")
	list, _ := f.Readdir(-1)
	_ = f.Close()
	var homeItemList []HalfBlockStruct
	size := len(list)

	for index := 0; index < size; index++ {
		buf := make([]byte, list[index].Size())
		filePath := "homeItems\\" + list[index].Name()
		file, _ := os.Open(filePath)
		fReader := bufio.NewReader(file)
		_, _ = fReader.Read(buf)
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)
		var item HalfBlockStruct
		item.Image1 = imgBase64Str
		item.Name1 = list[index].Name()

		index++
		if index != size {
			buf = make([]byte, list[index].Size())
			filePath = "homeItems\\" + list[index].Name()
			file, _ = os.Open(filePath)
			fReader = bufio.NewReader(file)
			_, _ = fReader.Read(buf)
			imgBase64Str = base64.StdEncoding.EncodeToString(buf)
			item.Image2 = imgBase64Str
			item.Name2 = list[index].Name()
		} else {
			item.Image2 = ""
			item.Name2 = ""
		}
		homeItemList = append(homeItemList, item)
	}
	finalMessage, _ := json.Marshal(homeItemList)
	_, _ = fmt.Fprint(writer, string(finalMessage))
}

func getBanners(writer http.ResponseWriter, request *http.Request) {
	f, _ := os.Open("Banner")
	list, _ := f.Readdir(-1)
	_ = f.Close()
	var banners []Banner
	for index, fInfo := range list {
		buf := make([]byte, fInfo.Size())
		filePath := "banner\\" + list[index].Name()
		file, _ := os.Open(filePath)
		fReader := bufio.NewReader(file)
		_, _ = fReader.Read(buf)
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)
		banners = append(banners, Banner{Image: imgBase64Str, Name: string(rune(index))})
	}
	finalMessage, _ := json.Marshal(banners)
	_, _ = fmt.Fprint(writer, string(finalMessage))

}

func search(writer http.ResponseWriter, request *http.Request) {
	input := request.Header.Get("Get")
	testTrie.Search(input)

	// TODO reply after looking from data base

	f, _ := os.Open("homeItems")
	list, _ := f.Readdir(-1)
	_ = f.Close()
	var itemList []FullBlockStruct
	size := len(list)

	for index := 0; index < size; index++ {
		buf := make([]byte, list[index].Size())
		filePath := "homeItems\\" + list[index].Name()
		file, _ := os.Open(filePath)
		fReader := bufio.NewReader(file)
		_, _ = fReader.Read(buf)
		imgBase64Str := base64.StdEncoding.EncodeToString(buf)
		var item FullBlockStruct
		item.Image = imgBase64Str
		item.Name = list[index].Name()
		item.Company = "company"
		item.Price = 200
		item.MRP = 300
		itemList = append(itemList, item)
	}
	finalMessage, _ := json.Marshal(itemList)
	_, _ = fmt.Fprint(writer, string(finalMessage))

}

func suggestionRequest(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(request.Header)
	input := request.Header.Get("Get")
	fmt.Println("hello")
	fmt.Println(input)
	suggestions := testTrie.GetSuggestion(input)
	fmt.Println("suggestions := ", suggestions)
	jsonSuggestion, err := json.Marshal(suggestions)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("\n json suggestions := ", string(jsonSuggestion))

		_, _ = fmt.Fprint(writer, string(jsonSuggestion))
	}
}
