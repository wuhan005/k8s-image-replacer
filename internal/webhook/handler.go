// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package webhook

import (
	"github.com/slok/kubewebhook/v2/pkg/webhook"
	kwhmutating "github.com/slok/kubewebhook/v2/pkg/webhook/mutating"
	corev1 "k8s.io/api/core/v1"
)

func NewHandler() (webhook.Webhook, error) {
	imageSwapper := NewImageReplacer()
	mutator := kwhmutating.MutatorFunc(imageSwapper.Mutate)
	webHookConfig := kwhmutating.WebhookConfig{
		ID:      "k8s-image-replacer",
		Obj:     &corev1.Pod{},
		Mutator: mutator,
	}
	return kwhmutating.NewWebhook(webHookConfig)
}
