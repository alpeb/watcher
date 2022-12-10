package main

import (
	"flag"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	klog "k8s.io/klog/v2"
)

func updateEndpointSlice(_ interface{}, _ interface{}) {
	fmt.Println("updateEndpointSlice()")
}

func main() {
	level, err := log.ParseLevel("debug")
	if err != nil {
		log.Fatal("invalid log-level")
	}
	log.SetLevel(level)
	flag.Set("stderrthreshold", "INFO")
	flag.Set("logtostderr", "true")
	flag.Set("v", "12") // At 7 and higher, authorization tokens get logged.
	// pipe klog entries to logrus
	klog.SetOutput(log.StandardLogger().Writer())

	config, err := clientcmd.BuildConfigFromFlags("", "")
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	factory := informers.NewSharedInformerFactory(clientset, 10*time.Minute)
	informer := factory.Discovery().V1().EndpointSlices().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		UpdateFunc: updateEndpointSlice,
	})

	stop := make(chan struct{})
	defer close(stop)

	factory.Start(stop)
	if !cache.WaitForCacheSync(stop, informer.HasSynced) {
		panic("Failed to sync")
	}

	select {}
}
