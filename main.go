// metadata listens on a TCP port waiting for metadata and writes
// it to a MariaDB database by the same name.
package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"

	"gopkg.in/alecthomas/kingpin.v2"
)

var appVersion = "0.0.1"

type cmdlineArgs struct {
	IpAddress net.IP
	TcpPort   uint16
	MyConfig  string
	DbHost    string
	DbPort    int
	DbUser    string
	DbPass    string
	DbName    string
	Verbose   bool
	Debug     bool
}

const defaultAddress = "0.0.0.0" // all interfaces
const defaultPort = 52341

var cmdline cmdlineArgs

func main() {
	// var cartBits CartData
	var (
		ipAddress = kingpin.Flag("ip-addr", "The IP address on which to listen for metadata.").
				Short('a').
				Default(defaultAddress).
				HintOptions(string(net.ParseIP(defaultAddress))).
				IP()
		tcpPort = kingpin.Flag("port", "The TCP port on which to listen for metadata.").
			Short('p').
			Default(strconv.Itoa(defaultPort)).
			HintOptions(strconv.Itoa(defaultPort)).
			Uint16()
		myconfig = kingpin.Flag("myconfig", "The full path to a .my.cnf configuration file").
				Short('y').
				Default(os.Getenv("HOME") + "/.my.cnf").
				String()
		dbhost = kingpin.Flag("dbhost", "The name or IP address of the database host").
			String()
		dbport = kingpin.Flag("dbport", "The TCP Port number for the database host").
			Default("3306").
			HintOptions("3306").
			Int()
		dbuser = kingpin.Flag("dbuser", "The name of the database user").
			String()
		dbpass = kingpin.Flag("dbpass", "The password for the database user").
			String()
		dbname = kingpin.Flag("dbname", "The name of the database.").
			String()
		verbose = kingpin.Flag("verbose", "Be chatty when running").
			Short('v').
			Bool()
		debug = kingpin.Flag("debug", "Be very verbose about what is going on (implies -v). Also enables profiling.").
			Short('d').
			Bool()
	)
	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.CommandLine.Help = "Metadata listener and writer."
	kingpin.UsageTemplate(kingpin.DefaultUsageTemplate).Version(appVersion).Author("Broadcast Tool & Die, David Klann")
	kingpin.Parse()

	cmdline.IpAddress = *ipAddress
	cmdline.TcpPort = *tcpPort
	cmdline.MyConfig = *myconfig
	cmdline.DbUser = *dbuser
	cmdline.DbPass = *dbpass
	cmdline.DbName = *dbname
	cmdline.DbHost = *dbhost
	cmdline.DbPort = *dbport
	cmdline.Verbose = *verbose
	cmdline.Debug = *debug
	debugPrint(fmt.Sprintf("cmdline: %+#v\n", cmdline))

	if cmdline.Debug {
		cmdline.Verbose = true
		debugPrint("Debugging enabled\n")
		debugPrint("Profiling enabled...\n")
		go func() {
			log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	log.Println(logSprintf("starting app version: %s", appVersion))

	log.Println(logSprintf("about to listen on TCP port %d", cmdline.TcpPort))

	listener, err := net.Listen("tcp", cmdline.IpAddress.String()+":"+strconv.Itoa(int(cmdline.TcpPort)))
	if err != nil {
		log.Println(logSprintf("error trying to listen on '%s:%d' (%v)", cmdline.IpAddress, cmdline.TcpPort, err))
		os.Exit(1)
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Println(logSprintf("Error on listener.Accept(): %#v", err))
			os.Exit(2)
		}
		verbosePrint(fmt.Sprintf("accepted connection from %s", connection.RemoteAddr().String()))
		go handleConnection(connection)
	}
}
