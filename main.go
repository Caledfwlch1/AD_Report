// main.go
package main
  
import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"myfunc"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const IniFile = "AD_Report.ini"

type IniParameters struct {
	lisIP   string
	lisPort string
	logFile string
	logIP   string
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	go func() {
		myfunc.CLog.PrintLog(http.ListenAndServe("localhost:6060", nil))
	}()

	os.Chdir(os.Args[0][:strings.LastIndex(os.Args[0], string(os.PathSeparator))])
	if err := myfunc.Conf.ReadINI(); err != nil {
		myfunc.CLog.PrintLog(err)
		return
	}
	filog := myfunc.LoadLog()
	defer filog.Close()
	myfunc.CLog.PrintLog(os.Getwd())
	myfunc.CLog.PrintLog(os.Args)
	myfunc.CheckFiles()
	port := readFlags()

	//	go myfunc.RunTimeOut()

	http.HandleFunc("/", myfunc.Handler)
	http.ListenAndServe(port, nil)

} //  func main()

func PrintHelp() {
	fmt.Println("Вебсервер 'AD утилиты'. Поиск объектов, создание отчётов.")
	fmt.Println("")
	fmt.Println("AD_Report [-h]|[-p:порт]|[-i]|[-u]")
	fmt.Println("	-h	Описание")
	fmt.Println("	-p	Номер порта, поумолчанию -p:8080")
}

func readFlags() string {
	var port, _ = myfunc.Conf.GetValue("Default", "Port")

	if len(os.Args) == 1 {
		PrintHelp()
		return port
	}

	for i := 1; i < len(os.Args); i++ {
		if strings.Contains(fmt.Sprint(os.Args[i]), "-h") {
			PrintHelp()
			os.Exit(0)
			return port
		}
		if strings.Contains(fmt.Sprint(os.Args[i]), "-p") {
			port = os.Args[i][2:]
			myfunc.Conf.SetValue("Default", "Port", port)
			if err := goconfig.SaveConfigFile(&myfunc.Conf.ConfigFile, IniFile); err != nil {
				myfunc.CLog.PrintLog(err)
			}
		}
	}
	return port
}
