/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"slices"

	"k8s.io/apimachinery/pkg/runtime"
	types "k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	batchv1 "github.com/bonavadeur/seika/api/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/bonavadeur/seika/internal/bonalib"
)

var _ = bonalib.Baka()

// SeikaReconciler reconciles a Seika object
type SeikaReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.bonavadeur.io,resources=seikas,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.bonavadeur.io,resources=seikas/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.bonavadeur.io,resources=seikas/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods/status,verbs=get
//+kubebuilder:rbac:groups=core,resources=nodes,verbs=get;list;
//+kubebuilder:rbac:groups=core,resources=nodes/status,verbs=get

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Seika object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.0/pkg/reconcile
func (r *SeikaReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	// _ = log.FromContext(ctx)

	seika := &batchv1.Seika{}
	if err := r.Get(ctx, req.NamespacedName, seika); err != nil {
		bonalib.Warn("unable to fetch seika", err)
		return ctrl.Result{}, nil
	}

	nodeList := &corev1.NodeList{}
	if err := r.List(ctx, nodeList); err != nil {
		bonalib.Warn("unable to fetch nodeList", err)
		return ctrl.Result{}, nil
	}
	checkNodeRepurika(seika.Spec.Repurika, nodeList.Items)

	podCount := map[string]int{}
	podNames := map[string][]string{}
	for _, node := range nodeList.Items {
		podList := &corev1.PodList{}
		listOpts := []client.ListOption{
			client.InNamespace(seika.Namespace),
			client.MatchingLabels(seika.Spec.Selector.MatchLabels),
			client.MatchingFields{"spec.nodeName": node.Name},
		}
		if err := r.List(ctx, podList, listOpts...); err != nil {
			bonalib.Warn("Failed to list pods on node", node.Name, err)
			return ctrl.Result{}, err
		}
		podCount[node.Name] = len(podList.Items)
		for _, pod := range podList.Items {
			podNames[node.Name] = append(podNames[node.Name], pod.Name)
		}
	}
	bonalib.Succ("Reconcile", seika.Name, podCount)

	repurika := seika.Spec.Repurika
	for node, size := range repurika {
		countDiff := int(math.Abs(float64(podCount[node] - int(size))))
		if podCount[node] < size {
			bonalib.Log("need to create more pod in", node)
			for i := 0; i < countDiff; i++ {
				pod, err := r.createPodTemplate(seika, node)
				if err != nil {
					bonalib.Warn("Failed to define new PodTemplate", err)
					return ctrl.Result{}, nil
				}
				if err := r.Create(ctx, pod); err != nil {
					bonalib.Warn("Failed to create new Pod", err)
					return ctrl.Result{}, err
				}
				if err := r.Status().Update(ctx, seika); err != nil {
					bonalib.Warn("Failed to update Repurika Status", err)
					return ctrl.Result{}, nil
				}
			}
			return ctrl.Result{Requeue: true}, nil
		} else if podCount[node] > size {
			bonalib.Log("need to delete less pod in", node)
			for i := 0; i < countDiff; i++ {
				pod := &corev1.Pod{}
				podTobeDeletedName := podNames[node][i]
				if err := r.Get(ctx, types.NamespacedName{Name: podTobeDeletedName, Namespace: seika.Namespace}, pod); err != nil {
					bonalib.Warn("Pod doesnt exist", podTobeDeletedName, err)
					return ctrl.Result{}, err
				}
				if err := r.Delete(ctx, pod, client.GracePeriodSeconds(0)); err != nil {
					bonalib.Warn("Failed to delete the pod", podTobeDeletedName, err)
					return ctrl.Result{}, err
				}
				if err := r.Status().Update(ctx, seika); err != nil {
					bonalib.Warn("Failed to update Seika status", err)
					return ctrl.Result{}, err
				}
			}
			return ctrl.Result{Requeue: true}, nil
		}
	}

	seika.Status.Repurika = seika.Spec.Repurika
	if err := r.Status().Update(ctx, seika); err != nil {
		bonalib.Warn("Failed to update Repurika status", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func generatePodName() string {
	letters := []byte("1234567890abcdefghijklmnopqrstuvwxyz")
	ranStr := make([]byte, 5)
	for i := 0; i < 5; i++ {
		ranStr[i] = letters[rand.Intn(len(letters))]
	}
	str := string(ranStr)
	return str
}

func (r *SeikaReconciler) createPodTemplate(seika *batchv1.Seika, node string) (*corev1.Pod, error) {
	labels := map[string]string{}
	annotations := map[string]string{}
	for k, v := range seika.Spec.Selector.MatchLabels {
		labels[k] = v
	}
	for k, v := range seika.Spec.Template.Labels {
		labels[k] = v
	}
	labels["bonavadeur.io/seika-hostname"] = node

	for k, v := range seika.Spec.Template.Annotations {
		annotations[k] = v
	}

	spec := seika.Spec.Template.Spec
	spec.NodeSelector = map[string]string{
		"kubernetes.io/hostname": node,
	}

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:        fmt.Sprintf("%v-%v-%v", seika.Name, node, generatePodName()),
			Namespace:   seika.Namespace,
			Labels:      labels,
			Annotations: annotations,
		},
		Spec: spec,
	}

	if err := ctrl.SetControllerReference(seika, pod, r.Scheme); err != nil {
		return nil, err
	}

	return pod, nil
}

func checkNodeRepurika(repurika map[string]int, nodes []corev1.Node) {
	nodeNames := []string{}
	for _, node := range nodes {
		nodeNames = append(nodeNames, node.Name)
	}
	for node := range repurika {
		if slices.Contains(nodeNames, node) {
			continue
		} else {
			bonalib.Warn("Node not found: ", node)
			return
		}
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *SeikaReconciler) SetupWithManager(mgr ctrl.Manager) error {
	bonalib.Succ("SetupWithManager")

	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &corev1.Pod{}, "spec.nodeName", func(rawObj client.Object) []string {
		pod := rawObj.(*corev1.Pod)
		return []string{pod.Spec.NodeName}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Seika{}).
		Owns(&corev1.Pod{}).
		Complete(r)
}
