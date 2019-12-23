package menu

import (
	"github.com/herb-go/fetch"
	"github.com/jarlyyn/herb-go-experimental/wechatmp"
)

func CreateMenu(App *wechatmp.App, menu *Menu) error {
	result := &wechatmp.ResultAPIError{}
	return App.CallJSONApiWithAccessToken(wechatmp.APIMenuCreate, nil, menu, result)
}

func GetMenu(App *wechatmp.App) (*MenuResult, error) {
	menu := NewMenuResult()
	err := App.CallJSONApiWithAccessToken(wechatmp.APIMenuGet, nil, nil, menu)
	if fetch.CompareAPIErrCode(err, "46003") {
		return menu, nil
	}
	return menu, err
}
