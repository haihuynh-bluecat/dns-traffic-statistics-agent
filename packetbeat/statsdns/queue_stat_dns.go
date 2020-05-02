// Copyright 2020 BlueCat Networks (USA) Inc. and its affiliates
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package statsdns

import (
	"time"

	"github.com/elastic/beats/libbeat/logp"
	"github.com/elastic/beats/packetbeat/model"
)

type (
	// DNS data after decode
	QueryDNS struct {
		srcIP   string
		dstIP   string
	}

	QueueStatDNS struct {
		isActive    bool
		queries     chan *QueryDNS
		records     chan *model.Record
	}
)

func NewQueryDNS(srcIP, dstIP string) (queryDNS *QueryDNS) {
	queryDNS = &QueryDNS{
		srcIP:   srcIP,
		dstIP:   dstIP,
	}
	return
}

func NewQueueStatDNS() (queue *QueueStatDNS) {
	queue = &QueueStatDNS{
		queries:     make(chan *QueryDNS),
		records:     make(chan *model.Record),
		isActive:    false,
	}
	return
}

func (queue *QueueStatDNS) PushStatDNS(queryDNS *QueryDNS, record *model.Record) {
	if !queue.isActive{
		return
	}
	if queryDNS != nil {
		queue.queries <- queryDNS
	} else if record != nil {
		queue.records <- record
	}
}

func (queue *QueueStatDNS) SubStatDNS(flagActive *bool) {
	for queue.isActive {
		if *flagActive == false {
			continue
		}
		select {
		case query := <-queue.queries:
			if query == nil {
				continue
			}
			// CreateCounterMetric(query.srcIP, query.dstIP)
			IncreaseQueryCounter(query.srcIP, query.dstIP, QUERY)
			IncreaseQueryCounterForPerView(query.srcIP, query.dstIP, QUERY)
		case record := <-queue.records:
			if record == nil {
				continue
			}
			time.Sleep(5 * time.Microsecond)
			ReceivedMessage(record)
		}
	}
}

func (queue *QueueStatDNS) Stop() {
	logp.Info("QueueStatDNS Stop")
	queue.isActive = false
	close(queue.queries)
	close(queue.records)
}
