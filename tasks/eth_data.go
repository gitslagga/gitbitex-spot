package tasks

import "math/big"

type RequestResult struct {
	ID      int         `json:"id"`
	Version string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *Error      `json:"error,omitempty"`
	Data    string      `json:"data,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type BlockCountResp struct {
	Error  *Error `json:"error"`
	Result string `json:"result"`
	ID     string `json:"id"`
}

type ContractDecimalsResp struct {
	Result string `json:"result"`
}

type BlockHashResp struct {
	Error  *Error `json:"error"`
	Result string `json:"result"`
	ID     string `json:"id"`
}

type BlockFalseResp struct {
	Result BlockTransactionsHash `json:"result"`
	ID     string                `json:"id"`
}

type BlockTransactionsHash struct {
	Transactions []string `json:"transactions"`
}

type BlockTrueResp struct {
	Result BlockTransactionsDetail `json:"result"`
	ID     string                  `json:"id"`
}

type BlockTransactionsDetail struct {
	Transactions []*Transaction `json:"transactions"`
}

type RowTransactionResp struct {
	Result Transaction `json:"result"`
	ID     string      `json:"id"`
}

type Transaction struct {
	BlockNumber string `json:"blockNumber"`
	FromAddress string `json:"from"`
	ToAddress   string `json:"to"`
	Value       string `json:"value"`
	Txid        string `json:"hash"`
	Input       string `json:"input"`
}

//查看transaction receipt
type RowTransactionReceipt struct {
	Result TransactionReceiptResult `json:"result"`
	ID     string                   `json:"id"`
}

type TransactionReceiptResult struct {
	Status string                  `json:"status"`
	Logs   []TransactionReceiptLog `json:"logs"`
}

type TransactionReceiptLog struct {
	Data   string   `json:"data" form:"data"`
	Topics []string `json:"topics" form:"topics"` //topics[0]: method, topics[1]: payer, topics[2]: payee
}

type RawTransactionResp struct {
	Result Transaction `json:"result"`
	ID     string      `json:"id"`
}

type RawTransactionReceipt struct {
	Result TransactionReceipt `json:"result"`
	ID     string             `json:"id"`
}

type TransactionReceipt struct {
	Status string `json:"status" form:"status"`
}

type SendTransactionLockUnlockResp struct {
	Error  *Error `json:"error"`
	Result bool   `json:"result"`
}

type SendTransactionResp struct {
	Error *Error `json:"error"`
	Txid  string `json:"result"`
}

type TokenBalanceResp struct {
	Error  *Error `json:"error"`
	Amount string `json:"result"`
}

// TransactionParameters GO transaction to make more easy controll the parameters
type TransactionParameters struct {
	From     string
	To       string
	Nonce    *big.Int
	Gas      *big.Int
	GasPrice *big.Int
	Value    *big.Int
	Data     string
}

// RequestTransactionParameters JSON
type RequestTransactionParameters struct {
	From     string `json:"from"`
	To       string `json:"to,omitempty"`
	Nonce    string `json:"nonce,omitempty"`
	Gas      string `json:"gas,omitempty"`
	GasPrice string `json:"gasPrice,omitempty"`
	Value    string `json:"value,omitempty"`
	Data     string `json:"data,omitempty"`
}

type NewAddressResp struct {
	Error  *Error `json:"error"`
	Result string `json:"result"`
	ID     string `json:"id"`
}

type EthTransactionCountResp struct {
	Error  string `json:"error" form:"error"`
	Result string `json:"result" form:"result"`
	ID     string `json:"id" form:"id"`
}

type SignTransactionResp struct {
	Error  *Error                     `json:"error" form:"error"`
	Result *SignTransactionResultItem `json:"result" form:"result"`
}
type SignTransactionResultItem struct {
	Raw string `json:"raw"`
}

type SendRawTransactionResp struct {
	Error *Error `json:"error" form:"error"`
	Txid  string `json:"result" form:"result"`
}
