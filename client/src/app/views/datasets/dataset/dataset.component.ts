import { Component, OnInit } from "@angular/core";
import { FormBuilder, FormGroup, Validators } from "@angular/forms";
import { ActivatedRoute } from "@angular/router";
import { Dataset } from "src/app/models/dataset";
import { DatasetsService } from "src/app/services/datasets.service";
import { DevicesService } from "src/app/services/devices.service";

@Component({
  selector: "app-dataset",
  templateUrl: "./dataset.component.html",
  styleUrls: ["./dataset.component.scss"],
})
export class DatasetComponent implements OnInit {
  dataset: Dataset;
  constructor(
    private activeRoute: ActivatedRoute,
    private datasetService: DatasetsService
  ) {}

  ngOnInit(): void {
    this.activeRoute.params.subscribe((key) => {
      this.datasetService.LoadDataset(key.reference).subscribe((res) => {
        this.dataset = res;
      });
    });
  }
}
