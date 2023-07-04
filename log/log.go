/**
 * @Description:
 * @Version: 1.0.0
 * @Author: liteng
 * @Date: 2023/7/4 11:02 上午
 */
package log

import "github.com/xelabs/go-mysqlstack/xlog"

var (
	Logger *xlog.Log = xlog.NewStdLog(xlog.Level(xlog.INFO))
)