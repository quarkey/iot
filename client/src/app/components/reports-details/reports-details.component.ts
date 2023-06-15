import { Component, Input, OnInit } from '@angular/core';
import { Dataset } from 'src/app/models/dataset';
import { DatasetsService } from 'src/app/services/datasets.service';

@Component({
  selector: 'app-reports-details',
  templateUrl: './reports-details.component.html',
  styleUrls: ['./reports-details.component.scss'],
})
export class ReportsDetailsComponent implements OnInit {
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
