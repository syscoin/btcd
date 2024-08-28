// Copyright (c) 2013-2016 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package wire

import (
	"io"
	"encoding/binary"
)
const (
	MAX_GUID_LENGTH = 20
	MAX_RLP_SIZE = 4096
	HASH_SIZE = 32
	MAX_NEVM_BLOCK_SIZE = 33554432 // 32 MB
)

type MintSyscoinType struct {
	TxHash []byte
	BlockHash []byte
    TxPos uint16
    TxParentNodes []byte
    TxPath []byte
    TxRoot []byte
    ReceiptRoot []byte
    ReceiptPos uint16
    ReceiptParentNodes []byte
}

type SyscoinBurnToEthereumType struct {
	EthAddress []byte `json:"ethAddress,omitempty"`
}

// NEVMAddressEntry represents an entry with an address and collateralHeight.
type NEVMAddressEntry struct {
	Address          []byte
	CollateralHeight uint32
}
// NEVMAddressUpdateEntry represents an update entry with old and new addresses.
type NEVMAddressUpdateEntry struct {
	OldAddress []byte
	NewAddress []byte
}
// NEVMRemoveEntry represents an entry with an address to be removed.
type NEVMRemoveEntry struct {
	Address []byte
}
// NEVMAddressDiff holds the differences in the NEVM address list for masternodes.
type NEVMAddressDiff struct {
	AddedMNNEVM   []NEVMAddressEntry
	UpdatedMNNEVM []NEVMAddressUpdateEntry
	RemovedMNNEVM []NEVMRemoveEntry
}

type NEVMBlockWire struct {
	NEVMBlockHash []byte
	TxRoot        []byte
	ReceiptRoot   []byte
	NEVMBlockData []byte
	SYSBlockHash  []byte
	VersionHashes [][]byte
	Diff          NEVMAddressDiff
}

type NEVMDisconnectBlockWire struct {
	SYSBlockHash  []byte
	Diff          NEVMAddressDiff
}

func (a *NEVMAddressEntry) Deserialize(r io.Reader) error {
	var err error
	a.Address, err = ReadVarBytes(r, 0, HASH_SIZE, "Address")
	if err != nil {
		return err
	}
	a.CollateralHeight, err = binarySerializer.Uint32(r, littleEndian)
	if err != nil {
		return err
	}
	return nil
}

func (a *NEVMAddressEntry) Serialize(w io.Writer) error {
	err := WriteVarBytes(w, 0, a.Address)
	if err != nil {
		return err
	}
	err = binarySerializer.PutUint32(w, littleEndian, a.CollateralHeight)
	if err != nil {
		return err
	}
	return nil
}

func (a *NEVMAddressUpdateEntry) Deserialize(r io.Reader) error {
	var err error
	a.OldAddress, err = ReadVarBytes(r, 0, HASH_SIZE, "OldAddress")
	if err != nil {
		return err
	}
	a.NewAddress, err = ReadVarBytes(r, 0, HASH_SIZE, "NewAddress")
	if err != nil {
		return err
	}
	return nil
}

func (a *NEVMAddressUpdateEntry) Serialize(w io.Writer) error {
	err := WriteVarBytes(w, 0, a.OldAddress)
	if err != nil {
		return err
	}
	err = WriteVarBytes(w, 0, a.NewAddress)
	if err != nil {
		return err
	}
	return nil
}

func (a *NEVMRemoveEntry) Deserialize(r io.Reader) error {
	var err error
	a.Address, err = ReadVarBytes(r, 0, HASH_SIZE, "Address")
	if err != nil {
		return err
	}
	return nil
}

func (a *NEVMRemoveEntry) Serialize(w io.Writer) error {
	err := WriteVarBytes(w, 0, a.Address)
	if err != nil {
		return err
	}
	return nil
}

func (d *NEVMAddressDiff) Deserialize(r io.Reader) error {
	var err error

	// Deserialize AddedMNNEVM
	numAdded, err := ReadVarInt(r, 0)
	if err != nil {
		return err
	}
	d.AddedMNNEVM = make([]NEVMAddressEntry, numAdded)
	for i := range d.AddedMNNEVM {
		err = d.AddedMNNEVM[i].Deserialize(r)
		if err != nil {
			return err
		}
	}

	// Deserialize UpdatedMNNEVM
	numUpdated, err := ReadVarInt(r, 0)
	if err != nil {
		return err
	}
	d.UpdatedMNNEVM = make([]NEVMAddressUpdateEntry, numUpdated)
	for i := range d.UpdatedMNNEVM {
		err = d.UpdatedMNNEVM[i].Deserialize(r)
		if err != nil {
			return err
		}
	}

	// Deserialize RemovedMNNEVM
	numRemoved, err := ReadVarInt(r, 0)
	if err != nil {
		return err
	}
	d.RemovedMNNEVM = make([]NEVMRemoveEntry, numRemoved)
	for i := range d.RemovedMNNEVM {
		err = d.RemovedMNNEVM[i].Deserialize(r)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *NEVMAddressDiff) Serialize(w io.Writer) error {
	var err error

	// Serialize AddedMNNEVM
	err = WriteVarInt(w, 0, uint64(len(d.AddedMNNEVM)))
	if err != nil {
		return err
	}
	for i := range d.AddedMNNEVM {
		err = d.AddedMNNEVM[i].Serialize(w)
		if err != nil {
			return err
		}
	}

	// Serialize UpdatedMNNEVM
	err = WriteVarInt(w, 0, uint64(len(d.UpdatedMNNEVM)))
	if err != nil {
		return err
	}
	for i := range d.UpdatedMNNEVM {
		err = d.UpdatedMNNEVM[i].Serialize(w)
		if err != nil {
			return err
		}
	}

	// Serialize RemovedMNNEVM
	err = WriteVarInt(w, 0, uint64(len(d.RemovedMNNEVM)))
	if err != nil {
		return err
	}
	for i := range d.RemovedMNNEVM {
		err = d.RemovedMNNEVM[i].Serialize(w)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *NEVMBlockWire) Deserialize(r io.Reader) error {
	var err error

	// Deserialize NEVMBlockHash
	a.NEVMBlockHash = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.NEVMBlockHash)
	if err != nil {
		return err
	}

	// Deserialize TxRoot
	a.TxRoot = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.TxRoot)
	if err != nil {
		return err
	}

	// Deserialize ReceiptRoot
	a.ReceiptRoot = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.ReceiptRoot)
	if err != nil {
		return err
	}

	// Deserialize NEVMBlockData
	a.NEVMBlockData, err = ReadVarBytes(r, 0, MAX_NEVM_BLOCK_SIZE, "NEVMBlockData")
	if err != nil {
		return err
	}

	// Deserialize SYSBlockHash
	a.SYSBlockHash = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.SYSBlockHash)
	if err != nil {
		return err
	}

	// Deserialize VersionHashes
	numVH, err := ReadVarInt(r, 0)
	if err != nil {
		return err
	}
	a.VersionHashes = make([][]byte, numVH)
	for i := 0; i < int(numVH); i++ {
		a.VersionHashes[i], err = ReadVarBytes(r, 0, HASH_SIZE, "VersionHash")
		if err != nil {
			return err
		}
	}

	// Deserialize Diff
	a.Diff = NEVMAddressDiff{}
	err = a.Diff.Deserialize(r)
	if err != nil {
		return err
	}

	return nil
}


func (a *NEVMDisconnectBlockWire) Deserialize(r io.Reader) error {
	var err error

	// Deserialize SYSBlockHash
	a.SYSBlockHash = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.SYSBlockHash)
	if err != nil {
		return err
	}

	// Deserialize Diff
	a.Diff = NEVMAddressDiff{}
	err = a.Diff.Deserialize(r)
	if err != nil {
		return err
	}

	return nil
}

func (a *NEVMBlockWire) Serialize(w io.Writer) error {
	var err error
	_, err = w.Write(a.NEVMBlockHash[:])
	if err != nil {
		return err
	}
	_, err = w.Write(a.TxRoot[:])
	if err != nil {
		return err
	}
	_, err = w.Write(a.ReceiptRoot[:])
	if err != nil {
		return err
	}
	err = WriteVarBytes(w, 0, a.NEVMBlockData)
	if err != nil {
		return err
	}
	return nil
}


func (a *MintSyscoinType) Deserialize(r io.Reader) error {
	a.TxHash = make([]byte, HASH_SIZE)
	_, err := io.ReadFull(r, a.TxHash)
	if err != nil {
		return err
	}
	a.BlockHash = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.BlockHash)
	if err != nil {
		return err
	}
	a.TxPos, err = binarySerializer.Uint16(r, binary.LittleEndian)
	if err != nil {
		return err
	}
	a.TxParentNodes, err = ReadVarBytes(r, 0, MAX_RLP_SIZE, "TxParentNodes")
	if err != nil {
		return err
	}
	a.TxPath, err = ReadVarBytes(r, 0, MAX_RLP_SIZE, "TxPath")
	if err != nil {
		return err
	}
	a.ReceiptPos, err = binarySerializer.Uint16(r, binary.LittleEndian)
	if err != nil {
		return err
	}
	a.ReceiptParentNodes, err = ReadVarBytes(r, 0, MAX_RLP_SIZE, "ReceiptParentNodes")
	if err != nil {
		return err
	}
	a.TxRoot = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.TxRoot)
	if err != nil {
		return err
	}
	a.ReceiptRoot = make([]byte, HASH_SIZE)
	_, err = io.ReadFull(r, a.ReceiptRoot)
	if err != nil {
		return err
	}
	return nil
}

func (a *MintSyscoinType) Serialize(w io.Writer) error {
	_, err := w.Write(a.TxHash[:])
	if err != nil {
		return err
	}
	_, err = w.Write(a.BlockHash[:])
	if err != nil {
		return err
	}
	err = binarySerializer.PutUint16(w, binary.LittleEndian, a.TxPos)
	if err != nil {
		return err
	}
	err = WriteVarBytes(w, 0, a.TxParentNodes)
	if err != nil {
		return err
	}
	err = WriteVarBytes(w, 0, a.TxPath)
	if err != nil {
		return err
	}
	err = binarySerializer.PutUint16(w, binary.LittleEndian, a.ReceiptPos)
	if err != nil {
		return err
	}
	err = WriteVarBytes(w, 0, a.ReceiptParentNodes)
	if err != nil {
		return err
	}
	_, err = w.Write(a.TxRoot[:])
	if err != nil {
		return err
	}
	_, err = w.Write(a.ReceiptRoot[:])
	if err != nil {
		return err
	}
	return nil
}

func (a *SyscoinBurnToEthereumType) Deserialize(r io.Reader) error {
	var err error
	a.EthAddress, err = ReadVarBytes(r, 0, MAX_GUID_LENGTH, "ethAddress")
	if err != nil {
		return err
	}
	return nil
}

func (a *SyscoinBurnToEthereumType) Serialize(w io.Writer) error {
	err := WriteVarBytes(w, 0, a.EthAddress)
	if err != nil {
		return err
	}
	return nil
}