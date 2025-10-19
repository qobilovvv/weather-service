package geocoding

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Name      string `json:"name"`
	Country   string `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Client struct {
	httpClient *http.Client
}

func NewClient(httClient *http.Client) *Client {
	return &Client{
		httpClient: httClient,
	}
}

func (c *Client) GetCoords(city string) (Response, error) {
	res, err := c.httpClient.Get(
		fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1&language=ru&format=json", city),
	)

	if err != nil {
		return Response{}, err
	}
	
	if res.StatusCode != http.StatusOK {
		return Response{}, fmt.Errorf("status code %d", res.StatusCode)
	}
	
	defer res.Body.Close()
	
	var geoResponse struct {
		Results []Response `json:"results"`
	}
	
	err = json.NewDecoder(res.Body).Decode(&geoResponse)
	if err != nil {
		return Response{}, err
	}
	
	return geoResponse.Results[0], err
}
