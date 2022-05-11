/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package baby_record

import (
	"management_backend/src/ctrl/common"

	"github.com/gin-gonic/gin"
)

type RecordAdd struct {
	Type   int
	Amount int
}

type RecordDel struct {
	Ts int64 `json:"ts"`
}

type RecordMeta struct {
	Type   int   `json:"type"`
	Amount int   `json:"amount"`
	Time   int64 `json:"time"`
}

const Breastfeeding int = 1
const Shit int = 2
const KeyPrefix string = "baby_"

func (reqBody *RecordDel) IsLegal() bool {
	return true
}

func (reqBody *RecordAdd) IsLegal() bool {
	if reqBody.Type != Breastfeeding && reqBody.Type != Shit {
		return false
	}
	return true
}

func BindAddRecordHandler(ctx *gin.Context) *RecordAdd {
	var body = &RecordAdd{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindDeleteRecordHandler(ctx *gin.Context) *RecordDel {
	var body = &RecordDel{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
