package buildmermaid

import (
	"fmt"
	"io/ioutil"
	"log"
	"main/proxy/binary"
	"main/proxy/graph"
	"os"
	"regexp"
	"time"
)

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

// RegisterAndLogin godoc
// @Summary Register and login users
// @Description Register and login users using JWT tokens
// @Tags auth
// @Accept  json
// @Produce  json
// @Param username query string true "Username"
// @Param password query string true "Password"
// @Success 200 {string} string "Success"
// @Failure 401 {string} string "Unauthorized"
// @Router /auth/register [get]
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
