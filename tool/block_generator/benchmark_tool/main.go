package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"github.com/dappley/go-dappley/config"
	"github.com/dappley/go-dappley/config/pb"
	"github.com/dappley/go-dappley/consensus"
	"github.com/dappley/go-dappley/contract"
	"github.com/dappley/go-dappley/core"
	"github.com/dappley/go-dappley/logic"
	"github.com/dappley/go-dappley/network"
	"github.com/dappley/go-dappley/storage"
	"github.com/shirou/gopsutil/mem"
	"github.com/shirou/gopsutil/cpu"
	logger "github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"time"
)

const(
	nodeDbPath = "db/temp.db"
	reportFilePath = "report.csv"
	genesisAddrTest = "121yKAXeG4cw6uaGCBYjWk9yTWmMkhcoDD"
	genesisFilePathTest = "../conf/genesis.conf"
	testport1 = 10851
	testport2 = 10852
	cdBtwTest = time.Second*10
)

func main() {

	logger.SetFormatter(&logger.TextFormatter{
		FullTimestamp: true,
	})

	logger.SetLevel(logger.InfoLevel)

	logger.Warn("Please make sure you have moved the test database files to benchmark_tool/db folder!")
	logger.Info("Type enter to continue...")
	reader := bufio.NewReader(os.Stdin)
	reader.ReadString('\n')

	files, err := ioutil.ReadDir("./db")
	if err!= nil{
		logger.WithError(err).Panic("can not read db files")
	}
	if len(files)==0 {
		logger.Warn("No db files are found. Exiting...")
		return
	}

	var elapsed time.Duration
	var blkHeight uint64
	var numOfTx int

	filePath := "report_" + time.Now().String() + ".csv"

	initRecordFile(filePath)

	for _, file := range files {
		elapsed, blkHeight, numOfTx = runTest("./db/"+file.Name())
		recordResult(filePath, file.Name(), elapsed, blkHeight, numOfTx)
		logger.WithFields(logger.Fields{
			"time_elapsed" : elapsed,
			"ave_blk_time"  : elapsed/time.Duration(blkHeight),
			"ave_tx_time"  : elapsed/time.Duration(blkHeight)/time.Duration(numOfTx),
		}).Info("Test Finished")
		time.Sleep(cdBtwTest)
	}
}

func initRecordFile(filePath string){
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Panic("Open file failed while recording failed transactions")
	}
	w := csv.NewWriter(f)
	vmStat, err := mem.VirtualMemory()
	cpuStat, err := cpu.Info()
	w.Write([]string{
		"","","","","","","Memory",
		vmStat.String(),
	})
	w.Write([]string{
		"","","","","","","Cpu",
		cpuStat[0].String(),
	})
	w.Write([]string{
		"Timestamp",
		"Db file",
		"Number of blocks",
		"Number of transactions per block",
		"Total time elapsed",
		"Average block time",
		"Average transaction time",
	})
	w.Flush()
}

func recordResult(filePath string, dbfileName string, elapsed time.Duration, blkHeight uint64, numOfTx int){
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		logger.Panic("Open file failed while recording failed transactions")
	}
	w := csv.NewWriter(f)
	w.Write([]string{
		time.Now().String(),
		dbfileName,
		fmt.Sprint(blkHeight),
		fmt.Sprint(numOfTx),
		elapsed.String(),
		(elapsed/time.Duration(blkHeight)).String(),
		(elapsed/time.Duration(blkHeight)/time.Duration(numOfTx)).String(),
	})
	w.Flush()
}

func runTest(fileName string) (time.Duration, uint64, int){
	logger.WithFields(logger.Fields{
		"db_file_path" : fileName,
	}).Info("Test starts...")
	defer os.RemoveAll(nodeDbPath)
	db1 := storage.OpenDatabase(fileName)
	defer db1.Close()
	db2 := storage.OpenDatabase(nodeDbPath)
	defer db2.Close()

	bc, node1 := prepareNode(db1)
	bc2, node2 := prepareNode(db2)

	node1.Start(testport1)
	defer node1.Stop()
	node2.Start(testport2)
	defer node2.Stop()

	node1.GetPeerManager().AddAndConnectPeer(node2.GetInfo())

	blkHeight := bc.GetMaxHeight()
	tailBlock,_ := bc.GetTailBlock()
	numOfTx := len(tailBlock.GetTransactions())

	logger.WithFields(logger.Fields{
		"blk_height" : blkHeight,
		"num_of_tx"  : numOfTx,
	}).Info("Start Downloading...")

	time.Sleep(time.Second)

	start := time.Now()
	downloadBlocks(node2, bc2)
	elapsed := time.Since(start)

	logger.WithFields(logger.Fields{
		"time_elapsed" : elapsed,
		"ave_blk_time"  : elapsed/time.Duration(blkHeight),
		"ave_tx_time"  : elapsed/time.Duration(blkHeight)/time.Duration(numOfTx),
	}).Info("Downloading ends... Cleaning up files...")



	return elapsed, blkHeight, numOfTx
}



func prepareNode(db storage.Storage) (*core.Blockchain, *network.Node){
	genesisConf := &configpb.DynastyConfig{}
	config.LoadConfig(genesisFilePathTest, genesisConf)
	maxProducers := (int)(genesisConf.GetMaxProducers())
	dynasty := consensus.NewDynastyWithConfigProducers(genesisConf.GetProducers(), maxProducers)
	conss := consensus.NewDPOS()
	conss.SetDynasty(dynasty)
	txPoolLimit:=uint32(2000)
	bc, err := core.GetBlockchain(db, conss, txPoolLimit, vm.NewV8EngineManager(core.Address{}))
	if err != nil {
		bc, err = logic.CreateBlockchain(core.NewAddress(genesisAddrTest), db, conss, txPoolLimit, vm.NewV8EngineManager(core.Address{}))
		if err != nil {
			logger.Panic(err)
		}
	}

	bc.SetState(core.BlockchainInit)
	node := network.NewNode(bc, core.NewBlockPool(0))
	return bc, node
}

func downloadBlocks(node *network.Node, bc *core.Blockchain) {
	downloadManager := node.GetDownloadManager()
	finishChan := make(chan bool, 1)
	bc.SetState(core.BlockchainDownloading)
	downloadManager.StartDownloadBlockchain(finishChan)
	<-finishChan
	bc.SetState(core.BlockchainReady)
}