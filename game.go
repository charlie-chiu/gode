package gode

type Game interface {
	OnReady() []byte
	OnLogin() []byte

	// ↓ 對應到 php 核心
	//machineOccupyAuto($_iUserID)
	OnTakeMachine() []byte

	//onLoadInfo($_iUserID, $_iGameCode)
	OnLoadInfo() []byte

	//getMachineDetail($_iUserID, $_iGameCode)
	OnGetMachineDetail() []byte

	//creditExchange($_sSid, $_iGameCode, $_sBetBase, $_iCredit)
	OnCreditExchange() []byte

	//balanceExchange($_iUserID, $_iHallID, $_iGameCode)
	OnBalanceExchange() []byte

	//beginGame($_sSid, $_iGameCode, $_aBetInfo, $_iPlatform = null, $_iClient = 0, $_sOrderIp = null)
	BeginGame() []byte
}

type (
	SessionID string
	GameType  int
	GameCode  int
	UserID    int
)
