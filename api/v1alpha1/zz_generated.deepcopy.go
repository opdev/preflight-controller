//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2022.

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CheckOptions) DeepCopyInto(out *CheckOptions) {
	*out = *in
	if in.ContainerOptions != nil {
		in, out := &in.ContainerOptions, &out.ContainerOptions
		*out = new(ContainerOptions)
		(*in).DeepCopyInto(*out)
	}
	if in.OperatorOptions != nil {
		in, out := &in.OperatorOptions, &out.OperatorOptions
		*out = new(OperatorOptions)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CheckOptions.
func (in *CheckOptions) DeepCopy() *CheckOptions {
	if in == nil {
		return nil
	}
	out := new(CheckOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContainerOptions) DeepCopyInto(out *ContainerOptions) {
	*out = *in
	if in.CertificationProjectID != nil {
		in, out := &in.CertificationProjectID, &out.CertificationProjectID
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContainerOptions.
func (in *ContainerOptions) DeepCopy() *ContainerOptions {
	if in == nil {
		return nil
	}
	out := new(ContainerOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OperatorOptions) DeepCopyInto(out *OperatorOptions) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OperatorOptions.
func (in *OperatorOptions) DeepCopy() *OperatorOptions {
	if in == nil {
		return nil
	}
	out := new(OperatorOptions)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreflightCheck) DeepCopyInto(out *PreflightCheck) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreflightCheck.
func (in *PreflightCheck) DeepCopy() *PreflightCheck {
	if in == nil {
		return nil
	}
	out := new(PreflightCheck)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PreflightCheck) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreflightCheckJobGenerator) DeepCopyInto(out *PreflightCheckJobGenerator) {
	*out = *in
	in.pc.DeepCopyInto(&out.pc)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreflightCheckJobGenerator.
func (in *PreflightCheckJobGenerator) DeepCopy() *PreflightCheckJobGenerator {
	if in == nil {
		return nil
	}
	out := new(PreflightCheckJobGenerator)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreflightCheckList) DeepCopyInto(out *PreflightCheckList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]PreflightCheck, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreflightCheckList.
func (in *PreflightCheckList) DeepCopy() *PreflightCheckList {
	if in == nil {
		return nil
	}
	out := new(PreflightCheckList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *PreflightCheckList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreflightCheckSpec) DeepCopyInto(out *PreflightCheckSpec) {
	*out = *in
	if in.PreflightImage != nil {
		in, out := &in.PreflightImage, &out.PreflightImage
		*out = new(string)
		**out = **in
	}
	in.CheckOptions.DeepCopyInto(&out.CheckOptions)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreflightCheckSpec.
func (in *PreflightCheckSpec) DeepCopy() *PreflightCheckSpec {
	if in == nil {
		return nil
	}
	out := new(PreflightCheckSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PreflightCheckStatus) DeepCopyInto(out *PreflightCheckStatus) {
	*out = *in
	if in.Jobs != nil {
		in, out := &in.Jobs, &out.Jobs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PreflightCheckStatus.
func (in *PreflightCheckStatus) DeepCopy() *PreflightCheckStatus {
	if in == nil {
		return nil
	}
	out := new(PreflightCheckStatus)
	in.DeepCopyInto(out)
	return out
}
