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

package core

import (
	"sync"
	"time"

	metricsserver "sigs.k8s.io/controller-runtime/pkg/metrics/server"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/wave-k8s/wave/test/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

var _ = Describe("Wave hash Suite", func() {
	Context("calculateConfigHash", func() {
		var c client.Client
		var m utils.Matcher

		var mgrStopped *sync.WaitGroup
		var stopMgr chan struct{}

		const timeout = time.Second * 5

		var cm1 *corev1.ConfigMap
		var cm2 *corev1.ConfigMap
		var cm3 *corev1.ConfigMap
		var s1 *corev1.Secret
		var s2 *corev1.Secret
		var s3 *corev1.Secret

		var modified = "modified"

		BeforeEach(func() {
			mgr, err := manager.New(cfg, manager.Options{
				Metrics: metricsserver.Options{
					BindAddress: "0",
				},
			})
			Expect(err).NotTo(HaveOccurred())
			var cerr error
			c, cerr = client.New(cfg, client.Options{Scheme: scheme.Scheme})
			Expect(cerr).NotTo(HaveOccurred())
			m = utils.Matcher{Client: c}

			stopMgr, mgrStopped = StartTestManager(mgr)

			cm1 = utils.ExampleConfigMap1.DeepCopy()
			cm2 = utils.ExampleConfigMap2.DeepCopy()
			cm3 = utils.ExampleConfigMap3.DeepCopy()
			s1 = utils.ExampleSecret1.DeepCopy()
			s2 = utils.ExampleSecret2.DeepCopy()
			s3 = utils.ExampleSecret3.DeepCopy()

			m.Create(cm1).Should(Succeed())
			m.Create(cm2).Should(Succeed())
			m.Create(cm3).Should(Succeed())
			m.Create(s1).Should(Succeed())
			m.Create(s2).Should(Succeed())
			m.Create(s3).Should(Succeed())

			m.Get(cm1, timeout).Should(Succeed())
			m.Get(cm2, timeout).Should(Succeed())
			m.Get(cm3, timeout).Should(Succeed())
			m.Get(s1, timeout).Should(Succeed())
			m.Get(s2, timeout).Should(Succeed())
			m.Get(s3, timeout).Should(Succeed())
		})

		AfterEach(func() {
			close(stopMgr)
			mgrStopped.Wait()

			utils.DeleteAll(cfg, timeout,
				&appsv1.DeploymentList{},
				&corev1.ConfigMapList{},
				&corev1.SecretList{},
			)
		})

		Context("all key children", func() {
			var configMaps map[types.NamespacedName]*corev1.ConfigMap
			var secrets map[types.NamespacedName]*corev1.Secret
			var configMapsConfig []configMetadata
			var secretsConfig []configMetadata

			BeforeEach(func() {

				configMaps = map[types.NamespacedName]*corev1.ConfigMap{
					GetNamespacedNameFromObject(cm1): cm1,
					GetNamespacedNameFromObject(cm2): cm2,
				}

				secrets = map[types.NamespacedName]*corev1.Secret{
					GetNamespacedNameFromObject(s1): s1,
					GetNamespacedNameFromObject(s2): s2,
				}

				configMapsConfig = []configMetadata{
					{name: GetNamespacedNameFromObject(cm1), allKeys: true},
					{name: GetNamespacedNameFromObject(cm2), allKeys: true},
				}

				secretsConfig = []configMetadata{
					{name: GetNamespacedNameFromObject(s1), allKeys: true},
					{name: GetNamespacedNameFromObject(s2), allKeys: true},
				}
			})

			It("returns a different hash when an allKeys child's configmap string data is updated", func() {
				h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				cm1.Data["key1"] = modified

				h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(h2).NotTo(Equal(h1))
			})

			It("returns a different hash when an allKeys child's configmap binary data is updated", func() {
				h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				cm1.BinaryData["binary_key1"] = []byte(modified)

				h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(h2).NotTo(Equal(h1))
			})

			It("returns a different hash when an allKeys child's secret data is updated", func() {
				h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				s1.Data["key1"] = []byte("modified")

				h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(h2).NotTo(Equal(h1))
			})

			It("returns the same hash when a child's metadata is updated", func() {
				h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				s1.Annotations = map[string]string{"new": "annotations"}

				h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(h2).To(Equal(h1))
			})
		})

		Context("single-field children", func() {
			var configMaps map[types.NamespacedName]*corev1.ConfigMap
			var secrets map[types.NamespacedName]*corev1.Secret
			var configMapsConfig []configMetadata
			var secretsConfig []configMetadata

			BeforeEach(func() {

				configMaps = map[types.NamespacedName]*corev1.ConfigMap{
					GetNamespacedNameFromObject(cm1): cm1,
					GetNamespacedNameFromObject(cm2): cm2,
				}

				secrets = map[types.NamespacedName]*corev1.Secret{
					GetNamespacedNameFromObject(s1): s1,
					GetNamespacedNameFromObject(s2): s2,
				}

				configMapsConfig = []configMetadata{
					{name: GetNamespacedNameFromObject(cm1), allKeys: false, keys: map[string]struct{}{
						"key1":        {},
						"binary_key1": {},
					},
					},
					{name: GetNamespacedNameFromObject(cm2), allKeys: true},
				}

				secretsConfig = []configMetadata{
					{name: GetNamespacedNameFromObject(s1), allKeys: false, keys: map[string]struct{}{
						"key1": {},
					},
					},
					{name: GetNamespacedNameFromObject(s2), allKeys: true},
				}
			})

			It("returns a different hash when a single-field child's data is updated", func() {

				h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				cm1.Data["key1"] = modified

				h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				cm1.BinaryData["binary_key1"] = []byte(modified)

				h3, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				s1.Data["key1"] = []byte("modified")

				h4, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(h2).NotTo(Equal(h1))
				Expect(h3).NotTo(Equal(h1))
				Expect(h3).NotTo(Equal(h2))
				Expect(h4).NotTo(Equal(h1))
				Expect(h4).NotTo(Equal(h2))
				Expect(h4).NotTo(Equal(h3))
			})

			It("returns the same hash when a single-field child's data is updated but not for that field", func() {

				h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				cm1.Data["key3"] = modified
				cm1.BinaryData["binary_key3"] = []byte("modified")
				s1.Data["key3"] = []byte("modified")

				h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig, secretsConfig)
				Expect(err).NotTo(HaveOccurred())

				Expect(h2).To(Equal(h1))
			})
		})

		It("returns the same hash independent of child ordering", func() {
			configMaps := map[types.NamespacedName]*corev1.ConfigMap{
				GetNamespacedNameFromObject(cm1): cm1,
				GetNamespacedNameFromObject(cm2): cm2,
				GetNamespacedNameFromObject(cm3): cm3,
			}

			secrets := map[types.NamespacedName]*corev1.Secret{
				GetNamespacedNameFromObject(s1): s1,
				GetNamespacedNameFromObject(s2): s2,
				GetNamespacedNameFromObject(s3): s3,
			}

			configMapsConfig1 := []configMetadata{
				{name: GetNamespacedNameFromObject(cm1), allKeys: true},
				{name: GetNamespacedNameFromObject(cm2), allKeys: true},
				{name: GetNamespacedNameFromObject(cm3), allKeys: false, keys: map[string]struct{}{
					"key1": {},
					"key2": {},
				},
				},
			}

			secretsConfig1 := []configMetadata{
				{name: GetNamespacedNameFromObject(s1), allKeys: true},
				{name: GetNamespacedNameFromObject(s2), allKeys: true},
				{name: GetNamespacedNameFromObject(s3), allKeys: false, keys: map[string]struct{}{
					"key1": {},
					"key2": {},
				},
				},
			}

			configMapsConfig2 := []configMetadata{
				{name: GetNamespacedNameFromObject(cm1), allKeys: true},
				{name: GetNamespacedNameFromObject(cm2), allKeys: true},
				{name: GetNamespacedNameFromObject(cm3), allKeys: false, keys: map[string]struct{}{
					"key1": {},
					"key2": {},
				},
				},
				{name: GetNamespacedNameFromObject(cm1), allKeys: true},
			}

			secretsConfig2 := []configMetadata{
				{name: GetNamespacedNameFromObject(s3), allKeys: false, keys: map[string]struct{}{
					"key2": {},
				},
				},
				{name: GetNamespacedNameFromObject(s2), allKeys: true},
				{name: GetNamespacedNameFromObject(s1), allKeys: true},
				{name: GetNamespacedNameFromObject(s3), allKeys: false, keys: map[string]struct{}{
					"key1": {},
				},
				},
			}

			h1, err := calculateConfigHash(configMaps, secrets, configMapsConfig1, secretsConfig1)
			Expect(err).NotTo(HaveOccurred())
			h2, err := calculateConfigHash(configMaps, secrets, configMapsConfig2, secretsConfig2)
			Expect(err).NotTo(HaveOccurred())

			Expect(h2).To(Equal(h1))
		})

	})

	Context("setConfigHash", func() {
		var deploymentObject *appsv1.Deployment

		BeforeEach(func() {
			deploymentObject = utils.ExampleDeployment.DeepCopy()
		})

		It("sets the hash annotation to the provided value", func() {
			setConfigHash(deploymentObject, "1234")

			podAnnotations := deploymentObject.Spec.Template.GetAnnotations()
			Expect(podAnnotations).NotTo(BeNil())

			hash, ok := podAnnotations[ConfigHashAnnotation]
			Expect(ok).To(BeTrue())
			Expect(hash).To(Equal("1234"))
		})

		It("leaves existing annotations in place", func() {
			// Add an annotation to the pod spec
			podAnnotations := deploymentObject.Spec.Template.GetAnnotations()
			if podAnnotations == nil {
				podAnnotations = make(map[string]string)
			}
			podAnnotations["existing"] = "annotation"
			deploymentObject.Spec.Template.SetAnnotations(podAnnotations)

			// Set the config hash
			setConfigHash(deploymentObject, "1234")

			// Check the existing annotation is still in place
			podAnnotations = deploymentObject.Spec.Template.GetAnnotations()
			Expect(podAnnotations).NotTo(BeNil())

			hash, ok := podAnnotations["existing"]
			Expect(ok).To(BeTrue())
			Expect(hash).To(Equal("annotation"))
		})
	})

})
