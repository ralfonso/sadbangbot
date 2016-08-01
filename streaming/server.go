package streaming

import (
	"fmt"
	"net/url"

	"github.com/ChimeraCoder/anaconda"
	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
)

type FilterFunc func(anaconda.Tweet) bool

type streamingServer struct {
	ctx       context.Context
	api       *anaconda.TwitterApi
	stream    *anaconda.Stream
	filter    FilterFunc
	workQueue chan interface{}
	Publisher func(string, url.Values) (anaconda.Tweet, error)
}

func NewStreamingServer(ctx context.Context, api *anaconda.TwitterApi, filter FilterFunc) *streamingServer {
	return &streamingServer{
		ctx:       ctx,
		api:       api,
		filter:    filter,
		workQueue: make(chan interface{}, 100),
		Publisher: api.PostTweet,
	}
}

func (s *streamingServer) Start() {
	v := url.Values{}
	v.Add("track", "sad")
	v.Add("lang", "en")
	s.stream = s.api.PublicStreamFilter(v)

	log.Info("streaming server started")

	go func() {
		for {
			select {
			case <-s.ctx.Done():
				return
			case msg := <-s.workQueue:
				s.handler(msg)
			}
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			s.stream.Stop()
			return
		case msg := <-s.stream.C:
			s.workQueue <- msg
		}
	}
}

func (s *streamingServer) handler(v interface{}) {
	if msg, ok := v.(anaconda.Tweet); ok {
		if s.filter(msg) {
			_ = s.reply(msg)
		}
	}
}

func (s *streamingServer) reply(msg anaconda.Tweet) error {
	tweet, err := s.Publisher(fmt.Sprintf("This must be hard for you. I'm sorry. %s", statusLink(msg)), nil)
	if err != nil {
		log.WithError(err).Error("unable to post tweet reply")
		return err
	}
	log.WithFields(log.Fields{
		"tweet.text": tweet.Text,
	}).Info("posted tweet reply")
	return nil
}

func statusLink(msg anaconda.Tweet) string {
	return fmt.Sprintf("https://twitter.com/%s/status/%s", msg.User.ScreenName, msg.IdStr)
}
