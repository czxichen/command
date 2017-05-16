package watchdog

import (
	"fmt"
	"time"

	log "github.com/golang/glog"
)

var (
	restartDelay      = 2 * time.Second
	restartBackoff    = 5 * time.Second
	restartBackoffMax = 60 * time.Second
)

type Watchdog struct {
	services map[string]*Service
	shutdown chan bool
}

func NewWatchdog() *Watchdog {
	return &Watchdog{
		services: make(map[string]*Service),
		shutdown: make(chan bool),
	}
}

//关闭服务
func (w *Watchdog) Shutdown() {
	select {
	case w.shutdown <- true:
	default:
	}
}

//添加服务,如果存在
func (w *Watchdog) AddService(name, binary string) (*Service, error) {
	if _, ok := w.services[name]; ok {
		return nil, fmt.Errorf("Service %q already exists", name)
	}

	svc := newService(name, binary)
	w.services[name] = svc

	return svc, nil
}

//启动服务
func (w *Watchdog) Walk() {
	log.Info("Seesaw watchdog starting...")

	w.mapDependencies()

	for _, svc := range w.services {
		go svc.run()
	}
	<-w.shutdown
	for _, svc := range w.services {
		go svc.stop()
	}
	for _, svc := range w.services {
		stopped := <-svc.stopped
		svc.stopped <- stopped
	}
}

//设置依赖关系
func (w *Watchdog) mapDependencies() {
	for name := range w.services {
		svc := w.services[name]
		for depName := range svc.dependencies {
			dep, ok := w.services[depName]
			if !ok {
				log.Fatalf("Failed to find dependency %q for service %q", depName, name)
			}
			svc.dependencies[depName] = dep //依赖谁
			dep.dependents[svc.name] = svc  //谁依赖它
		}
	}
}
