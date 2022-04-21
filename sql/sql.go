package sql

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/mattn/go-sqlite3"
)

type Db struct {
	db *sqlx.DB
}

func NewDb(filename string) (*Db, error) {
	db, err := sqlx.Open("sqlite3", "./data.db?_txlock=IMMEDIATE&_journal_mode=WAL")
	if err != nil {
		return nil, fmt.Errorf("NewDb: %w", err)
	}
	d := new(Db)
	d.db = db

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS danmaku(
		danmaku_id INTEGER PRIMARY KEY AUTOINCREMENT,
		send_mode INT NOT NULL,
		send_font_size INT NOT NULL,
		danmaku_color INT NOT NULL,
		time INT NOT NULL,
		dmid INT NOT NULL,
		msg_type INT NOT NULL,
		bubble TEXT NOT NULL,
		content TEXT NOT NULL,
		mid INT NOT NULL,
		uname TEXT NOT NULL,
		room_admin INT NOT NULL,
		vip INT NOT NULL,
		svip INT NOT NULL,
		rank INT NOT NULL,
		mobile_verify INT NOT NULL,
		uname_color TEXT NOT NULL,
		medal_name TEXT NOT NULL,
		up_name TEXT NOT NULL,
		medal_level INT NOT NULL,
		user_level INT NOT NULL
	);
	
	CREATE TABLE IF NOT EXISTS hot(
		time INT NOT NULL,
		hot INT,
		watched INT
	);
	
	CREATE TABLE IF NOT EXISTS gift(
		gift_id INTEGER PRIMARY KEY AUTOINCREMENT,
		uname TEXT NOT NULL,
		uid INT NOT NULL,
		rnd TEXT NOT NULL,
		gift_name TEXT NOT NULL,
		gift_num INT NOT NULL,
		gift_price INT NOT NULL,
		time INT NOT NULL,
		num INT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS sc(
		sc_id INTEGER PRIMARY KEY AUTOINCREMENT,
		id INT NOT NULL,
		uname TEXT NOT NULL,
		uid INT NOT NULL,
		time INT NOT NULL,
		start_time INT NOT NULL,
		message TEXT NOT NULL,
		price INT NOT NULL
	);
	CREATE TABLE IF NOT EXISTS delsc(
		delsc_id INTEGER PRIMARY KEY AUTOINCREMENT,
		id INT NOT NULL,
		time INT NOT NULL
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("NewDb: %w", err)
	}

	return d, nil
}

func (db *Db) Close() error {
	return db.db.Close()
}

func (db *Db) InsertDanmaku(ctx context.Context, t *Danmaku) error {
	err := insert(ctx, t, db, `INSERT INTO "danmaku" ("send_mode", "send_font_size", "danmaku_color", "time", "dmid", "msg_type", "bubble", "content", "mid", "uname", "room_admin", "vip", "svip", "rank", "mobile_verify", "uname_color", "medal_name", "up_name", "medal_level", "user_level") VALUES (:send_mode, :send_font_size, :danmaku_color, :time, :dmid, :msg_type, :bubble, :content, :mid, :uname, :room_admin, :vip, :svip, :rank, :mobile_verify, :uname_color, :medal_name, :up_name, :medal_level, :user_level);`)
	if err != nil {
		return fmt.Errorf("db.InsertDanmaku: %w", err)
	}
	return nil
}

func (db *Db) InsertHot(ctx context.Context, hot, watched int64) error {
	h := Hot{
		Time:    time.Now().Unix(),
		Hot:     hot,
		Watched: watched,
	}
	err := insert(ctx, &h, db, `INSERT INTO "hot" ("time", "hot", "watched") VALUES (:time, :hot, :watched);`)
	if err != nil {
		return fmt.Errorf("db.InsertHot: %w", err)
	}
	return nil
}

func (db *Db) InsertGift(ctx context.Context, gift *Gift) error {
	err := insert(ctx, gift, db, `INSERT INTO gift ("uname", "uid", "rnd", "gift_name", "gift_num", "gift_price", "time", "num") VALUES (:uname, :uid, :rnd, :gift_name, :gift_num, :gift_price, :time, :num);`)
	if err != nil {
		return fmt.Errorf("db.InsertHot: %w", err)
	}
	return nil
}

func (db *Db) InsertSC(ctx context.Context, sc *Sc) error {
	err := insert(ctx, sc, db, `INSERT INTO sc ("id", "uname", "uid", "time", "start_time", "message", "price") VALUES (:id, :uname, :uid, :time, :start_time, :message, :price);`)
	if err != nil {
		return fmt.Errorf("db.InsertHot: %w", err)
	}
	return nil
}

func (db *Db) InsertDelSC(ctx context.Context, id int64) error {
	err := insert(ctx, map[string]interface{}{"id": id, "time": time.Now().Unix()}, db, `INSERT INTO delsc ("id", "time") VALUES (:id, :time);`)
	if err != nil {
		return fmt.Errorf("db.InsertHot: %w", err)
	}
	return nil
}

func insert[D any](ctx context.Context, d D, db *Db, query string) error {
	_, err := db.db.NamedExecContext(ctx, query, d)
	if err != nil {
		e := sqlite3.Error{}
		if errors.As(err, &e) {
			if e.Code == sqlite3.ErrConstraint {
				log.Println(err)
				return nil
			}
			if e.Code == sqlite3.ErrBusy || e.Code == sqlite3.ErrLocked {
				log.Println(err)
				time.Sleep(1 * time.Second)
				return insert(ctx, d, db, query)
			}
		}
		return fmt.Errorf("insert: %w", err)
	}
	return nil
}
