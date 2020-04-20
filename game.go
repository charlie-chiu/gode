package gode

type Game interface {
	// 對應 php 核心的API

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
	HallID    int
)
