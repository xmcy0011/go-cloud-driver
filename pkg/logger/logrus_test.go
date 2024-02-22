package logger

import (
	"os"
)

func Example_basic() {
	os.Setenv("LOGOUT", "1")

	l := NewLogger()
	// l.(*serverLog).logger.SetReportCaller(false)
	// l.(*serverLog).logger.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})

	l.Debugln("debug")
	l.Infoln("info")
	l.Infof("hello %s", "world")
	l.Errorln("error")
	l.Warnln("warn")

	// Output:
	// {"level":"info","msg":"info"}
	// {"level":"info","msg":"hello world"}
	// {"level":"error","msg":"error"}
	// {"level":"warning","msg":"warn"}
}

// func Example_fields() {
// 	os.Setenv("LOGOUT", "1")

// 	l := NewLogger()
// 	l.(*serverLog).logger.SetReportCaller(false)
// 	l.(*serverLog).logger.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})

// 	logger := l.WithField("userId", "admin")
// 	logger.Infoln("login")
// 	logger.Infoln("logout")

// 	// instead
// 	// logger.Infof("userId: %d, login", "admin")
// 	// logger.Infof("userId: %d, logout", "admin")

// 	// Output:
// 	// {"level":"info","msg":"login","userId":"admin"}
// 	// {"level":"info","msg":"logout","userId":"admin"}
// }

// func Example_trace() {
// 	os.Setenv("LOGOUT", "1")

// 	l := NewLogger()
// 	l.(*serverLog).logger.SetReportCaller(false)
// 	l.(*serverLog).logger.SetFormatter(&logrus.JSONFormatter{DisableTimestamp: true})

// 	ctx := context.WithValue(context.Background(), key("requestId"), "fd9415de-4f2f-4a59-a464-9dcd206c3d01")

// 	driver := func(ctx context.Context) {
// 		l.Trace(ctx).Infoln("driver in")
// 	}
// 	logics := func(ctx context.Context) {
// 		l.Trace(ctx).Infoln("logics in")
// 	}
// 	dbaccess := func(ctx context.Context) {
// 		l.Trace(ctx).Infoln("dbaccess in")
// 	}
// 	driven := func(ctx context.Context) {
// 		l.Trace(ctx).Infoln("driven in")
// 	}
// 	driverout := func(ctx context.Context) {
// 		l.Trace(ctx).Infoln("driver out")
// 	}

// 	driver(ctx)
// 	logics(ctx)
// 	dbaccess(ctx)

// 	ctx2 := WithRequestId(context.Background(), "fd9415de-4f2f-4a59-a464-9dcd206c3d01")
// 	driven(ctx2)
// 	driverout(ctx2)

// 	// Output:
// 	// {"level":"info","msg":"driver in","requestId":"fd9415de-4f2f-4a59-a464-9dcd206c3d01"}
// 	// {"level":"info","msg":"logics in","requestId":"fd9415de-4f2f-4a59-a464-9dcd206c3d01"}
// 	// {"level":"info","msg":"dbaccess in","requestId":"fd9415de-4f2f-4a59-a464-9dcd206c3d01"}
// 	// {"level":"info","msg":"driven in","requestId":"fd9415de-4f2f-4a59-a464-9dcd206c3d01"}
// 	// {"level":"info","msg":"driver out","requestId":"fd9415de-4f2f-4a59-a464-9dcd206c3d01"}
// }
