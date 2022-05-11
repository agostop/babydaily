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

	"management_backend/src/ctrl/common"
	leveldbhandle "management_backend/src/db/leveldb"

	"github.com/syndtr/goleveldb/leveldb"
)

func BabyRecordHandle(ctx *gin.Context) {
	params := BindAddRecordHandler(ctx)
	if params == nil || !params.IsLegal() {
		common.ConvergeFailureResponse(ctx, common.ErrorParamWrong)
		return
	}

	ts := time.Now().Unix()
	key := getKeyWithPrefix(ts)
	val, err := json.Marshal(&RecordMeta{Amount: params.Amount, Time: ts, Type: params.Type})
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorMarshalParameters)
		return
	}

	db := leveldbhandle.GetHandleInstance()
	batch := &leveldb.Batch{}
	batch.Put([]byte(key), []byte(val))

	err = db.BatchPut(batch)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.ErrorHandleFailure)
		return
	}

	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func BabyRecordListHandle(ctx *gin.Context) {
	db := leveldbhandle.GetHandleInstance()
	rs, err := db.IteratorWithPrefix(KeyPrefix)
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
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
	rd := BindDeleteRecordHandler(ctx)
	db := leveldbhandle.GetHandleInstance()
	err := db.Delete(getKeyWithPrefix(rd.Ts))
	if err != nil {
		common.ConvergeFailureResponse(ctx, common.InternalError)
		return
	}
	common.ConvergeDataResponse(ctx, common.NewStatusResponse(), nil)
}

func getKeyWithPrefix(ts int64) string {
	return fmt.Sprintf("%v%v", KeyPrefix, ts)
}
