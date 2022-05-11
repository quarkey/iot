import { Component, Input, OnInit } from "@angular/core";
import { Device } from "src/app/models/device";

@Component({
  selector: "app-device-card",
  templateUrl: "./device-card.component.html",
  styleUrls: ["./device-card.component.scss"],
})
export class DeviceCardComponent implements OnInit {
  @Input() device: Device;
  constructor() {}

  ngOnInit(): void {}
}
