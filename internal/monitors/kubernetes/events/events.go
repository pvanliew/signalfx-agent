package events

import (
	"strings"
	"time"

	"github.com/signalfx/golib/event"
	"github.com/signalfx/signalfx-agent/internal/core/common/kubernetes"
	"github.com/signalfx/signalfx-agent/internal/core/config"
	"github.com/signalfx/signalfx-agent/internal/monitors"
	"github.com/signalfx/signalfx-agent/internal/monitors/kubernetes/leadership"
	"github.com/signalfx/signalfx-agent/internal/monitors/types"
	"github.com/signalfx/signalfx-agent/internal/utils"

	log "github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	k8s "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
)

const monitorType = "kubernetes-events"

// MONITOR(kubernetes-events): This monitor sends Kubernetes events as SignalFx
// events.  Upon startup, it will send all of the events that K8s has that are
// still persisted and then send any new events that come in.  The various
// agents perform leader election amongs themselves to decide which instance
// will send events, unless the `alwaysClusterReporter` config option is set to
// true.
//
// You can see the types of events happening in your cluster with `kubectl get
// events -o yaml --all-namespaces`.

var logger = log.WithFields(log.Fields{"monitorType": monitorType})

func init() {
	monitors.Register(monitorType, func() interface{} { return &Monitor{} }, &Config{})
}

// EventInclusionSpec specifies a type of event to send
type EventInclusionSpec struct {
	Reason             string `yaml:"reason"`
	InvolvedObjectKind string `yaml:"involvedObjectKind"`
}

// Config for the K8s event monitor
type Config struct {
	config.MonitorConfig
	// Configuration of the Kubernetes API client
	KubernetesAPI *kubernetes.APIConfig `yaml:"kubernetesAPI" default:"{}"`
	// A list of event types to send events for.  Only events matching these
	// items will be sent.
	WhitelistedEvents []EventInclusionSpec `yaml:"whitelistedEvents"`
	// If true, all events from Kubernetes will be sent.  Please don't use this
	// option unless you really want to act on all possible K8s events.
	SendAllEvents bool `yaml:"_sendAllEvents"`
	// Whether to always send events from this agent instance or to do leader
	// election to only send from one agent instance.
	AlwaysClusterReporter bool `yaml:"alwaysClusterReporter"`
}

// Monitor for K8s Cluster Metrics.  Also handles syncing certain properties
// about pods.
type Monitor struct {
	Output        types.Output
	stopper       chan struct{}
	sendAllEvents bool
	whitelistSet  map[EventInclusionSpec]bool
}

// Configure the monitor and kick off event syncing
func (m *Monitor) Configure(conf *Config) error {
	k8sClient, err := kubernetes.MakeClient(conf.KubernetesAPI)
	if err != nil {
		return err
	}

	m.sendAllEvents = conf.SendAllEvents
	m.whitelistSet = make(map[EventInclusionSpec]bool, len(conf.WhitelistedEvents))
	for i := range conf.WhitelistedEvents {
		spec := conf.WhitelistedEvents[i]
		spec.InvolvedObjectKind = strings.ToLower(spec.InvolvedObjectKind)
		spec.Reason = strings.ToLower(spec.Reason)
		m.whitelistSet[spec] = true
	}

	m.stopper = make(chan struct{})

	return m.start(k8sClient, conf.AlwaysClusterReporter)
}

func (m *Monitor) start(k8sClient *k8s.Clientset, alwaysReport bool) error {
	var syncStopper chan struct{}

	runSync := func() {
		syncStopper = make(chan struct{})
		syncEvents(k8sClient, cache.ResourceEventHandlerFuncs{
			AddFunc: func(obj interface{}) {
				ev := obj.(*v1.Event)
				m.handleNewEvent(ev)
			},
		}, syncStopper)
	}

	var leaderCh <-chan bool
	var unregister func()
	if alwaysReport {
		logger.Info("This instance will send K8s events")
		runSync()
	} else {
		var err error
		leaderCh, unregister, err = leadership.RequestLeaderNotification(k8sClient.CoreV1())
		if err != nil {
			return err
		}
	}

	go func() {
		for {
			select {
			case isLeader := <-leaderCh:
				if isLeader {
					logger.Info("This instance is now the leader and will send events")
					runSync()
				} else {
					logger.Info("No longer leader")
					close(syncStopper)
					syncStopper = nil
				}
			case <-m.stopper:
				logger.Info("Stopping k8s event syncing")
				if unregister != nil {
					unregister()
				}
				if syncStopper != nil {
					close(syncStopper)
				}
				return
			}
		}
	}()
	return nil
}

func (m *Monitor) shouldSendEvent(ev *v1.Event) bool {
	// Always ignore any events older than 1 minute so we don't cause an event
	// flood upon agent restarts.  This doesn't eliminate the possibility of
	// duplicated events but should limit them to a fairly narrow time window.
	if ev.LastTimestamp.Time.Before(time.Now().Add(-1 * time.Minute)) {
		return false
	}

	if m.sendAllEvents {
		return true
	}

	return m.whitelistSet[EventInclusionSpec{
		Reason:             strings.ToLower(ev.Reason),
		InvolvedObjectKind: strings.ToLower(ev.InvolvedObject.Kind),
	}]
}

func (m *Monitor) handleNewEvent(ev *v1.Event) {
	if m.shouldSendEvent(ev) {
		sfxEvent := k8sEventToSignalFxEvent(ev)
		m.Output.SendEvent(sfxEvent)
	}
}

func k8sEventToSignalFxEvent(ev *v1.Event) *event.Event {
	dims := map[string]string{
		"kubernetes_kind":      ev.InvolvedObject.Kind,
		"kubernetes_namespace": ev.InvolvedObject.Namespace,
		"obj_field_path":       ev.InvolvedObject.FieldPath,
	}

	// Reuse the existing pod dimensions that we send for metrics
	if ev.InvolvedObject.Kind == "Pod" {
		dims["kubernetes_pod_name"] = ev.InvolvedObject.Name
		dims["kubernetes_pod_uid"] = string(ev.InvolvedObject.UID)
	} else {
		dims["kubernetes_name"] = ev.InvolvedObject.Name
		dims["kubernetes_uid"] = string(ev.InvolvedObject.UID)
	}

	properties := utils.RemoveEmptyMapValues(map[string]string{
		"message":                     ev.Message,
		"source_component":            ev.Source.Component,
		"source_host":                 ev.Source.Host,
		"kubernetes_event_type":       ev.Type,
		"kubernetes_resource_version": ev.InvolvedObject.ResourceVersion,
	})

	return event.NewWithProperties(
		ev.Reason,
		event.AGENT,
		utils.RemoveEmptyMapValues(dims),
		utils.StringMapToInterfaceMap(properties),
		ev.LastTimestamp.Time)
}

func syncEvents(clientset *k8s.Clientset, handlers cache.ResourceEventHandlerFuncs, stopper chan struct{}) {
	client := clientset.CoreV1().RESTClient()
	watchList := cache.NewListWatchFromClient(client, "events", v1.NamespaceAll, fields.Everything())

	_, controller := cache.NewInformer(watchList, &v1.Event{}, 0, handlers)

	go controller.Run(stopper)
}

// Shutdown the monitor and stop any syncing
func (m *Monitor) Shutdown() {
	close(m.stopper)
}
