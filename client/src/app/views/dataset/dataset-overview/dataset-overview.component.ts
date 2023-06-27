import { Component, Input, OnInit } from '@angular/core';
import { Dataset } from 'src/app/models/dataset';
import { DatasetsService } from 'src/app/services/datasets.service';

@Component({
  selector: 'app-dataset-overview',
  templateUrl: './dataset-overview.component.html',
  styleUrls: ['./dataset-overview.component.scss'],
})
export class DatasetOverviewComponent implements OnInit {
  constructor(private datasetService: DatasetsService) {}
  @Input() dataset: Dataset;
  report: any;
  ngOnInit(): void {
    this.datasetService.GetTemperatureReport(this.dataset.reference, 'd1c0').subscribe((res) => {
      if (res) {
        this.report = res;
      }
    });
  }
}
