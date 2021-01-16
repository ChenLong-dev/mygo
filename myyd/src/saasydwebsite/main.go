/**
* @Author: cl
* @Date: 2021/1/16 15:39
 */
package main
import (
	"github.com/ChenLong-dev/gobase/mbase"
	"github.com/ChenLong-dev/gobase/mg"
	"github.com/ChenLong-dev/gobase/mlog"
	"google.golang.org/grpc"
	"math"
	"myyd/src/saasydwebsite/config"
	ws "myyd/src/saasydwebsite/website"
	"net"
)

type server struct{}

func main() {
	listenAddr := config.Conf.AvailConf.ListenAddr
	mongoHost := config.Conf.Mongo.Host
	mongoUsername := config.Conf.Mongo.UserName
	mongoPassword := config.Conf.Mongo.Password
	mongoDbname := config.Conf.Mongo.DbName

	mbase.Init()

	mlog.Infof("param: listenAddr:%s, mongoHost:%s, mongoUsername:%s, mongoPassword:%s, mongoDbname:%s\n",
		listenAddr, mongoHost, mongoUsername, mongoPassword, mongoDbname)

	if err := initMongoDB(mongoHost, mongoUsername, mongoPassword, mongoDbname); err != nil {
		return
	}

	listen, err := net.Listen("tcp", listenAddr)
	if err != nil {
		mlog.Infof("listen is failed: %v\n", err)
		return
	}
	var options = []grpc.ServerOption{
		grpc.MaxRecvMsgSize(math.MaxInt32),
		grpc.MaxSendMsgSize(1073741824),
	}
	s := grpc.NewServer(options...)

	ws.RegisterWebsiteServer(s, &server{})

	if err := s.Serve(listen); err != nil {
		mlog.Fatalf("failed to serve: %v", err)
	}

}

func initMongoDB(host, username, password, dbname string) error {
	if err := mg.Connect(host, username, password, dbname); err != nil {
		mlog.Errorf("init MongoDB is failed, err:%+v\n", err)
		return err
	}
	mlog.Info("init MongoDB is success ...")
	return nil
}

