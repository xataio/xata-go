// This file was auto-generated by Fern from our API Definition.

package api

type UpsertRecordWithIdResponse map[string]interface{}

//type UpsertRecordWithIdResponse struct {
//	typeName                     string
//	Record                       *Record
//	UpsertRecordWithIdResponseId *UpsertRecordWithIdResponseId
//}
//
//func NewUpsertRecordWithIdResponseFromRecord(value *Record) *UpsertRecordWithIdResponse {
//	return &UpsertRecordWithIdResponse{typeName: "record", Record: value}
//}
//
//func NewUpsertRecordWithIdResponseFromUpsertRecordWithIdResponseId(value *UpsertRecordWithIdResponseId) *UpsertRecordWithIdResponse {
//	return &UpsertRecordWithIdResponse{typeName: "upsertRecordWithIdResponseId", UpsertRecordWithIdResponseId: value}
//}
//
//func (u *UpsertRecordWithIdResponse) UnmarshalJSON(data []byte) error {
//	valueRecord := new(Record)
//	if err := json.Unmarshal(data, &valueRecord); err == nil {
//		u.typeName = "record"
//		u.Record = valueRecord
//		return nil
//	}
//	valueUpsertRecordWithIdResponseId := new(UpsertRecordWithIdResponseId)
//	if err := json.Unmarshal(data, &valueUpsertRecordWithIdResponseId); err == nil {
//		u.typeName = "upsertRecordWithIdResponseId"
//		u.UpsertRecordWithIdResponseId = valueUpsertRecordWithIdResponseId
//		return nil
//	}
//	return fmt.Errorf("%s cannot be deserialized as a %T", data, u)
//}
//
//func (u UpsertRecordWithIdResponse) MarshalJSON() ([]byte, error) {
//	switch u.typeName {
//	default:
//		return nil, fmt.Errorf("invalid type %s in %T", u.typeName, u)
//	case "record":
//		return json.Marshal(u.Record)
//	case "upsertRecordWithIdResponseId":
//		return json.Marshal(u.UpsertRecordWithIdResponseId)
//	}
//}
//
//type UpsertRecordWithIdResponseVisitor interface {
//	VisitRecord(*Record) error
//	VisitUpsertRecordWithIdResponseId(*UpsertRecordWithIdResponseId) error
//}
//
//func (u *UpsertRecordWithIdResponse) Accept(v UpsertRecordWithIdResponseVisitor) error {
//	switch u.typeName {
//	default:
//		return fmt.Errorf("invalid type %s in %T", u.typeName, u)
//	case "record":
//		return v.VisitRecord(u.Record)
//	case "upsertRecordWithIdResponseId":
//		return v.VisitUpsertRecordWithIdResponseId(u.UpsertRecordWithIdResponseId)
//	}
//}
