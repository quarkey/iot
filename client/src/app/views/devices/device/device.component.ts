import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { Device } from 'src/app/models/device';
import { DevicesService } from 'src/app/services/devices.service';

@Component({
  selector: 'app-device',
  templateUrl: './device.component.html',
  styleUrls: ['./device.component.scss'],
})
export class DeviceComponent implements OnInit {
  device: Device;
  constructor(
    private activeRoute: ActivatedRoute,
    private deviceService: DevicesService
  ) {}

  ngOnInit(): void {
    this.activeRoute.params.subscribe((key) => {
      this.deviceService.LoadDevice(key.arduino_key).subscribe((res) => {
        this.device = res;
      });
    });
  }
}
