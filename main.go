package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	biliapi "github.com/iyear/biligo"
	live "github.com/iyear/biligo-live"
	"github.com/xmdhs/bilivesave/sql"
)

var (
	RoomID int64
)

func init() {
	flag.Int64Var(&RoomID, "room", 0, "room id")
	flag.Parse()
}

func main() {
	c := biliapi.NewCommClient(&biliapi.CommSetting{
		Client: &http.Client{
			Timeout: time.Second * 10,
		},
	})
	rd, err := c.LiveGetRoomInfoByID(RoomID)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	db, err := sql.NewDb(strconv.FormatInt(rd.RoomID, 10) + ".db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	for {
		do(ctx, rd.RoomID, db)
		log.Println("sleep 3s, reconnect")
		time.Sleep(3 * time.Second)
	}
}

func do(ctx context.Context, roomID int64, db *sql.Db) {
	l := live.NewLive(true, 30*time.Second, 50, func(err error) {
		log.Println("panic:", err)
	})

	if err := l.Conn(websocket.DefaultDialer, live.WsDefaultHost); err != nil {
		log.Fatal(err)
		return
	}

	ctx, stop := context.WithCancel(ctx)
	defer stop()

	go func() {
		if err := l.Enter(ctx, roomID, "", 12345678); err != nil {
			log.Println(err)
			stop()
			return
		}
	}()

	rev(ctx, l, db)
}

func rev(ctx context.Context, l *live.Live, db *sql.Db) {
	ch := make(chan struct{}, runtime.NumCPU())
	for {
		select {
		case tp := <-l.Rev:
			if tp.Error != nil {
				log.Println(tp.Error)
				continue
			}
			go func() {
				ch <- struct{}{}
				handle(ctx, tp.Msg, db)
				<-ch
			}()
		case <-ctx.Done():
			log.Println("rev func stopped")
			return
		}
	}
}

var logger = log.New(os.Stdout, "", log.LstdFlags)

func handle(ctx context.Context, msg live.Msg, db *sql.Db) {
	// 使用 msg.(type) 进行事件跳转和处理，常见事件基本都完成了解析(Parse)功能，不常见的功能有一些实在太难抓取
	// 更多注释和说明等待添加
	switch msg := msg.(type) {
	// 心跳回应直播间人气值
	case *live.MsgHeartbeatReply:
		hot := msg.GetHot()
		logger.Printf("hot: %d\n", hot)
		db.InsertHot(ctx, int64(hot), 0)
	// 弹幕消息
	case *live.MsgDanmaku:
		dm, err := msg.Parse()
		if err != nil {
			log.Println(err)
			return
		}
		logger.Printf("弹幕: %s (%d:%s) 【%s】| %d\n", dm.Content, dm.MID, dm.Uname, dm.MedalName, dm.Time)
		d := sql.Danmaku{
			SendMode:     dm.SendMode,
			SendFontSize: dm.SendFontSize,
			DanmakuColor: dm.DanmakuColor,
			Time:         dm.Time,
			Dmid:         dm.DMID,
			MsgType:      dm.MsgType,
			Bubble:       dm.Bubble,
			Content:      dm.Content,
			Mid:          dm.MID,
			Uname:        dm.Uname,
			RoomAdmin:    dm.RoomAdmin,
			Vip:          dm.Vip,
			Svip:         dm.SVip,
			Rank:         dm.Rank,
			MobileVerify: dm.MobileVerify,
			UnameColor:   dm.UnameColor,
			MedalName:    dm.MedalName,
			UpName:       dm.UpName,
			MedalLevel:   dm.MedalLevel,
			UserLevel:    dm.UserLevel,
		}
		err = db.InsertDanmaku(ctx, &d)
		if err != nil {
			log.Println(err)
		}

	// 礼物消息
	case *live.MsgSendGift:
		g, err := msg.Parse()
		if err != nil {
			log.Println(err)
			return
		}
		logger.Printf("%s: %s %d个%s\n", g.Action, g.Uname, g.Num, g.GiftName)
		gift := sql.Gift{
			Uname:     g.Uname,
			Uid:       g.UID,
			Rnd:       g.Rnd,
			GiftName:  g.GiftName,
			GiftNum:   g.Num,
			GiftPrice: g.Price,
			GiftID:    g.GiftID,
			Action:    g.Action,
			Time:      g.Timestamp,
			Num:       g.Num,
		}
		err = db.InsertGift(ctx, &gift)
		if err != nil {
			log.Println(err)
			return
		}
	case *live.MsgSuperChatMessage:
		sc, err := msg.Parse()
		if err != nil {
			log.Println(err)
			return
		}
		logger.Printf("sc: %s (%d:%s) 【%s】| %d\n", sc.Message, sc.UID, sc.UserInfo.Uname, sc.MedalInfo.MedalName, sc.Ts)

		scm := sql.Sc{
			ID:        sc.ID,
			Uname:     sc.UserInfo.Uname,
			Uid:       sc.UID,
			Time:      sc.Time,
			StartTime: sc.StartTime,
			Message:   sc.Message,
			Price:     sc.Price,
		}
		err = db.InsertSC(ctx, &scm)
		if err != nil {
			log.Println(err)
			return
		}

	case *live.MsgSuperChatMessageDelete:
		l, err := msg.GetList()
		if err != nil {
			log.Println(err)
			return
		}
		logger.Printf("sc delete: %v\n", l)
		for _, v := range l {
			err = db.InsertDelSC(ctx, v)
			if err != nil {
				log.Println(err)
			}
		}

	case *live.MsgWatChed:
		w, err := msg.Parse()
		if err != nil {
			log.Println(err)
			return
		}
		logger.Printf("看过: %d\n", w.Num)
		err = db.InsertHot(ctx, 0, int64(w.Num))
		if err != nil {
			log.Println(err)
		}
	// General 表示live未实现的CMD命令，请自行处理raw数据。也可以提issue更新这个CMD
	case *live.MsgGeneral:
		logger.Println("unknown msg type|raw:", string(msg.Raw()))
	default:
		logger.Printf("not case msg type %v|raw: %v\n", msg.Cmd(), string(msg.Raw()))
	}

}
