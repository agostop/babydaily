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

type RecordOfId struct {
	Id int `json:"id"`
}

type RecordMeta struct {
	Id     int   `json:"id"`
	Type   int   `json:"type"`
	Amount int   `json:"amount"`
	Time   int64 `json:"time"`
}

type RecordMaintain struct {
	Type int
}

const CleanAll int = 1
const GetAllData int = 2
const RmPrimary int = 3
const Upgrade int = 4

const Breastfeeding int = 1
const Shit int = 2
const KeyPrefix string = "baby_"
const PrimaryKey string = "primarykey"

func (reqBody *RecordOfId) IsLegal() bool {
	return true
}

func (reqBody *RecordAdd) IsLegal() bool {
	if reqBody.Type != Breastfeeding && reqBody.Type != Shit {
		return false
	}
	return true
}

func (reqBody *RecordMaintain) IsLegal() bool {
	if reqBody.Type != CleanAll && reqBody.Type != GetAllData {
		return false
	}
	return true
}

func (reqBody *RecordMeta) IsLegal() bool {
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

func BindEditMetaRecordHandler(ctx *gin.Context) *RecordMeta {
	var body = &RecordMeta{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindTimeOfRecordHandler(ctx *gin.Context) *RecordOfId {
	var body = &RecordOfId{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}

func BindMaintainHandler(ctx *gin.Context) *RecordMaintain {
	var body = &RecordMaintain{}
	if err := common.BindBody(ctx, body); err != nil {
		return nil
	}
	return body
}
