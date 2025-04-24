package core

import (
	"context"
	"fmt"
	"log/slog"
	"slices"
	"sync"
)

type Service struct {
	mx          sync.Mutex
	updateNow   bool
	log         *slog.Logger
	db          DB
	xkcd        XKCD
	words       Words
	concurrency int
}

func NewService(
	log *slog.Logger, db DB, xkcd XKCD, words Words, concurrency int,
) (*Service, error) {
	if concurrency < 1 {
		return nil, fmt.Errorf("wrong concurrency specified: %d", concurrency)
	}
	return &Service{
		log:         log,
		db:          db,
		xkcd:        xkcd,
		words:       words,
		concurrency: concurrency,
	}, nil
}

func (s *Service) Update(ctx context.Context) (err error) {
	s.log.Info("Update start")
	defer func() {
		s.log.Info("Update finish")
	}()

	if !s.mx.TryLock() {
		return ErrAlreadyExists
	}
	defer s.mx.Unlock()

	s.updateNow = true
	defer func() {
		s.updateNow = false
	}()

	lastID, err := s.xkcd.LastID(ctx)
	if err != nil {
		s.log.Error("failed to fetch last comic ID from XKCD", "error", err)
		return err
	}

	IDs, err := s.db.IDs(ctx)
	if err != nil {
		s.log.Error("failed to retrieve existing comic IDs from database", "error", err)
		return err
	}

	slices.Sort(IDs)

	in1 := make(chan int)

	go func() {
		for i := 1; i <= lastID; i++ {
			if ctx.Err() != nil {
				break
			}
			in1 <- i
		}

		close(in1)
	}()

	var wg sync.WaitGroup
	wg.Add(s.concurrency)

	out1 := make([]chan XKCDInfo, s.concurrency)
	for i := 0; i < s.concurrency; i++ {
		out1[i] = make(chan XKCDInfo)
		// XKCD GET() Получить JSON
		go func() {
			for id := range in1 {
				if ctx.Err() != nil {
					break
				}
				if _, ok := slices.BinarySearch(IDs, id); !ok {
					xkcd, err := s.xkcd.Get(ctx, id)
					if err != nil {
						s.log.Error("failed to fetch comic from XKCD API", "comic_id", id, "error", err)
						continue
					}

					out1[i] <- xkcd
				}
			}
			close(out1[i])
		}()

	}

	in2 := merge1(out1...)

	out2 := make([]chan Comics, s.concurrency)
	for i := 0; i < s.concurrency; i++ {
		out2[i] = make(chan Comics)
		// WORDS.NORM() Нормализация
		go func() {
			for xkcd := range in2 {
				words, err := s.words.Norm(ctx, xkcd.Description)
				if err != nil {
					s.log.Error("failed to process comic keywords", "comic_id", xkcd.ID, "error", err)
					continue
				}

				out2[i] <- Comics{ID: xkcd.ID, URL: xkcd.URL, Words: words}
			}

			close(out2[i])
		}()
	}

	in3 := merge2(out2...)

	for i := 0; i < s.concurrency; i++ {
		// DATABASE.ADD() Сохранение в БД
		go func() {
			defer wg.Done()

			for comics := range in3 {
				err = s.db.Add(ctx, comics)
				if err != nil {
					s.log.Error("failed to persist comic in database", "comic_id", comics.ID, "error", err)
					continue
				}
			}

		}()
	}

	wg.Wait()

	return nil
}

func merge1(in ...chan XKCDInfo) chan XKCDInfo {
	out := make(chan XKCDInfo)

	var wg sync.WaitGroup
	wg.Add(len(in))
	for _, ch := range in {
		go func() {
			defer wg.Done()
			for xkcd := range ch {
				out <- xkcd
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func merge2(in ...chan Comics) chan Comics {
	out := make(chan Comics)

	var wg sync.WaitGroup
	wg.Add(len(in))

	for _, ch := range in {
		go func() {
			defer wg.Done()
			for comics := range ch {
				out <- comics
			}
		}()
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}

func (s *Service) Stats(ctx context.Context) (ServiceStats, error) {
	lastID, err := s.xkcd.LastID(ctx)
	comicsTotal := lastID - 1 // т.к. ресурс под индексом 404 not found

	if err != nil {
		s.log.Error("failed to fetch last comic ID from XKCD", "error", err)
		return ServiceStats{}, err
	}

	DBstat, err := s.db.Stats(ctx)
	if err != nil {
		s.log.Error("failed to retrieve database statistics", "error", err)
		return ServiceStats{}, err
	}

	s.log.Info("Stats", "Result", ServiceStats{DBStats: DBstat, ComicsTotal: comicsTotal})

	return ServiceStats{DBStats: DBstat, ComicsTotal: comicsTotal}, nil
}

func (s *Service) Status(ctx context.Context) ServiceStatus {
	if s.updateNow {
		return StatusRunning
	} else {
		return StatusIdle
	}
}

func (s *Service) Drop(ctx context.Context) error {
	return s.db.Drop(ctx)
}
