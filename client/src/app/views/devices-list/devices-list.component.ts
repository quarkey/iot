import { Component, OnInit } from "@angular/core";
import { Observable } from "rxjs";
import { Device } from "src/app/models/device";
import { DevicesService } from "src/app/services/devices.service";
import { DialogsService } from "src/app/services/dialogs.service";

@Component({
  selector: "app-devices",
  templateUrl: "./devices-list.component.html",
  styleUrls: ["./devices-list.component.scss"],
})
export class DevicesListComponent implements OnInit {
  devices$: Observable<Device[]>;
  devices: Device[];
  loading: boolean = true;

  constructor(
    private deviceService: DevicesService,
    private dialogService: DialogsService
  ) {}

  ngOnInit(): void {
    this.deviceService.LoadDevices().subscribe((res) => {
      console.log("fetching devices");
      if (res) {
        this.devices = res;
        this.loading = false;
      }
    });
  }
  newDeviceDialog() {
    this.dialogService.openNewDeviceDialog();
  }
}
