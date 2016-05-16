package myfunc

import (
	"fmt"
	"github.com/gogits/gogs/modules/ldap"
	"html/template"
	"net/http"
	//"runtime"
	"strconv"
	"strings"
	"sync"
	//"time"
)

// тип для вывода результата поиска по AD на web-страницу
type OutObject struct {
	NameObject string
	SAMName    string
	TypeObject string
	ExtAttrib  []string
}

// тип для обработки страниц поиска
type Page struct {
	Title         string
	OutputString  []OutObject
	Col           int
	SearchOUweb   string
	PSearchString string
}

type Job struct {
	idx     int
	strName string
	strGrp  []string
}

var ldapConnect *ldap.Conn
var SSearchString string

type TypeConnForRequest struct {
	ldapConn *ldap.Conn
	srchReq  *ldap.SearchRequest
	srchRes  *ldap.SearchResult
	wg       sync.WaitGroup
	jobIn    chan Job
	jobOut   chan Job
}

//	вывод web-страниц, обработка и отображение результатов
func Handler(w http.ResponseWriter, r *http.Request) {
	var result []*ldap.Entry
	var vConnReq TypeConnForRequest
	var length int
	var OutputString []OutObject
	var index = 1
	var err error
	snameUser := ""
	snameComp := ""
	Focus := HelloScreen
	sFindType := "login"
	groupMember := "memberOf"

	sFindType = r.FormValue("FindType")
	SSearchString = r.FormValue("SearchString")
	snameUser = r.FormValue("User")
	snameComp = r.FormValue("Computer")
	ColObjects := r.FormValue("ColObjects")
	ssearchOU := r.FormValue("SearchOU")
	chkName := r.FormValue("ChkName")
	chkSAM := r.FormValue("ChkSAM")
	chkType := r.FormValue("ChkType")
	chkGroup := r.FormValue("ChkGroup")

	if SSearchString == "" {
		SSearchString = "*"
	}

	if ssearchOU == "" {
		ssearchOU = "npp"
	}

	if sFindType != "" && index <= 1 {
		Filter = makeFilter(snameUser, snameComp, SSearchString, sFindType)
		if sFindType == "group" {
			groupMember = "member"
		}

		if ssearchOU == "npp" {
			Base_dn = "DC=npp"
		} else {
			Base_dn = "OU=" + ssearchOU + ",DC=npp"
		}

		CLog.PrintLog("Base_dn = " + Base_dn)

		IntColObjects, _ = strconv.Atoi(ColObjects)

		vConnReq.ldapConn, result, err = WorkLDAP() // LDAP-запрос в AD
		//		fmt.Println("err = ", err.Error())
		if err != nil && !(strings.Contains(err.Error(), "Size Limit Exceeded") || strings.Contains(err.Error(), "Time Limit Exceeded")) {
			OutputString = append(OutputString, OutObject{fmt.Sprintln(err), "", "", []string{""}})
			length = 0
		} else {
			length = len(result)
		}
		if vConnReq.ldapConn != nil {
			defer vConnReq.Close()
		}
		index++

		if length == 0 {
			OutputString = append(OutputString, OutObject{"Нет данных.", "", "", []string{""}})
		} else {
			for i := 0; i < length; i++ {
				OutputString = append(OutputString, fillOutputString(result[i], chkName, chkSAM, chkType, chkGroup, groupMember))
				/*				OutputString = append(OutputString,
								OutObject{stringToUTF8(result[i].DN),
									result[i].GetAttributeValue("sAMAccountName"),
									modifyAttribute(result[i].GetAttributeValue("objectCategory")),
									arrAttrib}) */ // result[i].GetAttributeValues("memberOf") или []string{""}
			}
			//			_ = vConnReq.getExtendedAttribute(OutputString)
		}
	}

	p := &Page{Title: Focus, OutputString: OutputString} // новая страница

	if r.Method == "POST" {
		Focus = ReportScreen
	}
	p.Col = length
	p.PSearchString = SSearchString
	if Base_dn == "npp" {
		p.SearchOUweb = "npp"
	} else {
		p.SearchOUweb = ssearchOU
	}
	renderTemplate(w, p, Focus) // рендерим ее

}

func renderTemplate(w http.ResponseWriter, p *Page, FocusIternal string) {
	var err error
	TemplateHello := template.Must(template.ParseFiles(HelloScreen))   // шаблон для первой страницы
	TemplateReport := template.Must(template.ParseFiles(ReportScreen)) // шаблон для первой страницы и страницы с результатом поиска

	if FocusIternal == HelloScreen {
		err = TemplateHello.ExecuteTemplate(w, FocusIternal, p)
	} else {
		err = TemplateReport.ExecuteTemplate(w, FocusIternal, p)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// формирование строки поиска для LDAP-запроса
func makeFilter(userName, computerName, objectName, sFType string) string {
	var outFilter string

	if sFType == "group" {
		outFilter = "(objectCategory=group)"
	} else {
		if userName == "true" && computerName != "true" {
			outFilter = "(objectCategory=person)"
		}
		if computerName == "true" && userName != "true" {
			outFilter = "(objectCategory=computer)"
		}
		if userName == "true" && computerName == "true" {
			outFilter = "(|(objectCategory=person)(objectCategory=computer))"
		}
		if userName != "true" && computerName != "true" {
			outFilter = ""
		}
	}
	if objectName != "*" {
		if sFType == "sam" {
			outFilter = fmt.Sprintf("(&(sAMAccountName=%s)%s)", objectName, outFilter)
		} else {
			outFilter = fmt.Sprintf("(&(CN=%s)%s)", objectName, outFilter)
		}
	}
	return outFilter
}

// преобразование результата для вывода на страницу - из длиной строки типа объекта в короткую
func modifyAttribute(inputString string) string {
	switch {
	case strings.Contains(inputString, "Person"):
		return "user"
	case strings.Contains(inputString, "Computer"):
		return "computer"
	case strings.Contains(inputString, "Group"):
		return "group"
	default:
		return inputString
	}
	return ""
}

func (v *TypeConnForRequest) Close() {
	defer v.ldapConn.Close()
}

func fillOutputString(vRez *ldap.Entry, sName, sSAM, sType, sGroup, sGM string) (vRezStr OutObject) {

	if sName == "" && sSAM == "" {
		sName = "true"
	}
	if sName == "true" {
		vRezStr.NameObject = vRez.DN
	}
	if sSAM == "true" {
		vRezStr.SAMName = vRez.GetAttributeValue("sAMAccountName")
	}
	if sType == "true" {
		vRezStr.TypeObject = modifyAttribute(vRez.GetAttributeValue("objectCategory"))
	}
	if sGroup == "true" {
		arrAttrib := vRez.GetAttributeValues(sGM)
		for i, arr := range arrAttrib {
			arrAttrib[i] = arr
		}
		vRezStr.ExtAttrib = arrAttrib
	}
	fmt.Println("fillOutputString: ", vRezStr)
	fmt.Println("fillOutputString: ", "sname="+sName, "ssam="+sSAM, "stype="+sType, "sgroup="+sGroup, sGM)
	return vRezStr
}
