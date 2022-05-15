import { Component, OnInit } from "@angular/core";
import { Observable } from "rxjs";
import { Dataset } from "src/app/models/dataset";
import { DatasetsService } from "src/app/services/datasets.service";
import { DevicesService } from "src/app/services/devices.service";
import { DialogsService } from "src/app/services/dialogs.service";
import { GeneralService } from "src/app/services/general.service";

@Component({
  selector: "app-datasets",
  templateUrl: "./dataset-list.component.html",
  styleUrls: ["./dataset-list.component.scss"],
})
export class DatasetsListComponent implements OnInit {
  datasets$: Observable<Dataset[]>;
  datasets: Dataset[];
  loading: boolean = true;

  constructor(
    private datasetsService: DatasetsService,
    private dialogService: DialogsService
  ) {}

  ngOnInit(): void {
    this.datasetsService.LoadDatasets().subscribe((res) => {
      if (res) {
        this.datasets = res;
        this.loading = false;
      }
    });
  }
  newDatasetDialog() {
    this.dialogService.openNewDatasetDialog();
  }
}
