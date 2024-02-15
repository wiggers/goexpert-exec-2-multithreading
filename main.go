package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Brasilapi struct {
	Cep    string `json:"cep"`
	Estado string `json:"state"`
	Cidade string `json:"city"`
	Bairro string `json:"neighborhood"`
	Rua    string `json:"street"`
	Origem string
}

type Viacepapi struct {
	Cep    string `json:"cep"`
	Estado string `json:"uf"`
	Cidade string `json:"localidade"`
	Bairro string `json:"bairro"`
	Rua    string `json:"logradouro"`
	Origem string
}

func main() {

	c1 := make(chan Brasilapi)
	c2 := make(chan Viacepapi)

	ctx := context.Background()
	ctx, cancelCtx := context.WithTimeout(ctx, 1*time.Second)
	defer cancelCtx()

	go link1(ctx, c1)
	go link2(ctx, c2)

	select {
	case msg1 := <-c1:
		fmt.Printf("%+v", msg1)

	case msg2 := <-c2:
		fmt.Printf("%+v", msg2)

	case <-time.After(time.Second * 1):
		fmt.Printf("timeout")
	}
}

func link1(ctx context.Context, c1 chan Brasilapi) {
	resp, err := call(ctx, "https://brasilapi.com.br/api/cep/v1/01153000")
	if err != nil {
		return
	}
	data := Brasilapi{Origem: "BrasilApi"}
	json.Unmarshal(resp, &data)
	if err != nil {
		return
	}

	c1 <- data
}

func link2(ctx context.Context, c2 chan Viacepapi) {
	resp, err := call(ctx, "http://viacep.com.br/ws/01153000/json/")
	if err != nil {
		return
	}
	data := Viacepapi{Origem: "ViaCep"}
	err = json.Unmarshal(resp, &data)
	if err != nil {
		return
	}

	c2 <- data
}

func call(ctx context.Context, address string) ([]byte, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", address, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
