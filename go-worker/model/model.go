package model

import "strings"

type Address struct {
	ZipCode  string `json:"zipCode"`
	Address  string `json:"address"`
	District string `json:"district"`
	City     string `json:"city"`
	State    string `json:"state"`
}

type ViaCepAddress struct {
	ZipCode  string `json:"cep"`
	Address  string `json:"logradouro"`
	District string `json:"bairro"`
	City     string `json:"localidade"`
	State    string `json:"uf"`
}

func (v *ViaCepAddress) NormalizdZipCode() {
	v.ZipCode = strings.Replace(v.ZipCode, "-", "", 1)
}
