package leveldb

import (
	"fmt"
	loggers "management_backend/src/logger"
	"os"
	"sync"

	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type LevelDbHandle struct {
	db        *leveldb.DB
	writeLock sync.Mutex
}

var (
	log = loggers.GetLogger(loggers.ModuleWeb)
)

var instance *LevelDbHandle

func GetHandleInstance() *LevelDbHandle {
	dbFolder := "./db"
	if instance != nil {
		return instance
	}

	err := createDbPathIfNotExist(dbFolder)
	if err != nil {
		panic(fmt.Sprintf("Error create dir %s by leveldbHandle: %s", dbFolder, err))
	}

	db, err := leveldb.OpenFile(dbFolder, nil)
	if err != nil {
		panic(fmt.Sprintf("the db file open failed: %s", err))
	}

	instance = &LevelDbHandle{
		db: db,
	}

	return instance
}

func (h *LevelDbHandle) Put(key string, value []byte) error {
	if value == nil {
		log.Warn("the value is nil.")
		return errors.New("the value is nil.")
	}
	err := h.db.Put([]byte(key), []byte(value), &opt.WriteOptions{Sync: false})
	if err != nil {
		log.Errorf("writing failed. key [%#v]", key)
		return errors.Wrapf(err, "error writing leveldb. key [%#v]", key)
	}

	return err
}

func (h *LevelDbHandle) Delete(key string) error {
	err := h.db.Delete([]byte(key), &opt.WriteOptions{Sync: false})
	if err != nil {
		log.Errorf("deleting key failed, key: [%#v]", key)
		return errors.Wrapf(err, "error deleting leveldb, key: [%#v]", key)
	}
	return err
}

func (h *LevelDbHandle) Get(key string) ([]byte, error) {
	value, err := h.db.Get([]byte(key), nil)
	if err == leveldb.ErrNotFound {
		value = nil
		err = nil
	}
	if err != nil {
		log.Errorf("getting leveldbprovider key [%#v], err:%s", key, err.Error())
		return nil, errors.Wrapf(err, "error getting leveldbhandle key [%#v]", []byte(key))
	}
	return value, nil
}

func (h *LevelDbHandle) BatchPut(batch *leveldb.Batch) error {

	h.writeLock.Lock()
	defer h.writeLock.Unlock()

	if err := h.db.Write(batch, nil); err != nil {
		log.Errorf("write batch to leveldb failed.")
		return errors.Wrap(err, "error write batch to leveldb.")
	}

	return nil
}

func (h *LevelDbHandle) IteratorWithPrefix(prefix string) ([]string, error) {
	if len(prefix) == 0 {
		return nil, errors.Errorf("iterator prefix should not be empty key.")
	}

	result := []string{}

	r := util.BytesPrefix([]byte(prefix))
	keyRange := &util.Range{Start: r.Start, Limit: r.Limit}
	it := h.db.NewIterator(keyRange, nil)
	defer it.Release()
	for it.Next() {
		result = append(result, string(it.Value()))
	}
	return result, nil
}

func createDbPathIfNotExist(path string) error {
	_, err := os.Stat(path)
	if err == nil {
		return nil
	}

	if os.IsNotExist(err) {
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
