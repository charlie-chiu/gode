package gode

type Game interface {
	// 對應 php 核心的API ，以5145步步高昇為例

	//machineOccupyAuto($_iUserID)
	OnTakeMachine(uid UserID) []byte

	//onLoadInfo($_iUserID, $_iGameCode)
	OnLoadInfo(uid UserID, gc GameCode) []byte

	//getMachineDetail($_iUserID, $_iGameCode)
	OnGetMachineDetail(uid UserID, gc GameCode) []byte

	//creditExchange($_sSid, $_iGameCode, $_sBetBase, $_iCredit)
	OnCreditExchange(sid SessionID, gc GameCode, bb string, credit int) []byte

	//balanceExchange($_iUserID, $_iHallID, $_iGameCode)
	OnBalanceExchange(uid UserID, hid HallID, gc GameCode) []byte

	//beginGame($_sSid, $_iGameCode, $_aBetInfo, $_iPlatform = null, $_iClient = 0, $_sOrderIp = null)
	BeginGame(sid SessionID, gc GameCode, betInfo string) []byte

	//machineLeave($_iUserID, $_iHallID, $_iGameCode)
	OnLeaveMachine(uid UserID, hid HallID, gameCode GameCode) []byte
}

type (
	SessionID string
	GameType  int
	GameCode  int
	UserID    int
	HallID    int
)
