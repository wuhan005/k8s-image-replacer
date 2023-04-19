// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webhook

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	kwhmodel "github.com/slok/kubewebhook/v2/pkg/model"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/wuhan005/k8s-image-replacer/internal/conf"
)

type ImageReplacer struct {
}

func NewImageReplacer() *ImageReplacer {
	return &ImageReplacer{}
}

// Mutate implements the Mutator interface.
func (p *ImageReplacer) Mutate(ctx context.Context, _ *kwhmodel.AdmissionReview, object metav1.Object) (*kwhmutating.MutatorResult, error) {
	pod, ok := object.(*corev1.Pod)
	if !ok {
		return &kwhmutating.MutatorResult{}, nil
	}

	logrus.WithContext(ctx).WithField("pod_name", pod.Name).Info("Mutating pod")

	for i, container := range pod.Spec.InitContainers {
		imageName := container.Image
		newImage := replaceImage(imageName)

		logrus.WithContext(ctx).WithField("old_image", imageName).WithField("new_image", newImage).Info("Replacing image")
		pod.Spec.InitContainers[i].Image = newImage
	}

	for i, container := range pod.Spec.Containers {
		imageName := container.Image
		newImage := replaceImage(imageName)

		logrus.WithContext(ctx).WithField("old_image", imageName).WithField("new_image", newImage).Info("Replacing image")
		pod.Spec.Containers[i].Image = newImage
	}

	return &kwhmutating.MutatorResult{MutatedObject: pod}, nil
}

const dockerhubRegistryPrefix = "registry-1.docker.io/"

func replaceImage(image string) string {
	splitCount := strings.Count(image, "/")
	if splitCount < 2 && conf.ImageReplacer.DockerRegistry == "" {
		return image
	}

	// Check if it is a docker hub image, e.g. nginx, wuhan005/Elaina.
	switch splitCount {
	case 0:
		image = "library/" + image
		fallthrough
	case 1:
		image = dockerhubRegistryPrefix + image
	}

	registry := strings.Split(image, "/")[0]

	if strings.HasPrefix(image, dockerhubRegistryPrefix) {
		return strings.Replace(image, registry, conf.ImageReplacer.DockerRegistry, 1)
	}

	for fromRegistry, toRegistry := range conf.ImageReplacer.ReplacePolicy {
		replacePrefix := fromRegistry + "/"
		if strings.HasPrefix(image, replacePrefix) {
			image = strings.Replace(image, fromRegistry, toRegistry, 1)
			break
		}
	}

	return image
}
