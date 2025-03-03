// Code generated by fastssz. DO NOT EDIT.
// Hash: f12c072ef432812fbe3b6a0d5414463072a6b8969ae1b79cddd5d18a845ec0a7
// Version: 0.1.3
package types

import (
	ssz "github.com/ferranbt/fastssz"
)

// MarshalSSZ ssz marshals the CommitteeMember object
func (c *CommitteeMember) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(c)
}

// MarshalSSZTo ssz marshals the CommitteeMember object to a target array
func (c *CommitteeMember) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf
	offset := int(515)

	// Field (0) 'OperatorID'
	dst = ssz.MarshalUint64(dst, uint64(c.OperatorID))

	// Field (1) 'CommitteeID'
	dst = append(dst, c.CommitteeID[:]...)

	// Field (2) 'SSVOperatorPubKey'
	if size := len(c.SSVOperatorPubKey); size != 459 {
		err = ssz.ErrBytesLengthFn("CommitteeMember.SSVOperatorPubKey", size, 459)
		return
	}
	dst = append(dst, c.SSVOperatorPubKey...)

	// Field (3) 'FaultyNodes'
	dst = ssz.MarshalUint64(dst, c.FaultyNodes)

	// Offset (4) 'Committee'
	dst = ssz.WriteOffset(dst, offset)

	// Field (5) 'DomainType'
	dst = append(dst, c.DomainType[:]...)

	// Field (4) 'Committee'
	if size := len(c.Committee); size > 13 {
		err = ssz.ErrListTooBigFn("CommitteeMember.Committee", size, 13)
		return
	}
	for ii := 0; ii < len(c.Committee); ii++ {
		if dst, err = c.Committee[ii].MarshalSSZTo(dst); err != nil {
			return
		}
	}

	return
}

// UnmarshalSSZ ssz unmarshals the CommitteeMember object
func (c *CommitteeMember) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size < 515 {
		return ssz.ErrSize
	}

	tail := buf
	var o4 uint64

	// Field (0) 'OperatorID'
	c.OperatorID = OperatorID(ssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'CommitteeID'
	copy(c.CommitteeID[:], buf[8:40])

	// Field (2) 'SSVOperatorPubKey'
	if cap(c.SSVOperatorPubKey) == 0 {
		c.SSVOperatorPubKey = make([]byte, 0, len(buf[40:499]))
	}
	c.SSVOperatorPubKey = append(c.SSVOperatorPubKey, buf[40:499]...)

	// Field (3) 'FaultyNodes'
	c.FaultyNodes = ssz.UnmarshallUint64(buf[499:507])

	// Offset (4) 'Committee'
	if o4 = ssz.ReadOffset(buf[507:511]); o4 > size {
		return ssz.ErrOffset
	}

	if o4 != 515 {
		return ssz.ErrInvalidVariableOffset
	}

	// Field (5) 'DomainType'
	copy(c.DomainType[:], buf[511:515])

	// Field (4) 'Committee'
	{
		buf = tail[o4:]
		num, err := ssz.DivideInt2(len(buf), 467, 13)
		if err != nil {
			return err
		}
		c.Committee = make([]*Operator, num)
		for ii := 0; ii < num; ii++ {
			if c.Committee[ii] == nil {
				c.Committee[ii] = new(Operator)
			}
			if err = c.Committee[ii].UnmarshalSSZ(buf[ii*467 : (ii+1)*467]); err != nil {
				return err
			}
		}
	}
	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the CommitteeMember object
func (c *CommitteeMember) SizeSSZ() (size int) {
	size = 515

	// Field (4) 'Committee'
	size += len(c.Committee) * 467

	return
}

// HashTreeRoot ssz hashes the CommitteeMember object
func (c *CommitteeMember) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(c)
}

// HashTreeRootWith ssz hashes the CommitteeMember object with a hasher
func (c *CommitteeMember) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'OperatorID'
	hh.PutUint64(uint64(c.OperatorID))

	// Field (1) 'CommitteeID'
	hh.PutBytes(c.CommitteeID[:])

	// Field (2) 'SSVOperatorPubKey'
	if size := len(c.SSVOperatorPubKey); size != 459 {
		err = ssz.ErrBytesLengthFn("CommitteeMember.SSVOperatorPubKey", size, 459)
		return
	}
	hh.PutBytes(c.SSVOperatorPubKey)

	// Field (3) 'FaultyNodes'
	hh.PutUint64(c.FaultyNodes)

	// Field (4) 'Committee'
	{
		subIndx := hh.Index()
		num := uint64(len(c.Committee))
		if num > 13 {
			err = ssz.ErrIncorrectListSize
			return
		}
		for _, elem := range c.Committee {
			if err = elem.HashTreeRootWith(hh); err != nil {
				return
			}
		}
		hh.MerkleizeWithMixin(subIndx, num, 13)
	}

	// Field (5) 'DomainType'
	hh.PutBytes(c.DomainType[:])

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the CommitteeMember object
func (c *CommitteeMember) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(c)
}

// MarshalSSZ ssz marshals the Operator object
func (o *Operator) MarshalSSZ() ([]byte, error) {
	return ssz.MarshalSSZ(o)
}

// MarshalSSZTo ssz marshals the Operator object to a target array
func (o *Operator) MarshalSSZTo(buf []byte) (dst []byte, err error) {
	dst = buf

	// Field (0) 'OperatorID'
	dst = ssz.MarshalUint64(dst, uint64(o.OperatorID))

	// Field (1) 'SSVOperatorPubKey'
	if size := len(o.SSVOperatorPubKey); size != 459 {
		err = ssz.ErrBytesLengthFn("Operator.SSVOperatorPubKey", size, 459)
		return
	}
	dst = append(dst, o.SSVOperatorPubKey...)

	return
}

// UnmarshalSSZ ssz unmarshals the Operator object
func (o *Operator) UnmarshalSSZ(buf []byte) error {
	var err error
	size := uint64(len(buf))
	if size != 467 {
		return ssz.ErrSize
	}

	// Field (0) 'OperatorID'
	o.OperatorID = OperatorID(ssz.UnmarshallUint64(buf[0:8]))

	// Field (1) 'SSVOperatorPubKey'
	if cap(o.SSVOperatorPubKey) == 0 {
		o.SSVOperatorPubKey = make([]byte, 0, len(buf[8:467]))
	}
	o.SSVOperatorPubKey = append(o.SSVOperatorPubKey, buf[8:467]...)

	return err
}

// SizeSSZ returns the ssz encoded size in bytes for the Operator object
func (o *Operator) SizeSSZ() (size int) {
	size = 467
	return
}

// HashTreeRoot ssz hashes the Operator object
func (o *Operator) HashTreeRoot() ([32]byte, error) {
	return ssz.HashWithDefaultHasher(o)
}

// HashTreeRootWith ssz hashes the Operator object with a hasher
func (o *Operator) HashTreeRootWith(hh ssz.HashWalker) (err error) {
	indx := hh.Index()

	// Field (0) 'OperatorID'
	hh.PutUint64(uint64(o.OperatorID))

	// Field (1) 'SSVOperatorPubKey'
	if size := len(o.SSVOperatorPubKey); size != 459 {
		err = ssz.ErrBytesLengthFn("Operator.SSVOperatorPubKey", size, 459)
		return
	}
	hh.PutBytes(o.SSVOperatorPubKey)

	hh.Merkleize(indx)
	return
}

// GetTree ssz hashes the Operator object
func (o *Operator) GetTree() (*ssz.Node, error) {
	return ssz.ProofTree(o)
}
