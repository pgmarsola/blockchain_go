package blockchain

type TxOutput struct {
	Value     int
	PublicKey string
}

type TxInput struct {
	ID        []byte
	Output    int
	Signature string
}

func (in *TxInput) CanUnlock(data string) bool {
	return in.Signature == data
}

func (out *TxOutput) CanBeUnlocked(data string) bool {
	return out.PublicKey == data
}