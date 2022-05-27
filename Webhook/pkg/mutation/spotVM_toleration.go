package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	"reflect"
)

// minLifespanTolerations is a container for mininum lifespan mutation
type minLifespanTolerations struct {
	Logger logrus.FieldLogger
}

// minLifespanTolerations implements the podMutator interface
var _ podMutator = (*minLifespanTolerations)(nil)

// Name returns the minLifespanTolerations short name
func (mpl minLifespanTolerations) Name() string {
	return "min_lifespan"
}

// Mutate returns a new mutated pod according to lifespan tolerations rules
func (mpl minLifespanTolerations) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {

	mpl.Logger = mpl.Logger.WithField("mutation", mpl.Name())
	mpod := pod.DeepCopy()

	//todo possible issues with multiple taints
	mpl.Logger.WithField("min_lifespan", 0).
		Printf("no lifespan label found, applying default lifespan toleration")

	/*
			//https://docs.microsoft.com/en-us/azure/aks/spot-node-pool
		  - key: "kubernetes.azure.com/scalesetpriority"
		    operator: "Equal"
		    value: "spot"
		    effect: "NoSchedule"
	*/
	tn := []corev1.Toleration{{
		Key:      "kubernetes.azure.com/scalesetpriority",
		Operator: corev1.TolerationOpEqual,
		Value:    "spot",
		Effect:   corev1.TaintEffectNoSchedule,
	}}

	mpod.Spec.Tolerations = appendTolerations(tn, mpod.Spec.Tolerations)
	return mpod, nil

}

// appendTolerations appends existing to new without duplicating any tolerations
func appendTolerations(new, existing []corev1.Toleration) []corev1.Toleration {
	var toAppend []corev1.Toleration

	for _, n := range new {
		found := false
		for _, e := range existing {
			if reflect.DeepEqual(n, e) {
				found = true
			}
		}
		if !found {
			toAppend = append(toAppend, n)
		}
	}

	return append(existing, toAppend...)
}
