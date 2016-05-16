package myfunc

import (
	"fmt"
	"github.com/Unknwon/goconfig"
	"log"
	"os"
	"time"
)

const (
	LSTR         = "\r\n"
	IniFile      = "AD_Report.ini"
	HelloScreen  = "helloscreen.html"
	ReportScreen = "reportscreen.html"
)

type CommonLog struct {
	log.Logger
}

var CLog CommonLog

type CongSet struct {
	goconfig.ConfigFile
}

var Conf CongSet

// Функция проверки html-файлов.
func CheckFiles() {
	if _, err := os.Stat(HelloScreen); err != nil {
		CLog.PrintLog(err)
		os.Exit(1)
		//		return
	}
	if _, err := os.Stat(ReportScreen); err != nil {
		CLog.PrintLog(err)
		os.Exit(1)
		//		return
	}
}

// Создание/открытие лог-файла
func LoadLog() (logFileWr *os.File) {
	nameLogFile, _ := Conf.GetValue("Logging", "Log_File")
	if _, err := os.Stat(nameLogFile); err != nil {
		logFileWr, err = os.Create(nameLogFile)
		if err != nil {
			fmt.Println(err)
			return nil
		}
	}
	logFileWr, err := os.OpenFile(nameLogFile, os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	logFileWr.Seek(0, os.SEEK_END)
	c := log.New(logFileWr, "", 0)
	CLog.Logger = *c
	CLog.PrintLog(" ################   Запуск AD_Report.exe.   ################ ")
	return logFileWr
}

func (v *CommonLog) PrintLog(s ...interface{}) {
	var str string
	str, _ = os.Hostname()
	str = time.Now().String() + " ; " + str
	str += " ; " + fmt.Sprint(s...) // + LSTR
	v.Println(str)
	v.Output(2, "")
	return
}

func (c *CongSet) ReadINI() (err error) {
	if _, err := os.Stat(IniFile); err != nil {
		if err := CreateDefaultConfig(); err != nil {
			return err
		}
	}
	conf, err := goconfig.LoadConfigFile(IniFile)
	c.ConfigFile = *conf
	return nil
}

func CreateDefaultConfig() error {
	confFile, err := os.Create(IniFile)
	defer confFile.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}
	_, err = fmt.Fprintln(confFile, "# Ini-файл веб-сервера AD_Report."+
		LSTR+"[Default]"+
		LSTR+"IP_Addres = "+
		LSTR+"Port = :8080"+
		LSTR+
		LSTR+"[Logging]"+
		LSTR+"Log_File = AD_Report.log"+
		LSTR+"Host = "+
		LSTR+LSTR)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}
