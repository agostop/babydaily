/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package baby_record

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/syndtr/goleveldb/leveldb"

	"management_backend/src/ctrl/common"
	leveldbhandle "management_backend/src/db/leveldb"
	loggers "management_backend/src/logger"
)

var (
	log = loggers.GetLogger(loggers.ModuleWeb)
)

const PAGE_SIZE int = 7

func BabyRecordHandle(ctx *gin.Context) {
	params := BindAddRecordHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	db := leveldbhandle.GetHandleInstance()

	ts := time.Now().Unix()
	id, err := common.GetAutoIncInstance().Id()
	if err != nil {
		log.Errorf("error get key id.")
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}
	key := getKeyWithPrefix(id)
	val, err := json.Marshal(&RecordMeta{Amount: params.Amount, Time: ts, Type: params.Type, Id: id})
	if err != nil {
		log.Errorf("marshal failed. case: %s", err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalParameters)
		return
	}

	db.Put(key, val)
	if err != nil {
		log.Errorf("leveldb put failed. key: %v", err.Error())
		common.ConvergeFailureResponse(ctx, common.ErrorHandleFailure)
		return
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func BabyRecordEditHandle(ctx *gin.Context) {
	params := BindEditMetaRecordHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	db := leveldbhandle.GetHandleInstance()
	b, err := db.Get(getKeyWithPrefix(params.Id))
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}

	if b == nil {
		common.ConvergeDataResponse(ctx, common.ErrorRecordNotFound, nil)
		return
	}

	r := &RecordMeta{}
	err = json.Unmarshal(b, r)
	if err != nil {
		log.Errorf("unmarshal failed. %v", err.Error())
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}

	r.Amount = params.Amount
	r.Time = params.Time
	r.Type = params.Type

	newVal, err := json.Marshal(r)
	if err != nil {
		log.Errorf("marshal failed. %v", err.Error())
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}

	if err = db.Put(getKeyWithPrefix(r.Id), newVal); err != nil {
		log.Errorf("put key failed. %v", err.Error())
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func BabyRecordListHandle(ctx *gin.Context) {
	db := leveldbhandle.GetHandleInstance()
	rs, err := db.IteratorWithPrefix([]byte(KeyPrefix))
	if err != nil {
		log.Errorf("error get IteratorWithRange: %v", err.Error())
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}

	// 仅输出最新的PAGE_SIZE条
	if len(rs) > PAGE_SIZE {
		rs = rs[:7]
	}

	var records []RecordMeta
	for _, v := range rs {
		var r RecordMeta
		json.Unmarshal([]byte(v), &r)
		records = append(records, r)
	}
	results := NewRecordResultView(records)
	common.ConvergeListResponse(ctx, results, int64(len(results)), nil)
}

func BabyRecordDeleteHandle(ctx *gin.Context) {
	rd := BindTimeOfRecordHandler(ctx)
	db := leveldbhandle.GetHandleInstance()

	if err := db.Delete(getKeyWithPrefix(rd.Id)); err != nil {
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func BabyRecordMaintainHandle(ctx *gin.Context) {
	param := BindMaintainHandler(ctx)
	db := leveldbhandle.GetHandleInstance()
	if param.Type == CleanAll {
		if err := db.CleanAll(); err != nil {
			log.Errorf("clean all failed, %v", err.Error())
			common.ConvergeFailureResponse(ctx, common.InternalError)
			return
		}
		common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
		return

	} else if param.Type == GetAllData {
		res, _ := db.IteratorWithPrefix([]byte(KeyPrefix))
		var allData []RecordMeta
		for _, v := range res {
			var r RecordMeta
			err := json.Unmarshal([]byte(v), &r)
			if err != nil {
				log.Errorf("%v", err.Error())
				continue
			}
			allData = append(allData, r)
		}
		common.ConvergeListResponse(ctx, NewRecordResultView(allData), int64(len(allData)), nil)
		return
	} else if param.Type == RmPrimary {
		db.Delete([]byte(PrimaryKey))
		db.Delete([]byte("id"))
		common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
		return
	} else if param.Type == Upgrade {
		rs, err := db.IteratorWithPrefix([]byte(KeyPrefix))
		if err != nil {
			log.Errorf("get data use prefix failed. ")
			common.ConvergeFailureResponse(ctx, common.InternalError)
			return
		}

		batch := &leveldb.Batch{}
		maxId := len(rs)
		db.Put([]byte("id"), common.ConvertIntToByte(uint64(maxId)))
		for _, v := range rs {
			var r RecordMeta
			json.Unmarshal([]byte(v), &r)
			r.Id = maxId
			b2, _ := json.Marshal(r)
			batch.Put(getKeyWithPrefix(r.Id), b2)
			batch.Delete([]byte(fmt.Sprintf("%v%v", KeyPrefix, r.Time)))
			maxId--
		}
		db.BatchPut(batch)
		common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
	}

}

func getKeyWithPrefix(id int) []byte {
	b := common.ConvertIntToByte(uint64(id))
	return append([]byte(KeyPrefix), b...)
}
