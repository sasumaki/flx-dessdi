package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

/*
Example Output:
☁️  cloudevents.Event
Validation: valid
Context Attributes,
  specversion: 1.0
  type: dev.knative.eventing.samples.heartbeat
  source: https://knative.dev/eventing-contrib/cmd/heartbeats/#event-test/mypod
  id: 2b72d7bf-c38f-4a98-a433-608fbcdd2596
  time: 2019-10-18T15:23:20.809775386Z
  contenttype: application/json
Extensions,
  beats: true
  heart: yes
  the: 42
Data,
  "data": {
      "names": [],
      "tensor": {
        "shape": [
          1
        ],
        "values": [
          0
        ]
      }
    },
*/

// CloudData is dada
type CloudData struct {
	Data struct {
		Names  []interface{} `json:"names"`
		Tensor struct {
			Shape  []int `json:"shape"`
			Values []int `json:"values"`
		} `json:"tensor"`
	} `json:"data"`
	Meta struct {
		Tags struct {
			ModelURI     string `json:"model_uri"`
			ModelVersion string `json:"model_version"`
		} `json:"tags"`
	} `json:"meta"`
}

func display(event cloudevents.Event) {
	var DataJSON CloudData

	fmt.Printf("☁️  cloudevents.Event\n%s", event.String())
	dada := event.Data()
	fmt.Println(dada)

	err := json.Unmarshal(dada, &DataJSON)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(strconv.Itoa(DataJSON.Data.Tensor.Values[0]))

}

func main() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}
	fmt.Println("starting")
	log.Fatal(c.StartReceiver(context.Background(), display))
}
