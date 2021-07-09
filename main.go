/*
Copyright (c) 2021 TriggerMesh Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"os"
	"strconv"
	"time"

	cloudevents "github.com/cloudevents/sdk-go"
	log "github.com/sirupsen/logrus"
)

// Dce runs a CloudEvents receiver.
type Dce struct {
	EventData    string
	EventSource  string
	EventSubject string
	EventID      string
	EventType    string
	TargetURL    string
	Timeout      string
	ceClient     cloudevents.Client
}

func main() {
	d := &Dce{}
	d.new()
	i, err := strconv.Atoi(d.Timeout)
	if err != nil {
		log.Fatal("error converting TIMEOUT to a string: ", err)
	}
	for {
		d.sendData()
		time.Sleep(time.Duration(i) * time.Second)
	}
}

func (dce *Dce) new() {
	dce.EventID = os.Getenv("EVENT_ID")
	if dce.EventID == "" {
		log.Fatal("EVENT_ID must be set")
	}
	dce.EventSource = os.Getenv("EVENT_SOURCE")
	if dce.EventSource == "" {
		log.Fatal("EVENT_SOURCE must be set")
	}
	dce.EventSubject = os.Getenv("EVENT_SUBJECT")
	if dce.EventSubject == "" {
		log.Fatal("EVENT_SUBJECT must be set")
	}
	dce.EventType = os.Getenv("EVENT_TYPE")
	if dce.EventType == "" {
		log.Fatal("EVENT_TYPE must be set")
	}
	dce.EventData = os.Getenv("EVENT_DATA")
	if dce.EventData == "" {
		log.Fatal("EVENT_DATA must be set")
	}
	dce.Timeout = os.Getenv("TIMEOUT")
	if dce.Timeout == "" {
		log.Fatal("TIMEOUT must be set")
	}
	sink := os.Getenv("K_SINK")
	if sink == "" {
		log.Fatal("K_SINK must be set")
	}

	t, err := cloudevents.NewHTTPTransport(
		cloudevents.WithTarget(sink),
	)
	if err != nil {
		log.Fatal("cloudevents.NewHTTPTransport failed: ", err)
	}

	dce.ceClient, err = cloudevents.NewClient(t, cloudevents.WithTimeNow())
	if err != nil {
		log.Fatal("cloudevents.NewClient failed: ", err)
	}
}

func (dce *Dce) sendData() {
	log.Info("sending data..")

	eventToSend := cloudevents.NewEvent("1.0")
	eventToSend.SetType(dce.EventType)
	eventToSend.SetSource(dce.EventSource)
	eventToSend.SetSubject(dce.EventSubject)
	eventToSend.SetID(dce.EventID)
	eventToSend.SetDataContentType(cloudevents.ApplicationJSON)
	b := []byte(dce.EventData)
	err := eventToSend.SetData(b)
	if err != nil {
		log.Error(err)
		return
	}

	_, _, err = dce.ceClient.Send(context.Background(), eventToSend)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Sent cloud event sucessfully")
}
