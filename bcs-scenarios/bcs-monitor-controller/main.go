/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	// Import all Kubernetes client auth plugins (e.g. Azure, GCP, OIDC, etc.)
	// to ensure that exec-entrypoint and run can make use of them.
	_ "k8s.io/client-go/plugin/pkg/client/auth"

	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/Tencent/bk-bcs/bcs-common/common/blog"
	"github.com/Tencent/bk-bcs/bcs-common/common/http/httpserver"
	monitorextensionv1 "github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/api/v1"
	"github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/pkg/apiclient"
	"github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/pkg/controllers"
	"github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/pkg/fileoperator"
	"github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/pkg/httpsvr"
	"github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/pkg/option"
	"github.com/Tencent/bk-bcs/bcs-runtime/bcs-k8s/bcs-component/bcs-monitor-controller/pkg/render"
)

var (
	scheme   = runtime.NewScheme()
	setupLog = ctrl.Log.WithName("setup")
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	utilruntime.Must(monitorextensionv1.AddToScheme(scheme))

	utilruntime.Must(corev1.AddToScheme(scheme))
	// +kubebuilder:scaffold:scheme
}

func main() {
	opts := &option.ControllerOption{}

	var metricsAddr string
	var enableLeaderElection bool
	var probeAddr string
	var verbosity int
	var scenarioRefreshFreqSec int64
	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")

	// log config
	flag.IntVar(&verbosity, "v", 3, "log level for V logs")
	flag.StringVar(&opts.LogDir, "log_dir", "./logs", "If non-empty, write log files in this directory")
	flag.Uint64Var(&opts.LogMaxSize, "log_max_size", 500, "Max size (MB) per log file.")
	flag.IntVar(&opts.LogMaxNum, "log_max_num", 10, "Max num of log file.")
	flag.BoolVar(&opts.ToStdErr, "logtostderr", false, "log to standard error instead of files")
	flag.BoolVar(&opts.AlsoToStdErr, "alsologtostderr", false, "log to standard error as well as files")

	flag.StringVar(&opts.Address, "address", "0.0.0.0", "address for controller")
	flag.StringVar(&opts.ScenarioPath, "scenario_path", "/data/bcs", "Store scenario templates")
	flag.Int64Var(&scenarioRefreshFreqSec, "scenario_refresh_req", 60, "refresh frequency ")
	flag.UintVar(&opts.HttpServerPort, "http_server_port", 8088, "http server port")
	opts.ScenarioGitRefreshFreq = time.Second * time.Duration(scenarioRefreshFreqSec)
	opts.Verbosity = int32(verbosity)
	flag.Parse()

	blog.InitLogs(opts.LogConfig)
	defer blog.CloseLogs()
	ctrl.SetLogger(zap.New(zap.UseDevMode(false)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       "333fb49e.monitorextension.bkbcs.tencent.com",
		// LeaderElectionReleaseOnCancel defines if the leader should step down voluntarily
		// when the Manager ends. This requires the binary to immediately end when the
		// Manager is stopped, otherwise, this setting is unsafe. Setting this significantly
		// speeds up voluntary leader transitions as the new leader don't have to wait
		// LeaseDuration time first.
		//
		// In the default scaffold provided, the program ends immediately after
		// the manager stops, so would be fine to enable this option. However,
		// if you are doing or is intended to do any operation such as perform cleanups
		// after the manager stops then its usage might be unsafe.
		// LeaderElectionReleaseOnCancel: true,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	ctx := context.Background()
	monitorCli := apiclient.NewBkmApiClient()
	fileOp := fileoperator.NewFileOperator(mgr.GetClient())
	if err = (&controllers.MonitorRuleReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),

		Ctx:           ctx,
		FileOp:        fileOp,
		MonitorApiCli: monitorCli,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "MonitorRule")
		os.Exit(1)
	}
	if err = (&controllers.NoticeGroupReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),

		Ctx:           ctx,
		FileOp:        fileOp,
		MonitorApiCli: monitorCli,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "NoticeGroup")
		os.Exit(1)
	}
	if err = (&controllers.PanelReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),

		Ctx:           ctx,
		FileOp:        fileOp,
		MonitorApiCli: monitorCli,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Panel")
		os.Exit(1)
	}

	monitorRender, err := render.NewMonitorRender(scheme, mgr.GetClient(), opts)
	if err != nil {
		blog.Errorf("new render failed, err: %v", err)
		os.Exit(1)
	}
	if err = (&controllers.AppMonitorReconciler{
		Client: mgr.GetClient(),
		Scheme: mgr.GetScheme(),

		Ctx:    ctx,
		Render: monitorRender,
	}).SetupWithManager(mgr); err != nil {
		setupLog.Error(err, "unable to create controller", "controller", "Panel")
		os.Exit(1)
	}
	// +kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	err = initHttpServer(opts, mgr)
	if err != nil {
		blog.Errorf("init http server failed: %v", err.Error())
		os.Exit(1)
	}
	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

// initHttpServer init ingress controller http server
func initHttpServer(op *option.ControllerOption, mgr manager.Manager) error {
	server := httpserver.NewHttpServer(op.HttpServerPort, op.Address, "")

	// server.SetInsecureServer(op.Conf.InsecureAddress, op.Conf.InsecurePort)
	server.SetInsecureServer(op.Address, op.HttpServerPort)
	ws := server.NewWebService("", nil)
	httpServerClient := &httpsvr.HttpServerClient{
		Mgr: mgr,
	}
	httpsvr.InitRouters(ws, httpServerClient)

	router := server.GetRouter()
	webContainer := server.GetWebContainer()
	router.Handle("/{sub_path:.*}", webContainer)
	if err := server.ListenAndServeMux(false); err != nil {
		return fmt.Errorf("http ListenAndServe error %w", err)
	}

	return nil
}