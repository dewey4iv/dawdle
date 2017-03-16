package processor

import (
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/dewey4iv/dawdle"
	"github.com/dewey4iv/timestamps"
)

// New returns a new Processor
func New(opts ...Option) (dawdle.Processor, error) {
	p := Processor{
		internalDelay: time.Second,
		taskCh:        make(chan dawdle.Task),
		workerCh:      make(chan chan dawdle.Task),
		workerCount:   10,
		quitCh:        make(chan struct{}),
		exitCh:        make(chan struct{}),
	}

	for _, opt := range opts {
		if err := opt.Apply(&p); err != nil {
			return nil, err
		}
	}

	if err := p.setup(); err != nil {
		return nil, err
	}

	go p.loop()

	return &p, nil
}

// Processor is the default implmentation of a Processor
type Processor struct {
	store         dawdle.Store
	registrar     dawdle.Registrar
	internalDelay time.Duration
	taskCh        chan dawdle.Task
	workerCh      chan chan dawdle.Task
	workerCount   int
	workers       []*worker
	quitCh        chan struct{}
	exitCh        chan struct{}
}

func (p *Processor) setup() error {
	for i := 0; i < p.workerCount; i++ {
		p.workers = append(p.workers, newWorker(p.store, p.registrar, p.workerCh, p.quitCh, p.exitCh))
	}

	go func() {
		for {
			select {
			case <-p.quitCh:
				return
			default:
				task, err := p.store.Tasks().GetOne()
				if err != nil {
					if _, ok := err.(dawdle.ErrNoPendingTasks); ok {
						time.Sleep(p.internalDelay)
						continue
					}

					// REVIEW: doesn't feel right to continue
					// perhaps we should close - or wait
					// - or simply just log that it happened
					continue
				}

				p.taskCh <- *task
			}
		}
	}()

	return nil
}

func (p *Processor) loop() {
	for {
		task := <-p.taskCh

		select {
		case worker := <-p.workerCh:
			worker <- task
		case <-p.quitCh:
			// wait for all workers to return
			for _ = range p.workers {
				<-p.exitCh
			}

			return
		}
	}
}

func newWorker(store dawdle.Store, registrar dawdle.Registrar, workerCh chan chan dawdle.Task, quitCh chan struct{}, exitCh chan struct{}) *worker {
	w := &worker{
		store:     store,
		registrar: registrar,
		workerCh:  workerCh,
		quitCh:    quitCh,
		workCh:    make(chan dawdle.Task),
		exitCh:    exitCh,
	}

	go w.loop()

	return w
}

type worker struct {
	store     dawdle.Store
	registrar dawdle.Registrar
	workerCh  chan chan dawdle.Task
	workCh    chan dawdle.Task
	quitCh    chan struct{}
	exitCh    chan struct{}
}

func (w *worker) loop() {
	for {
		w.workerCh <- w.workCh

		select {
		case task := <-w.workCh:
			if err := Invoke(w.store, w.registrar, task); err != nil {
				logrus.Error(err)
				// TODO: handle this?
			}
		case <-w.quitCh:
			w.exitCh <- struct{}{}
			return
		}
	}
}

// Invoke takes a store, registrar and task and invokes Perform
func Invoke(store dawdle.Store, registrar dawdle.Registrar, task dawdle.Task) error {
	in := dawdle.NewInvocation(task.ID)

	// save regardless of the outcome
	defer func() {
		if err := store.Tasks().Save(task); err != nil {
			// TODO: log
		}

		if err := store.Invocations().Save(in); err != nil {
			// TODO: log
		}
	}()

	converter, err := registrar.Fetch(task.FuncName)
	if err != nil {
		handleFailure(&task, &in, err)
		return err
	}

	performer, err := converter(task)
	if err != nil {
		handleFailure(&task, &in, err)
		return err
	}

	if err = performer.Perform(); err != nil {
		handleFailure(&task, &in, err)
		return err
	}

	// edit before saving
	task.Status = dawdle.PassedTaskStatus
	task.Mark(timestamps.Updated)
	in.Result = true

	return nil
}

func handleFailure(t *dawdle.Task, in *dawdle.Invocation, err error) {
	t.Status = dawdle.FailedTaskStatus
	in.Result = false
	in.Err = err
}

// Option defines something that builds a Processor
type Option interface {
	Apply(*Processor) error
}

// WithStore sets the store
func WithStore(store dawdle.Store) Option {
	return &withStore{store}
}

type withStore struct {
	store dawdle.Store
}

func (opt *withStore) Apply(p *Processor) error {
	p.store = opt.store
	return nil
}

// WithRegistrar sets the registrar
func WithRegistrar(r dawdle.Registrar) Option {
	return &withRegistrar{r}
}

type withRegistrar struct {
	r dawdle.Registrar
}

func (opt *withRegistrar) Apply(p *Processor) error {
	p.registrar = opt.r

	return nil
}
