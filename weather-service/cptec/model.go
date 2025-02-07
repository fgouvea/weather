package cptec

import "encoding/xml"

type CitiesResponse struct {
	XMLName xml.Name `xml:"cidades"`
	Cities  []City   `xml:"cidade"`
}

type City struct {
	ID    string `xml:"id"`
	Name  string `xml:"nome"`
	State string `xml:"uf"`
}
