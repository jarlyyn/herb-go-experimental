package nsq

import (
	"net"
	"time"

	"github.com/jarlyyn/herb-go-experimental/messagequeue"
	nsq "github.com/nsqio/go-nsq"
)

type Config struct {
	Addr                            string
	Topic                           string
	Chanel                          string
	LookupAddr                      string
	LocalAddr                       string
	DialTimeoutInSecond             int64
	ReadTimeoutInSecond             int64
	WriteTimeoutInSecond            int64
	LookupdPollIntervalInSecond     int64
	LookupdPollJitter               float64
	MaxRequeueDelayInSecond         int64
	DefaultRequeueDelayInSecond     int64
	MaxBackoffDurationInSecond      int64
	BackoffMultiplierInSecond       int64
	MaxAttempts                     uint16
	LowRdyIdleTimeoutInSecond       int64
	LowRdyTimeoutInSecond           int64
	RDYRedistributeIntervalInSecond int64
	ClientID                        string
	Hostname                        string
	UserAgent                       string
	HeartbeatIntervalInSecond       int64
	SampleRate                      int32
	Deflate                         bool
	DeflateLevel                    int
	Snappy                          bool
	OutputBufferSize                int64
	OutputBufferTimeoutInSecond     int64
	MaxInFlight                     int
	MsgTimeoutInSecond              int64
	AuthSecret                      string
}

func (c *Config) ApplyTo(q *Queue) error {
	q.Addr = c.Addr
	q.Topic = c.Topic
	q.Chanel = c.Chanel
	q.LookupAddr = c.LookupAddr
	if c.LocalAddr != "" {
		addr, err := net.ResolveTCPAddr("tcp", c.LocalAddr)
		if err != nil {
			return err
		}
		q.Config.LocalAddr = addr
	}
	if c.DialTimeoutInSecond != 0 {
		q.Config.DialTimeout = time.Duration(c.DialTimeoutInSecond) * time.Second
	}
	if c.ReadTimeoutInSecond != 0 {
		q.Config.ReadTimeout = time.Duration(c.ReadTimeoutInSecond) * time.Second
	}
	if c.WriteTimeoutInSecond != 0 {
		q.Config.WriteTimeout = time.Duration(c.WriteTimeoutInSecond) * time.Second
	}
	if c.LookupdPollIntervalInSecond != 0 {
		q.Config.LookupdPollInterval = time.Duration(c.LookupdPollIntervalInSecond) * time.Second
	}
	if c.LookupdPollJitter != 0 {
		q.Config.LookupdPollJitter = c.LookupdPollJitter
	}
	if c.MaxRequeueDelayInSecond != 0 {
		q.Config.MaxRequeueDelay = time.Duration(c.MaxRequeueDelayInSecond) * time.Second
	}
	if c.DefaultRequeueDelayInSecond != 0 {
		q.Config.DefaultRequeueDelay = time.Duration(c.DefaultRequeueDelayInSecond) * time.Second
	}
	if c.MaxBackoffDurationInSecond != 0 {
		q.Config.MaxBackoffDuration = time.Duration(c.MaxBackoffDurationInSecond) * time.Second
	}
	if c.BackoffMultiplierInSecond != 0 {
		q.Config.BackoffMultiplier = time.Duration(c.BackoffMultiplierInSecond) * time.Second
	}
	if c.MaxAttempts != 0 {
		q.Config.MaxAttempts = c.MaxAttempts
	}
	if c.LowRdyIdleTimeoutInSecond != 0 {
		q.Config.LowRdyIdleTimeout = time.Duration(c.LowRdyIdleTimeoutInSecond) * time.Second
	}
	if c.LowRdyTimeoutInSecond != 0 {
		q.Config.LowRdyTimeout = time.Duration(c.LowRdyTimeoutInSecond) * time.Second
	}
	if c.RDYRedistributeIntervalInSecond != 0 {
		q.Config.RDYRedistributeInterval = time.Duration(c.RDYRedistributeIntervalInSecond) * time.Second
	}
	if c.ClientID != "" {
		q.Config.ClientID = c.ClientID
	}
	if c.Hostname != "" {
		q.Config.Hostname = c.Hostname
	}
	if c.UserAgent != "" {
		q.Config.UserAgent = c.UserAgent
	}
	if c.HeartbeatIntervalInSecond != 0 {
		q.Config.HeartbeatInterval = time.Duration(c.HeartbeatIntervalInSecond) * time.Second
	}
	if c.SampleRate != 0 {
		q.Config.SampleRate = c.SampleRate
	}
	if c.Deflate != false {
		q.Config.Deflate = c.Deflate
	}
	if c.DeflateLevel != 0 {
		q.Config.DeflateLevel = c.DeflateLevel
	}
	if c.Snappy != false {
		q.Config.Snappy = c.Snappy
	}
	if c.OutputBufferSize != 0 {
		q.Config.OutputBufferSize = c.OutputBufferSize
	}
	if c.OutputBufferTimeoutInSecond != 0 {
		q.Config.OutputBufferTimeout = time.Duration(c.OutputBufferTimeoutInSecond) * time.Second
	}
	if c.MaxInFlight != 0 {
		q.Config.MaxInFlight = c.MaxInFlight
	}
	if c.MsgTimeoutInSecond != 0 {
		q.Config.MsgTimeout = time.Duration(c.MsgTimeoutInSecond) * time.Second
	}
	if c.AuthSecret != "" {
		q.Config.AuthSecret = c.AuthSecret
	}
	return nil
}
func NewConfig() *Config {
	return &Config{}
}

type Queue struct {
	Addr       string
	Topic      string
	Chanel     string
	LookupAddr string
	Config     *nsq.Config
	Producer   *nsq.Producer
	Consumer   *nsq.Consumer
	consumer   func([]byte) messagequeue.ConsumerStatus
	recover    func()
}

func (q *Queue) SetRecover(r func()) {
	q.recover = r
}
func (q *Queue) Hanlder(message *nsq.Message) error {
	q.consumer(message.Body)
	return nil
}
func (q *Queue) Start() error {
	var err error
	q.Producer, err = nsq.NewProducer(q.Addr, q.Config)
	if err != nil {
		return err
	}
	q.Consumer, err = nsq.NewConsumer(q.Topic, q.Chanel, q.Config)
	if err != nil {
		return err
	}
	q.Consumer.AddHandler(nsq.HandlerFunc(q.Hanlder))
	if q.LookupAddr != "" {
		err = q.Consumer.ConnectToNSQLookupd(q.LookupAddr)
	} else {
		err = q.Consumer.ConnectToNSQD(q.Addr)
	}
	if err != nil {
		return err
	}
	return nil

}
func (q *Queue) Close() error {
	q.Producer.Stop()
	q.Consumer.Stop()
	return nil
}
func (q *Queue) ProduceMessages(messages ...[]byte) (sent []bool, err error) {
	sent = make([]bool, len(messages))
	for k := range messages {
		err := q.Producer.Publish(q.Topic, messages[k])
		if err != nil {
			return sent, err
		}
		sent[k] = true
	}
	return sent, nil
}
func (q *Queue) SetConsumer(c func([]byte) messagequeue.ConsumerStatus) {
	q.consumer = c
}

func NewQueue() *Queue {
	return &Queue{
		recover: func() {},
		Config:  nsq.NewConfig(),
	}
}

func QueueFactory(conf messagequeue.Config, prefix string) (messagequeue.Driver, error) {
	c := NewConfig()
	var err error
	err = conf.Get(prefix+"Addr", &c.Addr)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Topic", &c.Topic)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Chanel", &c.Chanel)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"LookupAddr", &c.LookupAddr)
	if err != nil {
		return nil, err
	}
	q := NewQueue()
	err = c.ApplyTo(q)
	if err != nil {
		return nil, err
	}
	return q, nil
}

func init() {
	messagequeue.Register("nsq", QueueFactory)
}
