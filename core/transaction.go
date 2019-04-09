package core

import (
	"bytes"
	"fmt"
	"time"
)

const (
	TransactionIDLength     = 32
	TimestampLength         = 4
	PublicKeyLength         = 64
	SignatureLength         = 32
	MetaDataLength          = 8
	TransactionHeaderLength = 268
)

// The output field contains the following 3 entries:
// (1) The total number of transactions generated by the requester that have been accepted by the requestee.
// (2) The total number of transactions rejected by the requestee.
// (3) The hash of the public key that the requester will use for its next transaction.
type Output struct {
	Accepted          uint64
	Rejected          uint64
	NextPublicKeyHash []byte
}

type TransactionHeader struct {
	TransactionID      []byte // SHA256(requesterPK, requesteePK, timestamp) : 256-bits : 32-bytes
	Timestamp          uint32 // Unix timestamp                              : 32-bits  : 4-byte
	PrevTransactionID  []byte // Previous transaction ID                     : 256-bits : 32-bytes
	RequesterPublicKey []byte // Base58 encoding of requester public key     : 512-bits : 64-bytes
	RequesterSignature []byte // Base58 encoding of requester signature      : 256-bits : 32-bytes
	RequesteePublicKey []byte // Base58 encoding of requestee public key     : 512-bits : 64-bytes
	RequesteeSignature []byte // Base58 encoding of requestee signature      : 256-bits : 32-bytes
	MetaLength         uint64 // Meta data length                            : 64-bits  : 8-bytes
}

type Transaction struct {
	Header TransactionHeader // Header
	Meta   []byte            // Meta data field
}

// We use requesterID, requesteeID, timestamp to identify a transaction in blocks.
func NewTransaction(from, to, meta []byte) *Transaction {
	time := uint32(time.Now().Unix())
	timeBuf := UInt32ToBytes(time)
	rawid := JoinBytes(timeBuf, from, to)
	transaction := Transaction{
		Header: TransactionHeader{
			TransactionID:      SHA256(rawid),
			Timestamp:          time,
			RequesterPublicKey: from,
			RequesteePublicKey: to,
			MetaLength:         uint64(len(meta))},
		Meta: meta}
	return &transaction
}

func (h TransactionHeader) EqualWith(temp TransactionHeader) bool {
	if !bytes.Equal(StripBytes(h.TransactionID, 0), StripBytes(temp.TransactionID, 0)) {
		return false
	}

	if h.Timestamp != temp.Timestamp {
		return false
	}

	if !bytes.Equal(StripBytes(h.PrevTransactionID, 0), StripBytes(temp.PrevTransactionID, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(h.RequesterPublicKey, 0), StripBytes(temp.RequesterPublicKey, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(h.RequesterSignature, 0), StripBytes(temp.RequesterSignature, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(h.RequesteePublicKey, 0), StripBytes(temp.RequesteePublicKey, 0)) {
		return false
	}

	if !bytes.Equal(StripBytes(h.RequesteeSignature, 0), StripBytes(temp.RequesteeSignature, 0)) {
		return false
	}

	if h.MetaLength != temp.MetaLength {
		return false
	}

	return true
}

func (h *TransactionHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(FitBytesIntoSpecificWidth(h.TransactionID, TransactionIDLength))
	buf.Write(UInt32ToBytes(h.Timestamp))
	buf.Write(FitBytesIntoSpecificWidth(h.PrevTransactionID, TransactionIDLength))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesterPublicKey, PublicKeyLength))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesterSignature, SignatureLength))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesteePublicKey, PublicKeyLength))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesteeSignature, SignatureLength))
	buf.Write(UInt64ToBytes(h.MetaLength))
	return buf.Bytes(), nil
}

func (h *TransactionHeader) UnMarshalBinary(data []byte) error {
	if len(data) != TransactionHeaderLength {
		return fmt.Errorf("Invalid transaction header")
	}

	buf := bytes.NewBuffer(data)

	h.TransactionID = buf.Next(TransactionIDLength)
	timestamp, err := BytesToUInt32(buf.Next(TimestampLength))
	if err != nil {
		return err
	}

	h.Timestamp = timestamp
	h.PrevTransactionID = buf.Next(TransactionIDLength)
	h.RequesterPublicKey = buf.Next(PublicKeyLength)
	h.RequesterSignature = buf.Next(SignatureLength)
	h.RequesteePublicKey = buf.Next(PublicKeyLength)
	h.RequesteeSignature = buf.Next(SignatureLength)
	metalen, err := BytesToUInt64(buf.Next(MetaDataLength))
	if err != nil {
		return err
	}

	h.MetaLength = metalen
	return nil
}
