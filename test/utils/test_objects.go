/*
Copyright 2018 Pusher Ltd. and Wave Contributors

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

package utils

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var labels = map[string]string{
	"app": "example",
}

var annotations = map[string]string{
	"wave.pusher.com/extra-configmaps": "ns1/test-cm1,ns2/test-cm2,local-cm1",
	"wave.pusher.com/extra-secrets":    "ns1/test-secret1,ns2/test-secret2,local-secret1",
}

var trueValue = true

var podTemplate = &corev1.PodTemplateSpec{
	ObjectMeta: metav1.ObjectMeta{
		Labels: labels,
	},
	Spec: corev1.PodSpec{
		SchedulerName: "default-scheduler",
		Volumes: []corev1.Volume{
			{
				Name: "secret1",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "example1",
					},
				},
			},
			{
				Name: "secret-optional",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "volume-optional",
						Optional:   &trueValue,
					},
				},
			},
			{
				Name: "configmap1",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "example1",
						},
					},
				},
			},
			{
				Name: "configmap-optional",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: "volume-optional",
						},
						Optional: &trueValue,
					},
				},
			},
			{
				Name: "projection1",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{{
							Secret: &corev1.SecretProjection{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example5",
								},
								Optional: &trueValue,
							},
						}},
					},
				},
			},
			{
				Name: "projection2",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{{
							ConfigMap: &corev1.ConfigMapProjection{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example5",
								},
								Optional: &trueValue,
							},
						}},
					},
				},
			},
			{
				Name: "projection3",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{{
							Secret: &corev1.SecretProjection{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example6",
								},
								Items: []corev1.KeyToPath{
									{
										Key:  "key1",
										Path: "example6_key1.txt",
									},
									{
										Key:  "key3",
										Path: "example6_key3.txt",
									},
								},
							},
						}},
					},
				},
			},
			{
				Name: "projection4",
				VolumeSource: corev1.VolumeSource{
					Projected: &corev1.ProjectedVolumeSource{
						Sources: []corev1.VolumeProjection{{
							ConfigMap: &corev1.ConfigMapProjection{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example6",
								},
								Items: []corev1.KeyToPath{
									{
										Key:  "key1",
										Path: "example6_key1.txt",
									},
									{
										Key:  "key3",
										Path: "example6_key3.txt",
									},
								},
							},
						}},
					},
				},
			},
		},
		Containers: []corev1.Container{
			{
				Name:  "container1",
				Image: "container1",
				Env: []corev1.EnvVar{
					{
						Name: "example1_key1",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example1",
								},
								Key: "key1",
							},
						},
					},
					{
						Name: "example1_key1_new_name",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example1",
								},
								Key: "key1",
							},
						},
					},
					{
						Name: "example3_key1",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example3",
								},
								Key: "key1",
							},
						},
					},
					{
						Name: "example3_key4",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example3",
								},
								Key:      "key4",
								Optional: &trueValue,
							},
						},
					},
					{
						Name: "example4_key1",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example4",
								},
								Key:      "key1",
								Optional: &trueValue,
							},
						},
					},
					{
						Name: "example1_secret_key1",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example1",
								},
								Key: "key1",
							},
						},
					},
					{
						Name: "example3_secret_key1",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example3",
								},
								Key: "key1",
							},
						},
					},
					{
						Name: "example3_secret_key4",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example3",
								},
								Key:      "key4",
								Optional: &trueValue,
							},
						},
					},
					{
						Name: "example4_secret_key1",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example4",
								},
								Key:      "key1",
								Optional: &trueValue,
							},
						},
					},
				},
				EnvFrom: []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "example1",
							},
						},
					},
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "envfrom-optional",
							},
							Optional: &trueValue,
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "example1",
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "envfrom-optional",
							},
							Optional: &trueValue,
						},
					},
				},
			},
			{
				Name:  "container2",
				Image: "container2",
				Env: []corev1.EnvVar{
					{
						Name: "env_optional_key2",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "env-optional",
								},
								Key:      "key2",
								Optional: &trueValue,
							},
						},
					},
					{
						Name: "example3_key2",
						ValueFrom: &corev1.EnvVarSource{
							ConfigMapKeyRef: &corev1.ConfigMapKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example3",
								},
								Key: "key2",
							},
						},
					},
					{
						Name: "example3_secret_key2",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "example3",
								},
								Key: "key2",
							},
						},
					},
					{
						Name: "env_optional_secret_key2",
						ValueFrom: &corev1.EnvVarSource{
							SecretKeyRef: &corev1.SecretKeySelector{
								LocalObjectReference: corev1.LocalObjectReference{
									Name: "env-optional",
								},
								Key:      "key2",
								Optional: &trueValue,
							},
						},
					},
				},
				EnvFrom: []corev1.EnvFromSource{
					{
						ConfigMapRef: &corev1.ConfigMapEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "example2",
							},
						},
					},
					{
						SecretRef: &corev1.SecretEnvSource{
							LocalObjectReference: corev1.LocalObjectReference{
								Name: "example2",
							},
						},
					},
				},
			},
		},
	},
}

// ExampleDeployment is an example Deployment object for use within test suites
var ExampleDeployment = &appsv1.Deployment{
	ObjectMeta: metav1.ObjectMeta{
		Name:        "example",
		Namespace:   "default",
		Labels:      labels,
		Annotations: annotations,
	},
	Spec: appsv1.DeploymentSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: *podTemplate,
	},
}

// ExampleStatefulSet is an example StatefulSet object for use within test suites
var ExampleStatefulSet = &appsv1.StatefulSet{
	ObjectMeta: metav1.ObjectMeta{
		Name:        "example",
		Namespace:   "default",
		Labels:      labels,
		Annotations: annotations,
	},
	Spec: appsv1.StatefulSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: *podTemplate,
	},
}

// ExampleDaemonSet is an example DaemonSet object for use within test suites
var ExampleDaemonSet = &appsv1.DaemonSet{
	ObjectMeta: metav1.ObjectMeta{
		Name:        "example",
		Namespace:   "default",
		Labels:      labels,
		Annotations: annotations,
	},
	Spec: appsv1.DaemonSetSpec{
		Selector: &metav1.LabelSelector{
			MatchLabels: labels,
		},
		Template: *podTemplate,
	},
}

// ExampleConfigMap1 is an example ConfigMap object for use within test suites
var ExampleConfigMap1 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example1",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example1:key1",
		"key2": "example1:key2",
		"key3": "example1:key3",
	},
	BinaryData: map[string][]byte{
		"binary_key1": []byte("example1:binary_key1"),
		"binary_key2": []byte("example1:binary_key2"),
		"binary_key3": []byte("example1:binary_key3"),
	},
}

// ExampleConfigMap2 is an example ConfigMap object for use within test suites
var ExampleConfigMap2 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example2",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example2:key1",
		"key2": "example2:key2",
		"key3": "example2:key3",
	},
}

// ExampleConfigMap3 is an example ConfigMap object for use within test suites
var ExampleConfigMap3 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example3",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example3:key1",
		"key2": "example3:key2",
		"key3": "example3:key3",
	},
}

// ExampleConfigMap4 is an example ConfigMap object for use within test suites
var ExampleConfigMap4 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example4",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example4:key1",
		"key2": "example4:key2",
		"key3": "example4:key3",
	},
}

// ExampleConfigMap5 is an example ConfigMap object for use within test suites
var ExampleConfigMap5 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example5",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example5:key1",
		"key2": "example5:key2",
		"key3": "example5:key3",
	},
}

// ExampleConfigMap6 is an example ConfigMap object for use within test suites
var ExampleConfigMap6 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example6",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example6:key1",
		"key2": "example6:key2",
		"key3": "example6:key3",
	},
}

// ExampleConfigMap6 is an example ConfigMap object for use within test suites
var ExampleConfigMap6WithoutKey3 = &corev1.ConfigMap{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example6",
		Namespace: "default",
		Labels:    labels,
	},
	Data: map[string]string{
		"key1": "example6:key1",
		"key2": "example6:key2",
	},
}

// ExampleSecret1 is an example Secret object for use within test suites
var ExampleSecret1 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example1",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example1:key1",
		"key2": "example1:key2",
		"key3": "example1:key3",
	},
}

// ExampleSecret2 is an example Secret object for use within test suites
var ExampleSecret2 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example2",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example2:key1",
		"key2": "example2:key2",
		"key3": "example2:key3",
	},
}

// ExampleSecret3 is an example Secret object for use within test suites
var ExampleSecret3 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example3",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example3:key1",
		"key2": "example3:key2",
		"key3": "example3:key3",
	},
}

// ExampleSecret4 is an example Secret object for use within test suites
var ExampleSecret4 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example4",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example4:key1",
		"key2": "example4:key2",
		"key3": "example4:key3",
	},
}

// ExampleSecret5 is an example Secret object for use within test suites
var ExampleSecret5 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example5",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example5:key1",
		"key2": "example5:key2",
		"key3": "example5:key3",
	},
}

// ExampleSecret6 is an example Secret object for use within test suites
var ExampleSecret6 = &corev1.Secret{
	ObjectMeta: metav1.ObjectMeta{
		Name:      "example6",
		Namespace: "default",
		Labels:    labels,
	},
	StringData: map[string]string{
		"key1": "example6:key1",
		"key2": "example6:key2",
		"key3": "example6:key3",
	},
}
