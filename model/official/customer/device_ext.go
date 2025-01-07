package customer

import "github.com/webx-top/echo"

var (
	DevicePlatforms       = echo.NewKVData()
	DeviceScenses         = echo.NewKVData()
	DefaultDevicePlatform = `pc`
	DefaultDeviceScense   = `web`
)

func init() {
	DevicePlatforms.Add(`ios`, echo.T(`iOS手机`))
	DevicePlatforms.Add(`android`, echo.T(`安卓手机`))
	DevicePlatforms.Add(`pc`, echo.T(`个人电脑`))
	DevicePlatforms.Add(`micro-program`, echo.T(`小程序`))

	DeviceScenses.Add(`app`, echo.T(`App应用`))
	DeviceScenses.Add(`web`, echo.T(`网页`))
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

func (d *DeviceInfo) UnsetSession(c echo.Context) {
	c.Session().Delete(`deviceInfo`)
}

func GetDeviceInfo(c echo.Context) *DeviceInfo {
	d, _ := c.Session().Get(`deviceInfo`).(*DeviceInfo)
	return d
}

func HasSignedInOtherDevice(c echo.Context) bool {
	return c.Internal().Bool(`hasSignedInOtherDevice`)
}
