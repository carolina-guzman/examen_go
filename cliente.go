package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"net"
	"os"
)

var chat string
var bandera bool

type Archivo struct {
	Nombre string
	Info []byte
	Tamanio int64
}

func getString() string {
	var stdin *bufio.Reader
	var line []rune
	var c rune
	var err error

	stdin = bufio.NewReader(os.Stdin)

	fmt.Printf("Mensaje: ")

	for {
		c, _, err = stdin.ReadRune()
		if err == io.EOF || c == '\n' {
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error\n")
			os.Exit(1)
		}
		line = append(line, c)
	}

	return string(line[:len(line)])
}


var conexion net.Conn
var usuario string

func cliente() {
	c, err := net.Dial("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	conexion = c
	validar := usuario
	err2 := gob.NewEncoder(conexion).Encode(validar)
	if err2 != nil {
		fmt.Println(err)
		return
	}
	for bandera {
		var x string
		texto := gob.NewDecoder(c).Decode(&x)
		if texto != nil {
			fmt.Println(texto)
		}
		if x == "File" {
			var recibido Archivo
			docError := gob.NewDecoder(c).Decode(&recibido)
			f, docError := os.Create(recibido.Nombre)
			if docError != nil {
				fmt.Println("Error: ", docError)
			} else {
				f.Write(recibido.Info)
				f.Sync()
				f.Close()
				chat += "Documento recibido -> " + recibido.Nombre + "\n"
				fmt.Println("Documento3 -> " + recibido.Nombre + "\n")
			}

		} else {
			chat += x + "\n"
			fmt.Println( x + "\n")
		}
	}
}

func main() {
	var op int64
	bandera = true
	fmt.Print("Username: ")
	fmt.Scanln(&usuario)
	go cliente()
	fmt.Println("-------Bienvenido ", usuario, "-------")
	for {
		fmt.Println("1: Enviar Mensaje")
		fmt.Println("2: Enviar Archivo")
		fmt.Println("3: Ver chat")
		fmt.Println("4: Salir")
		fmt.Scanln(&op)
		switch op {
		case 1:

			msj := usuario + ": "
			msj += getString()
			err := gob.NewEncoder(conexion).Encode(msj)
			if err != nil {
				fmt.Println(err)
				break
			}
			break
		case 2:
			var n string
			var ruta string
			fmt.Print("Ingrese la ruta del archivo a enviar: ")
			fmt.Scanln(&ruta)
			fmt.Print("Ingrese el nombre del archivo: ")
			fmt.Scanln(&n)
			g, err := os.Open(ruta)
			if err != nil {
				fmt.Println(err)
			} else {
				err2 := gob.NewEncoder(conexion).Encode("f")
				if err2 != nil {
					fmt.Println(err2)
				}
				data, _ := g.Stat()
				var size int64
				size = data.Size()
				array := make([]byte, data.Size())
				g.Read(array)
				object := Archivo{n, array, size}
				errorObject := gob.NewEncoder(conexion).Encode(object)
				if errorObject != nil {
					fmt.Println(errorObject)
				}
			}
			break
		case 3:
			fmt.Println(chat)
			break
		case 4:
			bandera = false
			msj := "-1"
			err := gob.NewEncoder(conexion).Encode(msj)
			if err != nil {
				fmt.Println(err)
				break
			}
			conexion.Close()
			return
			break
		}
	}
	var x string
	fmt.Scanln(&x)
}
