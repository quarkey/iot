import { Component, Input, OnInit } from "@angular/core";
import { Controller } from "src/app/models/controllers";

@Component({
  selector: "app-controller-card",
  templateUrl: "./controller-card.component.html",
  styleUrls: ["./controller-card.component.scss"],
})
export class ControllerCardComponent implements OnInit {
  @Input() citem: Controller;

  constructor() {}

  ngOnInit(): void {}
}
