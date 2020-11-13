// +build !ignore_autogenerated

// Copyright 2020 Oracle and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ndb) DeepCopyInto(out *Ndb) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ndb.
func (in *Ndb) DeepCopy() *Ndb {
	if in == nil {
		return nil
	}
	out := new(Ndb)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Ndb) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NdbList) DeepCopyInto(out *NdbList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Ndb, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NdbList.
func (in *NdbList) DeepCopy() *NdbList {
	if in == nil {
		return nil
	}
	out := new(NdbList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *NdbList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NdbMgmdSpec) DeepCopyInto(out *NdbMgmdSpec) {
	*out = *in
	if in.NodeCount != nil {
		in, out := &in.NodeCount, &out.NodeCount
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NdbMgmdSpec.
func (in *NdbMgmdSpec) DeepCopy() *NdbMgmdSpec {
	if in == nil {
		return nil
	}
	out := new(NdbMgmdSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NdbMysqldSpec) DeepCopyInto(out *NdbMysqldSpec) {
	*out = *in
	if in.NodeCount != nil {
		in, out := &in.NodeCount, &out.NodeCount
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NdbMysqldSpec.
func (in *NdbMysqldSpec) DeepCopy() *NdbMysqldSpec {
	if in == nil {
		return nil
	}
	out := new(NdbMysqldSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NdbNdbdSpec) DeepCopyInto(out *NdbNdbdSpec) {
	*out = *in
	if in.NoOfReplicas != nil {
		in, out := &in.NoOfReplicas, &out.NoOfReplicas
		*out = new(int32)
		**out = **in
	}
	if in.NodeCount != nil {
		in, out := &in.NodeCount, &out.NodeCount
		*out = new(int32)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NdbNdbdSpec.
func (in *NdbNdbdSpec) DeepCopy() *NdbNdbdSpec {
	if in == nil {
		return nil
	}
	out := new(NdbNdbdSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NdbSpec) DeepCopyInto(out *NdbSpec) {
	*out = *in
	in.Mgmd.DeepCopyInto(&out.Mgmd)
	in.Ndbd.DeepCopyInto(&out.Ndbd)
	in.Mysqld.DeepCopyInto(&out.Mysqld)
	if in.Config != nil {
		in, out := &in.Config, &out.Config
		*out = new(v1.LocalObjectReference)
		**out = **in
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NdbSpec.
func (in *NdbSpec) DeepCopy() *NdbSpec {
	if in == nil {
		return nil
	}
	out := new(NdbSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NdbStatus) DeepCopyInto(out *NdbStatus) {
	*out = *in
	in.LastUpdate.DeepCopyInto(&out.LastUpdate)
	if in.ReceivedConfigHash != nil {
		in, out := &in.ReceivedConfigHash, &out.ReceivedConfigHash
		*out = make([]byte, len(*in))
		copy(*out, *in)
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NdbStatus.
func (in *NdbStatus) DeepCopy() *NdbStatus {
	if in == nil {
		return nil
	}
	out := new(NdbStatus)
	in.DeepCopyInto(out)
	return out
}
