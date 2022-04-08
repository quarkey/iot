import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { MatDialogRef } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { DevicesService } from 'src/app/services/devices.service';

@Component({
  selector: 'app-new-device',
  templateUrl: './new-device.component.html',
  styleUrls: ['./new-device.component.scss'],
})
export class NewDeviceDialogComponent implements OnInit {
  constructor(
    private formBuilder: FormBuilder,
    private deviceService: DevicesService,
    private dialogRef: MatDialogRef<NewDeviceDialogComponent>,
    private router: Router
  ) {}
  form: FormGroup;
  ngOnInit(): void {
    this.form = this.formBuilder.group({
      title: ['', Validators.required],
      description: ['', Validators.required],
    });
  }
  addNewDevice() {
    const title = this.form.get('title').value;
    const desc = this.form.get('description').value;
    this.deviceService.AddNewDevice(title, desc).subscribe((res) => {
      this.router.navigate([`/sensors/${res.arduino_key}`]);
      this.dialogRef.close();
    });
  }
}
