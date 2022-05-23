package util

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"math/rand"
	"ptc-Game/common/pkg/logiclog"
	"time"
)

func RandomDigit() string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	return fmt.Sprintf("%06v", rnd.Int31n(1000000))
}

func GetUid(c iris.Context) (uint, bool) {

	uid, err := c.Values().GetInt("userinfo.uid")
	if err != nil || uid < 0 {
		logiclog.CtxLogger(c).Errorf("getUid err: %+v", err)
		return 0, false
	}
	return uint(uid), true
}
