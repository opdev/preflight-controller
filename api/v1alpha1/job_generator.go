package v1alpha1

import (
	"path"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kubeconfigSecretKey  = "kubeconfig"
	kubeconfigVolumeName = "kubeconfig"

	dockerConfigJSONSecretKey  = "config.json"
	dockerConfigJSONVolumeName = "dockerconfigjson"

	defaultVolumeMountPath = "/preflight"

	pyxisAPITokenSecretKey = "pyxisapitoken"
)

// PreflightCheckJobGenerator generates a batchv1.Job for a
// given PreflightCheck
type PreflightCheckJobGenerator struct {
	pc PreflightCheck
}

// Generate returns a batchv1.Job based on the input g.PreflightCheck.
func (g *PreflightCheckJobGenerator) Generate() batchv1.Job {
	var backoffLimit int32 = 0

	// the baseContainer we're working with. Needs modifications to args, volumeMounts, etc.
	container := &corev1.Container{
		Name:         "preflight",
		Image:        g.preflightImage(),
		Args:         []string{"check", g.checkCommand(), g.pc.Spec.Image},
		VolumeMounts: []corev1.VolumeMount{},
	}

	// base set of volumes we'll use in the PodSpecTemplate
	volumes := []corev1.Volume{}

	// configure the runtime environment based on the spec.
	container, volumes = g.configureDockerConfigJSON(container, volumes)
	container = g.configureLogLevel(container)

	// configure check container options if necessary.
	if g.pc.Spec.CheckOptions.ContainerOptions != nil {
		container = g.configurePyxisAPIToken(container)
		container = g.configureCertificationProjectID(container)
	}

	// configure check operator options if necessary.
	if g.pc.Spec.CheckOptions.OperatorOptions != nil {
		container = g.configureIndexImage(container)
	}

	return batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: g.pc.GetName() + "-",
			Namespace:    g.pc.GetNamespace(),
			Labels: map[string]string{
				"spawning-resource-name": g.pc.GetName(),
			},
		},
		Spec: batchv1.JobSpec{
			BackoffLimit: &backoffLimit, // only run once, then fail.
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					GenerateName: g.pc.GetName() + "-",
					Namespace:    g.pc.GetNamespace(),
					Labels: map[string]string{
						"spawning-resource-name": g.pc.GetName(),
					},
				},
				Spec: corev1.PodSpec{
					RestartPolicy: corev1.RestartPolicyNever,
					Containers: []corev1.Container{
						*container,
					},
					Volumes: volumes,
				},
			},
		},
	}
}

// preflightImage returns the value for Spec.preflightImage in pc,
// or the default value of quay.io/opdev/preflight:stable.
func (g *PreflightCheckJobGenerator) preflightImage() string {
	if g.pc.Spec.PreflightImage != nil {
		return *g.pc.Spec.PreflightImage
	}

	return "quay.io/opdev/preflight:stable"
}

// configureLogLevel add the user's log level preference to the pod, or defaults to Info
// if none is provided.
func (g *PreflightCheckJobGenerator) configureLogLevel(c *corev1.Container) *corev1.Container {
	loglevel := "Info"

	if g.pc.Spec.LogLevel != "" {
		loglevel = g.pc.Spec.LogLevel
	}

	c.Env = append(c.Env, corev1.EnvVar{
		Name:  "PFLT_LOGLEVEL",
		Value: loglevel,
	})

	return c
}

func (g *PreflightCheckJobGenerator) configureCertificationProjectID(c *corev1.Container) *corev1.Container {
	if g.pc.Spec.CheckOptions.ContainerOptions != nil {
		if g.pc.Spec.CheckOptions.ContainerOptions.CertificationProjectID != nil {
			c.Env = append(c.Env, corev1.EnvVar{
				Name:  "PFLT_CERTIFICATION_PROJECT_ID",
				Value: *g.pc.Spec.CheckOptions.ContainerOptions.CertificationProjectID,
			})
		}
	}

	return c
}

func (g *PreflightCheckJobGenerator) configureIndexImage(c *corev1.Container) *corev1.Container {
	if g.pc.Spec.CheckOptions.OperatorOptions != nil {
		c.Env = append(c.Env, corev1.EnvVar{
			Name:  "PFLT_INDEXIMAGE",
			Value: g.pc.Spec.CheckOptions.OperatorOptions.IndexImage,
		})
	}

	return c
}

// checkCommand resolves the command to provide to `preflight check` based on
// the PreflightCheck.Spec.CheckOptions. If the OperatorOptions are set, this will
// return the operator, and otherwise fall back to container.
func (g *PreflightCheckJobGenerator) checkCommand() string {
	cmd := "container"

	if g.pc.Spec.CheckOptions.OperatorOptions != nil {
		cmd = "operator"
	}

	return cmd
}

// configurePyxisAPIToken will add the pyxisAPIToken value to the environment from the secret ref.
func (g *PreflightCheckJobGenerator) configurePyxisAPIToken(c *corev1.Container) *corev1.Container {
	if g.pc.Spec.CheckOptions.ContainerOptions == nil {
		// The job being executed is not a container check.
		return c
	}

	if g.pc.Spec.CheckOptions.ContainerOptions.PyxisAPITokenSecretRef == "" {
		// The user didn't include a secret to check.
		return c
	}

	// Add the secret to the volume list.
	notOptional := false

	// Project the mountpoint as the environment variable.
	c.Env = append(c.Env, corev1.EnvVar{
		Name: "PFLT_PYXIS_API_TOKEN",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: g.pc.Spec.CheckOptions.ContainerOptions.PyxisAPITokenSecretRef,
				},
				Key:      pyxisAPITokenSecretKey,
				Optional: &notOptional,
			},
		},
	})

	return c
}

// configureKubeconfig will mount the provided kubeconfig secret to the v and configure the KUBECONFIG
// environment variable to c
func (g *PreflightCheckJobGenerator) configureKubeconfig(c *corev1.Container, v []corev1.Volume) (*corev1.Container, []corev1.Volume) {
	if g.pc.Spec.CheckOptions.OperatorOptions == nil {
		// The spec indicates that we're not running operator checks
		return c, v
	}

	if g.pc.Spec.CheckOptions.OperatorOptions.KubeconfigSecretRef == "" {
		// The user didn't include a secret to check.
		return c, v
	}

	// Add the secret to the volume list.
	notOptional := false
	v = append(v, corev1.Volume{
		Name: kubeconfigVolumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: g.pc.Spec.CheckOptions.OperatorOptions.KubeconfigSecretRef,
				Optional:   &notOptional,
			},
		},
	})

	// Mount the secret to the container.
	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      kubeconfigVolumeName,
		ReadOnly:  true,
		MountPath: defaultVolumeMountPath,
	})

	// Project the mountpoint as the environment variable.
	c.Env = append(c.Env, corev1.EnvVar{
		Name:  "KUBECONFIG",
		Value: path.Join(defaultVolumeMountPath, kubeconfigSecretKey),
	})

	return c, v
}

// configureDockerConfigJSON will add the dockerConfigSecretRef as a volume to v, as
// a volumeMount to c, and as the PFLT_DOCKERCONFIG environment variable to c.
func (g *PreflightCheckJobGenerator) configureDockerConfigJSON(c *corev1.Container, v []corev1.Volume) (*corev1.Container, []corev1.Volume) {
	if g.pc.Spec.DockerConfigSecretRef == "" {
		// The user didn't include a secret to check.
		return c, v
	}

	// Add the secret to the volume list.
	notOptional := false
	v = append(v, corev1.Volume{
		Name: dockerConfigJSONVolumeName,
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: g.pc.Spec.DockerConfigSecretRef,
				Optional:   &notOptional,
			},
		},
	})

	// Mount the secret to the container.
	c.VolumeMounts = append(c.VolumeMounts, corev1.VolumeMount{
		Name:      dockerConfigJSONVolumeName,
		ReadOnly:  true,
		MountPath: defaultVolumeMountPath},
	)

	// Project the mountpoint as the environment variable.
	c.Env = append(c.Env, corev1.EnvVar{
		Name:  "PFLT_DOCKERCONFIG",
		Value: path.Join(defaultVolumeMountPath, dockerConfigJSONSecretKey),
	})

	return c, v
}
