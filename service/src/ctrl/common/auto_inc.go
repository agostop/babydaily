package common

import (
	"encoding/binary"
	leveldbhandle "management_backend/src/db/leveldb"

	"github.com/pkg/errors"
)

type AutoInc struct {
	start     int
	step      int
	queue     chan int
	isRunning bool
	db        *leveldbhandle.LevelDbHandle
}

const MAX_QUEUE_SIZE int = 3

var (
	instance *AutoInc
	idKey    string = "id"
)

func GetAutoIncInstance() *AutoInc {
	return instance
}

func InitAutoInc(db *leveldbhandle.LevelDbHandle) {
	b, err := db.Get([]byte(idKey))
	if err != nil {
		log.Errorf("error get id in leveldb.")
		panic("failed init id")
	}

	if len(b) == 0 {
		b = ConvertIntToByte(1)
	}

	currentId := getValueFromInt(b)
	instance = NewAutoInc(currentId, 1, db)
}

func NewAutoInc(start int, step int, db *leveldbhandle.LevelDbHandle) *AutoInc {

	a := &AutoInc{
		start: start,
		step:  step,
		queue: make(chan int, MAX_QUEUE_SIZE),
		db:    db,
	}
	a.isRunning = true
	go a.process()
	return a
}

func (a *AutoInc) process() {
	for i := a.start; a.isRunning; i = i + a.step {
		a.queue <- i
	}
}

func (a *AutoInc) Id() (int, error) {
	id := <-a.queue
	if err := a.db.Put([]byte(idKey), ConvertIntToByte(uint64(id))); err != nil {
		return -1, errors.Errorf("error put number into leveldb")
	}
	return id, nil
}

func (a *AutoInc) Close() {
	a.isRunning = false
	close(a.queue)
}

func getValueFromInt(b []byte) int {
	u := binary.BigEndian.Uint64(b)
	return int(u)
}
