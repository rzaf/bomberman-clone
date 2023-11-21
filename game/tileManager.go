package game

import (
	"fmt"
	"github.com/rzaf/bomberman-clone/core"
	"log"
	"os"
	"sync"

	ray "github.com/gen2brain/raylib-go/raylib"
)

const (
	TILE_LENGTH = 50
)

type GameMap struct {
	sync.Mutex
	tiles          []TileInterface
	extraTiles     []TileInterface
	XCount, YCount int
	data           []byte
	P1Pos, P2Pos   ray.Vector2
	Name           string
	EditorMode     bool
}

func CopyGameMap(dst *GameMap, src *GameMap) {
	dst.tiles = make([]TileInterface, len(src.tiles))
	copy(dst.tiles, src.tiles)
	dst.extraTiles = make([]TileInterface, len(src.extraTiles))
	copy(dst.extraTiles, src.extraTiles)
	dst.XCount = src.XCount
	dst.YCount = src.YCount
	dst.P1Pos = src.P1Pos
	dst.P2Pos = src.P2Pos
	dst.Name = src.Name
	dst.EditorMode = src.EditorMode
}

type tileManager struct {
	*GameMap
}

var (
	TileManager tileManager
)

func (tm *tileManager) Change(m *GameMap) *tileManager {
	lastX := tm.XCount
	lastY := tm.YCount
	lastTiles := tm.tiles

	m.tiles = make([]TileInterface, m.XCount*m.YCount)
	tm.GameMap = m
	for i := 0; i < m.XCount; i++ {
		for j := 0; j < m.YCount; j++ {
			if i == 0 || j == 0 || i == m.XCount-1 || j == m.YCount-1 {
				tm.Add(NewWallTile(i, j, false))
			} else {
				if i < lastX-1 && j < lastY-1 {
					tm.Add(lastTiles[j*lastX+i])
				} else {
					tm.Add(NewFloorTile(i, j))
				}
			}
		}
	}
	return tm
}

func (m *GameMap) Init() {
	m.tiles = make([]TileInterface, m.XCount*m.YCount)
	for i := 0; i < m.XCount; i++ {
		for j := 0; j < m.YCount; j++ {
			if i == 0 || j == 0 || i == m.XCount-1 || j == m.YCount-1 {
				m.Add(NewWallTile(i, j, false))
			} else {
				m.Remove(i, j)
				m.Add(NewFloorTile(i, j))
			}
		}
	}
}

func (m *GameMap) Length() (int, int) {
	return m.XCount, m.YCount
}

func (m *GameMap) Add(ti TileInterface) {
	t := m.TryLock()
	if t {
		defer m.Unlock()
	}
	_, isBomb := ti.(*BombTile)
	i, j := ti.Index()
	if isBomb {
		m.extraTiles = append(m.extraTiles, ti)
	} else {
		m.tiles[j*m.XCount+i] = ti
	}
}

func (m *GameMap) Remove(i, j int) {
	t := m.TryLock()
	if t {
		defer m.Unlock()
	}
	if m.tiles[j*m.XCount+i] != nil {
		switch t := m.tiles[j*m.XCount+i].(type) {
		case *WallTile:
			CollisionManager.Remove(t)
		case *Tile:
			CollisionManager.Remove(t)
		case *FloorTile:
			CollisionManager.Remove(t)
		}
	}
	m.tiles[j*m.XCount+i] = nil
}

func (m *GameMap) Get(i, j int) TileInterface {
	t := m.TryLock()
	if t {
		defer m.Unlock()
	}
	return m.tiles[j*m.XCount+i]
}

func (tm *tileManager) DrawInEditor() {
	tm.Lock()
	defer tm.Unlock()
	for _, tile := range TileManager.tiles {
		if tile != nil {
			var rec ray.Rectangle
			switch t := tile.(type) {
			case *WallTile:
				rec = ray.NewRectangle(float32(t.X), float32(t.Y), float32(t.Width), float32(t.Height))
			case *FloorTile:
				rec = ray.NewRectangle(float32(t.X), float32(t.Y), float32(t.Width), float32(t.Height))
			}

			tile.Draw()
			ray.DrawRectangleLinesEx(rec, 1, ray.White)
		}
	}
	if tm.P1Pos.X != 0 {
		core.GetTexture("anims").Crop(ray.NewRectangle(64, 0, 17, 25)).DrawAt(ray.NewRectangle(tm.P1Pos.X-TILE_LENGTH/2+8, tm.P1Pos.Y-TILE_LENGTH/2, 34, 50))
	}
	if tm.P2Pos.X != 0 {
		core.GetTexture("anims").Crop(ray.NewRectangle(167, 0, 17, 25)).DrawAt(ray.NewRectangle(tm.P2Pos.X-TILE_LENGTH/2+8, tm.P2Pos.Y-TILE_LENGTH/2, 34, 50))
	}
}

func (tm *tileManager) Draw() {
	t := tm.TryLock()
	if t {
		defer tm.Unlock()
	}
	for _, tile := range tm.tiles {
		if tile != nil {
			tile.Draw()
		}
	}
	for _, tile := range tm.extraTiles {
		if tile != nil {
			tile.Draw()
		}
	}
}

func (tm *tileManager) Update() {
	t := tm.TryLock()
	if t {
		defer tm.Unlock()
	}
	for _, tile := range tm.extraTiles {
		if tile != nil {
			switch t := tile.(type) {
			case *BombTile:
				t.Update()
			}
		}
	}
	for _, tile := range tm.tiles {
		if tile != nil {
			switch t := tile.(type) {
			case *WallTile:
				t.Update()
			case *UpgradeTile:
				t.Update()
			}
		}
	}
}

func (m *GameMap) Reload() *GameMap {
	data := m.data
	nWidth := 0
	nHeight := 1
	i := 0
	if data[len(data)-1] == '\n' {
		data[len(data)-1] = '0'
	}
	for _, b := range data {
		i++
		if b == '\n' {
			if nWidth != 0 && nWidth != i-1 {
				log.Fatalf("all lines should have same length! ")
			}
			nWidth = i - 1
			i = 0
			nHeight++
		}
	}
	m.tiles = make([]TileInterface, nWidth*nHeight)
	m.extraTiles = make([]TileInterface, 0, 10)
	m.XCount = nWidth
	m.YCount = nHeight
	m.P1Pos = ray.NewVector2(0, 0)
	m.P2Pos = ray.NewVector2(0, 0)
	i = 0
	j := 0
	for _, b := range data {
		switch b {
		case '#':
			// fmt.Println(i, j, "#")
			m.tiles[j*nWidth+i] = NewWallTile(i, j, false)
		case '*':
			// fmt.Println(i, j, "*")
			m.tiles[j*nWidth+i] = NewFloorTile(i, j)
		case '%':
			// fmt.Println(i, j, "%")
			m.tiles[j*nWidth+i] = NewWallTile(i, j, true)
		case '\n':
			// fmt.Println(i, j, "n")
			i = -1
			j++
		case '1':
			m.tiles[j*nWidth+i] = NewFloorTile(i, j)
			m.P1Pos = ray.NewVector2(float32(i)*TILE_LENGTH+TILE_LENGTH/2, float32(j)*TILE_LENGTH+TILE_LENGTH/2)
			// fmt.Println("1", i, j)
		case '2':
			m.tiles[j*nWidth+i] = NewFloorTile(i, j)
			m.P2Pos = ray.NewVector2(float32(i)*TILE_LENGTH+TILE_LENGTH/2, float32(j)*TILE_LENGTH+TILE_LENGTH/2)
			// fmt.Println("2", i, j)
		}
		i++
	}
	return m
}

func (m *GameMap) LoadFrom(data []byte, name string) *GameMap {
	m.Name = name
	m.EditorMode = false
	m.data = data
	m.Reload()
	return m
}

func (m *GameMap) LoadFromFile(path string) *GameMap {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("openning %s failed:%s", path, err.Error())
	}
	return m.LoadFrom(data, path)
}

func (m *GameMap) SaveInto(path string) {
	fmt.Printf("%d,%d\n", m.XCount*m.YCount, len(m.tiles))
	data := make([]byte, 0, (m.XCount+1)*(m.YCount))
	i2 := 0
	for i := 0; i < (m.XCount+1)*m.YCount; i++ {
		if i%(m.XCount+1) == m.XCount {
			data = append(data, '\n')
			fmt.Println()
		} else {
			switch t := m.tiles[i2].(type) {
			case *WallTile:
				if t.isDestoryable {
					fmt.Printf("%c", '%')
					data = append(data, '%')
				} else {
					fmt.Printf("%c", '#')
					data = append(data, '#')
				}
			case *FloorTile:
				fmt.Printf("%c", '*')
				data = append(data, '*')
			default:
				log.Fatalln("")
			}
			i2++
		}
	}

	if m.P1Pos.X != 0 {
		i := int(m.P1Pos.X / TILE_LENGTH)
		j := int(m.P1Pos.Y / TILE_LENGTH)
		data[j*m.XCount+i+j] = '1'
	}
	if m.P2Pos.X != 0 {
		i := int(m.P2Pos.X / TILE_LENGTH)
		j := int(m.P2Pos.Y / TILE_LENGTH)
		data[j*m.XCount+i+j] = '2'
	}
	err := os.WriteFile(path+m.Name, data, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
