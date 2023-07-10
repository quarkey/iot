import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
// import { MatLegacyDialogRef as MatDialogRef } from "@angular/material/legacy-dialog";
import { MatDialogRef } from '@angular/material/dialog';

import { Router } from '@angular/router';
import { Device } from 'src/app/models/device';
import { DatasetsService } from 'src/app/services/datasets.service';
import { DevicesService } from 'src/app/services/devices.service';

@Component({
  selector: 'app-new-dataset-dialog',
  templateUrl: './new-dataset.component.html',
  styleUrls: ['./new-dataset.component.scss'],
})
export class NewDatasetDialogComponent implements OnInit {
  form: FormGroup;
  deviceList: Device[];
  ref: string;

  constructor(
    private formBuilder: FormBuilder,
    private dialogRef: MatDialogRef<NewDatasetDialogComponent>,
    private datasetService: DatasetsService,
    private deviceService: DevicesService,
    private router: Router
  ) {}

  ngOnInit(): void {
    this.deviceService.LoadDevices().subscribe((res) => {
      if (res) {
        this.deviceList = res;
      }
    });
    this.form = this.formBuilder.group({
      title: ['', Validators.required],
      description: ['', Validators.required],
      sensor_id: ['', Validators.required],
    });
    this.form.get('sensor_id').valueChanges.subscribe((sensor_id) => {
      this.deviceList.filter((x) => {
        if (sensor_id == x.id) {
          this.ref = x.arduino_key;
        }
      });
    });
  }
  addNewDataset() {
    var newDataset = { ...this.form.value, reference: this.ref } as Device;
    this.datasetService.newDataset(newDataset).subscribe((res) => {
      if (res) {
        this.router.navigate([`/datasets/${this.ref}`]);
        this.dialogRef.close();
      }
    });
  }
}
