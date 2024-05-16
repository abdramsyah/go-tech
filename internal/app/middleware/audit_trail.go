package middleware

import (
	"emoney-backoffice/internal/app/constant"
	"emoney-backoffice/internal/app/dto"
	"emoney-backoffice/internal/app/util"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (m *CustomMiddleware) AuditTrailMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		//Get route name
		routeName := "UNKNOWN"
		for _, r := range c.Echo().Routes() {
			if r.Method == c.Request().Method && r.Path == c.Path() {
				routeName = r.Name
			}
		}
		actx, err := util.NewAppContext(c)
		if err != nil {
			resp := dto.FailedHttpResponse("", constant.ErrUnauthorized, nil)
			return c.JSON(http.StatusUnauthorized, resp)
		}
		adminID := actx.GetAdminID()
		urlString := c.Request().URL.String()

		//Get request body
		//bodyBytes, _ := io.ReadAll(c.Request().Body)
		//bodyBytesCopy := io.NopCloser(bytes.NewBuffer(bodyBytes))
		//var prettyJSON bytes.Buffer
		//if err := json.Indent(&prettyJSON, bodyBytes, "", "\t"); err == nil {
		//	fmt.Println("Request body V2: ", string(prettyJSON.Bytes()))
		//}
		//c.Request().Body = bodyBytesCopy

		req := dto.AuditTrailRequest{
			AdminID:   uint(adminID),
			URL:       fmt.Sprintf("[%s] %s", c.Request().Method, urlString),
			RouteName: routeName,
		}
		_ = m.Service.AuditTrail.Create(c, &req)

		return next(c)
	}
}
