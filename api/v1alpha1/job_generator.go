package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	kubeconfigSecretKey            = "kubeconfig"
	dockerConfigJSONSecretKey      = "config.json"
	defaultDockerConfigSecretName  = "dockerconfig"
	defaultPyxisAPITokenSecretName = "pyxisapi"
)

// PreflightCheckJobGenerator generates a batchv1.Job for a
// given PreflightCheck
type PreflightCheckJobGenerator struct {
	pc PreflightCheck
}

// preflightImage returns the value for Spec.preflightImage in pc,
// or the default value of quay.io/opdev/preflight:stable.
func (g *PreflightCheckJobGenerator) preflightImage() string {
	img := "quay.io/opdev/preflight:stable"
	if g.pc.Spec.PreflightImage != nil {
		img = *g.pc.Spec.PreflightImage
	}

	return img
}

// envLogLevel returns the value for Spec.LogLevel in pc, or
// a default value of Info.
func (g *PreflightCheckJobGenerator) envLogLevel() corev1.EnvVar {
	if g.pc.Spec.LogLevel != nil {
		return corev1.EnvVar{
			Name:  "PFLT_LOGLEVEL",
			Value: *g.pc.Spec.LogLevel,
		}
	}

	return corev1.EnvVar{
		Name:  "PFLT_LOGLEVEL",
		Value: "Info", // defaulting to Info
	}
}

// envCertificationProjectID returns the value for
// Spec.CheckContainerOtions.CertificationProjectID in g.pc. Because this
// is optional, this returns a pointer.
//
// a nil value indicates that this is not set in g.pc.
func (g *PreflightCheckJobGenerator) envCertificationProjectID() *corev1.EnvVar {
	if g.pc.Spec.CheckOptions.ContainerOptions != nil {
		if g.pc.Spec.CheckOptions.ContainerOptions.CertificationProjectID != nil {
			return &corev1.EnvVar{
				Name:  "PFLT_CERTIFICATION_PROJECT_ID",
				Value: *g.pc.Spec.CheckOptions.ContainerOptions.CertificationProjectID,
			}
		}
	}

	return nil
}

func (g *PreflightCheckJobGenerator) envIndexImage() *corev1.EnvVar {
	if g.pc.Spec.CheckOptions.OperatorOptions != nil {
		return &corev1.EnvVar{
			Name:  "PFLT_INDEXIMAGE",
			Value: g.pc.Spec.CheckOptions.OperatorOptions.IndexImage,
		}
	}

	return nil
}

// envKubeconfig returns a kubeconfig environment variable that references
// a volume-mounted secret. If nil, it's assumed that we don't need this
// environment variable.
func (g *PreflightCheckJobGenerator) envKubeconfig() *corev1.EnvVar {
	isOptional := false

	if g.pc.Spec.CheckOptions.OperatorOptions != nil {
		return &corev1.EnvVar{
			Name: "KUBECONFIG",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: "kubeconfig",
					},
					Key:      kubeconfigSecretKey,
					Optional: &isOptional,
				},
			},
		}
	}

	return nil
}

// volumeKubeconfig returns the volume configuration for the preflight kubecofig
// based on the secret reference provided by the user.
func (g *PreflightCheckJobGenerator) volumeKubeconfig() *corev1.Volume {
	isOptional := false

	if g.pc.Spec.CheckOptions.OperatorOptions != nil {
		return &corev1.Volume{
			Name: "kubeconfig",
			VolumeSource: corev1.VolumeSource{
				Secret: &corev1.SecretVolumeSource{
					SecretName: g.pc.Spec.CheckOptions.OperatorOptions.KubeconfigSecretRef,
					Optional:   &isOptional,
				},
			},
		}
	}

	return nil
}

// volumeDockerConfigJSON returns the volume configuration for the DockerConfigJSON
// based on the secret reference provided by the user, or a default value.
func (g *PreflightCheckJobGenerator) volumeDockerConfigJSON() corev1.Volume {
	isOptional := true

	secretName := defaultDockerConfigSecretName
	if g.pc.Spec.DockerConfigSecretRef != nil {
		secretName = *g.pc.Spec.DockerConfigSecretRef
	}

	return corev1.Volume{
		Name: "dockerconfigjson",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
				Optional:   &isOptional,
			},
		},
	}
}

// volumePyxisAPIToken returns the volume configuration for the PyxisAPIToken
// based on the secret reference provided by the user, or a default value.
func (g *PreflightCheckJobGenerator) volumePyxisAPIToken() corev1.Volume {
	isOptional := true

	secretName := defaultPyxisAPITokenSecretName
	if g.pc.Spec.CheckOptions.ContainerOptions != nil {
		if g.pc.Spec.CheckOptions.ContainerOptions.PyxisAPITokenSecretRef != nil {
			secretName = *g.pc.Spec.CheckOptions.ContainerOptions.PyxisAPITokenSecretRef
		}
	}

	return corev1.Volume{
		Name: "pyxisapitoken",
		VolumeSource: corev1.VolumeSource{
			Secret: &corev1.SecretVolumeSource{
				SecretName: secretName,
				Optional:   &isOptional,
			},
		},
	}
}

// Volumes returns all volumes that would be configured for the job
// based on the PreflightCheck.Spec.
func (g *PreflightCheckJobGenerator) Volumes() []corev1.Volume {
	vols := []corev1.Volume{
		g.volumeDockerConfigJSON(),
		g.volumePyxisAPIToken(),
	}

	// Kubeconfig only applies to the operator check, so
	// only add it if that's the requested command.
	if v := g.volumeKubeconfig(); v != nil {
		vols = append(vols, *v)
	}

	return vols
}

// Env returns the environment variables that would be configured
// for the job based on the PreflightCheck.Spec.
func (g *PreflightCheckJobGenerator) Env() []corev1.EnvVar {
	envs := []corev1.EnvVar{
		g.envLogLevel(),
		g.envDockerConfigValue(),
	}

	// Add optional environments.
	if e := g.envCertificationProjectID(); e != nil {
		envs = append(envs, *e)
	}

	if e := g.envKubeconfig(); e != nil {
		envs = append(envs, *e)
	}

	if e := g.envIndexImage(); e != nil {
		envs = append(envs, *e)
	}

	return envs
}

// envDockerConfigValue will return the corev1.EnvVar for PFLT_DOCKERCONFIG.
func (g *PreflightCheckJobGenerator) envDockerConfigValue() corev1.EnvVar {
	isOptional := true
	return corev1.EnvVar{
		Name: "PFLT_DOCKERCONFIG",
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: &corev1.SecretKeySelector{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: g.dockerConfigSecretRef(),
				},
				Key:      dockerConfigJSONSecretKey,
				Optional: &isOptional, // TODO if it doesn't exist, do we still get the env with an empty string?
			},
		},
	}
}

// dockerConfigSecretRef returns a dockerConfigSecret name from the
// PreflightCheck.Spec if set, or a default value if not.
func (g *PreflightCheckJobGenerator) dockerConfigSecretRef() string {
	ref := defaultDockerConfigSecretName
	if g.pc.Spec.DockerConfigSecretRef != nil {
		ref = *g.pc.Spec.DockerConfigSecretRef
	}

	return ref
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

// TODO(komish): Wire up credentials so that the environment variables refer to
// the path on the filesystem where the mounted secret volumes exist.

// Generate returns a batchv1.Job based on the input g.PreflightCheck.
func (g *PreflightCheckJobGenerator) Generate() batchv1.Job {
	var backoffLimit int32 = 0

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
						{
							Name:  "preflight",
							Image: g.preflightImage(),
							Args:  []string{"check", g.checkCommand(), g.pc.Spec.Image},
							Env:   g.Env(),
						},
					},
					Volumes: g.Volumes(),
				},
			},
		},
	}
}
