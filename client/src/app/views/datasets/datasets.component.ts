import { Component, OnInit } from "@angular/core";
import { Observable } from "rxjs";
import { Dataset } from "src/app/models/dataset";
import { DatasetsService } from "src/app/services/datasets.service";
import { DevicesService } from "src/app/services/devices.service";
import { DialogsService } from "src/app/services/dialogs.service";

@Component({
  selector: "app-datasets",
  templateUrl: "./datasets.component.html",
  styleUrls: ["./datasets.component.scss"],
})
export class DatasetsComponent implements OnInit {
  datasets$: Observable<Dataset[]>;
  datasets: Dataset[];

  constructor(
    private datasetsService: DatasetsService,
    private dialogService: DialogsService
  ) {}

  ngOnInit(): void {
    this.datasetsService.LoadDatasets().subscribe((res) => {
      this.datasets = res;
    });
  }
  newDatasetDialog() {
    this.dialogService.openNewDatasetDialog();
  }
}
