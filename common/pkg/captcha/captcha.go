package captcha

import (
	bc "github.com/mojocn/base64Captcha"
	"math/rand"
	"time"
)

func init() {
	//init rand seed
	rand.Seed(time.Now().UnixNano())
}

type Type string

const (
	Str   Type = "string"
	Digit Type = "digit"
)

type Captcha interface {
	Generate() (id, content, b64s string, err error)
	Verify(id, answer string, clear bool) (match bool)
}

// New return a Captcha interface
func New(type_ Type) Captcha {
	var captchaDriver Captcha
	switch type_ {
	case Digit:
		captchaDriver = newCaptcha(bc.DefaultDriverDigit, bc.DefaultMemStore)
	case Str:
		//stringDriver := bc.NewDriverString(80, 240, 20, 100, 2, "", nil, nil)
		//captchaDriver = newCaptcha()
	default:
		captchaDriver = newCaptcha(bc.DefaultDriverDigit, bc.DefaultMemStore)
	}
	return captchaDriver
}

// captcha implement Captcha interface
type captcha struct {
	Driver bc.Driver
	Store  bc.Store
}

//NewCaptcha creates a captcha instance from driver and store
func newCaptcha(driver bc.Driver, store bc.Store) *captcha {
	return &captcha{Driver: driver, Store: store}
}

//Generate generates a random id, base64 image string or an error if any
func (c *captcha) Generate() (id, content, b64s string, err error) {
	id, content, answer := c.Driver.GenerateIdQuestionAnswer()
	item, err := c.Driver.DrawCaptcha(content)
	if err != nil {
		return "", "", "", err
	}
	c.Store.Set(id, answer)
	b64s = item.EncodeB64string()
	return
}

//Verify by a given id key and remove the captcha value in store,
//return boolean value.
//if you has multiple captcha instances which share a same store.
//You may want to call `store.Verify` method instead.
func (c *captcha) Verify(id, answer string, clear bool) (match bool) {
	match = c.Store.Get(id, clear) == answer
	return
}
