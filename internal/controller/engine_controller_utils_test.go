/*
Copyright Coraza Kubernetes Operator contributors.

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
	"testing"

	"github.com/stretchr/testify/assert"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	wafv1alpha1 "github.com/networking-incubator/coraza-kubernetes-operator/api/v1alpha1"
)

func TestEngineMatchesLabels(t *testing.T) {
	podLabels := map[string]string{
		"app":                                    "gateway",
		"gateway.networking.k8s.io/gateway-name": "my-gw",
	}

	t.Run("nil driver returns false", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: nil},
		}
		assert.False(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("nil istio returns false", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{Istio: nil}},
		}
		assert.False(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("nil wasm returns false", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{Wasm: nil},
			}},
		}
		assert.False(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("nil workload selector returns false", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{WorkloadSelector: nil},
				},
			}},
		}
		assert.False(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("matching labels returns true", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "gateway"},
						},
					},
				},
			}},
		}
		assert.True(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("non-matching labels returns false", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "other"},
						},
					},
				},
			}},
		}
		assert.False(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("subset of pod labels still matches", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{
								"app":                                    "gateway",
								"gateway.networking.k8s.io/gateway-name": "my-gw",
							},
						},
					},
				},
			}},
		}
		assert.True(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("selector requires label pod does not have", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "gateway", "extra": "label"},
						},
					},
				},
			}},
		}
		assert.False(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("empty selector matches everything", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{},
					},
				},
			}},
		}
		assert.True(t, engineMatchesLabels(engine, podLabels))
	})

	t.Run("nil pod labels with non-empty selector returns false", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{
							MatchLabels: map[string]string{"app": "gateway"},
						},
					},
				},
			}},
		}
		assert.False(t, engineMatchesLabels(engine, nil))
	})

	t.Run("matchExpressions selector works", func(t *testing.T) {
		engine := &wafv1alpha1.Engine{
			Spec: wafv1alpha1.EngineSpec{Driver: &wafv1alpha1.DriverConfig{
				Istio: &wafv1alpha1.IstioDriverConfig{
					Wasm: &wafv1alpha1.IstioWasmConfig{
						WorkloadSelector: &metav1.LabelSelector{
							MatchExpressions: []metav1.LabelSelectorRequirement{
								{Key: "app", Operator: metav1.LabelSelectorOpIn, Values: []string{"gateway", "proxy"}},
							},
						},
					},
				},
			}},
		}
		assert.True(t, engineMatchesLabels(engine, podLabels))
	})
}
