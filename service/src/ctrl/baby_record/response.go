/*
 Copyright (C) THL A29 Limited, a Tencent company. All rights reserved.
   SPDX-License-Identifier: Apache-2.0
*/
package baby_record

type RecordResultView struct {
	Type   int   `json:"type"`
	Amount int   `json:"amount"`
	Time   int64 `json:"time"`
}

func NewRecordResultView(out []RecordMeta) []interface{} {

	var resultView []interface{}
	for _, v := range out {
		r := &RecordResultView{
			Type:   v.Type,
			Amount: v.Amount,
			Time:   v.Time,
		}
		resultView = append(resultView, r)
	}

	return resultView
}
