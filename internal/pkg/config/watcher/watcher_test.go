package watcher

import (
	"context"
	"io/ioutil"
	"testing"

	_ "github.com/tumblr/k8s-sidecar-injector/internal/pkg/testing"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
)

var (
	testConfig = Config{
		Namespace: "default",
		ConfigMapLabels: map[string]string{
			"thing": "fake",
		},
	}
)

func TestGet(t *testing.T) {
	data, err := ioutil.ReadFile("test/fixtures/sidecars/sidecar-test.yaml")
	if err != nil {
		t.Fatal(err)
	}

	client := fake.NewSimpleClientset(
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ok",
				Namespace: "default",
				Labels:    testConfig.ConfigMapLabels,
			},
			Data: map[string]string{
				"sidecar-test": string(data),
			},
		},
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "ok1",
				Namespace: "default",
				Labels:    testConfig.ConfigMapLabels,
			},
			Data: map[string]string{
				"sidecar-test": string(data),
			},
		},
		// This configmap shouldn't be retrieved because it's missing
		// the label selector
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nolabels",
				Namespace: "default",
			},
			Data: map[string]string{
				"sidecar-test": string(data),
			},
		},
		// This configmap shouldn't be retrieved because it doesn't have
		// any valid data
		&v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "nodata",
				Namespace: "default",
				Labels:    testConfig.ConfigMapLabels,
			},
		},
	)
	sharedInformer, err := newSharedInformer(client, testConfig.Namespace, testConfig.ConfigMapLabels)
	if err != nil {
		t.Fatal(err)
	}
	w := K8sConfigMapWatcher{
		Config:         testConfig,
		sharedInformer: sharedInformer,
	}

	ctx := context.Background()

	go w.Watch(ctx, make(chan interface{}, 10))
	if !w.WaitForCacheSync(ctx) {
		t.Fatalf("expected cache to sync")
	}

	messages, err := w.Get(ctx)
	if err != nil {
		t.Fatal(err.Error())
	}
	if len(messages) != 2 {
		t.Fatalf("expected 2 messages, but got %d", len(messages))
	}
}
