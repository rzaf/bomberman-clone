package game

import "sync"

// "github.com/rzaf/bomberman-clone/states/running"

// var State StateI

type StateI interface {
	Update()
	Draw()
	OnEnter()
	OnExit()
	OnWindowResized()
}

type GameState uint8

func (s *GameState) Change(newState GameState) {
	stateLocker.Lock()
	*s = newState
	stateLocker.Unlock()
}

func (s *GameState) Get() GameState {
	stateLocker.Lock()
	defer stateLocker.Unlock()
	return *s
}

var (
	LastState, State *GameState
	stateLocker      sync.Mutex
)

const (
	MENU GameState = iota
	OFFLINE_BATTLE
	ONLINE_BATTLE
	BATTLE_MENU
	ONLINE_MENU
	SETTING
	WIN
	PAUSED
	EDITOR
	QUIT
)
