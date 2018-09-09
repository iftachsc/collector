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
func getScsiLunsRoutine(wg *sync.WaitGroup, c *govmomi.Client, ctx context.Context,
	resultSink chan map[string]types.ScsiLun, errorSink chan error) {

	defer wg.Done()
	scsiLuns, err := vmware.GetScsiLunDisks(c, ctx)
	if err != nil {
		errorSink <- err
		return
	}
	resultSink <- scsiLuns
}

func getVmsRoutine(wg *sync.WaitGroup, c *govmomi.Client, ctx context.Context,
	resultSink chan []contracts.VsphereVM, errorSink chan error) {

	defer wg.Done()
	vms, err := vmware.GetVM(c, ctx)
	if err != nil {
		errorSink <- err
		return
	}
	resultSink <- vms
}

func getZadataLuns(wg *sync.WaitGroup, c *govmomi.Client, ctx context.Context,
	resultSink chan []contracts.VsphereVM, errorSink chan error) {
}

func getStorageVolumes(wg *sync.WaitGroup, sf contracts.StoageFiler, ctx context.Context,
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
