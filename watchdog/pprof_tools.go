package watchdog

import (
	"net/http"
	_ "net/http/pprof"
	"single-ticker/config"
	"single-ticker/utils/logger"
)

func StartPprofNet(cfg *config.Config) {
	if cfg.PprofListenAddress == "" {
		logger.Info("[Watchdog] No Need Start Pprof Net")
		return
	}
	go func() {
		http.ListenAndServe(cfg.PprofListenAddress, nil)
	}()
	logger.Info("[Watchdog] Start Pprof Net")
}
