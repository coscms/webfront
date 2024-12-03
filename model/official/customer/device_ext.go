package customer

import "github.com/webx-top/echo"

var (
	DevicePlatforms       = echo.NewKVData()
	DeviceScenses         = echo.NewKVData()
	DefaultDevicePlatform = `pc`
	DefaultDeviceScense   = `web`
)

func init() {
	DevicePlatforms.Add(`ios`, `iOS手机`)
	DevicePlatforms.Add(`android`, `安卓手机`)
	DevicePlatforms.Add(`pc`, `个人电脑`)
	DevicePlatforms.Add(`micro-program`, `小程序`)

	DeviceScenses.Add(`app`, `App应用`)
	DeviceScenses.Add(`web`, `网页`)
}

type DeviceInfo struct {
	Scense   string // 场景
	Platform string // 系统平台
	DeviceNo string // 设备编号
}

func (d *DeviceInfo) Init(c echo.Context) {
	d.Platform = c.Header(`X-Platform`)
	d.Scense = c.Header(`X-Scense`)
	d.DeviceNo = c.Header(`X-Device-Id`)
}

func (d *DeviceInfo) SetSession(c echo.Context) {
	c.Session().Set(`deviceInfo`, d)
}

func GetDeviceInfo(c echo.Context) *DeviceInfo {
	d, _ := c.Session().Get(`deviceInfo`).(*DeviceInfo)
	return d
}
