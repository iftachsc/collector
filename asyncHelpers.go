package main

import (
	"context"
	"sync"

	"github.com/iftachsc/contracts"
	"github.com/iftachsc/vmware"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vim25/types"
)

//All functions here receive a WaitGroup, a channel for results, and channel for errors
func getScsiLunsRoutine(ctx context.Context, wg *sync.WaitGroup, c *govmomi.Client,
	resultSink chan map[string]types.ScsiLun, errorSink chan error) {

	defer wg.Done()
	scsiLuns, err := vmware.GetScsiLunDisks(ctx, c)
	if err != nil {
		errorSink <- err
		return
	}
	resultSink <- uniqueLunsMapByUUID(scsiLuns)
}

func getVmsRoutine(ctx context.Context, wg *sync.WaitGroup, c *govmomi.Client,
	resultSink chan []contracts.VsphereVM, errorSink chan error) {

	defer wg.Done()
	vms, err := vmware.GetVM(c, ctx)
	if err != nil {
		errorSink <- err
		return
	}
	resultSink <- vms
}

func getZadaraLuns(ctx context.Context, wg *sync.WaitGroup, c *govmomi.Client,
	resultSink chan []contracts.VsphereVM, errorSink chan error) {
}

func getStorageVolumes(ctx context.Context, wg *sync.WaitGroup, sf contracts.StoageFiler,
	resultSink chan []contracts.StorageVolume, errorSink chan error) {
	defer wg.Done()
	volumes, err := sf.GetVolumes()

	if err != nil {
		errorSink <- err
		return
	}
	resultSink <- volumes
}

// func addVsphereTask(wg *sync.WaitGroup, c *govmomi.Client, ctx context.Context, ,
// 	f vsphereTask, errorSink chan error) {
// 	before := time.Now()
// 	defer wg.Done()
// 	result, err := f(c, ctx)
// 	if err != nil {
// 		errorSink <- err
// 		return
// 	}
// 	resultSink <- result
// 	println("Task took", int(time.Now().Sub(before).Seconds()), "seconds")
// }

//removes duplicate luns by the uuid property and return a map
func uniqueLunsMapByUUID(luns []types.ScsiLun) map[string]types.ScsiLun {
	//u := []types.ScsiLun{}
	m := make(map[string]types.ScsiLun)

	for _, lun := range luns {
		if _, ok := m[lun.Uuid]; !ok {
			m[lun.Uuid] = lun
		}
	}

	return m
}
