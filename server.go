package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)


type Archivo struct {
	Nombre string
	Info []byte
	Tamanio int64
}

var conversacion string
var servidorActivo bool
var clientes = make(map[string]net.Conn)



func server() {
	s, err := net.Listen("tcp", ":9000")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Ahora estás conectado! ")
	for {
		c, err := s.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go handleCliente(c)

		if !servidorActivo {
			s.Close()
			break
		}
	}
}

func handleCliente(c net.Conn) {
	var usuario string
	errC := gob.NewDecoder(c).Decode(&usuario)
	if errC != nil {
		fmt.Println(errC)
	}
	conversacion += usuario + " se ha conectado \n"

	clientes[usuario] = c
	enviarUsuarios(c, usuario+" está conectado!")
	for {
		var x string
		validacion := gob.NewDecoder(c).Decode(&x)
		if validacion != nil {
			fmt.Println(validacion)
		}
		if x == "-1" {
			nombre := obtenerUsuario(c)
			delete(clientes, nombre)
			conversacion += nombre + "se ha desconectado\n"
			enviarUsuarios(c, nombre+" se ha desconectado\n")
			break
		} else if x == "f" {
			var recibido Archivo
			error2 := gob.NewDecoder(c).Decode(&recibido)
			if error2 != nil {
				fmt.Println(error2)
			}
			f, ErrorF := os.Create(recibido.Nombre)
			if ErrorF != nil {
				fmt.Println(ErrorF)
			} else {
				f.Write(recibido.Info)
				f.Sync()
				f.Close()
				conversacion += "Recibido: " + recibido.Nombre + "\n"
				enviarArchivoUsuarios(c, recibido)
			}
		} else {
			conversacion += x + "\n"
			fmt.Println( x + "\n")
			enviarUsuarios(c, x)
		}
	}
}

func obtenerUsuario(c net.Conn) string {
	usuario := ""
	for key, value := range clientes {
		if value == c {
			usuario = key
			break
		}
	}
	return usuario
}



func enviarUsuarios(c net.Conn, msj string) {
	for _, element := range clientes {
		if element != c {
			err := gob.NewEncoder(element).Encode(&msj)
			if err != nil {
				fmt.Println(err)
			}
		}
	}
}


func enviarArchivoUsuarios(c net.Conn, fileReceived Archivo) {
	for _, element := range clientes {
		if element != c {
			err := gob.NewEncoder(element).Encode("File")
			if err != nil {
				fmt.Println(err)
			}
			errorObject := gob.NewEncoder(element).Encode(fileReceived)
			if errorObject != nil {
				fmt.Println(errorObject)
			}
		}
	}
}

func escribirArchivo() {
	n:= "backup.txt"
	f, err := os.Create(n)
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()
	bn, err2 := f.WriteString(conversacion)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println("Escrito: ", bn)
	f.Sync()
}

func main() {
	servidorActivo = true
	go server()
	var op int
	for {
		fmt.Println("1. Mostrar Mensajes")
		fmt.Println("2. Hacer backup")
		fmt.Println("3. Salir")
		fmt.Scanln(&op)
		switch op {
		case 1:
			fmt.Println(conversacion)
			break
		case 2:
			escribirArchivo()
			break
		case 3:
			servidorActivo = false
			return
			break
		}
	}
	var x string
	fmt.Scanln(&x)
}
