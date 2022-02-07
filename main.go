/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"fmt"
	"os"

	api "github.com/test-project/pkg"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

var (
	setupLog = ctrl.Log.WithName("setup")
)

type reconciler struct {
	client.Client
	scheme *runtime.Scheme
}

func (r *reconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithValues("chaospod", req.NamespacedName)
	log.V(1).Info("reconciling chaos pod")

	log.Info(fmt.Sprintf("%+v\n\n%+v\n", ctx, req))

	var chaospod api.ChaosPod
	if err := r.Get(ctx, req.NamespacedName, &chaospod); err != nil {
		log.Error(err, "unable to get chaosctl")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func main() {

	sch := runtime.NewScheme()
	err := corev1.AddToScheme(sch)
	if err != nil {
		fmt.Errorf("err while adding to scheme %w", err)
		os.Exit(1)
	}

	cfg, err := ctrl.GetConfig()
	if err != nil {
		fmt.Errorf("err while getting config %w", err)
	}

	cl, err := client.New(cfg, client.Options{Scheme: sch})
	if err != nil {
		fmt.Errorf("err while creating client %w", err)
	}

	// create a namespace
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-2",
			ClusterName: "admin_org",
		},
		TypeMeta: metav1.TypeMeta{
			Kind:       "Namespace",
			APIVersion: "v1",
		},
	}

	err = cl.Create(context.TODO(), ns)
	if err != nil {
		fmt.Println("error while creating namespace")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	list := &corev1.NamespaceList{}
	err = cl.List(context.Background(), list)
	if err != nil {
		fmt.Errorf("error listing namespaces %w", err)
	}

	for _, l := range list.Items {
		fmt.Println(l.Name)
	}
	fmt.Println("created namespace")

	// create a configMap
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-cm",
			Namespace:   "test-2",
			ClusterName: "admin_org",
		},
		Data: map[string]string{
			"test-key": "test-value",
		},
	}

	fmt.Println("successfully created cm in cluster", cm.ClusterName)
}
