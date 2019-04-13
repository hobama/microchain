package core

import (
	"bytes"
	"fmt"
	"time"
)

var (
	TransactionIDBufferSize     = 32
	TimestampBufferSize         = 4
	PublicKeyBufferSize         = 64
	SignatureBufferSize         = 32
	MetaDataLenBufferSize       = 8
	PublicKeyHashBufferSize     = 32
	TXOutputBufferSize          = 48
	TXOutputLenBufferSize       = 8
	TransactionHeaderBufferSize =
	    2*TransactionIDBufferSize + // TransactionID + PrevTransactionID
		TimestampBufferSize       + // Timestamp
		2*PublicKeyBufferSize     + // RequesterPublicKey + RequesteePublicKey
		2*SignatureBufferSize     + // RequesterSignature + RequesteeSignature
		MetaDataLenBufferSize     + // MetaDataLen
		TXOutputLenBufferSize       // TXOutputLen
)

// The output field contains the following 3 entries:
// (1) The total number of transactions generated by the requester that have been accepted by the requestee.
// (2) The total number of transactions rejected by the requestee.
// (3) The hash of the public key that the requester will use for its next transaction.
type TXOutput struct {
	Accepted          uint64
	Rejected          uint64
	NextPublicKeyHash []byte
}
type TXOutputs []TXOutput

type TransactionHeader struct {
	TransactionID      []byte // SHA256(requesterPK, requesteePK, timestamp) : 256-bits : 32-bytes
	Timestamp          uint32 // Unix timestamp                              : 32-bits  : 4-byte
	PrevTransactionID  []byte // Previous transaction ID                     : 256-bits : 32-bytes
	RequesterPublicKey []byte // Base58 encoding of requester public key     : 512-bits : 64-bytes
	RequesterSignature []byte // Base58 encoding of requester signature      : 256-bits : 32-bytes
	RequesteePublicKey []byte // Base58 encoding of requestee public key     : 512-bits : 64-bytes
	RequesteeSignature []byte // Base58 encoding of requestee signature      : 256-bits : 32-bytes
	MetaLength         uint64 // Meta data length                            : 64-bits  : 8-bytes
	OutputLength       uint64 // TXOutput length                             : 64-bits  : 8-bytes
}

type Transaction struct {
	Header TransactionHeader // Header
	Meta   []byte            // Meta data field
	Output TXOutputs         // TXOutput        // TXOutput
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

// Test if two TXOutputs are equal.
func (txo TXOutput) EqualWith(temp TXOutput) bool {
	if txo.Accepted != temp.Accepted {
		return false
	}

	if txo.Rejected != temp.Rejected {
		return false
	}

	if !bytes.Equal(StripBytes(txo.NextPublicKeyHash, 0), StripBytes(temp.NextPublicKeyHash, 0)) {
		return false
	}

	return true
}

// Serialize TXOutput into bytes.
func (txo TXOutput) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(UInt64ToBytes(txo.Accepted))
	buf.Write(UInt64ToBytes(txo.Rejected))
	buf.Write(FitBytesIntoSpecificWidth(txo.NextPublicKeyHash, PublicKeyHashBufferSize))

	return buf.Bytes(), nil
}

// Read TXOutput from bytes.
func (txo *TXOutput) UnmarshalBinary(data []byte) error {
	if len(data) != TXOutputBufferSize {
		return fmt.Errorf("Invalid TXOutput")
	}

	buf := bytes.NewBuffer(data)

	acc, err := BytesToUInt64(buf.Next(8))
	if err != nil {
		return err
	}

	txo.Accepted = acc

	rej, err := BytesToUInt64(buf.Next(8))
	if err != nil {
		return err
	}

	txo.Rejected = rej

	txo.NextPublicKeyHash = buf.Bytes()

	return nil
}

// Test if two TXOutputLists are equal
func (txos TXOutputs) EqualWith(temp TXOutputs) bool {
	if len(txos) != len(temp) {
		return false
	}

	for i, txo := range txos {
		if !txo.EqualWith(temp[i]) {
			return false
		}
	}

	return true
}

// Append new TXOutput to XOutputs.
func (txos TXOutputs) Append(txo TXOutput) TXOutputs {
	return append(txos, txo)
}

// Serialize TXOutputs into bytes.
func (txos TXOutputs) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, txo := range txos {
		b, err := txo.MarshalBinary()
		if err != nil {
			return nil, fmt.Errorf("Cannot serialize TXOutputs into bytes.")
		}

		buf.Write(b)
	}

	return buf.Bytes(), nil
}

// Read TXOutputs from bytes.
func (txos *TXOutputs) UnmarshalBinary(data []byte) error {
	if len(data) % TXOutputBufferSize != 0 {
		return fmt.Errorf("Invalid TXOutputs.")
	}

	buf := bytes.NewBuffer(data)

	for buf.Len() != 0 {
		txo := new(TXOutput)

		err := txo.UnmarshalBinary(buf.Next(int(TXOutputBufferSize)))
		if err != nil {
			return err
		}

		*txos = txos.Append(*txo)
	}

	return nil
}

// Test if two transaction headers are equal.
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

	if h.OutputLength != temp.OutputLength {
		return false
	}

	return true
}

// Serialize transaction header into bytes.
func (h TransactionHeader) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	buf.Write(FitBytesIntoSpecificWidth(h.TransactionID, TransactionIDBufferSize))
	buf.Write(UInt32ToBytes(h.Timestamp))
	buf.Write(FitBytesIntoSpecificWidth(h.PrevTransactionID, TransactionIDBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesterPublicKey, PublicKeyBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesterSignature, SignatureBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesteePublicKey, PublicKeyBufferSize))
	buf.Write(FitBytesIntoSpecificWidth(h.RequesteeSignature, SignatureBufferSize))
	buf.Write(UInt64ToBytes(h.MetaLength))
	buf.Write(UInt64ToBytes(h.OutputLength))

	return buf.Bytes(), nil
}

// Read tansaction header from bytes.
func (h *TransactionHeader) UnmarshalBinary(data []byte) error {
	if len(data) != TransactionHeaderBufferSize {
		return fmt.Errorf("Invalid transaction header")
	}

	buf := bytes.NewBuffer(data)

	h.TransactionID = StripBytes(buf.Next(TransactionIDBufferSize), 0)
	timestamp, err := BytesToUInt32(buf.Next(TimestampBufferSize))
	if err != nil {
		return err
	}

	h.Timestamp = timestamp
	h.PrevTransactionID = StripBytes(buf.Next(TransactionIDBufferSize), 0)
	h.RequesterPublicKey = StripBytes(buf.Next(PublicKeyBufferSize), 0)
	h.RequesterSignature = StripBytes(buf.Next(SignatureBufferSize), 0)
	h.RequesteePublicKey = StripBytes(buf.Next(PublicKeyBufferSize), 0)
	h.RequesteeSignature = StripBytes(buf.Next(SignatureBufferSize), 0)

	metalen, err := BytesToUInt64(buf.Next(MetaDataLenBufferSize))
	if err != nil {
		return err
	}

	h.MetaLength = metalen

	outputlen, err := BytesToUInt64(buf.Next(TXOutputLenBufferSize))
	if err != nil {
		return err
	}

	h.OutputLength = outputlen

	return nil
}

// Test if two transactions are equal.
func (t Transaction) EqualWith(temp Transaction) bool {
	if !t.Header.EqualWith(temp.Header) {
		return false
	}

	if !bytes.Equal(t.Meta, temp.Meta) {
		return false
	}

	return true
}

// Serialize transaction into bytes.
func (t Transaction) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	h, err := t.Header.MarshalBinary()
	if err != nil {
		return nil, err
	}

	buf.Write(h)
	buf.Write(t.Meta)

	return buf.Bytes(), nil
}

// Read tansaction from bytes.
func (t *Transaction) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)
	if err := t.Header.UnmarshalBinary(buf.Next(TransactionHeaderBufferSize)); err != nil {
		return err
	}

	t.Meta = buf.Bytes()

	return nil
}

// Calculate SHA256 sum of transaction.
func (t *Transaction) Hash() []byte {
	h, _ := t.Header.MarshalBinary()
	return SHA256(h)
}

// We have 2 ways to represent transactions in a block
// (1) TransactionsList: Consume lower memory, but low performace
// (2) TransactionsMap: High performance, but consume more memory
type TransactionsList []Transaction
type TransactionsMap struct {
	mapping map[string]Transaction
	order   []string
}

// Test if given tansaction is contained in the list.
func (list TransactionsList) Contains(tr Transaction) (bool, int) {
	for i, t := range list {
		if bytes.Equal(tr.Header.TransactionID, t.Header.TransactionID) {
			return true, i
		}
	}

	return false, 0
}

// Test if transaction with given id is contained in the list.
func (list TransactionsList) ContainsByID(id []byte) (bool, int) {
	for i, t := range list {
		if bytes.Equal(id, t.Header.TransactionID) {
			return true, i
		}
	}

	return false, 0
}

// Append new transaction to transaction list.
func (list TransactionsList) Append(tr Transaction) TransactionsList {
	return append(list, tr)
}

// Insert new transaction to transaction list.
func (list TransactionsList) Insert(tr Transaction) TransactionsList {
	for i, t := range list {
		if t.Header.Timestamp >= tr.Header.Timestamp {
			return append(append(list[:i], tr), list[i:]...)
		}
	}

	return list.Append(tr)
}

// Serialize transactions into bytes.
func (list TransactionsList) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)

	for _, t := range list {
		trBytes, err := t.MarshalBinary()
		if err != nil {
			return nil, err
		}
		buf.Write(trBytes)
	}

	return buf.Bytes(), nil
}

// Read tansaction header from bytes.
func (list *TransactionsList) UnmarshalBinary(data []byte) error {
	buf := bytes.NewBuffer(data)

	for buf.Len() != 0 {
		t := new(Transaction)
		err := t.Header.UnmarshalBinary(buf.Next(int(TransactionHeaderBufferSize)))
		if err != nil {
			return err
		}

		t.Meta = buf.Next(int(t.Header.MetaLength))
		*list = list.Append(*t)
	}

	return nil
}

// Test if two transaction lists are equal.
func (list TransactionsList) EqualWith(temp TransactionsList) bool {
	if len(list) != len(temp) {
		return false
	}

	for i, t := range list {
		if !t.EqualWith(temp[i]) {
			return false
		}
	}

	return true
}
