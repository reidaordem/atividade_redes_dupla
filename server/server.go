package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println("Erro ao iniciar o servidor:", err)
		os.Exit(1)
	}
	defer listener.Close()
	fmt.Println("Servidor ouvindo na porta 9000...")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erro ao aceitar conexão:", err)
			continue
		}
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()

    	fmt.Printf("[CONECTOU] %s\n", addr)

	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')

		if err != nil {
			if err == io.EOF {
				fmt.Printf("Cliente %s desconectou.\n", conn.RemoteAddr().String())
			} else {
				fmt.Println("Erro na leitura:", err)
			}
			break
		}

		fmt.Printf("[%s]: %s", conn.RemoteAddr().String(), message)

		operacao := operacoes(message, conn)

		if operacao {
			fmt.Printf("[ENCERROU] %s\n", addr)
			break
		}
	}
}

func operacoes(mensagem string, conn net.Conn) bool {
	slice_operacao := strings.Fields(mensagem)

	if len(slice_operacao) != 3 {
		if len(slice_operacao) == 1 && slice_operacao[0] == "SAIR" {
			fmt.Fprintf(conn, "(conexão encerrada)\n")
			return true
		}
		fmt.Fprintf(conn, "ERRO: formato invalido(use: OPERACAO NUM1 NUM2)\n")
		return false
	}

	NUM1, err1 := strconv.Atoi(slice_operacao[1])
	NUM2, err2 := strconv.Atoi(slice_operacao[2])
	var resultado int

	if err1 != nil || err2 != nil {
		fmt.Fprintf(conn, "ERRO: formato invalido(use: OPERACAO NUM1 NUM2)\n")
		return false
	}

	if slice_operacao[0] == "DIV" && NUM2 == 0 {
		fmt.Fprintf(conn, "ERRO: divisao por zero\n")
		return false
	}

	switch slice_operacao[0] {

	case "SOMA":
		resultado = NUM1 + NUM2
		fmt.Fprintf(conn, "RESULTADO %d\n", resultado)
	case "SUB":
		resultado = NUM1 - NUM2
		fmt.Fprintf(conn, "RESULTADO %d\n", resultado)
	case "MUL":
		resultado = NUM1 * NUM2
		fmt.Fprintf(conn, "RESULTADO %d\n", resultado)
	case "DIV":
		resultado = NUM1 / NUM2
		fmt.Fprintf(conn, "RESULTADO %d\n", resultado)
	case "SAIR":
		fmt.Fprintf(conn, "(conexão encerrada)\n")

		return true

	default:
		fmt.Fprintf(conn, "ERRO: comando desconhecido\n")

	}
	return false
}