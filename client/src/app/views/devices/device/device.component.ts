import { Component, OnInit } from "@angular/core";
import {
  FormBuilder,
  FormControl,
  FormGroup,
  Validators,
} from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { Device } from "src/app/models/device";
import { DevicesService } from "src/app/services/devices.service";

@Component({
  selector: "app-device",
  templateUrl: "./device.component.html",
  styleUrls: ["./device.component.scss"],
})
export class DeviceComponent implements OnInit {
  device: Device;
  form: FormGroup;

  constructor(
    private activeRoute: ActivatedRoute,
    private deviceService: DevicesService,
    private formBuilder: FormBuilder
  ) {}

  ngOnInit(): void {
    this.form = this.formBuilder.group({
      title: ["", Validators.required],
      description: ["", Validators.required],
    });
    this.activeRoute.params.subscribe((key) => {
      this.deviceService.LoadDevice(key.arduino_key).subscribe((res) => {
        if (res) {
          this.device = res;
          this.form.patchValue(res);
        }
      });
    });
  }
  updateDevice() {
    var dat: Device = {
      ...this.form.value,
      arduino_key: this.device.arduino_key,
    };

    this.deviceService.UpdateDevice(dat).subscribe((res) => {
      this.form.markAsPristine();
    });
  }
}
