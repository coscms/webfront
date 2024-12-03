package customer

import (
	"time"

	"github.com/coscms/webfront/dbschema"
	"github.com/webx-top/echo"
)

func NewCustomerOptions(customerM *dbschema.OfficialCustomer, noClear ...bool) *CustomerOptions {
	var customer *dbschema.OfficialCustomer
	if customerM == nil {
		customer = dbschema.NewOfficialCustomer(nil)
	} else {
		if len(noClear) > 0 && noClear[0] {
			customer = customerM
		} else {
			_customer := ClearPasswordData(customerM)
			customer = &_customer
		}
	}
	return &CustomerOptions{OfficialCustomer: customer}
}

type CustomerOptions struct {
	*dbschema.OfficialCustomer
	MaxAge     time.Duration // 登录状态有效时长
	SignInType string        // 登录方式

	// device information
	DeviceInfo

	// for qrcode scan login
	SessionID string
	IPAddress string
}

func (co *CustomerOptions) ApplyOptions(options ...CustomerOption) *CustomerOptions {
	for _, option := range options {
		option(co)
	}
	return co
}

type CustomerOption func(*CustomerOptions)

func CustomerName(name string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Name = name
	}
}
func CustomerPassword(password string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Password = password
	}
}
func CustomerMobile(mobile string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Mobile = mobile
	}
}
func CustomerEmail(email string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Email = email
	}
}
func CustomerMaxAgeSeconds(maxAgeSeconds int) CustomerOption {
	return func(c *CustomerOptions) {
		c.MaxAge = time.Duration(maxAgeSeconds) * time.Second
	}
}
func CustomerMaxAge(maxAge time.Duration) CustomerOption {
	return func(c *CustomerOptions) {
		c.MaxAge = maxAge
	}
}
func CustomerSignInType(signInType string) CustomerOption {
	return func(c *CustomerOptions) {
		c.SignInType = signInType
	}
}
func CustomerScense(scense string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Scense = scense
	}
}
func CustomerPlatform(platform string) CustomerOption {
	return func(c *CustomerOptions) {
		c.Platform = platform
	}
}
func CustomerDeviceNo(deviceNo string) CustomerOption {
	return func(c *CustomerOptions) {
		c.DeviceNo = deviceNo
	}
}
func CustomerSessionID(sessionID string) CustomerOption {
	return func(c *CustomerOptions) {
		c.SessionID = sessionID
	}
}
func CustomerIPAddress(ipAddr string) CustomerOption {
	return func(c *CustomerOptions) {
		c.IPAddress = ipAddr
	}
}
func CustomerDeviceInfo(deviceInfo DeviceInfo) CustomerOption {
	return func(c *CustomerOptions) {
		c.DeviceInfo = deviceInfo
	}
}
func GenerateOptionsFromHeader(c echo.Context, maxAge ...int) []CustomerOption {
	d := DeviceInfo{}
	d.Init(c)
	co := []CustomerOption{
		CustomerDeviceInfo(d),
	}
	if len(maxAge) > 0 {
		co = append(co, CustomerMaxAgeSeconds(maxAge[0]))
	}
	return co
}
