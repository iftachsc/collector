package main

import (
	"context"
	"strings"
	"sync"

	"github.com/iftachsc/contracts"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"
)

type collection struct {
	VMware *govmomi.Client
}

func createCollections(db *mongo.Database) {
	return
}

func startCollection(uuid string, dbClient *mongo.Client, db *mongo.Database, c *govmomi.Client, ctx context.Context, locationUuid string) (Shuki, error) {

	collectionSource := contracts.Location{
		contracts.IaasHost
		StoageFilers}

	var wg = new(sync.WaitGroup)

	vmsChan := make(chan []contracts.VsphereVM, 1)
	scsiLunsChan := make(chan map[string]types.ScsiLun, 1)
	//storageVolumesChan := make(chan contracts.StorageVolume)

	errors := make(chan error, 2)

	wg.Add(2)

	go getScsiLunsRoutine(wg, c, ctx, scsiLunsChan, errors)
	go getVmsRoutine(wg, c, ctx, vmsChan, errors)

	//println("TIME VMS:", int(time.Now().Sub(before).Seconds()), "seconds")
	go func() {
		wg.Wait()
		close(scsiLunsChan)
		close(vmsChan)
		close(errors)
	}()

	for err := range errors {
		// here error happend u could exit your caller function
		println(err.Error())
		//return
	}
	vms := <-vmsChan
	scsiLuns := <-scsiLunsChan

	//set rdm mappings to storage system objects
	for _, vm := range vms {
		for _, disk := range vm.Disks {
			if disk.IsRdm() {
				rdm := disk.(*contracts.VsphereRdmDisk)
				lun := scsiLuns[rdm.LunUuid]
				rdm.SetVendor(strings.TrimSpace(lun.Vendor))
				
			}
		}
	}

	return Shuki{vms, nil}, nil
}

type Shuki struct {
	Vms  []contracts.VsphereVM `json:"vms"`
	Luns []types.ScsiLun       `json:"luns"`
}
