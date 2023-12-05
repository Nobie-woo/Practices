package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {

	fmt.Print("Введите путь до файла и запрос (напр. 'set.txt SADD myset 28'): ")
	reader := bufio.NewReader(os.Stdin)
	query, _ := reader.ReadString('\n')
	fmt.Println("Вы ввели:", query)

	conn, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		fmt.Println("Ошибка подключения к серверу: ", err)
		return
	}
	defer conn.Close()

	_, err = fmt.Fprintf(conn, "%s\n", query)
	if err != nil {
		fmt.Println("ошибка записи: ", err)
		return
	}
}
