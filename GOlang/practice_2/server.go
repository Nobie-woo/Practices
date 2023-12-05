package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Node struct {
	data string
	next *Node
}

type Stack struct {
	head *Node
}

type Queue struct {
	head *Node
	tail *Node
}

type Set struct {
	array [256]string
}

type KeyValue struct {
	key       string
	value     string
	doesExist string
}

type HashTable struct {
	table [256]*KeyValue
}

//---------------------------------------------------------------------------------------------------------

func writeToFile(filePath string, content string) error { // ЗАПИСЬ В ФАЙЛ
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	err := os.WriteFile(filePath, []byte(content), 0644)
	return err
}

func readFile(filePath string) []byte { // ЧТЕНИЕ ФАЙЛА
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil
	}
	return content
}

func deleteStrFromFile(filePath string, content string) {
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()
	file := readFile(filePath)
	if file == nil {
		fmt.Println("Не удалось считать файл")
	} else {
		lines := strings.Split(string(file), "\n")
		newLines := ""
		isFound := false
		for i := 0; i < len(lines); i++ {
			if lines[i] == content && isFound == false {
				isFound = true
			} else if lines[i] != "" {
				newLines += lines[i] + "\n"
			}
		}
		writeToFile(filePath, newLines)
	}
}

//----------------------------------------------SET-------------------------------------------------------

func (set *Set) readSet(filePath string) {
	file := readFile(filePath)
	if file == nil {
		fmt.Println("Не удалось считать файл")
	} else {
		lines := strings.Split(string(file), "\n")
		for i := 0; i < len(lines); i++ {
			if lines[i] != "" {
				set.add(lines[i])
			}
		}
	}
}

func (set *Set) writeSet(filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println("Не удалось создать файл")
		return
	}
	defer file.Close()
	for i := 0; i < 256; i++ {
		if set.array[i] != "" {
			file.Write([]byte(set.array[i]))
			file.Write([]byte("\n"))
		}
	}
}

func (set *Set) add(value string) {
	index := toint(value, 256)
	if set.array[index] == "" {
		set.array[index] = value
	} else {
		for i := index % 256; i < 256; i = (i + 1) % 256 {
			if set.array[i] == value {
				fmt.Println("Элемент уже есть")
				return
			} else if set.array[i] == "" {
				set.array[i] = value
				return
			} else {
				fmt.Println("Множество заполнено")
				return
			}
		}
	}
}

func (set *Set) remove(value string) (string, error) {
	index := toint(value, 256)
	for i := index % 256; i < 256; i = (i + 1) % 256 {
		if set.array[i] == "" {
			return "", errors.New("элемент не найден")
		} else if set.array[i] == value {
			set.array[i] = ""
			return "", errors.New("всё путем :)")
		}
	}
	return "", errors.New("все ячейки заполнены, но элемент всё равно не найден")
}

func (set *Set) ismember(value string) (string, error) {
	index := toint(value, 256)
	for i := index % 256; i < 256; i = (i + 1) % 256 {
		if set.array[i] == "" {
			return "", errors.New("элемент не обнаружен")
		} else if set.array[i] == value {
			return "Элемент обнаружен", nil
		}
	}
	return "", errors.New("так и не обнаружили :(")
}

//-----------------------------------------------STACK----------------------------------------------------

func (stack *Stack) readStack(filePath string) {
	file := readFile(filePath)
	if file == nil {
		fmt.Println("Не удалось считать файл")
	} else {
		lines := strings.Split(string(file), "\n")
		for i := 0; i < len(lines); i++ {
			if lines[i] != "" {
				stack.push(lines[i])
			}
		}
	}
}

func (stack *Stack) writeStack(filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Не удалось создать файл")
		return
	}
	defer file.Close()
	file.WriteString(stack.head.data)
	file.WriteString("\n")
}

func (stack *Stack) push(value string) {
	node := &Node{data: value}
	if stack.head == nil {
		stack.head = node
	} else {
		node.next = stack.head
		stack.head = node
	}
}

func (stack *Stack) pop() (string, error) {
	if stack.head == nil {
		return "", errors.New("ошибочка вышла (указатель нулевой)")
	} else {
		value := stack.head.data
		stack.head = stack.head.next
		return value, nil
	}
}

//-----------------------------------------------QUEUE-----------------------------------------------------

func (queue *Queue) readQueue(filePath string) {
	file := readFile(filePath)
	if file == nil {
		fmt.Println("Не удалось считать файл")
	} else {
		lines := strings.Split(string(file), "\n")
		for i := 0; i < len(lines); i++ {
			if lines[i] != "" {
				queue.enqueue(lines[i])
			}
		}
	}
}

func (queue *Queue) writeQueue(filePath string) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Не удалось создать файл")
		return
	}
	defer file.Close()
	file.WriteString(queue.tail.data)
	file.WriteString("\n")
}

func (queue *Queue) enqueue(value string) {
	node := &Node{data: value}
	if queue.head == nil {
		queue.head = node
		queue.tail = node
	} else {
		queue.tail.next = node
		queue.tail = node
	}
}

func (queue *Queue) dequeue() (string, error) {
	if queue.head == nil {
		return "", errors.New("ошибочка вышла (указатель нулевой)")
	} else {
		value := queue.head.data
		queue.head = queue.head.next
		return value, nil

	}
}

//-----------------------------------------------HASH------------------------------------------------------

func toint(key string, size int) int {
	hash := 0
	for i := 0; i < len(key); i++ {
		hash = hash + int(key[i])
	}
	return hash % size
}

func (hashTable *HashTable) readHash(filePath string) {
	file := readFile(filePath)
	if file == nil {
		fmt.Println("Не удалось считать файл")
	} else {
		lines := strings.Split(string(file), "\n")
		for i := 0; i < len(lines); i++ {
			line := strings.Split(string(lines[i]), " ")
			if len(line) == 1 {
				return
			}
			hashTable.insert(line[0], line[1], line[2])
		}
	}
}

func (HashTable *HashTable) writeHash(filePath string) {
	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Не удалось создать файл")
		return
	}
	defer file.Close()
	for i := 0; i < 256; i++ {
		if HashTable.table[i] != nil {
			file.Write([]byte(HashTable.table[i].key))
			file.Write([]byte(" "))
			file.Write([]byte(HashTable.table[i].value))
			file.Write([]byte(" "))
			file.Write([]byte(HashTable.table[i].doesExist))
			file.Write([]byte("\n"))
		}
	}
}

func (hashTable *HashTable) insert(key string, value string, doesExist string) {
	newkeyValue := &KeyValue{key, value, doesExist}
	index := toint(key, 256)
	if hashTable.table[index] == nil {
		hashTable.table[index] = newkeyValue
	} else {
		for i := (index + 1) % 256; i < 256; i = (i + 1) % 256 {
			if hashTable.table[i] == nil {
				hashTable.table[i] = newkeyValue
				return
			} else if hashTable.table[i].key == key {
				hashTable.table[i].value = value
				hashTable.table[i].doesExist = "1"
				return
			} else if hashTable.table[i].doesExist == "0" {
				continue
			} else {
				fmt.Println("Места больше нет")
			}
		}
	}
}

func (hashTable *HashTable) delete(key string) (string, error) {
	index := toint(key, 256)
	for i := (index + 1) % 256; i < 256; i = (i + 1) % 256 {
		if hashTable.table[i] == nil || (hashTable.table[i].doesExist == "0" && hashTable.table[i].key == key) {
			return "", errors.New("ключ не найден")
		} else if hashTable.table[i].key == key && hashTable.table[i].doesExist == "1" {
			hashTable.table[i].doesExist = "0"
			return hashTable.table[i].value, errors.New("всё путем :)")
		}
	}
	return "", errors.New("все ячейки заполнены, но ключ все равно не найден")
}

func (hashTable *HashTable) get(key string) (string, error) {
	index := toint(key, 256)
	for i := (index + 1) % 256; i < 256; i = (i + 1) % 256 {
		if hashTable.table[i] == nil {
			return "", errors.New("элемент не обнаружен")
		} else if hashTable.table[i].key == key && hashTable.table[i].doesExist == "1" {
			return hashTable.table[i].value, nil
		}
	}
	return "", errors.New("так и не обнаружили :(")
}

func editDatabase(query string) {
	set := &Set{}
	stack := &Stack{}
	queue := &Queue{}
	hash := &HashTable{}
	str := strings.Fields(query)
	file := str[0]

	switch str[1] {

	case "SADD":
		set.readSet(file)
		set.add(str[3])
		set.writeSet(file)

	case "SREM":
		set.readSet(file)
		value, err := set.remove(str[3])
		if err == nil {
			fmt.Println(value)
			deleteStrFromFile(file, value)
		} else {
			fmt.Println(err)
		}

	case "SISMEMBER":
		set.readSet(file)
		value, err := set.ismember(str[3])
		if err == nil {
			fmt.Println(value)
		} else {
			fmt.Println(err)
		}

	case "SPUSH":
		stack.readStack(file)
		stack.push(str[3])
		stack.writeStack(file)

	case "SPOP":
		stack.readStack(file)
		value, err := stack.pop()
		if err == nil {
			fmt.Println(value)
			deleteStrFromFile(file, value)
		} else {
			fmt.Println(err)
		}

	case "QPUSH":
		queue.readQueue(file)
		queue.enqueue(str[3])
		queue.writeQueue(file)

	case "QPOP":
		queue.readQueue(file)
		value, err := queue.dequeue()
		if err == nil {
			fmt.Println(value)
			deleteStrFromFile(file, value)
		} else {
			fmt.Println(err)
		}

	case "HSET":
		hash.readHash(file)
		hash.insert(str[3], str[4], "1")
		hash.writeHash(file)

	case "HDEL":
		hash.readHash(file)
		hash.delete(str[3])
		hash.writeHash(file)

	case "HGET":
		hash.readHash(file)
		value, err := hash.get(str[3])
		if err == nil {
			fmt.Println(value)
		} else {
			fmt.Println(err)
		}
	default:
		fmt.Println("Ошибка в написании команды!")
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	var query string
	for {
		buf, err := bufio.NewReader(conn).ReadBytes('\n')
		if err != nil {
			break
		}
		fmt.Print("Запрос: ", string(buf))
		query = string(buf)
		break
	}
	editDatabase(query)
}

func main() {
	listener, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			conn.Close()
			continue
		}
		go handleConnection(conn)
	}
}
