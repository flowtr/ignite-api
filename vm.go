package main

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/weaveworks/ignite/pkg/apis/ignite"
	"github.com/weaveworks/ignite/pkg/providers"
	"github.com/weaveworks/libgitops/pkg/runtime"
	"net/http"
)

func GetVMS(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vmList, err := providers.Client.VMs().List()
	if err != nil {
		err := json.NewEncoder(w).Encode(vmList)
		if err != nil {
			return
		}
	}
}

func GetVM(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	vm, err := providers.Client.VMs().Get(runtime.UID(params["id"]))
	if err == nil {
		err := json.NewEncoder(w).Encode(vm)
		if err != nil {
			return
		}
	} else {
		return
	}
}

func withSpec(obj *ignite.VM, spec ignite.VMSpec) *ignite.VM {
	obj.Spec = spec
	return obj
}

func CreateVM(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var data ignite.VMSpec
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		return
	}
	vm := withSpec(providers.Client.VMs().New(), data)
	err = providers.Client.VMs().Set(vm)
	if err != nil {
		err := json.NewEncoder(w).Encode(err)
		if err != nil {
			return
		}
	} else {
		err := json.NewEncoder(w).Encode(vm)
		if err != nil {
			return
		}
	}
}
