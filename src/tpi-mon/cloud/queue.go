package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	uuid "github.com/satori/go.uuid"
)

var qMessageSep = []byte(":")

type qMessage struct {
	queue      *queue
	routingKey string
	id         []byte
	expires    time.Time
	data       []byte
}

type queue struct {
	id          string
	redisClient *redis.Client
}

func newQueue(redisHost string, redisPort uint16) (*queue, error) {
	addr := fmt.Sprintf("%s:%d", redisHost, redisPort)
	redisClient := redis.NewClient(&redis.Options{
		Addr:       addr,
		MaxRetries: 8,
	})

	if err := redisClient.Ping().Err(); err != nil {
		return nil, err
	}

	q := &queue{
		id:          uuid.NewV4().String(),
		redisClient: redisClient,
	}

	return q, nil
}

func (q *queue) publish(routingKey string, data []byte) error {
	return q.publishEx(routingKey, data, time.Time{})
}

func (q *queue) publishEx(routingKey string, data []byte, expires time.Time) error {
	msg := newQMessage(q, routingKey, data, expires)
	serialized := msg.marshal()
	return q.redisClient.RPush(routingKey, serialized).Err()
}

func (q *queue) startConsumeLoop(routingKey string, process func(qMessage) error) {

	go func() {
		for {

			processingList := q.processingQueueName(routingKey)
			if err := q.redisClient.BRPopLPush(routingKey, processingList, time.Duration(0)).Err(); err != nil {
				// TODO: need better error handling
				panic(err)
			}

			for {
				vals, err := q.redisClient.LRange(processingList, 0, 128).Result()

				if err != nil {
					// TODO: need better error handling
					panic(err)
				}

				if len(vals) == 0 {
					break
				}

				for _, val := range vals {
					msg := newQMessageUnmarshalled(q, routingKey, []byte(val))

					if !msg.expires.IsZero() && msg.expires.Before(time.Now()) {
						logger.Printf("Discarding expired msg %v\n", msg)
						continue
					}

					if err := process(msg); err != nil {
						//TODO: a good way to requeue a failed msg
						logger.Panicf("queue:%v failed to process message %v: %v\n", routingKey, msg, err)
					}

					msg.ack()
				}

			}
		}
	}()
}

func (q *queue) processingQueueName(routingKey string) string {
	return "processing:" + q.id + ":" + routingKey
}

func (msg qMessage) ack() error {
	processingList := msg.queue.processingQueueName(msg.routingKey)
	serialized := msg.marshal()
	nRemoved, err := msg.queue.redisClient.LRem(processingList, 1, serialized).Result()
	if err != nil {
		return err
	}

	if nRemoved == 0 {
		panic(fmt.Errorf("Invalid state: processed message was already removed: %v", msg))
	}

	return nil
}

func newQMessage(q *queue, routingKey string, data []byte, expires time.Time) qMessage {
	id, _ := uuid.NewV4().MarshalBinary()
	return qMessage{
		queue:      q,
		routingKey: routingKey,

		id:      id,
		data:    data,
		expires: expires,
	}
}

func newQMessageUnmarshalled(q *queue, routingKey string, unmarshalled []byte) qMessage {
	msg := qMessage{
		queue:      q,
		routingKey: routingKey,
	}

	msg.unmarshal(unmarshalled)
	return msg
}

func (msg qMessage) marshal() []byte {
	expiresBytes, err := msg.expires.MarshalBinary()
	if err != nil {
		panic(fmt.Errorf("Invalid expiry %v: %v", msg.expires, err))
	}

	serialized := bytes.Join([][]byte{
		msg.id,
		expiresBytes,
		msg.data,
	}, qMessageSep)

	return serialized
}

func (msg qMessage) unmarshal(marshalled []byte) {

	fields := bytes.SplitN(marshalled, qMessageSep, 3)
	if len(fields) != 3 {
		panic(fmt.Errorf("bad data: no enough separators: %v", marshalled))
	}

	id := fields[0]

	var expires time.Time
	if err := expires.UnmarshalBinary(fields[1]); err != nil {
		panic(fmt.Errorf("bad data: invalid expires: %v", fields[1]))
	}

	data := fields[2]

	msg.id = id
	msg.data = data
	msg.expires = expires
}
