package taobao

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiHost = "taobao-datahub.p.rapidapi.com"
	apiKey  = "5900c7af08mshcaa0cb0046fed36p16afa7jsn2b8bf5a505d9" // Replace with your actual API key
)

func SearchImageOnTaobao(imageURL string, pageSize int) (*ApiResponse, error) {
	url := fmt.Sprintf("https://%s/item_search_image?imgUrl=%s&pageSize=%d", apiHost, imageURL, pageSize)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var resp ApiResponse
	err = json.Unmarshal(body, &resp)
	return &resp, err
}

func SearchItemDetails(itemID string) (*Response, error) {
	url := fmt.Sprintf("https://%s/item_detail?itemId=%s", apiHost, itemID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("x-rapidapi-key", apiKey)
	req.Header.Add("x-rapidapi-host", apiHost)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to do request: %w", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	var resp Response
	err = json.Unmarshal(body, &resp)
	return &resp, err
}

type ApiResponse struct {
	Result struct {
		Status struct {
			Code             int    `json:"code"`
			Attempt          int    `json:"attempt"`
			P                bool   `json:"p"`
			Data             string `json:"data"`
			ExecutionTime    string `json:"executionTime"`
			RequestTime      string `json:"requestTime"`
			RequestID        string `json:"requestId"`
			Endpoint         string `json:"endpoint"`
			APIVersion       string `json:"apiVersion"`
			FunctionsVersion string `json:"functionsVersion"`
			La               string `json:"la"`
			Pmu              int    `json:"pmu"`
			Mu               int    `json:"mu"`
		} `json:"status"`
		Settings struct {
			CatID     string `json:"catId"`
			Sort      string `json:"sort"`
			ImgRegion string `json:"imgRegion"`
			ImgURL    string `json:"imgUrl"`
			Region    string `json:"region"`
			Locale    string `json:"locale"`
			Currency  string `json:"currency"`
		} `json:"settings"`
		Base struct {
			SortValues    []string `json:"sortValues"`
			ImgRegion     string   `json:"imgRegion"`
			ImgRegionFull string   `json:"imgRegionFull"`
			PageSize      string   `json:"pageSize"`
			CategoryList  []struct {
				Name string `json:"name"`
				ID   any `json:"id"`
			} `json:"categoryList"`
		} `json:"base"`
		ResultList []struct {
			Item struct {
				ItemID  any `json:"itemId"`
				Title   string `json:"title"`
				Sales   any `json:"sales"`
				ItemURL string `json:"itemUrl"`
				Image   string `json:"image"`
				Sku     struct {
					Def struct {
						Price          string `json:"price"`
						PromotionPrice string `json:"promotionPrice"`
					} `json:"def"`
				} `json:"sku"`
			} `json:"item"`
			Delivery struct {
				ShippingFrom string `json:"shippingFrom"`
			} `json:"delivery"`
			Seller struct {
				SellerID   string `json:"sellerId"`
				StoreTitle string `json:"storeTitle"`
				StoreType  string `json:"storeType"`
			} `json:"seller"`
		} `json:"resultList"`
	} `json:"result"`
}

/*{
	"access_token": "50000300a02Moz1133a89f7jAqyRsEY7iVDIwDVQE8FFRBOUVgIFyo3FLyMYn9",
	"refresh_token": "500013016020q215e7ed7cdjqcuTfBbjsfLS0GJFvIZIS0nFYCRQFRwXKrzrm5",
	"user_id": "2100008076043",
	"account_platform": "seller_center",
	"refresh_expires_in": 5184000,
	"expires_in": 2592000,
	"seller_id": "2100008076043",
	"account": "hajyyevnazar123@gmail.com",
	"short_code": "",
	"code": "0",
	"request_id": "213bceab17181929135374222",
	"_trace_id_": "2102fcce17181929135341345ef1f8"
 }

 {'access_token': '50000301434bs7gvHiTaWRCfdBPhsWxcFDRoaiQUHl1HOR17987e03AxJ12hdw',
 'refresh_token': '50001301934e1yatTdAcFSdttUCJlykvDQ0HH1iQwjOxVS195203abLQQ1HkSH',
 'user_id': '2100008076043', 'account_platform': 'seller_center', 'refresh_expires_in': 5184000, 'expires_in': 2592000,
 'seller_id': '2100008076043', 'account': 'hajyyevnazar123@gmail.com', 'short_code': '', 'code': '0',
 'request_id': '2101751417182076200382241', '_trace_id_': '2101448e17182076200372719ea07b'}
*/

// Define the structs
type Response struct {
	Result Result `json:"result"`
}

type Result struct {
	Status   Status   `json:"status"`
	Settings Settings `json:"settings"`
	Item     Item     `json:"item"`
}

type Status struct {
	Code             int    `json:"code"`
	Data             string `json:"data"`
	ExecutionTime    string `json:"executionTime"`
	RequestTime      string `json:"requestTime"`
	RequestId        string `json:"requestId"`
	Endpoint         string `json:"endpoint"`
	ApiVersion       string `json:"apiVersion"`
	FunctionsVersion string `json:"functionsVersion"`
	La               string `json:"la"`
	Pmu              int    `json:"pmu"`
	Mu               int    `json:"mu"`
}

type Settings struct {
	Locale    string `json:"locale"`
	Currency  string `json:"currency"`
	ItemId    string `json:"itemId"`
	ItemIdStr string `json:"itemIdStr"`
}

type Item struct {
	ItemId      int64       `json:"itemId"`
	Title       string      `json:"title"`
	CatId       string      `json:"catId"`
	CatName     string      `json:"catName"`
	Sales       int         `json:"sales"`
	ItemUrl     string      `json:"itemUrl"`
	Images      []string    `json:"images"`
	Video       Video       `json:"video"`
	Properties  Properties  `json:"properties"`
	Description Description `json:"description"`
	Sku         Sku         `json:"sku"`
}

type Video struct {
	Id        string `json:"id"`
	Thumbnail string `json:"thumbnail"`
	Url       string `json:"url"`
}

type Properties struct {
	Cut  string     `json:"cut"`
	List []PropList `json:"list"`
}

type PropList struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Description struct {
	Url    string   `json:"url"`
	Images []string `json:"images"`
}

type Sku struct {
	Def   SkuDef    `json:"def"`
	Base  []SkuBase `json:"base"`
	Props []struct {
		PID    int       `json:"pid"`
		Name   string       `json:"name"`
		Values []PropValues `json:"values"`
	} `json:"props"`
}

type SkuDef struct {
	Quantity       int    `json:"quantity"`
	Price          string `json:"price"`
	PromotionPrice string `json:"promotionPrice"`
}

type SkuBase struct {
	PropPath       string `json:"propPath"`
	SkuId          int64  `json:"skuId"`
	Quantity       int    `json:"quantity"`
	Price          string `json:"price"`
	PromotionPrice string `json:"promotionPrice"`
}

type PropValues struct {
	Name string `json:"name"`
	VID  int  `json:"vid"`
	Image string `json:"image"`
}

func (s *SkuBase) UnmarshalJSON(data []byte) error {
	type Alias SkuBase
	aux := &struct {
		PropPath string `json:"propPath"`
		PropMap  string `json:"propMap"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Set PropPath to the value of propPath or propMap
	if aux.PropPath != "" {
		s.PropPath = aux.PropPath
	} else {
		s.PropPath = aux.PropMap
	}

	return nil
}