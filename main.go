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

type kBugger struct {
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
	d := &kBugger{}
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

func (kBugger *kBugger) new() {
	kBugger.EventID = os.Getenv("EVENT_ID")
	if kBugger.EventID == "" {
		log.Fatal("EVENT_ID must be set")
	}
	kBugger.EventSource = os.Getenv("EVENT_SOURCE")
	if kBugger.EventSource == "" {
		log.Fatal("EVENT_SOURCE must be set")
	}
	kBugger.EventSubject = os.Getenv("EVENT_SUBJECT")
	if kBugger.EventSubject == "" {
		log.Fatal("EVENT_SUBJECT must be set")
	}
	kBugger.EventType = os.Getenv("EVENT_TYPE")
	if kBugger.EventType == "" {
		log.Fatal("EVENT_TYPE must be set")
	}
	kBugger.EventData = os.Getenv("EVENT_DATA")
	if kBugger.EventData == "" {
		log.Fatal("EVENT_DATA must be set")
	}
	kBugger.Timeout = os.Getenv("TIMEOUT")
	if kBugger.Timeout == "" {
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

	kBugger.ceClient, err = cloudevents.NewClient(t, cloudevents.WithTimeNow())
	if err != nil {
		log.Fatal("cloudevents.NewClient failed: ", err)
	}
}

func (kBugger *kBugger) sendData() {
	log.Info("sending data..")

	eventToSend := cloudevents.NewEvent("1.0")
	eventToSend.SetType(kBugger.EventType)
	eventToSend.SetSource(kBugger.EventSource)
	eventToSend.SetSubject(kBugger.EventSubject)
	eventToSend.SetID(kBugger.EventID)
	eventToSend.SetDataContentType(cloudevents.ApplicationJSON)
	b := []byte(kBugger.EventData)
	err := eventToSend.SetData(b)
	if err != nil {
		log.Error(err)
		return
	}

	_, _, err = kBugger.ceClient.Send(context.Background(), eventToSend)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("Sent cloud event sucessfully")
}
