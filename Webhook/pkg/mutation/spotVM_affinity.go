package mutation

import (
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
)

// minLifespanTolerations is a container for mininum lifespan mutation
type affinity struct {
	Logger logrus.FieldLogger
}

// minLifespanTolerations implements the podMutator interface
var _ podMutator = (*affinity)(nil)

// Name returns the minLifespanTolerations short name
func (a affinity) Name() string {
	return "tolerations"
}

// Mutate returns a new mutated pod according to lifespan tolerations rules
func (a affinity) Mutate(pod *corev1.Pod) (*corev1.Pod, error) {

	a.Logger = a.Logger.WithField("mutation", a.Name())
	mpod := pod.DeepCopy()

	spotvmAffinity := corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			RequiredDuringSchedulingIgnoredDuringExecution: &corev1.NodeSelector{
				NodeSelectorTerms: []corev1.NodeSelectorTerm{{
					MatchExpressions: []corev1.NodeSelectorRequirement{{
						Key:      "kubernetes.azure.com/scalesetpriority",
						Operator: "in",
						Values:   []string{"Spot"},
					}},
					MatchFields: nil,
				}}},
			PreferredDuringSchedulingIgnoredDuringExecution: nil,
		},
	}
	//todo this will override affinities. You should merge match expressions
	mpod.Spec.Affinity = &spotvmAffinity
	return mpod, nil
}
