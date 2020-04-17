package gode

const PathPrefix = "/amfphp/json.php"
const Path5145 = "/casino.slot.line243.BuBuGaoSheng."

type Flash2dbPhpGame struct {
	GameType GameType
}

func NewFlash2dbPhpGame(gameType GameType) *Flash2dbPhpGame {
	path := map[GameType]string{
		5145: Path5145,
	}

	_ = PathPrefix + path[5145]

	return &Flash2dbPhpGame{
		GameType: gameType,
	}
}

func (Flash2dbPhpGame) OnTakeMachine() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) OnLoadInfo() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) OnGetMachineDetail() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) BeginGame() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) OnCreditExchange() []byte {
	panic("implement me")
}

func (Flash2dbPhpGame) OnBalanceExchange() []byte {
	panic("implement me")
}
