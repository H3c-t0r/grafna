package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"runtime/trace"
	"strconv"
	"syscall"
	"time"

	"github.com/grafana/grafana/pkg/api"
	"github.com/grafana/grafana/pkg/extensions"
	"github.com/grafana/grafana/pkg/infra/log"
	"github.com/grafana/grafana/pkg/infra/metrics"
	"github.com/grafana/grafana/pkg/server"
	_ "github.com/grafana/grafana/pkg/services/alerting/conditions"
	_ "github.com/grafana/grafana/pkg/services/alerting/notifiers"
	"github.com/grafana/grafana/pkg/setting"
)

// The following variables cannot be constants, since they can be overridden through the -X link flag
var version = "7.5.0"
var commit = "NA"
var buildBranch = "main"
var buildstamp string

type exitWithCode struct {
	reason string
	code   int
}

func (e exitWithCode) Error() string {
	return e.reason
}

func main() {
	var (
		configFile = flag.String("config", "", "path to config file")
		homePath   = flag.String("homepath", "", "path to grafana install/home path, defaults to working directory")
		pidFile    = flag.String("pidfile", "", "path to pid file")
		packaging  = flag.String("packaging", "unknown", "describes the way Grafana was installed")

		v           = flag.Bool("v", false, "prints current version and exits")
		vv          = flag.Bool("vv", false, "prints current version, all dependencies and exits")
		profile     = flag.Bool("profile", false, "Turn on pprof profiling")
		profileAddr = flag.String("profile-addr", "localhost", "Define custom address for profiling")
		profilePort = flag.Uint64("profile-port", 6060, "Define custom port for profiling")
		tracing     = flag.Bool("tracing", false, "Turn on tracing")
		tracingFile = flag.String("tracing-file", "trace.out", "Define tracing output file")
	)

	flag.Parse()

	if *v || *vv {
		fmt.Printf("Version %s (commit: %s, branch: %s)\n", version, commit, buildBranch)
		if *vv {
			fmt.Println("Dependencies:")
			if info, ok := debug.ReadBuildInfo(); ok {
				for _, dep := range info.Deps {
					fmt.Println(dep.Path, dep.Version)
				}
			}
		}
		os.Exit(0)
	}

	profileDiagnostics := newProfilingDiagnostics(*profile, *profileAddr, *profilePort)
	if err := profileDiagnostics.overrideWithEnv(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	traceDiagnostics := newTracingDiagnostics(*tracing, *tracingFile)
	if err := traceDiagnostics.overrideWithEnv(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	if profileDiagnostics.enabled {
		fmt.Println("diagnostics: pprof profiling enabled", "addr", profileDiagnostics.addr, "port", profileDiagnostics.port)
		runtime.SetBlockProfileRate(1)
		go func() {
			err := http.ListenAndServe(fmt.Sprintf("%s:%d", profileDiagnostics.addr, profileDiagnostics.port), nil)
			if err != nil {
				panic(err)
			}
		}()
	}

	if err := executeServer(*configFile, *homePath, *pidFile, *packaging, traceDiagnostics); err != nil {
		code := 1
		var ewc exitWithCode
		if errors.As(err, &ewc) {
			code = ewc.code
		}
		if code != 0 {
			fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		}

		os.Exit(code)
	}
}

func executeServer(configFile, homePath, pidFile, packaging string, traceDiagnostics *tracingDiagnostics) error {
	defer func() {
		if err := log.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to close log: %s\n", err)
		}
	}()

	if traceDiagnostics.enabled {
		fmt.Println("diagnostics: tracing enabled", "file", traceDiagnostics.file)
		f, err := os.Create(traceDiagnostics.file)
		if err != nil {
			panic(err)
		}
		defer func() {
			if err := f.Close(); err != nil {
				log.Error("Failed to write trace diagnostics", "path", traceDiagnostics.file, "err", err)
			}
		}()

		if err := trace.Start(f); err != nil {
			panic(err)
		}
		defer trace.Stop()
	}

	buildstampInt64, err := strconv.ParseInt(buildstamp, 10, 64)
	if err != nil || buildstampInt64 == 0 {
		buildstampInt64 = time.Now().Unix()
	}

	setting.BuildVersion = version
	setting.BuildCommit = commit
	setting.BuildStamp = buildstampInt64
	setting.BuildBranch = buildBranch
	setting.IsEnterprise = extensions.IsEnterprise
	setting.Packaging = validPackaging(packaging)

	metrics.SetBuildInformation(version, commit, buildBranch)

	s, err := server.Initialize(setting.CommandLineArgs{
		Config: configFile, HomePath: homePath, Args: flag.Args(),
	}, server.Options{
		PidFile: pidFile, Version: version, Commit: commit, BuildBranch: buildBranch,
	}, api.ServerOptions{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to start grafana. error: %s\n", err.Error())
		return err
	}

	ctx := context.Background()

	go listenToSystemSignals(ctx, s)

	if err := s.Run(); err != nil {
		code := s.ExitCode(err)
		return exitWithCode{
			reason: err.Error(),
			code:   code,
		}
	}

	return nil
}

func validPackaging(packaging string) string {
	validTypes := []string{"dev", "deb", "rpm", "docker", "brew", "hosted", "unknown"}
	for _, vt := range validTypes {
		if packaging == vt {
			return packaging
		}
	}
	return "unknown"
}

func listenToSystemSignals(ctx context.Context, s *server.Server) {
	signalChan := make(chan os.Signal, 1)
	sighupChan := make(chan os.Signal, 1)

	signal.Notify(sighupChan, syscall.SIGHUP)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-sighupChan:
			if err := log.Reload(); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to reload loggers: %s\n", err)
			}
		case sig := <-signalChan:
			ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
			defer cancel()
			if err := s.Shutdown(ctx, fmt.Sprintf("System signal: %s", sig)); err != nil {
				fmt.Fprintf(os.Stderr, "Timed out waiting for server to shut down\n")
			}
			return
		}
	}
}
