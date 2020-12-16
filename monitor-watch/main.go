package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

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

type clouddata struct {
	data struct {
		names  []string
		tensor struct {
			shape  []int
			values []int
		}
	}
}

func display(event cloudevents.Event) {
	var JSONAsInterface clouddata

	fmt.Printf("☁️  cloudevents.Event\n%s", event.String())
	dada := event.Data()
	fmt.Println(dada)
	if err := json.Unmarshal(dada, &JSONAsInterface); err != nil {
		fmt.Println(JSONAsInterface.data.tensor.values[0])
	}
}

func main() {
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal("Failed to create client, ", err)
	}
	fmt.Println("starting")
	log.Fatal(c.StartReceiver(context.Background(), display))
}
