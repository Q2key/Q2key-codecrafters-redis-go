package core

import (
	"fmt"
	"log"
	"os"
)

func (r *Redis) SendToReplicas(buff *[]byte) {
	*r.ReceivedBytes += len(*buff)
	for _, r := range *r.RepConnPool {
		r.Conn().Write(*buff)
	}
}

func (r *Redis) RegisterReplicaConn(conn Conn) {
	(*r.RepConnPool)[(conn).Id()] = conn
}

func (r *Redis) UpdateReplica(id string, offset int) {
	rep := (*r.RepConnPool)[id]
	rep.SetOffset(offset)
	(*r.RepConnPool)[id] = rep
}

func (r *Redis) TryReadDb() {
	if r.Config.DbFileName == "" || r.Config.Dir == "" {
		return
	}

	path := fmt.Sprintf("%s/%s", r.Config.Dir, r.Config.DbFileName)

	db := NewRedisDB(path)
	if !db.IsFileExists(r.Config.DbFileName) {
		_ = os.Mkdir(r.Config.Dir, os.ModeDir)
	}

	if !db.IsFileExists(path) {
		err := db.Create()
		if err != nil {
			log.Fatal(err)
		}
	}

	err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range db.GetData() {
		r.Store.Set(k, v, STRING)
		exp, ok := db.GetExpires()[k]
		if ok {
			r.Store.SetExpiredAt(k, exp)
		}
	}
}
