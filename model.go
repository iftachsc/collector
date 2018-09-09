package main

import (
	"github.com/iftachsc/contracts"
)

// type CollectedDisk struct {
// 	Disk vmware.VsphereDisk
//
// }
type CollectedVm struct {
	contracts.VsphereVM
	//disks []CollectedDisk
}

type VirtualMachineWithBackend struct {
	contracts.VsphereVM
}
type CollectionResult struct {
	VirtualMachineWithBackend
}

type Collection struct {
	//ID     bson.		    `bson:"_id" json:"id"`
	LocationUuid string           `bson:"location_uuid" json:"location_uuid"`
	Status       string           `bson:"name" json:"name"`
	Error        error            `bson:"cover_image" json:"cover_image"`
	Result       CollectionResult `bson:"description" json:"description"`
}
