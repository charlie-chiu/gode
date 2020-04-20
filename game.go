package gode

import "encoding/json"

type Game interface {
	// 對應 php 核心的API ，以5145步步高昇為例

	//machineOccupyAuto($_iUserID)
	OnTakeMachine(uid UserID) json.RawMessage

	//onLoadInfo($_iUserID, $_iGameCode)
	OnLoadInfo(uid UserID, gc GameCode) json.RawMessage

	//getMachineDetail($_iUserID, $_iGameCode)
	OnGetMachineDetail(uid UserID, gc GameCode) json.RawMessage

	//creditExchange($_sSid, $_iGameCode, $_sBetBase, $_iCredit)
	OnCreditExchange(sid SessionID, gc GameCode, bb string, credit int) json.RawMessage

	//balanceExchange($_iUserID, $_iHallID, $_iGameCode)
	OnBalanceExchange(uid UserID, hid HallID, gc GameCode) json.RawMessage

	//beginGame($_sSid, $_iGameCode, $_aBetInfo, $_iPlatform = null, $_iClient = 0, $_sOrderIp = null)
	BeginGame(sid SessionID, gc GameCode, betInfo string) json.RawMessage

	//machineLeave($_iUserID, $_iHallID, $_iGameCode)
	OnLeaveMachine(uid UserID, hid HallID, gameCode GameCode) json.RawMessage
}

type (
	SessionID string
	GameType  int
	GameCode  int
	UserID    int
	HallID    int
)
