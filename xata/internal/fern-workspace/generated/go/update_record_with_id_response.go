// This file was auto-generated by Fern from our API Definition.

package api

//import (
//	json "encoding/json"
//	fmt "fmt"
//)

type UpdateRecordWithIdResponse map[string]interface{}

//
//type UpdateRecordWithIdResponse struct {
//	typeName                     string
//	Record                       *Record
//	UpdateRecordWithIdResponseId *UpdateRecordWithIdResponseId
//}
//
//func NewUpdateRecordWithIdResponseFromRecord(value *Record) *UpdateRecordWithIdResponse {
//	return &UpdateRecordWithIdResponse{typeName: "record", Record: value}
//}
//
//func NewUpdateRecordWithIdResponseFromUpdateRecordWithIdResponseId(value *UpdateRecordWithIdResponseId) *UpdateRecordWithIdResponse {
//	return &UpdateRecordWithIdResponse{typeName: "updateRecordWithIdResponseId", UpdateRecordWithIdResponseId: value}
//}
//
//func (u *UpdateRecordWithIdResponse) UnmarshalJSON(data []byte) error {
//	valueRecord := new(Record)
//	if err := json.Unmarshal(data, &valueRecord); err == nil {
//		u.typeName = "record"
//		u.Record = valueRecord
//		return nil
//	}
//	valueUpdateRecordWithIdResponseId := new(UpdateRecordWithIdResponseId)
//	if err := json.Unmarshal(data, &valueUpdateRecordWithIdResponseId); err == nil {
//		u.typeName = "updateRecordWithIdResponseId"
//		u.UpdateRecordWithIdResponseId = valueUpdateRecordWithIdResponseId
//		return nil
//	}
//	return fmt.Errorf("%s cannot be deserialized as a %T", data, u)
//}
//
//func (u UpdateRecordWithIdResponse) MarshalJSON() ([]byte, error) {
//	switch u.typeName {
//	default:
//		return nil, fmt.Errorf("invalid type %s in %T", u.typeName, u)
//	case "record":
//		return json.Marshal(u.Record)
//	case "updateRecordWithIdResponseId":
//		return json.Marshal(u.UpdateRecordWithIdResponseId)
//	}
//}
//
//type UpdateRecordWithIdResponseVisitor interface {
//	VisitRecord(*Record) error
//	VisitUpdateRecordWithIdResponseId(*UpdateRecordWithIdResponseId) error
//}
//
//func (u *UpdateRecordWithIdResponse) Accept(v UpdateRecordWithIdResponseVisitor) error {
//	switch u.typeName {
//	default:
//		return fmt.Errorf("invalid type %s in %T", u.typeName, u)
//	case "record":
//		return v.VisitRecord(u.Record)
//	case "updateRecordWithIdResponseId":
//		return v.VisitUpdateRecordWithIdResponseId(u.UpdateRecordWithIdResponseId)
//	}
//}
