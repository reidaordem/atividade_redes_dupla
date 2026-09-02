package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9000")
	if err != nil {
		fmt.Println("Erro ao conectar ao servidor:", err)
		os.Exit(1)
	}
	defer conn.Close()

	addr := conn.RemoteAddr().String()
	fmt.Printf("Conectado ao servidor %s\n", addr)
	
	reader := bufio.NewReader(os.Stdin)
	serverReader := bufio.NewReader(conn)
	for {
		fmt.Print("calc> ")
		text, _ := reader.ReadString('\n')

		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Erro ao enviar mensagem:", err)
			break
		}

		response, err := serverReader.ReadString('\n')
		if err != nil {
			fmt.Println("Servidor fechou a conexão.")
			break
		}
		fmt.Print(response)
		if response == "(conexão encerrada)\n" {
			break
		}
	}

}
