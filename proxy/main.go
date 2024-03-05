package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/webdeskltd/dadata.v2"
	"io/ioutil"
	"log"
	"main/proxy/binary"
	"main/proxy/graph"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

func main() {

	// Создаем реверсивный прокси для сервера Hugo
	hugoURL, err := url.Parse("http://hugo:1313")
	if err != nil {
		panic(err)
	}
	hugoProxy := httputil.NewSingleHostReverseProxy(hugoURL)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/address/search") {
			search(w, r)
		}
		if strings.HasPrefix(r.URL.Path, "/api/address/geo") {
			geocode(w, r)
		} else if strings.HasPrefix(r.URL.Path, "/api/") {
			w.Write([]byte("Hello from API"))
		} else {
			hugoProxy.ServeHTTP(w, r)
		}
	})
	//go WorkerCounter()
	//go WorkerBinary()
	//go WorkerGraph()
	// Запускаем сервер на порту 8080
	http.ListenAndServe(":8080", nil)
}

type ReverseProxy struct {
	host string
	port string
}

func NewReverseProxy(host, port string) *ReverseProxy {
	return &ReverseProxy{
		host: host,
		port: port,
	}
}

func (rp *ReverseProxy) ReverseProxy(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}

const content = ``

func WorkerTest() {
	t := time.NewTicker(1 * time.Second)
	var b byte = 0
	for {
		select {
		case <-t.C:
			fmt.Printf("Текущее время: %s\n", time.Now().String()[11:19])
			fmt.Printf("Текущая дата %s", time.Now().String()[0:7])
			err := writeToFile("./app/static/_index.md", fmt.Sprintf("%s%d", content, b))
			if err != nil {
				log.Fatal(err)
			}
			b++
		}
	}
}
func writeToFile(path string, data string) error {
	// Открываем файл для записи, флаг O_TRUNC обрежет файл до нуля и удалит предыдущее содержимое
	file, err := os.OpenFile(path, os.O_RDWR|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	time.Sleep(5 * time.Second)
	// Записываем данные в файл
	_, err = file.WriteString(data)
	if err != nil {
		return err
	}

	return nil
}

func writeCounter(counter byte) error {
	b := counter
	filePath := "./app/static/tasks/_index.md"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return err
	}
	text := string(content)
	re := regexp.MustCompile(`Счетчик:\s*\d+`)
	match := re.FindStringIndex(text)
	if match == nil {
		fmt.Println("Фраза 'Счетчик' не найдена в файле.")
		return err
	}
	с := fmt.Sprintf("Счетчик: %v", b)
	newText := text[:match[0]] + с + text[match[1]:]
	err = ioutil.WriteFile(filePath, []byte(newText), os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		return err
	}
	time.Sleep(5 * time.Second)
	filePath = "./app/static/tasks/_index.md"
	content, err = ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return err
	}
	text = string(content)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return nil
	}
	re = regexp.MustCompile(`Текущее\s*время\s*:\s*\d{4}-\d{2}-\d{2}\s*\d{2}:\d{2}:\d{2}`)
	match = re.FindStringIndex(text)
	if match == nil {
		fmt.Println("Фраза 'Текущее время' не найдена в файле.")
		return err
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	newTimeLine := fmt.Sprintf("Текущее время: %s", currentTime)
	newText = text[:match[0]] + newTimeLine + text[match[1]:]
	err = ioutil.WriteFile(filePath, []byte(newText), os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		return err
	}

	return nil
}

func WriteBinary(tree *binary.AVLTree, key int) error {
	filePath := "/app/static/tasks//binary.md"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return err
	}
	text := string(content)
	re := regexp.MustCompile(`{{< /columns >}}`)
	match := re.FindStringIndex(text)
	if match == nil {
		fmt.Println("Фраза 'Счетчик' не найдена в файле.")
		return err
	}
	tree.Insert(key)
	с := fmt.Sprintf(tree.ToMermaid())
	newText := text[:match[0]] + с
	err = ioutil.WriteFile(filePath, []byte(newText), os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		return err
	}
	fmt.Println(string(content))
	fmt.Println("Текущее время успешно обновлено в файле.")
	time.Sleep(1 * time.Second)
	fmt.Print("\033[H\033[2J")
	return nil

}

func WriteGraph() error {
	filePath := "/app/static/tasks//graph.md"
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("Ошибка чтения файла:", err)
		return err
	}
	text := string(content)
	re := regexp.MustCompile(`{{< mermaid >}}`)
	match := re.FindStringIndex(text)
	if match == nil {
		fmt.Println("Фраза '{{< mermaid >}}' не найдена в файле.")
		return err
	}
	с := fmt.Sprintf(graph.GenerateMermaid())
	newText := text[:match[0]] + с
	err = ioutil.WriteFile(filePath, []byte(newText), os.ModePerm)
	if err != nil {
		fmt.Println("Ошибка записи файла:", err)
		return err
	}
	fmt.Println(string(content))
	fmt.Println("Текущее время успешно обновлено в файле.")
	time.Sleep(1 * time.Second)
	fmt.Print("\033[H\033[2J")
	return nil

}

func WorkerCounter() {
	counter := byte(0)
	for range time.Tick(5 * time.Second) {
		err := writeCounter(counter)
		if err != nil {
			log.Fatal(err)
		}
		counter++
	}
}

func WorkerBinary() {
	tree := binary.GenerateTree(5)
	i := 6
	for range time.Tick(5 * time.Second) {
		if i > 100 {
			tree = binary.GenerateTree(5)
		}
		err := WriteBinary(tree, i)
		if err != nil {
			log.Fatal(err)
		}
		i++
	}
}

func WorkerGraph() {
	for range time.Tick(5 * time.Second) {
		err := WriteGraph()
		if err != nil {
			log.Fatal(err)
		}
	}
}

type SearchRequest struct {
	Query string `json:"query"`
}
type SearchResponse struct {
	Addresses []*Address `json:"addresses"`
}
type GeocodeRequest struct {
	Lat float32 `json:"lat"`
	Lng float32 `json:"lon"`
}
type GeocodeResponse struct {
	Addresses []*Address `json:"addresses"`
}

func search(w http.ResponseWriter, r *http.Request) {
	var searchRequest SearchRequest
	err := json.NewDecoder(r.Body).Decode(&searchRequest)
	if err != nil {
		log.Println("err502.1")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	log.Println(searchRequest.Query)

	addresses, err := api.CleanAddresses(searchRequest.Query)
	if err != nil {
		log.Println("err502.2")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	log.Println(addresses)
	var searchResponse SearchResponse
	searchResponse.Addresses = []*Address{{Lat: addresses[0].GeoLat, Lng: addresses[0].GeoLon}}
	err = json.NewEncoder(w).Encode(&searchResponse)
	if err != nil {
		log.Println("err502.3")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

type Address struct {
	Lat string `json:"lat"`
	Lng string `json:"lon"`
}

func geocode(w http.ResponseWriter, r *http.Request) {
	var geocodeRequest GeocodeRequest
	err := json.NewDecoder(r.Body).Decode(&geocodeRequest)
	if err != nil {
		log.Println("err1")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	api := dadata.NewDaData("602f4fabeedea0f000f4cee8ab9a5773d800f005", "f57d7df9064c22a9c4a7c61b90109cd44fd7f284")

	req := dadata.GeolocateRequest{
		Lat:          geocodeRequest.Lat,
		Lon:          geocodeRequest.Lng,
		Count:        5,
		RadiusMeters: 100,
	}
	addresses, err := api.GeolocateAddress(req)
	if err != nil {
		log.Println("err4")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var geocodeResponse GeocodeResponse
	geocodeResponse.Addresses = []*Address{{Lat: addresses[0].Data.City, Lng: addresses[0].Data.Street + " " + addresses[0].Data.House}}
	err = json.NewEncoder(w).Encode(&geocodeResponse)
	if err != nil {
		log.Println("err5")
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
