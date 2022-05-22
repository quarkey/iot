import { Component, Input, OnInit } from "@angular/core";
import { Controller } from "src/app/models/controllers";

@Component({
  selector: "app-controller-details",
  templateUrl: "./controller-details.component.html",
  styleUrls: ["./controller-details.component.scss"],
})
export class ControllerDetailsComponent implements OnInit {
  @Input() citem: Controller;

  constructor() {}

  ngOnInit(): void {}
}
