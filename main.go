package main

import (
	"fmt"
	"time"
	"os"
	"os/exec"
)

type proceso struct { //Estructura de los datos que se ocuparian.
	id int
	i chan uint64
	hecho chan bool
	cerrar bool
}

// En esta funcion se pide un numero entero el cual se utilizara para evaluar
func (proc *proceso) incremento_numero(numero uint64){ 
	var i uint64 = 0 // el proceso se iniciara en 0
	//Comienza la concurrencia
	for{
		if i == numero || proc.cerrar { //si un numero se repite o el proceso termino
			return
		}
		i++ //incrementa el valor
		proc.i <- i //paso de información hacia el canal
	}
}

func (proc *proceso) terminarProceso(){
	//Ayudara a poder terminar el proceso
	for{
		select{
		case <-proc.i://Segun cada cuando se iran a crear
			time.Sleep(time.Millisecond * 500)
		case <-proc.hecho:
			proc.cerrar = true //se levanta bandera para cerrar o terminar proceso
			return
		}
	}
}
//funcion para comenzar a contar
func (proc *proceso) empezarConteo(){
	contador := <- proc.i //paso del valor del contador
	fmt.Println(contador)
}

//func Proceso(id uint64) {
	//i := uint64(0);
	//for {
	//	fmt.Printf("id %d: %d \n", id, i)
	//	i = i + 1
	//	time.Sleep(time.Millisecond * 500)
	//}
//}

func main(){
	var opcion int
	enListar := make([]proceso,0,1) //slice para los procesos el cual se pasa la estructura y valores principales
	sigId := 0
	//var id uint64
	//chanel := make(chan uint64) // el chanel para la funcion

	for{
		fmt.Println("\t Administrador de procesos: \n\t 1) Agregar proceso \n\t 2) Mostrar procesos \n\t 3) terminar procesos \n\t 4) salir ")
		fmt.Scan(&opcion)
		if opcion == 1 {
			const maxCant = ^uint64(0) //Se utiliza de esta forma para tener un espacio delimetado
			proc := proceso{id: sigId, i: make(chan uint64), hecho: make(chan bool), cerrar: false}//pasamos los adatos a la estructura.
			enListar = append(enListar,proc) //paso de datos a un slice
			go proc.incremento_numero(maxCant)
			go proc.terminarProceso()
			sigId++

			//id = id + 1
			//go Proceso(id)
			//id <- 1 
			//m := <- id
			//go Proceso(m)
		}else if opcion == 2 {
			if len(enListar) > 0 {
				salir := make(chan bool)
				go func() {
					for{
						select {
						case <- salir:
							return //en caso de que se termine de mostrar retornar.
						default:
							for _, v := range enListar {
								fmt.Printf("id %d: ", v.id)
								v.empezarConteo()
							}
							time.Sleep(time.Second)
							out, _ := exec.Command("clear").Output()
							os.Stdout.Write(out)
						}
					}
				}()
				var a string
				fmt.Scan(&a)
				if(a == "2"){				 
     				salir <- true
      				//break
				}
			} else {
				fmt.Println("No hay procesos en ejecucion")
			}
		}else if opcion == 3 {
			var i int
			fmt.Print("Ingresa id : ")
			fmt.Scan(&i)
			var element *proceso
			for idx, e := range enListar {
				if e.id == i {
					element = &e
					enListar = append(enListar[:idx], enListar[idx+1:]...)
					break
				}
			}
			element.hecho <- false
		}else if opcion == 4 {
			break
		}
	}
}