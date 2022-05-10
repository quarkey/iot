import { Component, Input, OnInit } from "@angular/core";
import { Dataset } from "src/app/models/dataset";

@Component({
  selector: "app-dataset-card",
  templateUrl: "./dataset-card.component.html",
  styleUrls: ["./dataset-card.component.scss"],
})
export class DatasetCardComponent implements OnInit {
  @Input() dataset: Dataset;
  constructor() {}

  ngOnInit(): void {}
}
