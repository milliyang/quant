package main

import (
	"bytes"
	"code.google.com/p/go.net/html"
	"errors"
	"fmt"
	iconv "github.com/djimenez/iconv-go"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

var (
	// money.finance.sina.com.cn/corp/go.php/vMS_FuQuanMarketHistory/stockid/600031.html?year=2014&jidu=1
	URL_SINA = "http://money.finance.sina.com.cn/corp/go.php/vMS_FuQuanMarketHistory/stockid/%s.html?year=%d&jidu=%d"

	URL_SINA_ETF = "http://vip.stock.finance.sina.com.cn/corp/go.php/vMS_MarketHistory/stockid/%s/type/S.phtml?year=%d&jidu=%d"

	// shang zheng zhi shu:
	// http://vip.stock.finance.sina.com.cn/corp/go.php/vMS_MarketHistory/stockid/000001/type/S.phtml?year=2014&jidu=1

	Client *http.Client

	err_failed = errors.New("failed")

	httpDebug = false
)

func init() {
	Client = &http.Client{}
}

func HttpGet(ins Instructment, year, season int) ([]Bar, error) {

	bars := []Bar{}

	var url string
	if ins.Type == "ETF" {
		url = fmt.Sprintf(URL_SINA_ETF, ins.getSymbolNumber(), year, season)
	} else {
		url = fmt.Sprintf(URL_SINA, ins.getSymbolNumber(), year, season)
	}

	if true || httpDebug {
		fmt.Println("url>>", url)
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return bars, err
	} else {

		// req.Header.Add("Host", `money.finance.sina.com.cn`)
		// req.Header.Add("Connection", `keep-alive`)
		// req.Header.Add("Cache-Control", `max-age=0`)
		// req.Header.Add("Accept", `text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8`)
		// req.Header.Add("User-Agent", `Mozilla/5.0 (Macintosh; Intel Mac OS X 10_10_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/39.0.2171.95 Safari/537.36`)
		// req.Header.Add("Accept-Encoding", `gzip, deflate, sdch`)
		// req.Header.Add("Accept-Language", `en-US,en;q=0.8,zh-CN;q=0.6,zh;q=0.4,zh-TW;q=0.2,es;q=0.2`)
		// req.Header.Add("Cookie", `U_TRS1=00000016.a6ce2d8f.535a0a1b.606ca844; UOR=blog.sina.com.cn,blog.sina.com.cn,; SINAGLOBAL=112.95.26.155_1398409757.566653; vjuids=4907f1f43.14597b7d9c5.0.36b186d3; SGUID=1398409779998_35991861; visited_uss=gb_pwrd; __utma=269849203.1754309498.1415416252.1415416252.1415416252.1; __utmz=269849203.1415416252.1.1.utmcsr=games.sina.com.cn|utmccn=(referral)|utmcmd=referral|utmcct=/; sso_info=v02m6alo5qzta2alrGsmpeloZumnKadlqWkj5OEt46DpLaMk4iyjYOIwA==; lxlrtst=1417012390_o; ArtiFSize=14; lxlrttp=1419658073; U_TRS2=00000033.802334af.54a006ad.0463a74c; dpha=usrmdinst_4; dpvar=usrmdinst_1; Apache=198.177.58.51_1419773621.17797; ULV=1419773622982:19:10:1:198.177.58.51_1419773621.17797:1419677652830; FIN_ALL_VISITED=sh600031; FINA_V_S_2=sh600031,sz002047,sh600093,sh601669,sh000001; SR_SEL=1_511; SUS=SID-1789612242-1419776221-GZ-tcvhj-a335593792f5fa798c1e3ea05a4ac1f0; SUB=_2A255pGCNDeTxGedJ41sX8SzOzz6IHXVa0wTFrDV_PUNbvtBeLXDjkW9SAs4ZS9R3OoFoTnceAk7I3niI9A..; ALF=1451312221; SUE=es%3D6b0f0cf20713f6a27b90a157a387cc76%26ev%3Dv1%26es2%3D5e70dc85263b569e95782da2259e7e41%26rs0%3DgjAADm6Z8zqm2PfW5qmNXhnrlk5Ub478b3nVgxcMYBpTKthST787svzoQ8XehweIZDBnehMW9xTSMyi%252FCNRgeYuqSNJATZwZD58ypfT2DgtADO5nblIi55gcCxX6s3BXAKW0nbc2B2J1Epm5FhFgSzYqiob8lHMej%252BHdk2yAOhI%253D%26rv%3D0; SUP=cv%3D1%26bt%3D1419776221%26et%3D1419871621%26d%3D40c3%26i%3D8a7e%26us%3D1%26vf%3D0%26vt%3D0%26ac%3D2%26st%3D0%26lt%3D7%26uid%3D1789612242%26user%3Dmilliyang.cn%26ag%3D8%26name%3Dmilliyang%2540sina.cn%26nick%3Dmilliyang%26sex%3D2%26ps%3D0%26email%3Dmilliyang%2540sina.cn%26dob%3D1984-08-29%26ln%3D%26os%3D%26fmp%3D%26lcp%3D2012-07-06%252012%253A16%253A11; SUBP=0033WrSXqPxfM725Ws9jqgMF55529P9D9WF7qZ67B1-pvUMYppRMEBJb; USRMDE16=usrmdinst_32; _s_upa=3; vjlast=1419791087.1419791087.10`)

		response, err := Client.Do(req)
		if err != nil {
			return bars, err
		} else {
			data, err := ioutil.ReadAll(response.Body)
			defer response.Body.Close()
			if err != nil {
				return bars, err
			} else {
				out := make([]byte, len(data))
				out = out[:]
				iconv.Convert(data, out, "gb2312", "utf-8")

				if httpDebug {
					fmt.Println("httpGet>>", string(out))
				}

				bars = parseHtml2(strings.NewReader(string(out)))
				if bars != nil {
					return bars, nil
				} else {
					return bars, err_failed
				}

				// htmlStr := string(out)
				// idx0 := strings.Index(htmlStr, "<table id=\"FundHoldSharesTable\">")
				// idx1 := strings.Index(htmlStr, "</thead>") + len("</thead>")
				// idx2 := strings.LastIndex(htmlStr, "</table>")
				// fmt.Println(htmlStr[idx0 : len(htmlStr)-1])
				// fmt.Println("idx:", idx0, idx1, idx2)
				// newStr := htmlStr[idx0:idx1] + "<tbody>" + htmlStr[idx1:idx2] + "<tbody></table>"
			}
		}
	}
	return bars, err_failed
}

func nodeText(n *html.Node) string {
	var b bytes.Buffer
	err := html.Render(&b, n)
	if err != nil {
		panic(err.Error())
	}
	return b.String()
}

func parseHtml2(r io.Reader) []Bar {
	doc, err := html.Parse(r)
	if err != nil {
		log.Fatal(err)
	}

	var tableNode *html.Node

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode {
			if n.Data == "table" {
				for _, a := range n.Attr {
					if a.Key == "id" && a.Val == "FundHoldSharesTable" {
						tableNode = n
						return
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)

	if tableNode == nil {
		return []Bar{}
	}

	return parseTableNode(tableNode)
}

func parseTableNode(n *html.Node) []Bar {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if httpDebug {
			fmt.Println("parseTableNode", c.Data, c.Attr)
		}
		if c.Type == html.ElementNode {
			if c.Data == "tbody" {
				return parseTableTboby(c)
			} else if c.Data == "thead" {

			}
		}
	}
	return nil
}

func parseTableTboby(tbody *html.Node) []Bar {
	bars := []Bar{}
	firstRow := true
	for c := tbody.FirstChild; c != nil; c = c.NextSibling {
		// parse row
		if c.Type == html.ElementNode {
			if firstRow {
				for td := c.FirstChild; td != nil; td = td.NextSibling {
					if td.Type == html.ElementNode {
						// fmt.Println("td:", nodeText(td))
						if httpDebug {
							fmt.Print(td.FirstChild.FirstChild.FirstChild.Data, ",")
						}
					}
				}
				if httpDebug {
					fmt.Print("\n")
				}
				firstRow = false
			} else {
				isDate := true

				oneBar := Bar{}

				for td := c.FirstChild; td != nil; td = td.NextSibling {
					if td.Type == html.ElementNode {
						// fmt.Println("td:", nodeText(td))
						if isDate {
							isDate = false
							found := false

							// http://money.finance.sina.com.cn/corp/go.php/vMS_FuQuanMarketHistory/stockid/200016.html?year=2006&jidu=3
							// <td class="head">
							// <div align="center">
							// <a target="_blank" href="http://vip.stock.finance.sina.com.cn/quotes_service/view/vMS_tradehistory.php?symbol=sz200016&amp;date=2006-09-29">
							// 2006-09-29						</a>
							// </div>
							// </td>
							for tdx := td.FirstChild.FirstChild; tdx != nil; tdx = tdx.NextSibling {
								if tdx.Type == html.ElementNode {
									item := tdx.FirstChild.Data
									oneBar.Date = strings.TrimSpace(item)
									found = true
									break
								}
							}
							if !found {
								// http://money.finance.sina.com.cn/corp/go.php/vMS_FuQuanMarketHistory/stockid/200016.html?year=2005&jidu=3
								// <td class="head">
								// <div align="center">
								// 2005-09-30 </div>
								// </td>
								for tdx := td.FirstChild; tdx != nil; tdx = tdx.NextSibling {
									if tdx.Type == html.ElementNode {
										item := tdx.FirstChild.Data
										oneBar.Date = strings.TrimSpace(item)
										found = true
										break
									}
								}
							}
						} else {
							item := td.FirstChild.FirstChild.Data
							oneBar.Items = append(oneBar.Items, strings.TrimSpace(item))
						}
					}
				}
				if oneBar.selfCheck() {
					// Fxxking SINA
					// open, high, close, low, volumn, amouunt, factor ->
					// open, high, close, low, volumn, amount, factor
					oneBar.fxxkingSina()
					bars = append(bars, oneBar)
				} else {
					fmt.Println("http invalid bar")
					JsonPrint(oneBar)
					// debug
					os.Exit(2)
				}
			}
		}
	}
	return bars
}
