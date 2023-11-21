package running

import (
	"context"
	"fmt"
	"github.com/rzaf/bomberman-clone/game"
	"github.com/rzaf/bomberman-clone/pb"
	"io"
	"log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func DisconnectServer() {
	fmt.Println("HOST: disconnecting server")
	if serverGrpc != nil {
		serverGrpc.Stop()
	}
	if serverLis != nil {
		serverLis.Close()
	}
}

func CloseClient() {
	fmt.Println("GUEST: closing client")
	if clientGrpc != nil {
		clientGrpc.Close()
	}
}

type server struct {
	pb.GameServiceServer
}

func (s *server) GameInfo(req *pb.Empty, stream pb.GameService_GameInfoServer) error {
	fmt.Printf("request `%v` revieved \n", req)
	for {
		if isHostWaiting {
			path := mapNames[currentMapIndex].Text
			data, err := os.ReadFile("assets/maps/" + path)
			if err != nil {
				log.Fatalf("openning `%s` failed: %s", path, err.Error())
			}
			stream.Send(&pb.Info{
				// Map:    mapNames[currentMapIndex].Text,
				Map:    string(data),
				Rounds: int32(rounds),
				Time:   int32(roundTimeSeconds),
			})
			break
		}

	}
	isHostWaiting = false
	game.State.Change(game.ONLINE_BATTLE)

	p2.IsControllable = false

	return nil
}

func (s *server) PlayerInfo(stream pb.GameService_PlayerInfoServer) error {
	// recieving p2 info
	go func() {
		for {
			pi, err := stream.Recv()
			if err == io.EOF {
				break
			}
			// fmt.Printf("p2 stream recieved %v\n", pi)
			if pi != nil {
				p2.Lock()
				p2.Velocity.X = pi.Vel.X
				p2.Velocity.Y = pi.Vel.Y
				p2.Direction = game.Direction(pi.Direction)
				switch p2.Direction {
				case game.LEFT:
					p2.Animations.WalkingLeft.Index = int(pi.FrameIndex)
				case game.DOWN:
					p2.Animations.WalkingDown.Index = int(pi.FrameIndex)
				case game.RIGHT:
					p2.Animations.WalkingRight.Index = int(pi.FrameIndex)
				case game.UP:
					p2.Animations.WalkingUp.Index = int(pi.FrameIndex)
				}
				if pi.B.X != -1 {
					p2.PlaceBomb(int(pi.B.X), int(pi.B.Y))
				}
				p2.Unlock()
			}
		}
	}()

	// sending p1 info
	var frameI int32
	var bi, bj int
	for {
		if game.State.Get() == game.WIN {
			break
		}
		p1.Lock()
		switch p1.Direction {
		case game.LEFT:
			frameI = int32(p1.Animations.WalkingLeft.Index)
		case game.DOWN:
			frameI = int32(p1.Animations.WalkingDown.Index)
		case game.RIGHT:
			frameI = int32(p1.Animations.WalkingRight.Index)
		case game.UP:
			frameI = int32(p1.Animations.WalkingUp.Index)
		}
		bi = p1.LastBombX
		bj = p1.LastBombY
		if p1.LastBombX != -1 {
			p1.LastBombX = -1
			p1.LastBombY = -1
		}
		v := pb.Vec2{X: p1.Velocity.X, Y: p1.Velocity.Y}
		dir := int32(p1.Direction)
		p1.Unlock()

		var upgrades []*pb.Upgrade = nil
		game.SendingUpgradeLock.Lock()
		if game.SendingUpgrades != nil {
			for _, u := range game.SendingUpgrades {
				upgrades = append(upgrades, &pb.Upgrade{X: int32(u.I), Y: int32(u.J), Type: int32(u.Type)})
			}
			game.SendingUpgrades = nil
		}
		game.SendingUpgradeLock.Unlock()
		err := stream.Send(&pb.Player{
			Vel:        &v,
			Direction:  dir,
			FrameIndex: frameI,
			B:          &pb.Vec2{X: float32(bi), Y: float32(bj)},
			Upgrades:   upgrades,
		})
		if err != nil {
			fmt.Printf("error sending player1 Info(host) stream %v\n", err)
			game.State.Change(game.MENU)
			break
		}
		// if isWon {
		time.Sleep(15 * time.Millisecond)
	}
	return nil
}

var (
	clientGrpc *grpc.ClientConn
	serverGrpc *grpc.Server
	serverLis  net.Listener
)

const defaultPort = "50051"

func host() {
	serverLis, err := net.Listen("tcp4", "0.0.0.0:"+defaultPort)
	if err != nil {
		fmt.Printf("Failed to listen: %v", err)
		game.State.Change(game.MENU)
		return
	}
	serverGrpc = grpc.NewServer()
	pb.RegisterGameServiceServer(serverGrpc, &server{})
	fmt.Printf("server started\n")
	if err = serverGrpc.Serve(serverLis); err != nil {
		fmt.Printf("Failed to serve:%v", err)
		game.State.Change(game.MENU)
		return
	}
}

func connectToServer() {
	var err error
	addr := addrs[0] + "." + addrs[1] + "." + addrs[2] + "." + addrs[3]
	clientGrpc, err = grpc.Dial(addr+":"+defaultPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Printf("connection failed:%v", err)
		game.State.Change(game.MENU)
		return
	}
	fmt.Printf("connected to server \n")
}

func sendInfoReq() {
	c := pb.NewGameServiceClient(clientGrpc)

	res, err := c.GameInfo(context.Background(), &pb.Empty{})
	if err != nil {
		fmt.Printf("failed to create stream:%v", err)
		game.State.Change(game.MENU)
		return
	}
	for {
		info, err := res.Recv()
		if err == io.EOF {
			fmt.Printf("gameinfo stream ended:%v", err)
			game.State.Change(game.MENU)
			break
		}
		if err != nil {
			fmt.Printf("failed to recieve info:%v", err)
			game.State.Change(game.MENU)
			break
		}
		if info != nil {
			fmt.Printf("game info recieved: %v\n", info)
			// loadLevel(info.Map)
			game.TileManager.GameMap = &game.GameMap{}
			game.TileManager.GameMap.LoadFrom([]byte(info.Map), "host map")
			reset()
			p1.Wins = 0
			p2.Wins = 0

			rounds = int(info.Rounds)
			roundTimeSeconds = int(info.Time)
			isGuestWaiting = false
			game.State.Change(game.ONLINE_BATTLE)
			p1.IsControllable = false
			go sendPlayerReq()
			break
		}
	}
}

func sendPlayerReq() {
	c := pb.NewGameServiceClient(clientGrpc)

	stream, err := c.PlayerInfo(context.Background())
	if err != nil {
		fmt.Printf("failed to create stream:%v", err)
		game.State.Change(game.MENU)
		return
	}
	// recieving p1 info
	go func() {
		for {
			pi, err := stream.Recv()
			if err == io.EOF {
				break
			}
			// fmt.Printf("p1 stream recieved %v\n", pi)
			if pi != nil {
				p1.Lock()
				p1.Velocity.X = pi.Vel.X
				p1.Velocity.Y = pi.Vel.Y
				p1.Direction = game.Direction(pi.Direction)
				switch p1.Direction {
				case game.LEFT:
					p1.Animations.WalkingLeft.Index = int(pi.FrameIndex)
				case game.DOWN:
					p1.Animations.WalkingDown.Index = int(pi.FrameIndex)
				case game.RIGHT:
					p1.Animations.WalkingRight.Index = int(pi.FrameIndex)
				case game.UP:
					p1.Animations.WalkingUp.Index = int(pi.FrameIndex)
				}

				if pi.B.X != -1 {
					p1.PlaceBomb(int(pi.B.X), int(pi.B.Y))
				}
				p1.Unlock()
				if pi.Upgrades != nil {
					fmt.Println("adding host upgrades", pi.Upgrades)
					for _, u := range pi.Upgrades {
						i := int(u.X)
						j := int(u.Y)
						t := game.TileManager.Get(i, j)
						if w, isWall := t.(*game.WallTile); isWall {
							w.Remove()
						}
						game.AddUpgradeTile(i, j, game.UpgradeType(u.Type))
					}
				}
			}

		}
	}()

	// sending p2 info
	var frameI int32
	var bi, bj int
	for {
		if game.State.Get() == game.WIN {
			break
		}
		p2.Lock()
		switch p2.Direction {
		case game.LEFT:
			frameI = int32(p2.Animations.WalkingLeft.Index)
		case game.DOWN:
			frameI = int32(p2.Animations.WalkingDown.Index)
		case game.RIGHT:
			frameI = int32(p2.Animations.WalkingRight.Index)
		case game.UP:
			frameI = int32(p2.Animations.WalkingUp.Index)
		}
		bi = p2.LastBombX
		bj = p2.LastBombY
		if p2.LastBombX != -1 {
			p2.LastBombX = -1
			p2.LastBombY = -1
		}
		dir := int32(p2.Direction)
		v := pb.Vec2{X: p2.Velocity.X, Y: p2.Velocity.Y}
		p2.Unlock()
		err := stream.Send(&pb.Player{
			Vel:        &v,
			Direction:  dir,
			FrameIndex: frameI,
			B:          &pb.Vec2{X: float32(bi), Y: float32(bj)},
		})
		if err != nil {
			fmt.Printf("error sending player2 Info (guest) stream %v\n", err)
			game.State.Change(game.MENU)
			break
		}
		time.Sleep(15 * time.Millisecond)
	}
}
