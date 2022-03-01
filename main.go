package main

import (
	"fmt"
	"sync"
)

var stop bool = false

func produce(products chan<- int, max int, wg *sync.WaitGroup) {
	for i := 1; i <= max; i++ {
		products <- i
	}
	close(products)
	defer wg.Done()
}

func saleswoman(id int, warehouse <-chan int, bin chan<- int, typ string, wg *sync.WaitGroup) {
	for {
		if len(warehouse) < 1 {
			break
		}

		x := <-warehouse
		fmt.Printf("Kobieta %v dodaje do kosza %v %v \n", id, typ, x)

		bin <- x
	}

	stop = true
	defer wg.Done()
}

func messenger(id int, candels_price <-chan int, flower_bouquet_price <-chan int, wg *sync.WaitGroup) {
	for {
		if len(candels_price) < 2 || len(flower_bouquet_price) < 1 {
			if stop {
				break
			} else {
				continue
			}
		}

		candel_1 := <-candels_price
		candel_2 := <-candels_price
		flower_bouquet := <-flower_bouquet_price

		fmt.Printf("Poslaniec %v odbiera znicze %v i %v oraz wiazanke %v\n", id, candel_1, candel_2, flower_bouquet)
	}

	defer wg.Done()
}

func main() {
	candels := 100
	flower_bouquet := 50
	saleswoamns := 2
	messengers := 5

	candels_warehouse := make(chan int, candels)
	flower_bouquet_warehouse := make(chan int, flower_bouquet)

	wg1 := new(sync.WaitGroup)
	wg1.Add(2)

	go produce(candels_warehouse, candels, wg1)
	go produce(flower_bouquet_warehouse, flower_bouquet, wg1)

	wg1.Wait()

	candels_price := make(chan int, 10)
	flower_bouquet_price := make(chan int, 10)

	wg2 := new(sync.WaitGroup)
	wg2.Add(messengers)
	wg2.Add(saleswoamns * 2)

	for i := 1; i <= saleswoamns; i++ {
		go saleswoman(i, candels_warehouse, candels_price, "znicz", wg2)
		go saleswoman(i, flower_bouquet_warehouse, flower_bouquet_price, "wiazanke", wg2)
	}

	for i := 1; i <= messengers; i++ {
		go messenger(i, candels_price, flower_bouquet_price, wg2)
	}

	wg2.Wait()
}
